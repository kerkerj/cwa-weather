# cwa-weather — Project Plan

## Overview

A Go CLI tool + library for Taiwan's Central Weather Administration (CWA) Open Data API.
Provides structured weather data for any township in Taiwan, designed to be used by both humans and AI agents.

**GitHub repo**: `github.com/kerkerj/cwa-weather`

## Background & Motivation

- Taiwan's CWA provides a free Open Data API with township-level forecasts, real-time observations, typhoon tracking, and weather alerts
- Existing tools are either outdated Go libraries (`go-cwb`, last updated ~2023), MCP servers tied to Claude Desktop, or Home Assistant integrations
- **No standalone CLI exists** for querying Taiwan weather from the command line
- This fills the gap: a single binary, zero dependencies, usable by any AI agent framework or human

## Architecture

### Single repo, multiple use cases

```
github.com/kerkerj/cwa-weather/
├── cwa/                  # Go library (importable package)
│   ├── client.go         # HTTP client, auth, error handling
│   ├── client_test.go
│   ├── forecast.go       # Township forecast (F-D0047-*)
│   ├── forecast_test.go
│   ├── observe.go        # Real-time observation stations (O-A0003-001)
│   ├── observe_test.go
│   ├── typhoon.go        # Typhoon data (W-C0034-005) [Phase 2]
│   ├── alert.go          # Weather alerts (W-C0034-001) [Phase 2]
│   └── endpoints.go      # City → dataset ID mapping (22 cities)
├── cmd/
│   └── cwa-weather/
│       ├── main.go       # Entry point
│       ├── forecast.go   # cobra: forecast subcommand
│       └── observe.go    # cobra: observe subcommand
├── skill/                # Agent integration files
│   ├── SKILL.md          # OpenClaw agent instructions
│   └── AGENT.md          # Generic agent instructions (for Claude Code, Codex, etc.)
├── go.mod
├── go.sum
├── README.md
├── LICENSE               # MIT
├── .goreleaser.yml
└── .github/
    └── workflows/
        └── release.yml   # Tag → build binaries → GitHub Release
```

### Three ways to use

1. **As a CLI tool**: `cwa-weather forecast --city 新北市 --town 板橋區`
2. **As a Go library**: `import "github.com/kerkerj/cwa-weather/cwa"`
3. **As an agent skill**: Agent reads `SKILL.md` or `AGENT.md`, runs CLI commands

## CWA API Details

### Authentication

- API key from https://opendata.cwa.gov.tw (free, instant, no review)
- Passed via `CWA_API_KEY` environment variable
- CLI reads from env; library accepts key as constructor parameter

### Endpoints to implement

#### Phase 1 (MVP)

| Subcommand | CWA Endpoint | Description |
|---|---|---|
| `forecast` | `F-D0047-001` ~ `F-D0047-091` | Township-level forecast (hourly/3-hourly, 2-3 days) |
| `observe` | `O-A0003-001` | Real-time weather station observations |

#### Phase 2 (Future)

| Subcommand | CWA Endpoint | Description |
|---|---|---|
| `typhoon` | `W-C0034-005` | Typhoon tracking (path, wind speed, pressure) |
| `alert` | `W-C0034-001` | Weather alerts (rain, cold, wind warnings) |
| `overview` | `F-C0032-001` | 36-hour city-level forecast (quick overview) |

### City → Endpoint Mapping

Each city has its own dataset ID for township forecasts. Full mapping for all 22 cities:

```
臺北市: F-D0047-061 (3hr), F-D0047-063 (weekly)
新北市: F-D0047-069 (3hr), F-D0047-071 (weekly)
桃園市: F-D0047-005 (3hr), F-D0047-007 (weekly)
臺中市: F-D0047-073 (3hr), F-D0047-075 (weekly)
臺南市: F-D0047-077 (3hr), F-D0047-079 (weekly)
高雄市: F-D0047-065 (3hr), F-D0047-067 (weekly)
基隆市: F-D0047-049 (3hr), F-D0047-051 (weekly)
新竹縣: F-D0047-009 (3hr), F-D0047-011 (weekly)
新竹市: F-D0047-053 (3hr), F-D0047-055 (weekly)
苗栗縣: F-D0047-013 (3hr), F-D0047-015 (weekly)
彰化縣: F-D0047-017 (3hr), F-D0047-019 (weekly)
南投縣: F-D0047-021 (3hr), F-D0047-023 (weekly)
雲林縣: F-D0047-025 (3hr), F-D0047-027 (weekly)
嘉義縣: F-D0047-029 (3hr), F-D0047-031 (weekly)
嘉義市: F-D0047-057 (3hr), F-D0047-059 (weekly)
屏東縣: F-D0047-033 (3hr), F-D0047-035 (weekly)
宜蘭縣: F-D0047-001 (3hr), F-D0047-003 (weekly)
花蓮縣: F-D0047-037 (3hr), F-D0047-039 (weekly)
臺東縣: F-D0047-041 (3hr), F-D0047-043 (weekly)
澎湖縣: F-D0047-045 (3hr), F-D0047-047 (weekly)
金門縣: F-D0047-085 (3hr), F-D0047-087 (weekly)
連江縣: F-D0047-081 (3hr), F-D0047-083 (weekly)
```

Note: The CLI should also accept common aliases (e.g., `台北市` → `臺北市`, `台中市` → `臺中市`).

### CWA API Response Structure

The CWA API returns deeply nested JSON. Key insight: **each weather element uses a different key name** inside `ElementValue`.

Example response path for township forecast:
```
.records.Locations[0].Location[0].WeatherElement[].Time[].ElementValue[0]
```

Element value keys vary by element:
- 溫度 → `{ "Temperature": "19" }`
- 體感溫度 → `{ "ApparentTemperature": "20" }`
- 3小時降雨機率 → `{ "ProbabilityOfPrecipitation": "20" }`
- 天氣現象 → `{ "Weather": "陰", "WeatherCode": "07" }`
- 相對濕度 → `{ "RelativeHumidity": "85" }`
- 風速 → `{ "WindSpeed": "3" }`
- 風向 → `{ "WindDirection": "東北" }`
- 天氣預報綜合描述 → `{ "WeatherDescription": "..." }`

Some elements use `DataTime` (hourly), others use `StartTime`/`EndTime` (range-based).

**Available elements for township 3-hour forecast (F-D0047-069 etc.):**
- 溫度
- 露點溫度
- 相對濕度
- 體感溫度
- 舒適度指數
- 風速
- 風向
- 3小時降雨機率
- 天氣現象
- 天氣預報綜合描述

### Observation Station Response Structure

Endpoint: `O-A0003-001` (automatic weather stations)

```
.records.Station[] → {
  StationName, StationId,
  GeoInfo: { CountyName, TownName, ... },
  WeatherElement: { AirTemperature, Weather, RelativeHumidity, WindSpeed, ... }
}
```

Can filter by `StationName` or iterate by `CountyName`.

## CLI Design

### Philosophy: Thin Wrapper (Level 1)

- **Zero information loss**: Output ALL elements from the API response
- CLI only handles: auth, endpoint routing, parameter validation, JSON flattening
- CLI does NOT: pick fields, format for humans, translate
- Output is always JSON, designed for agent consumption
- The agent (or human with `jq`) decides what fields to use

### Commands

```bash
# Township forecast (Phase 1)
cwa-weather forecast --city "新北市" --town "板橋區"
cwa-weather forecast --city "臺北市"                  # All towns in city
cwa-weather forecast --city "台北市" --town "中正區"    # 台→臺 auto-converted

# Real-time observation (Phase 1)
cwa-weather observe --city "新北市"                    # All stations in city
cwa-weather observe --station "淡水"                   # Specific station

# Future (Phase 2)
cwa-weather typhoon                                    # Current typhoon info (if any)
cwa-weather alert                                      # Current weather alerts (if any)

# Utility
cwa-weather --help
cwa-weather --version
cwa-weather cities                                     # List all supported cities
```

### Output Format

```json
{
  "success": true,
  "command": "forecast",
  "location": {
    "city": "新北市",
    "town": "板橋區"
  },
  "source": {
    "dataset": "F-D0047-069",
    "api": "CWA Open Data",
    "url": "https://opendata.cwa.gov.tw"
  },
  "generated_at": "2026-02-28T18:00:00+08:00",
  "elements": {
    "溫度": [
      { "time": "2026-02-28T18:00:00+08:00", "value": { "Temperature": "19" } },
      { "time": "2026-02-28T19:00:00+08:00", "value": { "Temperature": "19" } }
    ],
    "體感溫度": [...],
    "3小時降雨機率": [...],
    "天氣現象": [...],
    "相對濕度": [...],
    "風速": [...],
    "風向": [...],
    "天氣預報綜合描述": [...]
  }
}
```

Error output:
```json
{
  "success": false,
  "error": "city not found: 新竹",
  "hint": "Did you mean 新竹市 or 新竹縣? Run 'cwa-weather cities' to see all options."
}
```

## Go Library Design

```go
package cwa

// Client is the CWA API client
type Client struct { ... }

// NewClient creates a new CWA API client
func NewClient(apiKey string, opts ...Option) *Client

// Forecast returns township-level weather forecast
// If town is empty, returns all towns in the city
func (c *Client) Forecast(ctx context.Context, city, town string) (*ForecastResult, error)

// Observe returns real-time observation data
// Either by city (all stations) or by station name
func (c *Client) Observe(ctx context.Context, opts ObserveOptions) (*ObserveResult, error)

// Cities returns all supported city names
func (c *Client) Cities() []string
```

## Dependencies

Minimal:
- `github.com/spf13/cobra` — CLI framework
- Standard library only for everything else (`net/http`, `encoding/json`)

## Release & Distribution

### goreleaser

Build for:
- `darwin/arm64` (macOS Apple Silicon)
- `darwin/amd64` (macOS Intel)
- `linux/arm64` (Raspberry Pi, ARM servers)
- `linux/amd64` (x86 servers, VPS)

### Installation methods

```bash
# 1. Go install
go install github.com/kerkerj/cwa-weather/cmd/cwa-weather@latest

# 2. Download binary from GitHub Releases
curl -L https://github.com/kerkerj/cwa-weather/releases/latest/download/cwa-weather_darwin_arm64 -o /usr/local/bin/cwa-weather

# 3. Homebrew (future)
brew install kerkerj/tap/cwa-weather
```

## Agent Skill Integration

### SKILL.md (OpenClaw format)

Placed in `skill/SKILL.md`, registered on ClawHub.
Tells OpenClaw agent:
- When to use which subcommand
- How to interpret the JSON output
- Common usage patterns (daily forecast, typhoon check, etc.)

### AGENT.md (Generic, for any agent)

Placed in `skill/AGENT.md`.
Framework-agnostic instructions that any LLM agent can follow.
Same content as SKILL.md but without OpenClaw-specific metadata.

## Implementation Phases

### Phase 1: MVP (this PR)
- [ ] Go module init
- [ ] `cwa/client.go` — HTTP client + auth + error handling
- [ ] `cwa/endpoints.go` — City → dataset ID mapping (22 cities, with 台→臺 alias)
- [ ] `cwa/forecast.go` — Township forecast
- [ ] `cwa/observe.go` — Observation stations
- [ ] `cmd/cwa-weather/` — CLI with `forecast`, `observe`, `cities` subcommands
- [ ] `README.md` — Installation, usage, examples
- [ ] `LICENSE` — MIT
- [ ] `.goreleaser.yml` + `.github/workflows/release.yml`
- [ ] `skill/SKILL.md` + `skill/AGENT.md`
- [ ] Tests for library functions

### Phase 2: Extended data
- [ ] `cwa/typhoon.go` + `typhoon` subcommand
- [ ] `cwa/alert.go` + `alert` subcommand
- [ ] `overview` subcommand (36-hour city-level)

### Phase 3: Polish
- [ ] Homebrew tap
- [ ] More comprehensive tests
- [ ] CI/CD pipeline (lint, test on PR)

## Notes for Implementation

1. **CWA API uses traditional Chinese city names** (臺 not 台). The CLI must handle both.
2. **Some elements use `DataTime`, others use `StartTime`/`EndTime`**. The flattened output should normalize this.
3. **Observation stations don't have one per township**. The closest station may be in a neighboring area. When querying by city, return all stations in that city and let the consumer pick.
4. **API rate limits**: CWA doesn't document specific rate limits for general members, but be reasonable. The CLI should not retry aggressively.
5. **Error responses from CWA**: When API key is invalid or endpoint doesn't exist, CWA returns `{"success": "false"}` with a message. Handle gracefully.

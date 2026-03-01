# cwa-weather

[![Go Report Card](https://goreportcard.com/badge/github.com/kerkerj/cwa-weather)](https://goreportcard.com/report/github.com/kerkerj/cwa-weather)
[![Go Reference](https://pkg.go.dev/badge/github.com/kerkerj/cwa-weather.svg)](https://pkg.go.dev/github.com/kerkerj/cwa-weather)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Built with Claude](https://img.shields.io/badge/Built%20with-Claude-blueviolet)](https://claude.ai)
[![Skill Scanner](https://img.shields.io/badge/Skill%20Scanner-SAFE-brightgreen)](https://github.com/cisco-ai-defense/skill-scanner)
![coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/kerkerj/6e8ceb0d7091bfa999901383f6486a29/raw/cwa-coverage.json)

[正體中文](README.zh-TW.md)

CWA Open Data API CLI and Go library for Taiwan weather data.

Queries forecast and real-time observation data from Taiwan's Central Weather Administration (CWA) Open Data platform.

## Installation

### Homebrew

```bash
brew tap kerkerj/tap
brew install cwa-weather
```

### Go

```bash
go install github.com/kerkerj/cwa-weather/cmd/cwa-weather@latest
```

Or download a binary from [GitHub Releases](https://github.com/kerkerj/cwa-weather/releases).

## Setup

Get an API key from https://opendata.cwa.gov.tw and export it:

```bash
export CWA_API_KEY=your-key
```

## Usage

### Forecast

```bash
# Township-level forecast
cwa-weather forecast --city 臺北市 --town 中正區

# City-level forecast (all towns)
cwa-weather forecast --city 台北市    # 台→臺 auto-converted

# Filter by weather elements
cwa-weather forecast --city 新北市 --town 板橋區 --element 溫度,天氣現象

# Show 3 days of forecast
cwa-weather forecast --city 臺北市 --town 中正區 --days 3

# Group by day
cwa-weather forecast --city 臺北市 --town 中正區 --days 3 --summary

# Filter by time range
cwa-weather forecast --city 臺北市 --time-from 2026-03-01T06:00:00

# Raw JSON output (pipe to jq for field extraction)
cwa-weather forecast --city 臺北市 --json
```

> **Tip**: Element names are defined by CWA API — run a command with `--json` and inspect the JSON keys to see all available names.

### Real-time Observation

```bash
# By city
cwa-weather observe --city 新北市

# By station name
cwa-weather observe --station 淡水

# Filter by weather elements
cwa-weather observe --city 新北市 --element AirTemperature,Weather
```

### 36-hour City-level Forecast (Overview)

```bash
# City-level 36-hour forecast
cwa-weather overview --city 臺北市

# Filter by weather elements
cwa-weather overview --city 臺北市 --element Wx,PoP

# Filter by time range
cwa-weather overview --city 臺北市 --time-from 2026-03-01T06:00:00 --time-to 2026-03-01T18:00:00
```

### Weather Alerts

```bash
# All active alerts
cwa-weather alert

# Alerts for a specific city
cwa-weather alert --city 臺北市
```

### Typhoon Tracking

```bash
# Current tropical cyclone info
cwa-weather typhoon

# Filter by tropical depression number and dataset
cwa-weather typhoon --td-no 03 --dataset ForecastData
```

### Marine Observation (海象)

```bash
# All marine stations
cwa-weather sea

# Specific station
cwa-weather sea --station 富貴角
```

### Generic Query

```bash
# Query any CWA endpoint by data ID
cwa-weather query F-D0047-069 -p LocationName=板橋區
```

### List Cities and Towns

```bash
# List all 22 cities
cwa-weather cities

# List towns in a city
cwa-weather cities --city 臺北市
```

## Library Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/kerkerj/cwa-weather/cwa"
)

func main() {
	client := cwa.NewClient("YOUR_API_KEY")
	ctx := context.Background()

	// Forecast
	forecast, _ := client.Forecast(ctx, "臺北市", "中正區")
	fmt.Println(forecast)

	// Forecast with element and time filters
	filtered, _ := client.Forecast(ctx, "臺北市", "", cwa.ForecastOption{
		Element:  "溫度,天氣現象",
		TimeFrom: "2026-03-01T06:00:00",
	})
	fmt.Println(filtered)

	// Observation
	obs, _ := client.Observe(ctx, cwa.ObserveByCity("新北市"))
	fmt.Println(obs)

	// Observation with element filter
	obsFiltered, _ := client.Observe(ctx, cwa.ObserveByCity("新北市"), cwa.ObserveWithElement("AirTemperature"))
	fmt.Println(obsFiltered)
}
```

## AI Agent Integration

This project includes skill files so AI agents can query Taiwan weather data via the CLI.

Requires `cwa-weather` CLI installed and `CWA_API_KEY` env var set.

### Claude Code (Plugin)

```
/plugin marketplace add kerkerj/cwa-weather
/plugin install cwa-weather@kerkerj-cwa-weather
```

Once installed, ask Claude things like:
- "What's the weather in Taipei?"
- "Any typhoons right now?"
- "Show me weather alerts"

### Other AI Agents

Point your agent to the skill files in [`plugins/cwa-weather/skills/`](plugins/cwa-weather/skills/) for command reference and usage instructions.

## Notes

- **Output**: Human-readable text by default. Use `--json` for raw JSON (pipe to `jq` for field extraction).
- **Supported cities**: All 22 Taiwan cities/counties.
- **台→臺 auto-conversion**: `台北市` is automatically converted to `臺北市` to match CWA's use of 正體字.

## License

MIT

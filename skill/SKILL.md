---
name: taiwan-weather
description: Query Taiwan weather data (forecast, observations) via CWA Open Data API
---

# Taiwan Weather (CWA Open Data)

## When to use
- User asks about Taiwan weather (forecast, temperature, rain, typhoon)
- Need real-time weather observation data for Taiwan locations

## Prerequisites
- `cwa-weather` CLI installed
- `CWA_API_KEY` environment variable set

## Commands

### Forecast (township-level)
`cwa-weather forecast --city "城市" --town "鄉鎮區"`
`cwa-weather forecast --city "城市" --element "溫度,天氣現象"`
`cwa-weather forecast --city "城市" --time-from "2026-03-01T06:00:00" --time-to "2026-03-01T18:00:00"`

**Forecast filter flags:**
- `--element` — filter weather elements (comma-separated, e.g. 溫度,天氣現象,降雨機率)
- `--time-from` — start time (yyyy-MM-ddThh:mm:ss)
- `--time-to` — end time (yyyy-MM-ddThh:mm:ss)

### Real-time Observation
`cwa-weather observe --city "城市"`
`cwa-weather observe --station "站名"`
`cwa-weather observe --city "城市" --element "AirTemperature,Weather"`

**Observe filter flags:**
- `--element` — filter weather elements (comma-separated, valid: Weather, VisibilityDescription, SunshineDuration, Now, WindDirection, WindSpeed, AirTemperature, RelativeHumidity, AirPressure, UVIndex, Max10MinAverage, GustInfo, DailyHigh, DailyLow)

### 36-hour City-level Forecast (Overview)
`cwa-weather overview --city "城市"`
`cwa-weather overview --city "城市" --element "Wx,PoP"`
`cwa-weather overview --city "城市" --time-from "2026-03-01T06:00:00" --time-to "2026-03-01T18:00:00"`

**Overview filter flags:**
- `--element` — filter weather elements (comma-separated, e.g. Wx,PoP,MinT,MaxT)
- `--time-from` — start time (yyyy-MM-ddThh:mm:ss)
- `--time-to` — end time (yyyy-MM-ddThh:mm:ss)

### Weather Alerts
`cwa-weather alert`
`cwa-weather alert --city "城市"`

**Alert filter flags:**
- `--city` — filter alerts by city name

### Typhoon Tracking
`cwa-weather typhoon`
`cwa-weather typhoon --td-no "03" --dataset "ForecastData"`

**Typhoon filter flags:**
- `--td-no` — tropical depression number
- `--dataset` — dataset type (e.g. ForecastData)

### Marine Observation (海象)
`cwa-weather sea`
`cwa-weather sea --station "站名"`

**Sea filter flags:**
- `--station` — filter by station name

### Generic Query (any CWA endpoint)
`cwa-weather query DATAID -p key=value`

### List Cities/Towns
`cwa-weather cities`
`cwa-weather cities --city "城市"`

## Output
All output is JSON. Use jq to extract specific fields.

## Notes
- City names use traditional Chinese (臺). Tool auto-converts 台→臺.
- Forecast returns all weather elements (temperature, rain probability, wind, humidity, etc.)
- Observation returns real-time station data (may not have one per township)

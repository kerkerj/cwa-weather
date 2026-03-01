---
name: cwa-weather
description: Use when user asks about Taiwan weather, forecast, observation, typhoon, alerts, or marine data. Requires cwa-weather CLI and CWA_API_KEY env var.
---

# cwa-weather

Taiwan CWA Open Data CLI. Output is always JSON — pipe to `jq` for extraction.

## Quick Reference

| Need | Command | Key flags |
|------|---------|-----------|
| Township forecast | `cwa-weather forecast --city 臺北市 --town 中正區` | `--element`, `--time-from`, `--time-to` |
| Observation | `cwa-weather observe --city 新北市` | `--station` (alt), `--element` |
| 36hr overview | `cwa-weather overview --city 臺北市` | `--element`, `--time-from`, `--time-to` |
| Alerts | `cwa-weather alert` | `--city` |
| Typhoon | `cwa-weather typhoon` | `--td-no`, `--dataset` |
| Marine | `cwa-weather sea` | `--station` |
| Any endpoint | `cwa-weather query DATAID -p key=value` | `-p` (repeatable) |
| List cities | `cwa-weather cities` | `--city` (show towns) |

Run `cwa-weather <command> --help` for all flags and details.

## Notes

- `台→臺` auto-converted (e.g. `台北市` → `臺北市`)
- `--element` values are defined by CWA API (not user's choice of language):
  - forecast/overview: `溫度`, `天氣現象`, `降雨機率` etc.
  - observe: `AirTemperature`, `Weather`, `WindSpeed` etc.
- To discover available elements: run command without `--element`, inspect the JSON keys
- `--element` accepts comma-separated values

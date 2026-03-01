# cwa-weather — Agent Instructions

Taiwan CWA Open Data CLI. Requires `CWA_API_KEY` env var. Output is always JSON.

## Command Selection

| User intent | Command |
|-------------|---------|
| Weather for a specific town (鄉鎮區) | `forecast --city X --town Y` |
| Weather for a city (all towns) | `forecast --city X` |
| Current conditions / temperature / rain now | `observe --city X` or `observe --station Y` |
| Quick city-level forecast (36hr) | `overview --city X` |
| Weather warnings / 特報 | `alert` or `alert --city X` |
| Typhoon info / 颱風 | `typhoon` |
| Sea conditions / 海象 / waves | `sea` or `sea --station X` |
| Anything else (83 CWA endpoints) | `query DATAID -p key=value` |
| List available cities | `cities` |

## Reduce output size

Use `--element` to filter (reduces ~80% tokens). Element names are CWA API-defined:
- forecast/overview: `--element 溫度,天氣現象,降雨機率`
- observe: `--element AirTemperature,Weather`
- Use `--time-from` / `--time-to` to narrow time range
- To discover available elements: run without `--element`, inspect JSON keys

## Key behaviors

- `台→臺` auto-converted (正體字) — accept either form from user
- Run `cwa-weather <command> --help` for full flag details
- All commands return full CWA JSON — use `jq` to extract specific fields
- No observation station per township; query by city returns all nearby stations

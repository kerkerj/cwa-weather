---
name: cwa-weather
description: Use when user asks about Taiwan weather, forecast, observation, typhoon, alerts, or marine data. Requires cwa-weather CLI and CWA_API_KEY env var.
license: MIT
---

# cwa-weather

Taiwan CWA Open Data CLI. Output is human-readable text by default. Use `--json` for raw JSON.

## Quick Reference

| Need | Command | Key flags |
|------|---------|-----------|
| Township forecast | `cwa-weather forecast --city 臺北市 --town 中正區` | `--element`, `--days N`, `--summary`, `--time-from`, `--time-to` |
| Observation | `cwa-weather observe --city 新北市` | `--station` (alt), `--element` |
| 36hr overview | `cwa-weather overview --city 臺北市` | `--element`, `--time-from`, `--time-to` |
| Alerts | `cwa-weather alert` | `--city` |
| Typhoon | `cwa-weather typhoon` | `--td-no`, `--dataset` |
| Marine | `cwa-weather sea` | `--station` |
| Any endpoint | `cwa-weather query DATAID -p key=value` | `-p` (repeatable) |
| List cities | `cwa-weather cities` | `--city` (show towns) |

Run `cwa-weather <command> --help` for all flags and details.

## Notes

- Default output is human-readable text
- Use `--json` flag for raw JSON output (pipe to `jq` for field extraction)
- If `CWA_API_KEY` is not set, commands exit with a clear error message
- `台→臺` auto-converted (e.g. `台北市` → `臺北市`)
- `observe`: `--city` and `--station` are mutually exclusive — use one or the other, not both
- `--element` accepts comma-separated values (CWA API-defined names)
- `forecast --days N` shows N days of data (1-3, default 1)
- `forecast --summary` groups entries by day instead of flat list
- `query` command always outputs JSON (generic endpoint, cannot format)

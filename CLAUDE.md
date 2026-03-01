# cwa-weather

CWA Open Data API CLI and Go library for Taiwan weather data.

## Build & Test

```bash
make build      # go build -o bin/cwa-weather ./cmd/cwa-weather
make test       # go test -v ./...
make lint       # golangci-lint run ./...
make sec        # gosec ./...
make check      # test + lint + sec (run before every commit)
```

## Conventions

- **Auth**: HTTP header `Authorization: {apikey}` (NOT query param — avoid log leakage)
- **台→臺**: Use `cwa.NormalizeCity()` for all user-facing city input
- **Tests**: `testify/assert` + `testify/require`, 3A pattern (Arrange/Act/Assert with comment separators), table-driven where appropriate
- **Test mocking**: `httptest.NewServer` + fixture files from `testdata/`
- **Error wrapping**: Always `fmt.Errorf("failed to X: %w", err)` with context
- **Response**: `json.RawMessage` for `records` field (each API has different structure)
- **Package**: Library tests use `package cwa_test` (external test package)

## Project Structure

- `cwa/` — Go library (Client, Forecast, Observe, Overview, Alert, Typhoon, Sea)
- `cmd/cwa-weather/` — CLI subcommands (cobra)
- `testdata/` — API fixture JSON files for httptest
- `skill/` — Agent skill files (SKILL.md for Claude Code skills, AGENT.md for agent mode)

## Adding a New Subcommand

1. Create `cwa/{name}.go` with `const {name}DatasetID` and method on `*Client`
2. Create `cwa/{name}_test.go` with httptest + fixture
3. Create `cmd/cwa-weather/{name}.go` with cobra command + `init()` registering to `rootCmd`
4. Add E2E help test in `cmd/cwa-weather/cmd_test.go`
5. Update `README.md`, `skill/SKILL.md`, `skill/AGENT.md`

# Human-Readable CLI Output Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** CLI 預設輸出人類可讀文字，解除 jq 依賴，讓 AI agent 可以直接讀取天氣資訊。

**Architecture:** 在 `cwa/` 層新增 typed records structs + Parse 方法，在 `cmd/` 層新增 text formatter。預設 text 輸出，`--json` flag 維持原本 raw JSON 行為。API key 檢查集中化，錯誤訊息附帶申請 URL。

**Tech Stack:** Go, Cobra PersistentFlags, `encoding/json`, `fmt` (no external dependencies)

---

## Key Design Decisions

- **預設文字輸出** — `--json` 才回傳 raw JSON（向後相容）
- **forecast text mode** — 只印「天氣預報綜合描述」element，最近 3 筆時間點
- **`query` 維持 JSON-only** — 通用指令，無法格式化
- **typed structs** — 每個 API 的 Records 定義 Go struct，比 `map[string]interface{}` 更安全
- **API key 錯誤** — 附帶 CWA 申請 URL，agent 看到就知道該停
- **format.go** — 共用格式化工具，不引入外部套件

---

## Task 1: Centralize API Key Check + Add --json Flag

**Files:**
- Modify: `cmd/cwa-weather/root.go`

**Step 1: Write the failing test**

```go
// cmd/cwa-weather/cmd_test.go — append to existing file
func TestCLI_NoAPIKey(t *testing.T) {
	// Arrange
	cmd := exec.Command("go", "run", ".", "forecast", "--city", "臺北市")
	cmd.Env = append(os.Environ(), "CWA_API_KEY=")

	// Act
	out, err := cmd.CombinedOutput()

	// Assert
	assert.Error(t, err)
	assert.Contains(t, string(out), "CWA_API_KEY")
	assert.Contains(t, string(out), "opendata.cwa.gov.tw")
}
```

**Step 2: Run test to verify it fails**

Run: `cd cmd/cwa-weather && go test -run TestCLI_NoAPIKey -v`
Expected: FAIL — current error message doesn't contain URL

**Step 3: Write minimal implementation**

Add to `cmd/cwa-weather/root.go`:

```go
var jsonOutput bool

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output raw JSON instead of human-readable text")
}

func getAPIKey() (string, error) {
	key := os.Getenv("CWA_API_KEY")
	if key == "" {
		return "", fmt.Errorf("CWA_API_KEY environment variable is not set — get a free key at https://opendata.cwa.gov.tw/userLogin")
	}
	return key, nil
}
```

**Step 4: Update all 7 commands to use `getAPIKey()`**

In each of `forecast.go`, `observe.go`, `overview.go`, `alert.go`, `typhoon.go`, `sea.go`, `query.go`, replace:

```go
apiKey := os.Getenv("CWA_API_KEY")
if apiKey == "" {
    return fmt.Errorf("CWA_API_KEY environment variable is required")
}
```

With:

```go
apiKey, err := getAPIKey()
if err != nil {
    return err
}
```

**Step 5: Run test to verify it passes**

Run: `cd cmd/cwa-weather && go test -run TestCLI_NoAPIKey -v`
Expected: PASS

**Step 6: Run all existing tests**

Run: `make check`
Expected: All 42 tests pass, lint + sec clean

**Step 7: Commit**

```bash
git add cmd/cwa-weather/root.go cmd/cwa-weather/forecast.go cmd/cwa-weather/observe.go cmd/cwa-weather/overview.go cmd/cwa-weather/alert.go cmd/cwa-weather/typhoon.go cmd/cwa-weather/sea.go cmd/cwa-weather/query.go cmd/cwa-weather/cmd_test.go
git commit -m "refactor: centralize API key check, add --json global flag"
```

---

## Task 2: Shared Format Helpers

**Files:**
- Create: `cmd/cwa-weather/format.go`

**Step 1: Create format.go with shared helpers**

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func printHeader(title string) {
	fmt.Println(title)
}

func printRow(key, value string) {
	fmt.Printf("  %-16s %s\n", key+":", value)
}
```

**Step 2: Update all 7+1 commands to use `printJSON()` instead of inline encoder**

In each command, replace:

```go
enc := json.NewEncoder(os.Stdout)
enc.SetIndent("", "  ")
return enc.Encode(resp)
```

With:

```go
return printJSON(resp)
```

**Step 3: Run all tests**

Run: `make check`
Expected: All pass — pure refactor, no behavior change

**Step 4: Commit**

```bash
git add cmd/cwa-weather/format.go cmd/cwa-weather/forecast.go cmd/cwa-weather/observe.go cmd/cwa-weather/overview.go cmd/cwa-weather/alert.go cmd/cwa-weather/typhoon.go cmd/cwa-weather/sea.go cmd/cwa-weather/query.go cmd/cwa-weather/cities.go
git commit -m "refactor: extract shared format helpers, DRY json encoder"
```

---

## Task 3: Forecast Records Type + Parse Method

**Files:**
- Create: `cwa/records_forecast.go`
- Create: `cwa/records_forecast_test.go`

**Step 1: Write the failing test**

```go
// cwa/records_forecast_test.go
package cwa_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseForecastRecords(t *testing.T) {
	// Arrange
	data, err := os.ReadFile("../testdata/F-D0047-069-板橋區.json")
	require.NoError(t, err)

	var resp cwa.Response
	require.NoError(t, json.Unmarshal(data, &resp))

	// Act
	rec, err := resp.ParseForecastRecords()

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, rec.Locations)
	assert.Equal(t, "新北市", rec.Locations[0].LocationsName)
	require.NotEmpty(t, rec.Locations[0].Location)
	assert.Equal(t, "板橋區", rec.Locations[0].Location[0].LocationName)
	require.NotEmpty(t, rec.Locations[0].Location[0].WeatherElement)
	assert.Equal(t, "溫度", rec.Locations[0].Location[0].WeatherElement[0].ElementName)
	require.NotEmpty(t, rec.Locations[0].Location[0].WeatherElement[0].Time)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./cwa/ -run TestParseForecastRecords -v`
Expected: FAIL — `ParseForecastRecords` not defined

**Step 3: Write minimal implementation**

```go
// cwa/records_forecast.go
package cwa

import (
	"encoding/json"
	"fmt"
)

type ForecastRecords struct {
	Locations []ForecastLocations `json:"Locations"`
}

type ForecastLocations struct {
	DatasetDescription string             `json:"DatasetDescription"`
	LocationsName      string             `json:"LocationsName"`
	Dataid             string             `json:"Dataid"`
	Location           []ForecastLocation `json:"Location"`
}

type ForecastLocation struct {
	LocationName   string            `json:"LocationName"`
	Geocode        string            `json:"Geocode"`
	WeatherElement []ForecastElement `json:"WeatherElement"`
}

type ForecastElement struct {
	ElementName string         `json:"ElementName"`
	Time        []ForecastTime `json:"Time"`
}

type ForecastTime struct {
	DataTime     string              `json:"DataTime"`
	StartTime    string              `json:"StartTime"`
	EndTime      string              `json:"EndTime"`
	ElementValue []map[string]string `json:"ElementValue"`
}

func (r *Response) ParseForecastRecords() (*ForecastRecords, error) {
	var rec ForecastRecords
	if err := json.Unmarshal(r.Records, &rec); err != nil {
		return nil, fmt.Errorf("failed to parse forecast records: %w", err)
	}
	return &rec, nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./cwa/ -run TestParseForecastRecords -v`
Expected: PASS

**Step 5: Commit**

```bash
git add cwa/records_forecast.go cwa/records_forecast_test.go
git commit -m "feat: add ForecastRecords type and ParseForecastRecords method"
```

---

## Task 4: Observe Records Type + Parse Method

**Files:**
- Create: `cwa/records_observe.go`
- Create: `cwa/records_observe_test.go`

**Step 1: Write the failing test**

```go
// cwa/records_observe_test.go
package cwa_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseObserveRecords(t *testing.T) {
	// Arrange
	data, err := os.ReadFile("../testdata/O-A0003-001-新北市.json")
	require.NoError(t, err)

	var resp cwa.Response
	require.NoError(t, json.Unmarshal(data, &resp))

	// Act
	rec, err := resp.ParseObserveRecords()

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, rec.Station)
	assert.Equal(t, "板橋", rec.Station[0].StationName)
	assert.NotEmpty(t, rec.Station[0].ObsTime.DateTime)
	assert.NotEmpty(t, rec.Station[0].WeatherElement.AirTemperature)
	assert.Equal(t, "新北市", rec.Station[0].GeoInfo.CountyName)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./cwa/ -run TestParseObserveRecords -v`
Expected: FAIL

**Step 3: Write minimal implementation**

```go
// cwa/records_observe.go
package cwa

import (
	"encoding/json"
	"fmt"
)

type ObserveRecords struct {
	Station []ObserveStation `json:"Station"`
}

type ObserveStation struct {
	StationName    string                `json:"StationName"`
	StationId      string                `json:"StationId"`
	ObsTime        ObserveObsTime        `json:"ObsTime"`
	GeoInfo        ObserveGeoInfo        `json:"GeoInfo"`
	WeatherElement ObserveWeatherElement `json:"WeatherElement"`
}

type ObserveObsTime struct {
	DateTime string `json:"DateTime"`
}

type ObserveGeoInfo struct {
	CountyName string `json:"CountyName"`
	TownName   string `json:"TownName"`
}

type ObserveWeatherElement struct {
	Weather          string          `json:"Weather"`
	AirTemperature   string          `json:"AirTemperature"`
	RelativeHumidity string          `json:"RelativeHumidity"`
	WindSpeed        string          `json:"WindSpeed"`
	WindDirection    string          `json:"WindDirection"`
	AirPressure      string          `json:"AirPressure"`
	SunshineDuration string          `json:"SunshineDuration"`
	UVIndex          string          `json:"UVIndex"`
	Now              ObserveNow      `json:"Now"`
	DailyExtreme     ObserveDailyExt `json:"DailyExtreme"`
}

type ObserveNow struct {
	Precipitation string `json:"Precipitation"`
}

type ObserveDailyExt struct {
	DailyHigh ObserveTempInfo `json:"DailyHigh"`
	DailyLow  ObserveTempInfo `json:"DailyLow"`
}

type ObserveTempInfo struct {
	TemperatureInfo struct {
		AirTemperature string `json:"AirTemperature"`
	} `json:"TemperatureInfo"`
}

func (r *Response) ParseObserveRecords() (*ObserveRecords, error) {
	var rec ObserveRecords
	if err := json.Unmarshal(r.Records, &rec); err != nil {
		return nil, fmt.Errorf("failed to parse observe records: %w", err)
	}
	return &rec, nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./cwa/ -run TestParseObserveRecords -v`
Expected: PASS

**Step 5: Commit**

```bash
git add cwa/records_observe.go cwa/records_observe_test.go
git commit -m "feat: add ObserveRecords type and ParseObserveRecords method"
```

---

## Task 5: Overview Records Type + Parse Method

**Files:**
- Create: `cwa/records_overview.go`
- Create: `cwa/records_overview_test.go`

**Step 1: Write the failing test**

```go
// cwa/records_overview_test.go
package cwa_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseOverviewRecords(t *testing.T) {
	// Arrange
	data, err := os.ReadFile("../testdata/F-C0032-001.json")
	require.NoError(t, err)

	var resp cwa.Response
	require.NoError(t, json.Unmarshal(data, &resp))

	// Act
	rec, err := resp.ParseOverviewRecords()

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, rec.Location)
	assert.NotEmpty(t, rec.Location[0].LocationName)
	require.NotEmpty(t, rec.Location[0].WeatherElement)
	assert.Equal(t, "Wx", rec.Location[0].WeatherElement[0].ElementName)
	require.NotEmpty(t, rec.Location[0].WeatherElement[0].Time)
	assert.NotEmpty(t, rec.Location[0].WeatherElement[0].Time[0].Parameter.ParameterName)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./cwa/ -run TestParseOverviewRecords -v`
Expected: FAIL

**Step 3: Write minimal implementation**

```go
// cwa/records_overview.go
package cwa

import (
	"encoding/json"
	"fmt"
)

type OverviewRecords struct {
	DatasetDescription string             `json:"datasetDescription"`
	Location           []OverviewLocation `json:"location"`
}

type OverviewLocation struct {
	LocationName   string            `json:"locationName"`
	WeatherElement []OverviewElement `json:"weatherElement"`
}

type OverviewElement struct {
	ElementName string         `json:"elementName"`
	Time        []OverviewTime `json:"time"`
}

type OverviewTime struct {
	StartTime string            `json:"startTime"`
	EndTime   string            `json:"endTime"`
	Parameter OverviewParameter `json:"parameter"`
}

type OverviewParameter struct {
	ParameterName  string `json:"parameterName"`
	ParameterValue string `json:"parameterValue"`
	ParameterUnit  string `json:"parameterUnit"`
}

func (r *Response) ParseOverviewRecords() (*OverviewRecords, error) {
	var rec OverviewRecords
	if err := json.Unmarshal(r.Records, &rec); err != nil {
		return nil, fmt.Errorf("failed to parse overview records: %w", err)
	}
	return &rec, nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./cwa/ -run TestParseOverviewRecords -v`
Expected: PASS

**Step 5: Commit**

```bash
git add cwa/records_overview.go cwa/records_overview_test.go
git commit -m "feat: add OverviewRecords type and ParseOverviewRecords method"
```

---

## Task 6: Alert Records Type + Parse Method

**Files:**
- Create: `cwa/records_alert.go`
- Create: `cwa/records_alert_test.go`

**Step 1: Write the failing test**

```go
// cwa/records_alert_test.go
package cwa_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAlertRecords(t *testing.T) {
	// Arrange
	data, err := os.ReadFile("../testdata/W-C0033-001-天氣特報.json")
	require.NoError(t, err)

	var resp cwa.Response
	require.NoError(t, json.Unmarshal(data, &resp))

	// Act
	rec, err := resp.ParseAlertRecords()

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, rec.Location)
	assert.NotEmpty(t, rec.Location[0].LocationName)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./cwa/ -run TestParseAlertRecords -v`
Expected: FAIL

**Step 3: Write minimal implementation**

```go
// cwa/records_alert.go
package cwa

import (
	"encoding/json"
	"fmt"
)

type AlertRecords struct {
	Location []AlertLocation `json:"location"`
}

type AlertLocation struct {
	LocationName     string           `json:"locationName"`
	Geocode          json.Number      `json:"geocode"`
	HazardConditions HazardConditions `json:"hazardConditions"`
}

type HazardConditions struct {
	Hazards []Hazard `json:"hazards"`
}

type Hazard struct {
	Info HazardInfo `json:"info"`
}

type HazardInfo struct {
	Phenomena    string `json:"phenomena"`
	Significance string `json:"significance"`
}

func (r *Response) ParseAlertRecords() (*AlertRecords, error) {
	var rec AlertRecords
	if err := json.Unmarshal(r.Records, &rec); err != nil {
		return nil, fmt.Errorf("failed to parse alert records: %w", err)
	}
	return &rec, nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./cwa/ -run TestParseAlertRecords -v`
Expected: PASS

**Step 5: Commit**

```bash
git add cwa/records_alert.go cwa/records_alert_test.go
git commit -m "feat: add AlertRecords type and ParseAlertRecords method"
```

---

## Task 7: Typhoon Records Type + Parse Method

**Files:**
- Create: `cwa/records_typhoon.go`
- Create: `cwa/records_typhoon_test.go`

**Step 1: Write the failing test**

```go
// cwa/records_typhoon_test.go
package cwa_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTyphoonRecords(t *testing.T) {
	// Arrange
	data, err := os.ReadFile("../testdata/W-C0034-005-颱風.json")
	require.NoError(t, err)

	var resp cwa.Response
	require.NoError(t, json.Unmarshal(data, &resp))

	// Act
	rec, err := resp.ParseTyphoonRecords()

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, rec.TropicalCyclones.TropicalCyclone)
	tc := rec.TropicalCyclones.TropicalCyclone[0]
	assert.NotEmpty(t, tc.TyphoonName)
	assert.NotEmpty(t, tc.CwaTyphoonName)
	require.NotNil(t, tc.AnalysisData)
	require.NotEmpty(t, tc.AnalysisData.Fix)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./cwa/ -run TestParseTyphoonRecords -v`
Expected: FAIL

**Step 3: Write minimal implementation**

```go
// cwa/records_typhoon.go
package cwa

import (
	"encoding/json"
	"fmt"
)

type TyphoonRecords struct {
	TropicalCyclones TropicalCyclones `json:"TropicalCyclones"`
}

type TropicalCyclones struct {
	TropicalCyclone []TropicalCyclone `json:"TropicalCyclone"`
}

type TropicalCyclone struct {
	Year           string       `json:"Year"`
	TyphoonName    string       `json:"TyphoonName"`
	CwaTyphoonName string       `json:"CwaTyphoonName"`
	CwaTdNo        string       `json:"CwaTdNo"`
	CwaTyNo        string       `json:"CwaTyNo"`
	AnalysisData   *CycloneData `json:"AnalysisData,omitempty"`
	ForecastData   *CycloneData `json:"ForecastData,omitempty"`
}

type CycloneData struct {
	Fix []CycloneFix `json:"Fix"`
}

type CycloneFix struct {
	DateTime            string `json:"DateTime"`
	CoordinateLongitude string `json:"CoordinateLongitude"`
	CoordinateLatitude  string `json:"CoordinateLatitude"`
	MaxWindSpeed        string `json:"MaxWindSpeed"`
	MaxGustSpeed        string `json:"MaxGustSpeed"`
	Pressure            string `json:"Pressure"`
	MovingSpeed         string `json:"MovingSpeed"`
	MovingDirection     string `json:"MovingDirection"`
}

func (r *Response) ParseTyphoonRecords() (*TyphoonRecords, error) {
	var rec TyphoonRecords
	if err := json.Unmarshal(r.Records, &rec); err != nil {
		return nil, fmt.Errorf("failed to parse typhoon records: %w", err)
	}
	return &rec, nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./cwa/ -run TestParseTyphoonRecords -v`
Expected: PASS

**Step 5: Commit**

```bash
git add cwa/records_typhoon.go cwa/records_typhoon_test.go
git commit -m "feat: add TyphoonRecords type and ParseTyphoonRecords method"
```

---

## Task 8: Sea Records Type + Parse Method

**Files:**
- Create: `cwa/records_sea.go`
- Create: `cwa/records_sea_test.go`

**Step 1: Write the failing test**

```go
// cwa/records_sea_test.go
package cwa_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSeaRecords(t *testing.T) {
	// Arrange
	data, err := os.ReadFile("../testdata/O-B0075-001-海象.json")
	require.NoError(t, err)

	var resp cwa.Response
	require.NoError(t, json.Unmarshal(data, &resp))

	// Act
	rec, err := resp.ParseSeaRecords()

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, rec.SeaSurfaceObs.Location)
	loc := rec.SeaSurfaceObs.Location[0]
	assert.NotEmpty(t, loc.Station.StationID)
	require.NotEmpty(t, loc.StationObsTimes.StationObsTime)
}
```

**Note:** Sea fixture 使用大寫 `Records` key — 如果 `resp.Records` 為 null，先修正 fixture 或在 `cwa/client.go` 的 `Response` struct 加上 case-insensitive handling。確認 fixture casing 後再決定。

**Step 2: Run test to verify it fails**

Run: `go test ./cwa/ -run TestParseSeaRecords -v`
Expected: FAIL

**Step 3: Write minimal implementation**

```go
// cwa/records_sea.go
package cwa

import (
	"encoding/json"
	"fmt"
)

type SeaRecords struct {
	SeaSurfaceObs SeaSurfaceObs `json:"SeaSurfaceObs"`
}

type SeaSurfaceObs struct {
	Location []SeaLocation `json:"Location"`
}

type SeaLocation struct {
	Station         SeaStation      `json:"Station"`
	StationObsTimes StationObsTimes `json:"StationObsTimes"`
}

type SeaStation struct {
	StationID   string `json:"StationID"`
	StationName string `json:"StationName"`
}

type StationObsTimes struct {
	StationObsTime []SeaObsTime `json:"StationObsTime"`
}

type SeaObsTime struct {
	DateTime        string             `json:"DateTime"`
	WeatherElements SeaWeatherElements `json:"WeatherElements"`
}

type SeaWeatherElements struct {
	TideHeight     string `json:"TideHeight"`
	TideLevel      string `json:"TideLevel"`
	SeaTemperature string `json:"SeaTemperature"`
	Temperature    string `json:"Temperature"`
}

func (r *Response) ParseSeaRecords() (*SeaRecords, error) {
	var rec SeaRecords
	if err := json.Unmarshal(r.Records, &rec); err != nil {
		return nil, fmt.Errorf("failed to parse sea records: %w", err)
	}
	return &rec, nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./cwa/ -run TestParseSeaRecords -v`
Expected: PASS (if fixture casing is compatible) or needs fixture fix

**Step 5: Commit**

```bash
git add cwa/records_sea.go cwa/records_sea_test.go
git commit -m "feat: add SeaRecords type and ParseSeaRecords method"
```

---

## Task 9: Forecast Text Formatter

**Files:**
- Modify: `cmd/cwa-weather/forecast.go`

**Step 1: Modify forecast RunE to support text output**

```go
RunE: func(cmd *cobra.Command, args []string) error {
    apiKey, err := getAPIKey()
    if err != nil {
        return err
    }

    client := cwa.NewClient(apiKey)

    var opts []cwa.ForecastOption
    if forecastElement != "" || forecastFrom != "" || forecastTo != "" {
        opts = append(opts, cwa.ForecastOption{
            Element:  forecastElement,
            TimeFrom: forecastFrom,
            TimeTo:   forecastTo,
        })
    }

    resp, err := client.Forecast(context.Background(), forecastCity, forecastTown, opts...)
    if err != nil {
        return fmt.Errorf("failed to get forecast: %w", err)
    }

    if jsonOutput {
        return printJSON(resp)
    }

    rec, err := resp.ParseForecastRecords()
    if err != nil {
        return err
    }

    for _, locs := range rec.Locations {
        for _, loc := range locs.Location {
            printHeader(fmt.Sprintf("%s（%s）", loc.LocationName, locs.LocationsName))
            for _, elem := range loc.WeatherElement {
                // text mode 預設只印綜合描述
                if forecastElement == "" && elem.ElementName != "天氣預報綜合描述" {
                    continue
                }
                fmt.Printf("  %s\n", elem.ElementName)
                maxTime := 3
                if len(elem.Time) < maxTime {
                    maxTime = len(elem.Time)
                }
                for _, t := range elem.Time[:maxTime] {
                    timeStr := t.DataTime
                    if t.StartTime != "" {
                        timeStr = t.StartTime + " ~ " + t.EndTime
                    }
                    val := firstValue(t.ElementValue)
                    fmt.Printf("    %-40s %s\n", timeStr, val)
                }
            }
        }
    }
    return nil
},
```

Add helper to `format.go`:

```go
func firstValue(ev []map[string]string) string {
	if len(ev) == 0 {
		return ""
	}
	for _, v := range ev[0] {
		return v
	}
	return ""
}
```

**Step 2: Run build and smoke test**

Run: `make build && ./bin/cwa-weather forecast --city 新北市 --town 板橋區`
Expected: Human-readable text output

Run: `./bin/cwa-weather forecast --city 新北市 --town 板橋區 --json`
Expected: Raw JSON (same as before)

**Step 3: Run all tests**

Run: `make check`
Expected: All pass

**Step 4: Commit**

```bash
git add cmd/cwa-weather/forecast.go cmd/cwa-weather/format.go
git commit -m "feat(forecast): add human-readable text output"
```

---

## Task 10: Observe Text Formatter

**Files:**
- Modify: `cmd/cwa-weather/observe.go`

**Step 1: Modify observe RunE**

Text output format:

```
板橋（板橋區, 新北市）  2026-02-28 23:50
  天氣: 陰    氣溫: 19.6°C    濕度: 95%
  風速: 1.3 m/s    風向: 10°    氣壓: 1011.0 hPa
  今日降雨: 1.0 mm    紫外線: 0
```

After the `if jsonOutput { return printJSON(resp) }` block:

```go
rec, err := resp.ParseObserveRecords()
if err != nil {
    return err
}

for _, stn := range rec.Station {
    printHeader(fmt.Sprintf("%s（%s, %s）  %s",
        stn.StationName, stn.GeoInfo.TownName, stn.GeoInfo.CountyName,
        stn.ObsTime.DateTime))
    we := stn.WeatherElement
    fmt.Printf("  天氣: %s    氣溫: %s°C    濕度: %s%%\n", we.Weather, we.AirTemperature, we.RelativeHumidity)
    fmt.Printf("  風速: %s m/s    風向: %s°    氣壓: %s hPa\n", we.WindSpeed, we.WindDirection, we.AirPressure)
    fmt.Printf("  今日降雨: %s mm    紫外線: %s\n", we.Now.Precipitation, we.UVIndex)
    fmt.Println()
}
return nil
```

**Step 2: Smoke test**

Run: `make build && ./bin/cwa-weather observe --city 新北市`

**Step 3: Commit**

```bash
git add cmd/cwa-weather/observe.go
git commit -m "feat(observe): add human-readable text output"
```

---

## Task 11: Overview Text Formatter

**Files:**
- Modify: `cmd/cwa-weather/overview.go`

**Step 1: Modify overview RunE**

Text output format — group elements by time period:

```
嘉義縣
  03/01 00:00 ~ 06:00
    天氣: 陰天    降雨: 20%    溫度: 17~20°C    舒適度: 稍有寒意至舒適
```

Build a map of time→elements, then print grouped.

**Step 2: Smoke test + commit**

```bash
git add cmd/cwa-weather/overview.go
git commit -m "feat(overview): add human-readable text output"
```

---

## Task 12: Alert Text Formatter

**Files:**
- Modify: `cmd/cwa-weather/alert.go`

**Step 1: Modify alert RunE**

```go
rec, err := resp.ParseAlertRecords()
if err != nil {
    return err
}

hasAlert := false
for _, loc := range rec.Location {
    if len(loc.HazardConditions.Hazards) > 0 {
        hasAlert = true
        printHeader(loc.LocationName)
        for _, h := range loc.HazardConditions.Hazards {
            fmt.Printf("  %s（%s）\n", h.Info.Phenomena, h.Info.Significance)
        }
    }
}
if !hasAlert {
    fmt.Println("目前無天氣特報。")
}
return nil
```

**Step 2: Smoke test + commit**

```bash
git add cmd/cwa-weather/alert.go
git commit -m "feat(alert): add human-readable text output"
```

---

## Task 13: Typhoon Text Formatter

**Files:**
- Modify: `cmd/cwa-weather/typhoon.go`

**Step 1: Modify typhoon RunE**

```go
rec, err := resp.ParseTyphoonRecords()
if err != nil {
    return err
}

tcs := rec.TropicalCyclones.TropicalCyclone
if len(tcs) == 0 {
    fmt.Println("目前無颱風資訊。")
    return nil
}
for _, tc := range tcs {
    printHeader(fmt.Sprintf("%s (%s) TD-%s", tc.CwaTyphoonName, tc.TyphoonName, tc.CwaTdNo))
    if tc.AnalysisData != nil && len(tc.AnalysisData.Fix) > 0 {
        fix := tc.AnalysisData.Fix[len(tc.AnalysisData.Fix)-1] // latest
        fmt.Printf("  位置: %s°E %s°N    氣壓: %s hPa\n", fix.CoordinateLongitude, fix.CoordinateLatitude, fix.Pressure)
        fmt.Printf("  最大風速: %s m/s    陣風: %s m/s\n", fix.MaxWindSpeed, fix.MaxGustSpeed)
        fmt.Printf("  移動: %s km/h 往 %s\n", fix.MovingSpeed, fix.MovingDirection)
    }
    fmt.Println()
}
return nil
```

**Step 2: Smoke test + commit**

```bash
git add cmd/cwa-weather/typhoon.go
git commit -m "feat(typhoon): add human-readable text output"
```

---

## Task 14: Sea Text Formatter

**Files:**
- Modify: `cmd/cwa-weather/sea.go`

**Step 1: Modify sea RunE**

```go
rec, err := resp.ParseSeaRecords()
if err != nil {
    return err
}

for _, loc := range rec.SeaSurfaceObs.Location {
    printHeader(fmt.Sprintf("Station %s", loc.Station.StationID))
    maxObs := 3
    obs := loc.StationObsTimes.StationObsTime
    if len(obs) < maxObs {
        maxObs = len(obs)
    }
    for _, o := range obs[:maxObs] {
        we := o.WeatherElements
        seaTemp := we.SeaTemperature
        if seaTemp == "None" || seaTemp == "" {
            seaTemp = "N/A"
        }
        fmt.Printf("  %s  潮高: %sm  %s  海溫: %s\n",
            o.DateTime, we.TideHeight, we.TideLevel, seaTemp)
    }
    fmt.Println()
}
return nil
```

**Step 2: Smoke test + commit**

```bash
git add cmd/cwa-weather/sea.go
git commit -m "feat(sea): add human-readable text output"
```

---

## Task 15: Cities Text Formatter

**Files:**
- Modify: `cmd/cwa-weather/cities.go`

**Step 1: Modify cities RunE**

Text output: plain list, one per line.

```go
if jsonOutput {
    return printJSON(result)
}

// text output
for _, name := range result {
    fmt.Println(name)
}
return nil
```

**Step 2: Add E2E tests**

```go
// cmd/cwa-weather/cmd_test.go — append
func TestCLI_Cities_TextOutput(t *testing.T) {
	// Arrange + Act
	out, err := exec.Command("go", "run", ".", "cities").CombinedOutput()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, string(out), "臺北市")
	assert.NotContains(t, string(out), "[") // not JSON
}
```

**Step 3: Run tests + commit**

Run: `make check`

```bash
git add cmd/cwa-weather/cities.go cmd/cwa-weather/cmd_test.go
git commit -m "feat(cities): add human-readable text output"
```

---

## Task 16: Update Plugin/Skill Files

**Files:**
- Modify: `plugins/cwa-weather/skills/cwa-weather/SKILL.md`
- Modify: `plugins/cwa-weather/agents/AGENT.md`

**Step 1: Update SKILL.md**

Change line 9 from:
```
Taiwan CWA Open Data CLI. Output is always JSON — pipe to `jq` for extraction.
```
To:
```
Taiwan CWA Open Data CLI. Output is human-readable text by default. Use `--json` for raw JSON.
```

Add to Notes section:
```
- Default output is human-readable text — no `jq` needed
- Use `--json` flag for raw JSON output (for scripting)
- If `CWA_API_KEY` is not set, commands exit with a clear error message
```

**Step 2: Update AGENT.md**

Change line 3 from:
```
Taiwan CWA Open Data CLI. Requires `CWA_API_KEY` env var. Output is always JSON.
```
To:
```
Taiwan CWA Open Data CLI. Requires `CWA_API_KEY` env var. Output is human-readable text by default.
```

Change line 31 from:
```
- All commands return full CWA JSON — use `jq` to extract specific fields
```
To:
```
- Default output is human-readable text — read directly, no jq needed
- Use `--json` for raw JSON when needed for data extraction
- If CWA_API_KEY is missing, stop and inform the user to set it
```

**Step 3: Commit**

```bash
git add plugins/cwa-weather/skills/cwa-weather/SKILL.md plugins/cwa-weather/agents/AGENT.md
git commit -m "docs: update skill/agent files for human-readable output"
```

---

## Task 17: Final Verification

**Step 1: Run full check**

```bash
make check
```

Expected: All tests pass, lint clean, sec clean

**Step 2: Smoke test all commands (requires CWA_API_KEY)**

```bash
./bin/cwa-weather forecast --city 新北市 --town 板橋區
./bin/cwa-weather forecast --city 新北市 --town 板橋區 --json
./bin/cwa-weather observe --city 新北市
./bin/cwa-weather overview
./bin/cwa-weather alert
./bin/cwa-weather typhoon
./bin/cwa-weather sea
./bin/cwa-weather cities
./bin/cwa-weather cities --json
CWA_API_KEY= ./bin/cwa-weather forecast --city 新北市
```

**Step 3: Verify no API key error message**

```bash
CWA_API_KEY= ./bin/cwa-weather forecast --city 臺北市 2>&1
# Expected: Error: CWA_API_KEY environment variable is not set — get a free key at https://opendata.cwa.gov.tw/userLogin
```

# cwa-weather MVP Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a Go CLI tool + library that queries any CWA Open Data API endpoint, with convenience wrappers for forecast and observe.

**Architecture:** Generic `Query(dataid, params)` core that hits any CWA REST endpoint and returns raw JSON. Convenience methods (`Forecast`, `Observe`) wrap the generic query with parameter validation and city→dataid routing. CLI uses Cobra with `query`, `forecast`, `observe`, `cities` subcommands.

**Tech Stack:** Go, Cobra, `encoding/json` with `json.RawMessage`, `embed` for static city/town data, `httptest` for testing.

---

## Key Design Decisions

- **83 endpoints** — all supported via generic `query` subcommand
- **Response parsing** — `json.RawMessage` for `records` (every API has different structure)
- **CWA casing inconsistency** — Go `encoding/json` is case-insensitive, no special handling needed
- **`success` field** — string `"true"`/`"false"`, not bool
- **台→臺** — global replacement before lookup
- **City/town data** — static `embed` from pre-fetched JSON
- **Observe filtering** — `--city` and `--station` mutually exclusive, both server-side query params
- **Output** — raw CWA JSON passthrough for `query`; structured JSON envelope for convenience commands
- **Testdata** — 12 fixture files already in `testdata/`

---

## Task 1: Go Module Init + Project Structure

**Files:**
- Create: `go.mod`
- Create: `cwa/client.go` (placeholder)
- Create: `cmd/cwa-weather/main.go` (placeholder)
- Create: `Makefile`

**Step 1: Init go module**

```bash
go mod init github.com/kerkerj/cwa-weather
```

**Step 2: Create Makefile**

```makefile
.PHONY: build test lint clean

build:
	go build -o bin/cwa-weather ./cmd/cwa-weather

test:
	go test -v ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/
```

**Step 3: Create placeholder files to verify structure compiles**

`cwa/client.go`:
```go
package cwa
```

`cmd/cwa-weather/main.go`:
```go
package main

func main() {}
```

**Step 4: Verify it builds**

Run: `go build ./...`
Expected: SUCCESS, no errors

**Step 5: Commit**

```bash
git add go.mod Makefile cwa/client.go cmd/cwa-weather/main.go
git commit -m "init: go module, project structure, Makefile"
```

---

## Task 2: Core HTTP Client + Generic Query

**Files:**
- Create: `cwa/client.go`
- Create: `cwa/client_test.go`

**Step 1: Write the failing test for NewClient**

`cwa/client_test.go`:
```go
package cwa_test

import (
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	// Arrange
	apiKey := "test-api-key"

	// Act
	c := cwa.NewClient(apiKey)

	// Assert
	assert.NotNil(t, c)
}

func TestNewClient_EmptyKey(t *testing.T) {
	// Arrange & Act
	c := cwa.NewClient("")

	// Assert
	assert.NotNil(t, c)
}
```

**Step 2: Run test to verify it fails**

Run: `go test -v ./cwa/...`
Expected: FAIL — `NewClient` not defined

**Step 3: Implement NewClient**

`cwa/client.go`:
```go
package cwa

import (
	"net/http"
)

const defaultBaseURL = "https://opendata.cwa.gov.tw/api/v1/rest/datastore"

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: http.DefaultClient,
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test -v ./cwa/...`
Expected: PASS

**Step 5: Write the failing test for Query**

Add to `cwa/client_test.go`:
```go
import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestQuery_Forecast(t *testing.T) {
	// Arrange
	fixture, err := os.ReadFile("../testdata/F-D0047-069-板橋區.json")
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/F-D0047-069", r.URL.Path)
		assert.Equal(t, "test-key", r.URL.Query().Get("Authorization"))
		assert.Equal(t, "JSON", r.URL.Query().Get("format"))
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture)
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Query(context.Background(), "F-D0047-069", map[string]string{
		"LocationName": "板橋區",
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
	assert.Equal(t, "F-D0047-069", resp.Result.ResourceID)
	assert.NotEmpty(t, resp.Records)
}

func TestQuery_InvalidAPIKey(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":"false","result":{"message":"Invalid API key"}}`))
	}))
	defer server.Close()

	c := cwa.NewClient("bad-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Query(context.Background(), "F-D0047-069", nil)

	// Assert
	assert.NoError(t, err) // HTTP succeeded, CWA returned error in body
	assert.Equal(t, "false", resp.Success)
}
```

**Step 6: Run test to verify it fails**

Run: `go test -v ./cwa/...`
Expected: FAIL — `Query`, `SetBaseURL`, response types not defined

**Step 7: Implement Query + response types**

Add to `cwa/client.go`:
```go
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Response struct {
	Success string          `json:"success"`
	Result  Result          `json:"result"`
	Records json.RawMessage `json:"records"`
}

type Result struct {
	ResourceID string  `json:"resource_id"`
	Fields     []Field `json:"fields"`
}

type Field struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func (c *Client) SetBaseURL(u string) {
	c.baseURL = u
}

func (c *Client) Query(ctx context.Context, dataID string, params map[string]string) (*Response, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%s", c.baseURL, dataID))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("Authorization", c.apiKey)
	q.Set("format", "JSON")
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}

	var resp Response
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}
```

**Step 8: Run tests to verify they pass**

Run: `go test -v ./cwa/...`
Expected: PASS

**Step 9: Commit**

```bash
git add cwa/client.go cwa/client_test.go go.mod go.sum
git commit -m "feat: core HTTP client with generic Query method"
```

---

## Task 3: City/Town Mapping + 台→臺 Alias

**Files:**
- Create: `cwa/endpoints.go`
- Create: `cwa/endpoints_test.go`
- Create: `cwa/cities.json` (embedded static data)

**Step 1: Write the failing tests**

`cwa/endpoints_test.go`:
```go
package cwa_test

import (
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
)

func TestGetDatasetID(t *testing.T) {
	tests := []struct {
		name     string
		city     string
		expected string
		wantErr  bool
	}{
		{
			name:     "standard city name",
			city:     "臺北市",
			expected: "F-D0047-061",
		},
		{
			name:     "simplified 台 alias",
			city:     "台北市",
			expected: "F-D0047-061",
		},
		{
			name:    "unknown city",
			city:    "不存在市",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			id, err := cwa.GetDatasetID(tt.city)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, id)
			}
		})
	}
}

func TestNormalizeCity(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"台北市", "臺北市"},
		{"台中市", "臺中市"},
		{"台南市", "臺南市"},
		{"台東縣", "臺東縣"},
		{"臺北市", "臺北市"}, // already correct
		{"新北市", "新北市"},  // no 台
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// Act
			result := cwa.NormalizeCity(tt.input)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCities(t *testing.T) {
	// Act
	cities := cwa.Cities()

	// Assert
	assert.Len(t, cities, 22)
	assert.Contains(t, cities, "臺北市")
	assert.Contains(t, cities, "新北市")
}
```

**Step 2: Run test to verify it fails**

Run: `go test -v ./cwa/...`
Expected: FAIL

**Step 3: Implement endpoints.go**

`cwa/endpoints.go`:
```go
package cwa

import (
	"fmt"
	"sort"
	"strings"
)

// cityDatasets maps city name to 3-day forecast dataset ID.
var cityDatasets = map[string]string{
	"宜蘭縣": "F-D0047-001",
	"桃園市": "F-D0047-005",
	"新竹縣": "F-D0047-009",
	"苗栗縣": "F-D0047-013",
	"彰化縣": "F-D0047-017",
	"南投縣": "F-D0047-021",
	"雲林縣": "F-D0047-025",
	"嘉義縣": "F-D0047-029",
	"屏東縣": "F-D0047-033",
	"臺東縣": "F-D0047-037",
	"花蓮縣": "F-D0047-041",
	"澎湖縣": "F-D0047-045",
	"基隆市": "F-D0047-049",
	"新竹市": "F-D0047-053",
	"嘉義市": "F-D0047-057",
	"臺北市": "F-D0047-061",
	"高雄市": "F-D0047-065",
	"新北市": "F-D0047-069",
	"臺中市": "F-D0047-073",
	"臺南市": "F-D0047-077",
	"連江縣": "F-D0047-081",
	"金門縣": "F-D0047-085",
}

// NormalizeCity replaces 台 with 臺 to match CWA naming.
func NormalizeCity(name string) string {
	return strings.ReplaceAll(name, "台", "臺")
}

// GetDatasetID returns the 3-day forecast dataset ID for a city.
func GetDatasetID(city string) (string, error) {
	city = NormalizeCity(city)
	id, ok := cityDatasets[city]
	if !ok {
		return "", fmt.Errorf("city not found: %s", city)
	}
	return id, nil
}

// Cities returns all supported city names sorted.
func Cities() []string {
	cities := make([]string, 0, len(cityDatasets))
	for k := range cityDatasets {
		cities = append(cities, k)
	}
	sort.Strings(cities)
	return cities
}
```

**Step 4: Run tests to verify they pass**

Run: `go test -v ./cwa/...`
Expected: PASS

**Step 5: Commit**

```bash
git add cwa/endpoints.go cwa/endpoints_test.go
git commit -m "feat: city-to-dataset mapping with 台→臺 normalization"
```

---

## Task 4: Static Town Data (embed)

**Files:**
- Create: `cwa/gen/fetch_towns.go` (one-time script to fetch town data)
- Create: `cwa/towns.json` (generated, embedded)
- Modify: `cwa/endpoints.go` — add `Towns(city)` function

**Step 1: Write a script to fetch all town names from CWA API**

`cwa/gen/fetch_towns.go`:
```go
//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
)

// Uses F-D0047-089 (全臺灣未來3天) to get all towns at once.
func main() {
	apiKey := os.Getenv("CWA_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "CWA_API_KEY required")
		os.Exit(1)
	}

	url := fmt.Sprintf(
		"https://opendata.cwa.gov.tw/api/v1/rest/datastore/F-D0047-089?Authorization=%s&format=JSON",
		apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Records struct {
			Locations []struct {
				LocationsName string `json:"LocationsName"`
				Location      []struct {
					LocationName string `json:"LocationName"`
				} `json:"Location"`
			} `json:"Locations"`
		} `json:"records"`
	}
	json.Unmarshal(body, &data)

	result := make(map[string][]string)
	for _, loc := range data.Records.Locations {
		var towns []string
		for _, l := range loc.Location {
			towns = append(towns, l.LocationName)
		}
		sort.Strings(towns)
		result[loc.LocationsName] = towns
	}

	out, _ := json.MarshalIndent(result, "", "  ")
	os.WriteFile("cwa/towns.json", out, 0644)
	fmt.Printf("Written %d cities\n", len(result))
}
```

**Step 2: Run the script to generate towns.json**

Run: `CWA_API_KEY=$CWA_API_KEY go run cwa/gen/fetch_towns.go`
Expected: `towns.json` created with 22 cities and their towns

Note: F-D0047-089 is the "全臺灣未來3天" endpoint. If it bundles all cities in one response, one call gets all towns. If not, we may need to iterate per city — verify from the response.

**Step 3: Write failing test for Towns()**

Add to `cwa/endpoints_test.go`:
```go
func TestTowns(t *testing.T) {
	// Act
	towns, err := cwa.Towns("臺北市")

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, towns, "中正區")
	assert.Contains(t, towns, "大安區")
}

func TestTowns_WithAlias(t *testing.T) {
	// Act
	towns, err := cwa.Towns("台北市")

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, towns)
}

func TestTowns_UnknownCity(t *testing.T) {
	// Act
	_, err := cwa.Towns("不存在市")

	// Assert
	assert.Error(t, err)
}
```

**Step 4: Run test to verify it fails**

Run: `go test -v ./cwa/...`
Expected: FAIL

**Step 5: Implement Towns() with embed**

Add to `cwa/endpoints.go`:
```go
import (
	_ "embed"
	"encoding/json"
)

//go:embed towns.json
var townsJSON []byte

var townsByCity map[string][]string

func init() {
	json.Unmarshal(townsJSON, &townsByCity)
}

// Towns returns all town names for a city.
func Towns(city string) ([]string, error) {
	city = NormalizeCity(city)
	towns, ok := townsByCity[city]
	if !ok {
		return nil, fmt.Errorf("city not found: %s", city)
	}
	return towns, nil
}
```

**Step 6: Run tests to verify they pass**

Run: `go test -v ./cwa/...`
Expected: PASS

**Step 7: Commit**

```bash
git add cwa/gen/fetch_towns.go cwa/towns.json cwa/endpoints.go cwa/endpoints_test.go
git commit -m "feat: static town data with go:embed"
```

---

## Task 5: Forecast Convenience Method

**Files:**
- Create: `cwa/forecast.go`
- Create: `cwa/forecast_test.go`

**Step 1: Write the failing test**

`cwa/forecast_test.go`:
```go
package cwa_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
)

func TestForecast(t *testing.T) {
	// Arrange
	fixture, err := os.ReadFile("../testdata/F-D0047-069-板橋區.json")
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/F-D0047-069", r.URL.Path)
		assert.Equal(t, "板橋區", r.URL.Query().Get("LocationName"))
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture)
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Forecast(context.Background(), "新北市", "板橋區")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

func TestForecast_CityAlias(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify it resolved 台北市 → 臺北市 → F-D0047-061
		assert.Equal(t, "/F-D0047-061", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":"true","result":{"resource_id":"F-D0047-061"},"records":{}}`))
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Forecast(context.Background(), "台北市", "")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

func TestForecast_UnknownCity(t *testing.T) {
	// Arrange
	c := cwa.NewClient("test-key")

	// Act
	_, err := c.Forecast(context.Background(), "不存在市", "")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "city not found")
}
```

**Step 2: Run test to verify it fails**

Run: `go test -v ./cwa/...`
Expected: FAIL

**Step 3: Implement Forecast**

`cwa/forecast.go`:
```go
package cwa

import "context"

// Forecast returns township-level weather forecast for a city.
// If town is empty, returns all towns in the city.
func (c *Client) Forecast(ctx context.Context, city, town string) (*Response, error) {
	datasetID, err := GetDatasetID(city)
	if err != nil {
		return nil, err
	}

	params := make(map[string]string)
	if town != "" {
		params["LocationName"] = NormalizeCity(town)
	}

	return c.Query(ctx, datasetID, params)
}
```

**Step 4: Run tests to verify they pass**

Run: `go test -v ./cwa/...`
Expected: PASS

**Step 5: Commit**

```bash
git add cwa/forecast.go cwa/forecast_test.go
git commit -m "feat: Forecast convenience method with city→dataid routing"
```

---

## Task 6: Observe Convenience Method

**Files:**
- Create: `cwa/observe.go`
- Create: `cwa/observe_test.go`

**Step 1: Write the failing test**

`cwa/observe_test.go`:
```go
package cwa_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
)

func TestObserve_ByCity(t *testing.T) {
	// Arrange
	fixture, err := os.ReadFile("../testdata/O-A0003-001-新北市.json")
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/O-A0003-001", r.URL.Path)
		assert.Equal(t, "新北市", r.URL.Query().Get("CountyName"))
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture)
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Observe(context.Background(), cwa.ObserveByCity("新北市"))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

func TestObserve_ByStation(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "淡水", r.URL.Query().Get("StationName"))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":"true","result":{"resource_id":"O-A0003-001"},"records":{}}`))
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Observe(context.Background(), cwa.ObserveByStation("淡水"))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

func TestObserve_NoOption(t *testing.T) {
	// Arrange
	c := cwa.NewClient("test-key")

	// Act
	_, err := c.Observe(context.Background(), nil)

	// Assert
	assert.Error(t, err)
}
```

**Step 2: Run test to verify it fails**

Run: `go test -v ./cwa/...`
Expected: FAIL

**Step 3: Implement Observe**

`cwa/observe.go`:
```go
package cwa

import (
	"context"
	"fmt"
)

const observeDatasetID = "O-A0003-001"

// ObserveOption configures an Observe query.
type ObserveOption func(params map[string]string)

// ObserveByCity filters observation stations by city.
func ObserveByCity(city string) ObserveOption {
	return func(params map[string]string) {
		params["CountyName"] = NormalizeCity(city)
	}
}

// ObserveByStation filters by station name.
func ObserveByStation(station string) ObserveOption {
	return func(params map[string]string) {
		params["StationName"] = station
	}
}

// Observe returns real-time observation data.
func (c *Client) Observe(ctx context.Context, opt ObserveOption) (*Response, error) {
	if opt == nil {
		return nil, fmt.Errorf("observe requires either city or station option")
	}

	params := make(map[string]string)
	opt(params)

	return c.Query(ctx, observeDatasetID, params)
}
```

**Step 4: Run tests to verify they pass**

Run: `go test -v ./cwa/...`
Expected: PASS

**Step 5: Commit**

```bash
git add cwa/observe.go cwa/observe_test.go
git commit -m "feat: Observe convenience method with city/station options"
```

---

## Task 7: CLI Root Command + Version

**Files:**
- Modify: `cmd/cwa-weather/main.go`
- Create: `cmd/cwa-weather/root.go`

**Step 1: Add cobra dependency**

Run: `go get github.com/spf13/cobra`

**Step 2: Implement root command**

`cmd/cwa-weather/root.go`:
```go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "cwa-weather",
	Short:   "CWA Open Data API CLI",
	Long:    "CLI tool for Taiwan's Central Weather Administration Open Data API.",
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

`cmd/cwa-weather/main.go`:
```go
package main

func main() {
	Execute()
}
```

**Step 3: Verify it builds and runs**

Run: `go build -o bin/cwa-weather ./cmd/cwa-weather && ./bin/cwa-weather --version`
Expected: `cwa-weather version dev`

**Step 4: Commit**

```bash
git add cmd/cwa-weather/ go.mod go.sum
git commit -m "feat: CLI root command with cobra"
```

---

## Task 8: CLI `query` Subcommand (Generic)

**Files:**
- Create: `cmd/cwa-weather/query.go`

**Step 1: Implement query subcommand**

`cmd/cwa-weather/query.go`:
```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/spf13/cobra"
)

var queryParams []string

var queryCmd = &cobra.Command{
	Use:   "query [dataid]",
	Short: "Query any CWA API endpoint by dataset ID",
	Long: `Query any CWA Open Data API endpoint directly.

Example:
  cwa-weather query F-D0047-069 --param LocationName=板橋區
  cwa-weather query O-A0003-001 --param CountyName=新北市
  cwa-weather query W-C0034-005`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := os.Getenv("CWA_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("CWA_API_KEY environment variable is required")
		}

		c := cwa.NewClient(apiKey)
		dataID := args[0]

		params := make(map[string]string)
		for _, p := range queryParams {
			parts := strings.SplitN(p, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid param format: %s (expected key=value)", p)
			}
			params[parts[0]] = parts[1]
		}

		resp, err := c.Query(context.Background(), dataID, params)
		if err != nil {
			return err
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(resp)
	},
}

func init() {
	queryCmd.Flags().StringArrayVarP(&queryParams, "param", "p", nil, "Query parameters (key=value)")
	rootCmd.AddCommand(queryCmd)
}
```

**Step 2: Build and test manually**

Run: `go build -o bin/cwa-weather ./cmd/cwa-weather && ./bin/cwa-weather query --help`
Expected: Help text with usage examples

**Step 3: Commit**

```bash
git add cmd/cwa-weather/query.go
git commit -m "feat: generic query subcommand for any CWA endpoint"
```

---

## Task 9: CLI `forecast` Subcommand

**Files:**
- Create: `cmd/cwa-weather/forecast.go`

**Step 1: Implement forecast subcommand**

`cmd/cwa-weather/forecast.go`:
```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/spf13/cobra"
)

var (
	forecastCity string
	forecastTown string
)

var forecastCmd = &cobra.Command{
	Use:   "forecast",
	Short: "Get township-level weather forecast",
	Long: `Get township-level weather forecast for a city/town.

Example:
  cwa-weather forecast --city 新北市 --town 板橋區
  cwa-weather forecast --city 台北市
  cwa-weather forecast --city 臺北市 --town 中正區`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if forecastCity == "" {
			return fmt.Errorf("--city is required")
		}

		apiKey := os.Getenv("CWA_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("CWA_API_KEY environment variable is required")
		}

		c := cwa.NewClient(apiKey)

		resp, err := c.Forecast(context.Background(), forecastCity, forecastTown)
		if err != nil {
			return err
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(resp)
	},
}

func init() {
	forecastCmd.Flags().StringVar(&forecastCity, "city", "", "City name (e.g. 臺北市, 台北市)")
	forecastCmd.Flags().StringVar(&forecastTown, "town", "", "Town name (e.g. 中正區)")
	rootCmd.AddCommand(forecastCmd)
}
```

**Step 2: Build and verify**

Run: `go build -o bin/cwa-weather ./cmd/cwa-weather && ./bin/cwa-weather forecast --help`
Expected: Help text

**Step 3: Commit**

```bash
git add cmd/cwa-weather/forecast.go
git commit -m "feat: forecast subcommand"
```

---

## Task 10: CLI `observe` Subcommand

**Files:**
- Create: `cmd/cwa-weather/observe.go`

**Step 1: Implement observe subcommand**

`cmd/cwa-weather/observe.go`:
```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/spf13/cobra"
)

var (
	observeCity    string
	observeStation string
)

var observeCmd = &cobra.Command{
	Use:   "observe",
	Short: "Get real-time weather station observations",
	Long: `Get real-time observation data from weather stations.

Example:
  cwa-weather observe --city 新北市
  cwa-weather observe --station 淡水`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if observeCity == "" && observeStation == "" {
			return fmt.Errorf("either --city or --station is required")
		}
		if observeCity != "" && observeStation != "" {
			return fmt.Errorf("--city and --station are mutually exclusive")
		}

		apiKey := os.Getenv("CWA_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("CWA_API_KEY environment variable is required")
		}

		c := cwa.NewClient(apiKey)

		var opt cwa.ObserveOption
		if observeCity != "" {
			opt = cwa.ObserveByCity(observeCity)
		} else {
			opt = cwa.ObserveByStation(observeStation)
		}

		resp, err := c.Observe(context.Background(), opt)
		if err != nil {
			return err
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(resp)
	},
}

func init() {
	observeCmd.Flags().StringVar(&observeCity, "city", "", "City name (e.g. 新北市)")
	observeCmd.Flags().StringVar(&observeStation, "station", "", "Station name (e.g. 淡水)")
	rootCmd.AddCommand(observeCmd)
}
```

**Step 2: Build and verify**

Run: `go build -o bin/cwa-weather ./cmd/cwa-weather && ./bin/cwa-weather observe --help`
Expected: Help text

**Step 3: Commit**

```bash
git add cmd/cwa-weather/observe.go
git commit -m "feat: observe subcommand"
```

---

## Task 11: CLI `cities` Subcommand

**Files:**
- Create: `cmd/cwa-weather/cities.go`

**Step 1: Implement cities subcommand**

`cmd/cwa-weather/cities.go`:
```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/spf13/cobra"
)

var citiesTown string

var citiesCmd = &cobra.Command{
	Use:   "cities",
	Short: "List supported cities and towns",
	Long: `List all supported cities, or list towns for a specific city.

Example:
  cwa-weather cities
  cwa-weather cities --town 臺北市`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if citiesTown != "" {
			towns, err := cwa.Towns(citiesTown)
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(towns)
		}

		cities := cwa.Cities()
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(cities)
	},
}

func init() {
	citiesCmd.Flags().StringVar(&citiesTown, "town", "", "Show towns for a specific city")
	rootCmd.AddCommand(citiesCmd)
}
```

**Step 2: Build and verify**

Run: `go build -o bin/cwa-weather ./cmd/cwa-weather && ./bin/cwa-weather cities`
Expected: JSON array of 22 cities

**Step 3: Commit**

```bash
git add cmd/cwa-weather/cities.go
git commit -m "feat: cities subcommand"
```

---

## Task 12: End-to-End Smoke Test

**Files:**
- Create: `cmd/cwa-weather/cmd_test.go`

**Step 1: Write integration test**

`cmd/cwa-weather/cmd_test.go`:
```go
package main

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCLI_Version(t *testing.T) {
	// Act
	out, err := exec.Command("go", "run", ".", "--version").CombinedOutput()

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, string(out), "cwa-weather version")
}

func TestCLI_Cities(t *testing.T) {
	// Act
	out, err := exec.Command("go", "run", ".", "cities").CombinedOutput()

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, string(out), "臺北市")
}

func TestCLI_QueryHelp(t *testing.T) {
	// Act
	out, err := exec.Command("go", "run", ".", "query", "--help").CombinedOutput()

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, string(out), "dataid")
}
```

**Step 2: Run tests**

Run: `go test -v ./...`
Expected: ALL PASS

**Step 3: Commit**

```bash
git add cmd/cwa-weather/cmd_test.go
git commit -m "test: CLI end-to-end smoke tests"
```

---

## Task 13: goreleaser + GitHub Actions

**Files:**
- Create: `.goreleaser.yml`
- Create: `.github/workflows/release.yml`
- Create: `.github/workflows/ci.yml`

**Step 1: Create .goreleaser.yml**

```yaml
version: 2
project_name: cwa-weather

builds:
  - main: ./cmd/cwa-weather
    binary: cwa-weather
    goos:
      - darwin
      - linux
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w -X main.version={{.Version}}

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Os }}_{{ .Arch }}"

checksum:
  name_template: "checksums.txt"

changelog:
  sort: asc
```

**Step 2: Create CI workflow**

`.github/workflows/ci.yml`:
```yaml
name: CI
on:
  push:
    branches: [main]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - run: go test -v ./...
      - run: go build ./...
```

**Step 3: Create release workflow**

`.github/workflows/release.yml`:
```yaml
name: Release
on:
  push:
    tags:
      - "v*"

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**Step 4: Verify goreleaser config**

Run: `goreleaser check` (if installed) or just verify YAML is valid

**Step 5: Commit**

```bash
git add .goreleaser.yml .github/
git commit -m "ci: goreleaser config and GitHub Actions workflows"
```

---

## Task 14: README + LICENSE + Agent Skills

**Files:**
- Create: `README.md`
- Create: `LICENSE`
- Create: `skill/SKILL.md`
- Create: `skill/AGENT.md`

**Step 1: Create LICENSE (MIT)**

Standard MIT license with `kerkerj` as author.

**Step 2: Create README.md**

Include:
- Project description
- Installation (go install, binary download)
- Usage examples for all subcommands (query, forecast, observe, cities)
- API key setup
- Library usage example
- Agent skill reference

**Step 3: Create skill/SKILL.md and skill/AGENT.md**

Agent instructions covering:
- Available subcommands and when to use each
- How to interpret JSON output
- Common usage patterns
- Error handling

**Step 4: Commit**

```bash
git add README.md LICENSE skill/
git commit -m "docs: README, LICENSE, agent skill files"
```

---

## Task Summary

| Task | Description | Depends On |
|------|-------------|------------|
| 1 | Go module init + structure | — |
| 2 | Core HTTP client + generic Query | 1 |
| 3 | City/town mapping + 台→臺 | 1 |
| 4 | Static town data (embed) | 3 |
| 5 | Forecast convenience method | 2, 3 |
| 6 | Observe convenience method | 2 |
| 7 | CLI root command + version | 1 |
| 8 | CLI query subcommand | 2, 7 |
| 9 | CLI forecast subcommand | 5, 7 |
| 10 | CLI observe subcommand | 6, 7 |
| 11 | CLI cities subcommand | 4, 7 |
| 12 | End-to-end smoke test | 8, 9, 10, 11 |
| 13 | goreleaser + CI | 7 |
| 14 | README + LICENSE + skills | 12 |

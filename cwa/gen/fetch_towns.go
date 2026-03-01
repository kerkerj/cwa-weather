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

const (
	baseURL = "https://opendata.cwa.gov.tw/api/v1/rest/datastore/"
	outFile = "cwa/towns.json"
)

// cityDatasets maps city name to its 3-hour forecast dataset ID.
// These per-city datasets contain town-level Location entries.
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

// apiResponse mirrors the relevant parts of the CWA per-city API response.
type apiResponse struct {
	Records struct {
		Locations []struct {
			LocationsName string `json:"LocationsName"`
			Location      []struct {
				LocationName string `json:"LocationName"`
			} `json:"Location"`
		} `json:"Locations"`
	} `json:"records"`
}

func main() {
	apiKey := os.Getenv("CWA_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "CWA_API_KEY environment variable is required")
		os.Exit(1)
	}

	result := make(map[string][]string, len(cityDatasets))

	// Sort city names for deterministic output.
	cities := make([]string, 0, len(cityDatasets))
	for city := range cityDatasets {
		cities = append(cities, city)
	}
	sort.Strings(cities)

	for _, city := range cities {
		datasetID := cityDatasets[city]
		reqURL := fmt.Sprintf("%s%s?format=JSON&elementName=T", baseURL, datasetID)

		fmt.Printf("Fetching %s (%s)...\n", city, datasetID)

		req, err := http.NewRequest(http.MethodGet, reqURL, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create request for %s: %v\n", city, err)
			os.Exit(1)
		}
		req.Header.Set("Authorization", apiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to fetch %s: %v\n", city, err)
			os.Exit(1)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			fmt.Fprintf(os.Stderr, "API returned status %d for %s: %s\n", resp.StatusCode, city, string(body))
			os.Exit(1)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read response for %s: %v\n", city, err)
			os.Exit(1)
		}

		var apiResp apiResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse JSON for %s: %v\n", city, err)
			os.Exit(1)
		}

		// Extract town names from the first Locations entry.
		townSet := make(map[string]bool)
		for _, loc := range apiResp.Records.Locations {
			for _, town := range loc.Location {
				townSet[town.LocationName] = true
			}
		}

		towns := make([]string, 0, len(townSet))
		for town := range townSet {
			towns = append(towns, town)
		}
		sort.Strings(towns)

		result[city] = towns
		fmt.Printf("  → %d towns\n", len(towns))
	}

	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal result: %v\n", err)
		os.Exit(1)
	}

	// Append newline for clean file ending.
	out = append(out, '\n')

	if err := os.WriteFile(outFile, out, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write %s: %v\n", outFile, err)
		os.Exit(1)
	}

	fmt.Printf("\nWritten %s with %d cities\n", outFile, len(result))
}

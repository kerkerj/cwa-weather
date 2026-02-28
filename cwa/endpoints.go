package cwa

import (
	"fmt"
	"sort"
	"strings"
)

// cityDatasets maps city name to 3-hour forecast dataset ID.
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

// NormalizeCity replaces 台 with 臺 to match CWA's official naming.
func NormalizeCity(name string) string {
	return strings.ReplaceAll(name, "台", "臺")
}

// GetDatasetID returns the dataset ID for the given city name.
// It normalizes 台→臺 before lookup.
func GetDatasetID(city string) (string, error) {
	normalized := NormalizeCity(city)

	id, ok := cityDatasets[normalized]
	if !ok {
		return "", fmt.Errorf("city not found: %s", city)
	}

	return id, nil
}

// Cities returns a sorted list of all 22 supported city names.
func Cities() []string {
	cities := make([]string, 0, len(cityDatasets))
	for city := range cityDatasets {
		cities = append(cities, city)
	}

	sort.Strings(cities)

	return cities
}

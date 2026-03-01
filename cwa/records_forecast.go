package cwa

import (
	"encoding/json"
	"fmt"
)

// ForecastRecords represents the parsed records from a forecast API response.
type ForecastRecords struct {
	Locations []ForecastLocations `json:"Locations"`
}

// ForecastLocations represents a group of forecast locations (e.g. a city).
type ForecastLocations struct {
	DatasetDescription string             `json:"DatasetDescription"`
	LocationsName      string             `json:"LocationsName"`
	Dataid             string             `json:"Dataid"`
	Location           []ForecastLocation `json:"Location"`
}

// ForecastLocation represents a single forecast location (e.g. a town/district).
type ForecastLocation struct {
	LocationName   string            `json:"LocationName"`
	Geocode        string            `json:"Geocode"`
	WeatherElement []ForecastElement `json:"WeatherElement"`
}

// ForecastElement represents a weather element in a forecast (e.g. temperature).
type ForecastElement struct {
	ElementName string         `json:"ElementName"`
	Time        []ForecastTime `json:"Time"`
}

// ForecastTime represents a time entry in a forecast element.
type ForecastTime struct {
	DataTime     string              `json:"DataTime"`
	StartTime    string              `json:"StartTime"`
	EndTime      string              `json:"EndTime"`
	ElementValue []map[string]string `json:"ElementValue"`
}

// ParseForecastRecords parses the raw Records field into ForecastRecords.
func (r *Response) ParseForecastRecords() (*ForecastRecords, error) {
	var rec ForecastRecords
	if err := json.Unmarshal(r.Records, &rec); err != nil {
		return nil, fmt.Errorf("failed to parse forecast records: %w", err)
	}
	return &rec, nil
}

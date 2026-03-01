package cwa

import (
	"encoding/json"
	"fmt"
)

// OverviewRecords represents the parsed records from an overview forecast API response.
type OverviewRecords struct {
	DatasetDescription string             `json:"datasetDescription"`
	Location           []OverviewLocation `json:"location"`
}

// OverviewLocation represents a location in the overview forecast.
type OverviewLocation struct {
	LocationName   string            `json:"locationName"`
	WeatherElement []OverviewElement `json:"weatherElement"`
}

// OverviewElement represents a weather element in the overview forecast.
type OverviewElement struct {
	ElementName string         `json:"elementName"`
	Time        []OverviewTime `json:"time"`
}

// OverviewTime represents a time entry in an overview forecast element.
type OverviewTime struct {
	StartTime string            `json:"startTime"`
	EndTime   string            `json:"endTime"`
	Parameter OverviewParameter `json:"parameter"`
}

// OverviewParameter represents a parameter value in the overview forecast.
type OverviewParameter struct {
	ParameterName  string `json:"parameterName"`
	ParameterValue string `json:"parameterValue"`
	ParameterUnit  string `json:"parameterUnit"`
}

// ParseOverviewRecords parses the raw Records field into OverviewRecords.
func (r *Response) ParseOverviewRecords() (*OverviewRecords, error) {
	var rec OverviewRecords
	if err := json.Unmarshal(r.Records, &rec); err != nil {
		return nil, fmt.Errorf("failed to parse overview records: %w", err)
	}
	return &rec, nil
}

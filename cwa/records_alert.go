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

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

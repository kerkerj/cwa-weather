package cwa

import (
	"encoding/json"
	"fmt"
)

// ObserveRecords represents the parsed records from an observation API response.
type ObserveRecords struct {
	Station []ObserveStation `json:"Station"`
}

// ObserveStation represents a single observation station.
type ObserveStation struct {
	StationName    string                `json:"StationName"`
	StationId      string                `json:"StationId"`
	ObsTime        ObserveObsTime        `json:"ObsTime"`
	GeoInfo        ObserveGeoInfo        `json:"GeoInfo"`
	WeatherElement ObserveWeatherElement `json:"WeatherElement"`
}

// ObserveObsTime represents the observation time.
type ObserveObsTime struct {
	DateTime string `json:"DateTime"`
}

// ObserveGeoInfo represents geographic information for a station.
type ObserveGeoInfo struct {
	CountyName string `json:"CountyName"`
	TownName   string `json:"TownName"`
}

// ObserveWeatherElement represents the weather data from an observation.
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

// ObserveNow represents current precipitation data.
type ObserveNow struct {
	Precipitation string `json:"Precipitation"`
}

// ObserveDailyExt represents daily temperature extremes.
type ObserveDailyExt struct {
	DailyHigh ObserveTempInfo `json:"DailyHigh"`
	DailyLow  ObserveTempInfo `json:"DailyLow"`
}

// ObserveTempInfo represents temperature information for daily extremes.
type ObserveTempInfo struct {
	TemperatureInfo struct {
		AirTemperature string `json:"AirTemperature"`
	} `json:"TemperatureInfo"`
}

// ParseObserveRecords parses the raw Records field into ObserveRecords.
func (r *Response) ParseObserveRecords() (*ObserveRecords, error) {
	var rec ObserveRecords
	if err := json.Unmarshal(r.Records, &rec); err != nil {
		return nil, fmt.Errorf("failed to parse observe records: %w", err)
	}
	return &rec, nil
}

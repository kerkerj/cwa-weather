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

// ObserveWithElement filters by weather element names (comma-separated).
func ObserveWithElement(element string) ObserveOption {
	return func(params map[string]string) {
		params["WeatherElement"] = element
	}
}

// Observe returns real-time observation data.
func (c *Client) Observe(ctx context.Context, opts ...ObserveOption) (*Response, error) {
	if len(opts) == 0 {
		return nil, fmt.Errorf("observe requires at least one option (city or station)")
	}

	params := make(map[string]string)
	for _, opt := range opts {
		opt(params)
	}

	// Must have either CountyName or StationName
	if params["CountyName"] == "" && params["StationName"] == "" {
		return nil, fmt.Errorf("observe requires either city or station option")
	}

	return c.Query(ctx, observeDatasetID, params)
}

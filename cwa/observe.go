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

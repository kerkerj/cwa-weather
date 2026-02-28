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

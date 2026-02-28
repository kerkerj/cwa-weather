package cwa

import "context"

const alertDatasetID = "W-C0033-001"

// Alert returns current weather alerts.
// If city is empty, returns alerts for all cities.
func (c *Client) Alert(ctx context.Context, city string) (*Response, error) {
	params := make(map[string]string)
	if city != "" {
		params["locationName"] = NormalizeCity(city)
	}

	return c.Query(ctx, alertDatasetID, params)
}

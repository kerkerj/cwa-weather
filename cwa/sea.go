package cwa

import "context"

const seaDatasetID = "O-B0075-001"

// Sea returns marine observation data (海象).
// If station is empty, returns data for all stations.
func (c *Client) Sea(ctx context.Context, station string) (*Response, error) {
	params := make(map[string]string)
	if station != "" {
		params["StationName"] = station
	}

	return c.Query(ctx, seaDatasetID, params)
}

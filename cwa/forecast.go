package cwa

import "context"

// ForecastOption configures optional server-side filters for a Forecast query.
type ForecastOption struct {
	Element  string // ElementName filter (comma-separated)
	TimeFrom string // timeFrom filter (yyyy-MM-ddThh:mm:ss)
	TimeTo   string // timeTo filter (yyyy-MM-ddThh:mm:ss)
}

// Forecast returns township-level weather forecast for a city.
// If town is empty, returns all towns in the city.
// Optional ForecastOption can be passed to filter by element or time range.
func (c *Client) Forecast(ctx context.Context, city, town string, opts ...ForecastOption) (*Response, error) {
	datasetID, err := GetDatasetID(city)
	if err != nil {
		return nil, err
	}

	params := make(map[string]string)
	if town != "" {
		params["LocationName"] = NormalizeCity(town)
	}
	if len(opts) > 0 {
		opt := opts[0]
		if opt.Element != "" {
			params["ElementName"] = opt.Element
		}
		if opt.TimeFrom != "" {
			params["timeFrom"] = opt.TimeFrom
		}
		if opt.TimeTo != "" {
			params["timeTo"] = opt.TimeTo
		}
	}

	return c.Query(ctx, datasetID, params)
}

package cwa

import "context"

const overviewDatasetID = "F-C0032-001"

// OverviewOption configures optional server-side filters for an Overview query.
type OverviewOption struct {
	Element  string // elementName filter (comma-separated, e.g. Wx,PoP,CI,MinT,MaxT)
	TimeFrom string // timeFrom filter (yyyy-MM-ddThh:mm:ss)
	TimeTo   string // timeTo filter (yyyy-MM-ddThh:mm:ss)
}

// Overview returns 36-hour city-level weather forecast.
// If city is empty, returns all cities.
// Optional OverviewOption can be passed to filter by element or time range.
func (c *Client) Overview(ctx context.Context, city string, opts ...OverviewOption) (*Response, error) {
	params := make(map[string]string)
	if city != "" {
		params["locationName"] = NormalizeCity(city)
	}
	if len(opts) > 0 {
		opt := opts[0]
		if opt.Element != "" {
			params["elementName"] = opt.Element
		}
		if opt.TimeFrom != "" {
			params["timeFrom"] = opt.TimeFrom
		}
		if opt.TimeTo != "" {
			params["timeTo"] = opt.TimeTo
		}
	}

	return c.Query(ctx, overviewDatasetID, params)
}

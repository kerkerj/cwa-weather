package cwa

import "context"

const typhoonDatasetID = "W-C0034-005"

// TyphoonOption configures optional server-side filters for a Typhoon query.
type TyphoonOption struct {
	CwaTdNo string // Tropical depression number
	Dataset string // AnalysisData or ForecastData
}

// Typhoon returns tropical cyclone tracking data.
// Optional TyphoonOption can be passed to filter by CwaTdNo or Dataset.
func (c *Client) Typhoon(ctx context.Context, opts ...TyphoonOption) (*Response, error) {
	params := make(map[string]string)
	if len(opts) > 0 {
		opt := opts[0]
		if opt.CwaTdNo != "" {
			params["CwaTdNo"] = opt.CwaTdNo
		}
		if opt.Dataset != "" {
			params["Dataset"] = opt.Dataset
		}
	}

	return c.Query(ctx, typhoonDatasetID, params)
}

package cwa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const defaultBaseURL = "https://opendata.cwa.gov.tw/api/v1/rest/datastore"

// Response represents the top-level CWA API response.
type Response struct {
	Success string          `json:"success"`
	Result  Result          `json:"result"`
	Records json.RawMessage `json:"records"`
}

// Result contains metadata about the API response.
type Result struct {
	ResourceID string  `json:"resource_id"`
	Fields     []Field `json:"fields"`
}

// Field describes a single field in the API result metadata.
type Field struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// Client is the CWA open data API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new CWA API client with the given API key.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetBaseURL overrides the base URL (useful for testing with httptest).
func (c *Client) SetBaseURL(u string) {
	c.baseURL = u
}

// Query sends a GET request to the CWA API for the given dataID and params.
func (c *Client) Query(ctx context.Context, dataID string, params map[string]string) (*Response, error) {
	reqURL, err := c.buildURL(dataID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (c *Client) buildURL(dataID string, params map[string]string) (string, error) {
	u, err := url.Parse(c.baseURL + "/" + dataID)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("Authorization", c.apiKey)
	q.Set("format", "JSON")

	for k, v := range params {
		q.Set(k, v)
	}

	u.RawQuery = q.Encode()

	return u.String(), nil
}

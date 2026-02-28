package cwa_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOverview(t *testing.T) {
	// Arrange
	fixture, err := os.ReadFile("../testdata/F-C0032-001.json")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/F-C0032-001", r.URL.Path)
		assert.Equal(t, "嘉義縣", r.URL.Query().Get("locationName"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Overview(context.Background(), "嘉義縣")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
	assert.Equal(t, "F-C0032-001", resp.Result.ResourceID)
}

func TestOverview_AllCities(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/F-C0032-001", r.URL.Path)
		assert.Empty(t, r.URL.Query().Get("locationName"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":"true","result":{"resource_id":"F-C0032-001","fields":[]},"records":{}}`))
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Overview(context.Background(), "")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

func TestOverview_WithElement(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Wx,PoP", r.URL.Query().Get("elementName"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":"true","result":{"resource_id":"F-C0032-001","fields":[]},"records":{}}`))
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Overview(context.Background(), "臺北市", cwa.OverviewOption{
		Element: "Wx,PoP",
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

func TestOverview_CityAlias(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "臺北市", r.URL.Query().Get("locationName"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":"true","result":{"resource_id":"F-C0032-001","fields":[]},"records":{}}`))
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Overview(context.Background(), "台北市")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

func TestOverview_WithTimeRange(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "2026-03-01T06:00:00", r.URL.Query().Get("timeFrom"))
		assert.Equal(t, "2026-03-01T18:00:00", r.URL.Query().Get("timeTo"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":"true","result":{"resource_id":"F-C0032-001","fields":[]},"records":{}}`))
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Overview(context.Background(), "臺北市", cwa.OverviewOption{
		TimeFrom: "2026-03-01T06:00:00",
		TimeTo:   "2026-03-01T18:00:00",
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

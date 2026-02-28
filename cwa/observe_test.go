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

func TestObserve_ByCity(t *testing.T) {
	// Arrange
	fixture, err := os.ReadFile("../testdata/O-A0003-001-新北市.json")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/O-A0003-001", r.URL.Path)
		assert.Equal(t, "新北市", r.URL.Query().Get("CountyName"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Observe(context.Background(), cwa.ObserveByCity("新北市"))

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

func TestObserve_ByStation(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/O-A0003-001", r.URL.Path)
		assert.Equal(t, "淡水", r.URL.Query().Get("StationName"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":"true","result":{"resource_id":"O-A0003-001","fields":[]},"records":{}}`))
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Observe(context.Background(), cwa.ObserveByStation("淡水"))

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

func TestObserve_NoOption(t *testing.T) {
	// Arrange
	c := cwa.NewClient("test-key")

	// Act
	resp, err := c.Observe(context.Background(), nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "observe requires either city or station option")
}

func TestObserve_CityAlias(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "臺北市", r.URL.Query().Get("CountyName"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":"true","result":{"resource_id":"O-A0003-001","fields":[]},"records":{}}`))
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Observe(context.Background(), cwa.ObserveByCity("台北市"))

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

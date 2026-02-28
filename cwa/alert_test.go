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

func TestAlert(t *testing.T) {
	// Arrange
	fixture, err := os.ReadFile("../testdata/W-C0033-001-天氣特報.json")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/W-C0033-001", r.URL.Path)
		assert.Equal(t, "臺北市", r.URL.Query().Get("locationName"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Alert(context.Background(), "臺北市")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
	assert.Equal(t, "W-C0033-001", resp.Result.ResourceID)
}

func TestAlert_AllCities(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/W-C0033-001", r.URL.Path)
		assert.Empty(t, r.URL.Query().Get("locationName"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":"true","result":{"resource_id":"W-C0033-001","fields":[]},"records":{}}`))
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Alert(context.Background(), "")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

func TestAlert_CityAlias(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "臺北市", r.URL.Query().Get("locationName"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":"true","result":{"resource_id":"W-C0033-001","fields":[]},"records":{}}`))
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Alert(context.Background(), "台北市")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

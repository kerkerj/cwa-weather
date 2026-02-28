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

func TestForecast(t *testing.T) {
	// Arrange
	fixture, err := os.ReadFile("../testdata/F-D0047-069-板橋區.json")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/F-D0047-069", r.URL.Path)
		assert.Equal(t, "板橋區", r.URL.Query().Get("LocationName"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fixture)
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Forecast(context.Background(), "新北市", "板橋區")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

func TestForecast_CityAlias(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/F-D0047-061", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":"true","result":{"resource_id":"F-D0047-061","fields":[]},"records":{}}`))
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Forecast(context.Background(), "台北市", "")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

func TestForecast_UnknownCity(t *testing.T) {
	// Arrange
	c := cwa.NewClient("test-key")

	// Act
	resp, err := c.Forecast(context.Background(), "火星市", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "city not found")
}

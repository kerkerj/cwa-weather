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

func TestSea(t *testing.T) {
	// Arrange
	fixture, err := os.ReadFile("../testdata/O-B0075-001-海象.json")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/O-B0075-001", r.URL.Path)
		assert.Equal(t, "富貴角", r.URL.Query().Get("StationName"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Sea(context.Background(), "富貴角")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

func TestSea_AllStations(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/O-B0075-001", r.URL.Path)
		assert.Empty(t, r.URL.Query().Get("StationName"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":"true","result":{"resource_id":"O-B0075-001","fields":[]},"records":{}}`))
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Sea(context.Background(), "")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

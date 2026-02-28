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

func TestNewClient(t *testing.T) {
	// Arrange & Act
	c := cwa.NewClient("test-api-key")

	// Assert
	assert.NotNil(t, c)
}

func TestNewClient_EmptyKey(t *testing.T) {
	// Arrange & Act
	c := cwa.NewClient("")

	// Assert
	assert.NotNil(t, c)
}

func TestQuery_Forecast(t *testing.T) {
	// Arrange
	fixture, err := os.ReadFile("../testdata/F-D0047-069-板橋區.json")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request path and params
		assert.Equal(t, "/F-D0047-069", r.URL.Path)
		assert.Equal(t, "test-key", r.Header.Get("Authorization"))
		assert.Equal(t, "JSON", r.URL.Query().Get("format"))
		assert.Equal(t, "板橋區", r.URL.Query().Get("LocationName"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Query(context.Background(), "F-D0047-069", map[string]string{
		"LocationName": "板橋區",
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
	assert.Equal(t, "F-D0047-069", resp.Result.ResourceID)
	assert.NotEmpty(t, resp.Result.Fields)
	assert.NotEmpty(t, resp.Records)
}

func TestQuery_InvalidAPIKey(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":"false","result":{"resource_id":"","fields":[]},"records":{}}`))
	}))
	defer server.Close()

	c := cwa.NewClient("bad-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Query(context.Background(), "F-D0047-069", nil)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "false", resp.Success)
}

func TestQuery_HTTPError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Query(context.Background(), "F-D0047-069", nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unexpected status code: 500")
}

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

func TestTyphoon(t *testing.T) {
	// Arrange
	fixture, err := os.ReadFile("../testdata/W-C0034-005-颱風.json")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/W-C0034-005", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Typhoon(context.Background())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
	assert.Equal(t, "W-C0034-005", resp.Result.ResourceID)
}

func TestTyphoon_WithCwaTdNo(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/W-C0034-005", r.URL.Path)
		assert.Equal(t, "03", r.URL.Query().Get("CwaTdNo"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":"true","result":{"resource_id":"W-C0034-005","fields":[]},"records":{}}`))
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Typhoon(context.Background(), cwa.TyphoonOption{
		CwaTdNo: "03",
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

func TestTyphoon_WithDataset(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/W-C0034-005", r.URL.Path)
		assert.Equal(t, "ForecastData", r.URL.Query().Get("Dataset"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":"true","result":{"resource_id":"W-C0034-005","fields":[]},"records":{}}`))
	}))
	defer server.Close()

	c := cwa.NewClient("test-key")
	c.SetBaseURL(server.URL)

	// Act
	resp, err := c.Typhoon(context.Background(), cwa.TyphoonOption{
		Dataset: "ForecastData",
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "true", resp.Success)
}

package cwa_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSeaRecords(t *testing.T) {
	// Arrange
	data, err := os.ReadFile("../testdata/O-B0075-001-海象.json")
	require.NoError(t, err)

	var resp cwa.Response
	require.NoError(t, json.Unmarshal(data, &resp))

	// Act
	rec, err := resp.ParseSeaRecords()

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, rec.SeaSurfaceObs.Location)
	loc := rec.SeaSurfaceObs.Location[0]
	assert.NotEmpty(t, loc.Station.StationID)
	require.NotEmpty(t, loc.StationObsTimes.StationObsTime)
}

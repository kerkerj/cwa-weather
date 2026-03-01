package cwa_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseObserveRecords(t *testing.T) {
	// Arrange
	data, err := os.ReadFile("../testdata/O-A0003-001-新北市.json")
	require.NoError(t, err)

	var resp cwa.Response
	require.NoError(t, json.Unmarshal(data, &resp))

	// Act
	rec, err := resp.ParseObserveRecords()

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, rec.Station)
	assert.Equal(t, "基隆", rec.Station[0].StationName)
	assert.NotEmpty(t, rec.Station[0].ObsTime.DateTime)
	assert.NotEmpty(t, rec.Station[0].WeatherElement.AirTemperature)
	assert.Equal(t, "基隆市", rec.Station[0].GeoInfo.CountyName)
}

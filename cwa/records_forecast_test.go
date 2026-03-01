package cwa_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseForecastRecords(t *testing.T) {
	// Arrange
	data, err := os.ReadFile("../testdata/F-D0047-069-板橋區.json")
	require.NoError(t, err)

	var resp cwa.Response
	require.NoError(t, json.Unmarshal(data, &resp))

	// Act
	rec, err := resp.ParseForecastRecords()

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, rec.Locations)
	assert.Equal(t, "新北市", rec.Locations[0].LocationsName)
	require.NotEmpty(t, rec.Locations[0].Location)
	assert.Equal(t, "板橋區", rec.Locations[0].Location[0].LocationName)
	require.NotEmpty(t, rec.Locations[0].Location[0].WeatherElement)
	assert.Equal(t, "溫度", rec.Locations[0].Location[0].WeatherElement[0].ElementName)
	require.NotEmpty(t, rec.Locations[0].Location[0].WeatherElement[0].Time)
}

package cwa_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseOverviewRecords(t *testing.T) {
	// Arrange
	data, err := os.ReadFile("../testdata/F-C0032-001.json")
	require.NoError(t, err)

	var resp cwa.Response
	require.NoError(t, json.Unmarshal(data, &resp))

	// Act
	rec, err := resp.ParseOverviewRecords()

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, rec.Location)
	assert.NotEmpty(t, rec.Location[0].LocationName)
	require.NotEmpty(t, rec.Location[0].WeatherElement)
	assert.Equal(t, "Wx", rec.Location[0].WeatherElement[0].ElementName)
	require.NotEmpty(t, rec.Location[0].WeatherElement[0].Time)
	assert.NotEmpty(t, rec.Location[0].WeatherElement[0].Time[0].Parameter.ParameterName)
}

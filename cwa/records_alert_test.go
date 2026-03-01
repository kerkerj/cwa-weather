package cwa_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAlertRecords(t *testing.T) {
	// Arrange
	data, err := os.ReadFile("../testdata/W-C0033-001-天氣特報.json")
	require.NoError(t, err)

	var resp cwa.Response
	require.NoError(t, json.Unmarshal(data, &resp))

	// Act
	rec, err := resp.ParseAlertRecords()

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, rec.Location)
	assert.NotEmpty(t, rec.Location[0].LocationName)
}

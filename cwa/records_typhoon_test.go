package cwa_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTyphoonRecords(t *testing.T) {
	// Arrange
	data, err := os.ReadFile("../testdata/W-C0034-005-颱風.json")
	require.NoError(t, err)

	var resp cwa.Response
	require.NoError(t, json.Unmarshal(data, &resp))

	// Act
	rec, err := resp.ParseTyphoonRecords()

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, rec.TropicalCyclones.TropicalCyclone)
	tc := rec.TropicalCyclones.TropicalCyclone[0]
	assert.NotEmpty(t, tc.TyphoonName)
	assert.NotEmpty(t, tc.CwaTyphoonName)
	require.NotNil(t, tc.AnalysisData)
	require.NotEmpty(t, tc.AnalysisData.Fix)
}

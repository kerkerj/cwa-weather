package cwa_test

import (
	"testing"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeCity(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"台北市 to 臺北市", "台北市", "臺北市"},
		{"台中市 to 臺中市", "台中市", "臺中市"},
		{"台南市 to 臺南市", "台南市", "臺南市"},
		{"台東縣 to 臺東縣", "台東縣", "臺東縣"},
		{"already correct 臺北市", "臺北市", "臺北市"},
		{"no 台 character", "新北市", "新北市"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			input := tt.input

			// Act
			got := cwa.NormalizeCity(input)

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetDatasetID(t *testing.T) {
	tests := []struct {
		name    string
		city    string
		wantID  string
		wantErr bool
	}{
		{"standard name 臺北市", "臺北市", "F-D0047-061", false},
		{"standard name 高雄市", "高雄市", "F-D0047-065", false},
		{"台 alias 台北市", "台北市", "F-D0047-061", false},
		{"台 alias 台中市", "台中市", "F-D0047-073", false},
		{"台 alias 台東縣", "台東縣", "F-D0047-037", false},
		{"unknown city", "火星市", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			city := tt.city

			// Act
			got, err := cwa.GetDatasetID(city)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "city not found")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, got)
			}
		})
	}
}

func TestTowns(t *testing.T) {
	// Act
	towns, err := cwa.Towns("臺北市")

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, towns, "中正區")
	assert.Contains(t, towns, "大安區")
}

func TestTowns_WithAlias(t *testing.T) {
	// Act
	towns, err := cwa.Towns("台北市")

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, towns)
}

func TestTowns_UnknownCity(t *testing.T) {
	// Act
	_, err := cwa.Towns("不存在市")

	// Assert
	assert.Error(t, err)
}

func TestCities(t *testing.T) {
	// Arrange
	// (no setup needed)

	// Act
	cities := cwa.Cities()

	// Assert
	assert.Len(t, cities, 22, "should have 22 cities")
	assert.Contains(t, cities, "臺北市")
	assert.Contains(t, cities, "新北市")
	assert.Contains(t, cities, "高雄市")

	// Verify sorted order
	for i := 1; i < len(cities); i++ {
		assert.LessOrEqual(t, cities[i-1], cities[i], "cities should be sorted")
	}
}

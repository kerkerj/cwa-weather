package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLI_Version(t *testing.T) {
	// Act
	out, err := exec.Command("go", "run", ".", "--version").CombinedOutput()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, string(out), "cwa-weather version")
}

func TestCLI_Cities(t *testing.T) {
	// Act
	out, err := exec.Command("go", "run", ".", "cities").CombinedOutput()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, string(out), "臺北市")
	assert.Contains(t, string(out), "新北市")
	assert.Contains(t, string(out), "高雄市")
}

func TestCLI_QueryHelp(t *testing.T) {
	// Act
	out, err := exec.Command("go", "run", ".", "query", "--help").CombinedOutput()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, string(out), "dataid")
}

func TestCLI_ForecastHelp(t *testing.T) {
	// Act
	out, err := exec.Command("go", "run", ".", "forecast", "--help").CombinedOutput()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, string(out), "--city")
	assert.Contains(t, string(out), "--town")
}

func TestCLI_ObserveHelp(t *testing.T) {
	// Act
	out, err := exec.Command("go", "run", ".", "observe", "--help").CombinedOutput()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, string(out), "--city")
	assert.Contains(t, string(out), "--station")
}

func TestCLI_OverviewHelp(t *testing.T) {
	// Act
	out, err := exec.Command("go", "run", ".", "overview", "--help").CombinedOutput()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, string(out), "--city")
}

func TestCLI_AlertHelp(t *testing.T) {
	// Act
	out, err := exec.Command("go", "run", ".", "alert", "--help").CombinedOutput()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, string(out), "--city")
}

func TestCLI_TyphoonHelp(t *testing.T) {
	// Act
	out, err := exec.Command("go", "run", ".", "typhoon", "--help").CombinedOutput()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, string(out), "--td-no")
}

func TestCLI_SeaHelp(t *testing.T) {
	// Act
	out, err := exec.Command("go", "run", ".", "sea", "--help").CombinedOutput()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, string(out), "--station")
}

func TestCLI_Cities_TextOutput(t *testing.T) {
	// Arrange + Act
	out, err := exec.Command("go", "run", ".", "cities").CombinedOutput()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, string(out), "臺北市")
	assert.NotContains(t, string(out), "[") // not JSON
}

func TestCLI_NoAPIKey(t *testing.T) {
	// Arrange
	cmd := exec.Command("go", "run", ".", "forecast", "--city", "臺北市")
	cmd.Env = append(os.Environ(), "CWA_API_KEY=")

	// Act
	out, err := cmd.CombinedOutput()

	// Assert
	assert.Error(t, err)
	assert.Contains(t, string(out), "CWA_API_KEY")
	assert.Contains(t, string(out), "opendata.cwa.gov.tw")
}

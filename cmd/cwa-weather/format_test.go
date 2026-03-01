package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureStdout redirects os.Stdout to a pipe and returns the written bytes after fn returns.
func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r) //nolint:errcheck
	return buf.String()
}

// ──────────────────────────────────────────────
// dash
// ──────────────────────────────────────────────

func TestDash(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "sentinel -99 replaced with dash",
			input: "-99",
			want:  "-",
		},
		{
			name:  "sentinel -99.0 replaced with dash",
			input: "-99.0",
			want:  "-",
		},
		{
			name:  "normal float passed through",
			input: "25.5",
			want:  "25.5",
		},
		{
			name:  "empty string passed through",
			input: "",
			want:  "",
		},
		{
			name:  "similar but non-sentinel value passed through",
			input: "-99.5",
			want:  "-99.5",
		},
		{
			name:  "zero value passed through",
			input: "0",
			want:  "0",
		},
		{
			name:  "positive integer passed through",
			input: "100",
			want:  "100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (test data already in tt)

			// Act
			got := dash(tt.input)

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}

// ──────────────────────────────────────────────
// formatElementValue
// ──────────────────────────────────────────────

func TestFormatElementValue(t *testing.T) {
	tests := []struct {
		name  string
		input []map[string]string
		want  string
	}{
		{
			name:  "empty slice returns empty string",
			input: []map[string]string{},
			want:  "",
		},
		{
			name: "single meaningful key returns its value",
			input: []map[string]string{
				{"Temperature": "23"},
			},
			want: "23",
		},
		{
			name: "key ending in Code is excluded",
			input: []map[string]string{
				{"Weather": "多雲", "WeatherCode": "04"},
			},
			want: "多雲",
		},
		{
			name: "key ending in Scale is excluded",
			input: []map[string]string{
				{"WindSpeed": "2", "BeaufortScale": "1"},
			},
			want: "2",
		},
		{
			name: "two meaningful keys joined with slash",
			input: []map[string]string{
				{"ComfortIndex": "22", "ComfortIndexDescription": "舒適"},
			},
			want: "22 / 舒適",
		},
		{
			name: "multiple maps in slice combined in order",
			input: []map[string]string{
				{"Temperature": "23"},
				{"Weather": "晴"},
			},
			want: "23 / 晴",
		},
		{
			name: "map with only Code and Scale keys returns empty",
			input: []map[string]string{
				{"WeatherCode": "04", "BeaufortScale": "3"},
			},
			want: "",
		},
		{
			name: "mixed: one map all excluded, one map has value",
			input: []map[string]string{
				{"WeatherCode": "04"},
				{"Weather": "晴天"},
			},
			want: "晴天",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (test data already in tt)

			// Act
			got := formatElementValue(tt.input)

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}

// ──────────────────────────────────────────────
// filterForecastByDays
// ──────────────────────────────────────────────

func makeStartTimeForecastTime(startTime, endTime, value string) cwa.ForecastTime {
	return cwa.ForecastTime{
		StartTime:    startTime,
		EndTime:      endTime,
		ElementValue: []map[string]string{{"Temperature": value}},
	}
}

func makeDataTimeForecastTime(dataTime, value string) cwa.ForecastTime {
	return cwa.ForecastTime{
		DataTime:     dataTime,
		ElementValue: []map[string]string{{"Temperature": value}},
	}
}

func TestFilterForecastByDays(t *testing.T) {
	// Arrange — fixed "now" at 2026-03-01 08:00 +08:00
	loc := time.FixedZone("CST", 8*60*60)
	now := time.Date(2026, 3, 1, 8, 0, 0, 0, loc)

	todayEntry := makeStartTimeForecastTime(
		"2026-03-01T06:00:00+08:00",
		"2026-03-01T12:00:00+08:00",
		"20",
	)
	tomorrowEntry := makeStartTimeForecastTime(
		"2026-03-02T06:00:00+08:00",
		"2026-03-02T12:00:00+08:00",
		"21",
	)
	dayAfterEntry := makeStartTimeForecastTime(
		"2026-03-03T06:00:00+08:00",
		"2026-03-03T12:00:00+08:00",
		"22",
	)
	beyondEntry := makeStartTimeForecastTime(
		"2026-03-04T06:00:00+08:00",
		"2026-03-04T12:00:00+08:00",
		"23",
	)
	badTimestampEntry := cwa.ForecastTime{
		StartTime:    "not-a-timestamp",
		ElementValue: []map[string]string{{"Temperature": "99"}},
	}
	dataTimeEntry := makeDataTimeForecastTime("2026-03-01T09:00:00+08:00", "25")

	allEntries := []cwa.ForecastTime{
		todayEntry, tomorrowEntry, dayAfterEntry, beyondEntry,
	}

	tests := []struct {
		name      string
		times     []cwa.ForecastTime
		days      int
		wantCount int
		wantFirst string // StartTime or DataTime of first result
	}{
		{
			name:      "days=1 includes only today",
			times:     allEntries,
			days:      1,
			wantCount: 1,
			wantFirst: "2026-03-01T06:00:00+08:00",
		},
		{
			name:      "days=2 includes today and tomorrow",
			times:     allEntries,
			days:      2,
			wantCount: 2,
			wantFirst: "2026-03-01T06:00:00+08:00",
		},
		{
			name:      "days=3 includes all three days",
			times:     allEntries,
			days:      3,
			wantCount: 3,
			wantFirst: "2026-03-01T06:00:00+08:00",
		},
		{
			name:      "days=0 clamps to 1",
			times:     allEntries,
			days:      0,
			wantCount: 1,
			wantFirst: "2026-03-01T06:00:00+08:00",
		},
		{
			name:      "days=5 clamps to 3",
			times:     allEntries,
			days:      5,
			wantCount: 3,
			wantFirst: "2026-03-01T06:00:00+08:00",
		},
		{
			name:      "bad timestamps are skipped",
			times:     []cwa.ForecastTime{badTimestampEntry, todayEntry},
			days:      1,
			wantCount: 1,
			wantFirst: "2026-03-01T06:00:00+08:00",
		},
		{
			name:      "uses DataTime when StartTime is empty",
			times:     []cwa.ForecastTime{dataTimeEntry},
			days:      1,
			wantCount: 1,
			wantFirst: "", // StartTime is empty; check via DataTime
		},
		{
			name:      "empty input returns empty result",
			times:     []cwa.ForecastTime{},
			days:      1,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (test data already in tt)

			// Act
			got := filterForecastByDays(tt.times, now, tt.days)

			// Assert
			require.Len(t, got, tt.wantCount)
			if tt.wantCount > 0 && tt.wantFirst != "" {
				assert.Equal(t, tt.wantFirst, got[0].StartTime)
			}
		})
	}
}

func TestFilterForecastByDays_UsesDataTime(t *testing.T) {
	// Arrange
	loc := time.FixedZone("CST", 8*60*60)
	now := time.Date(2026, 3, 1, 8, 0, 0, 0, loc)

	entry := makeDataTimeForecastTime("2026-03-01T09:00:00+08:00", "25")

	// Act
	got := filterForecastByDays([]cwa.ForecastTime{entry}, now, 1)

	// Assert
	require.Len(t, got, 1)
	assert.Equal(t, "2026-03-01T09:00:00+08:00", got[0].DataTime)
	assert.Empty(t, got[0].StartTime)
}

// ──────────────────────────────────────────────
// formatTimeRange
// ──────────────────────────────────────────────

func TestFormatTimeRange(t *testing.T) {
	tests := []struct {
		name  string
		start string
		end   string
		want  string
	}{
		{
			name:  "same day shows compact time range",
			start: "2026-03-01T06:00:00+08:00",
			end:   "2026-03-01T09:00:00+08:00",
			want:  "03/01 06:00~09:00",
		},
		{
			name:  "cross-day shows full date on both sides",
			start: "2026-03-01T21:00:00+08:00",
			end:   "2026-03-02T00:00:00+08:00",
			want:  "03/01 21:00 ~ 03/02 00:00",
		},
		{
			name:  "invalid start falls back to concatenation",
			start: "not-a-time",
			end:   "2026-03-01T09:00:00+08:00",
			want:  "not-a-time ~ 2026-03-01T09:00:00+08:00",
		},
		{
			name:  "invalid end falls back to concatenation",
			start: "2026-03-01T06:00:00+08:00",
			end:   "bad-end",
			want:  "2026-03-01T06:00:00+08:00 ~ bad-end",
		},
		{
			name:  "both invalid falls back to concatenation",
			start: "bad",
			end:   "also-bad",
			want:  "bad ~ also-bad",
		},
		{
			name:  "same day midnight boundary",
			start: "2026-03-01T00:00:00+08:00",
			end:   "2026-03-01T06:00:00+08:00",
			want:  "03/01 00:00~06:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (test data already in tt)

			// Act
			got := formatTimeRange(tt.start, tt.end)

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}

// ──────────────────────────────────────────────
// printForecastDetailed
// ──────────────────────────────────────────────

func TestPrintForecastDetailed_EmptyTimes(t *testing.T) {
	// Arrange
	elemName := "溫度"
	times := []cwa.ForecastTime{}

	// Act
	output := captureStdout(func() {
		printForecastDetailed(elemName, times)
	})

	// Assert
	assert.Empty(t, output, "empty times should produce no output")
}

func TestPrintForecastDetailed_WithData(t *testing.T) {
	// Arrange
	elemName := "溫度"
	times := []cwa.ForecastTime{
		{
			StartTime: "2026-03-01T06:00:00+08:00",
			EndTime:   "2026-03-01T12:00:00+08:00",
			ElementValue: []map[string]string{
				{"Temperature": "23"},
			},
		},
		{
			StartTime: "2026-03-01T12:00:00+08:00",
			EndTime:   "2026-03-01T18:00:00+08:00",
			ElementValue: []map[string]string{
				{"Temperature": "26"},
			},
		},
	}

	// Act
	output := captureStdout(func() {
		printForecastDetailed(elemName, times)
	})

	// Assert
	assert.Contains(t, output, "溫度", "output should contain element name")
	assert.Contains(t, output, "03/01 06:00~12:00", "output should contain formatted time range")
	assert.Contains(t, output, "23", "output should contain first temperature value")
	assert.Contains(t, output, "03/01 12:00~18:00", "output should contain second time range")
	assert.Contains(t, output, "26", "output should contain second temperature value")
}

func TestPrintForecastDetailed_UsesDataTime(t *testing.T) {
	// Arrange
	elemName := "天氣現象"
	times := []cwa.ForecastTime{
		{
			DataTime: "2026-03-01T09:00:00+08:00",
			ElementValue: []map[string]string{
				{"Weather": "晴"},
			},
		},
	}

	// Act
	output := captureStdout(func() {
		printForecastDetailed(elemName, times)
	})

	// Assert
	// When StartTime is empty, DataTime is used as-is (no formatting)
	assert.Contains(t, output, "2026-03-01T09:00:00+08:00")
	assert.Contains(t, output, "晴")
}

// ──────────────────────────────────────────────
// printForecastSummary
// ──────────────────────────────────────────────

func TestPrintForecastSummary_EmptyTimes(t *testing.T) {
	// Arrange
	loc := time.FixedZone("CST", 8*60*60)
	now := time.Date(2026, 3, 1, 8, 0, 0, 0, loc)

	// Act
	output := captureStdout(func() {
		printForecastSummary("溫度", []cwa.ForecastTime{}, now)
	})

	// Assert
	assert.Empty(t, output, "empty times should produce no output")
}

func TestPrintForecastSummary_TwoDays(t *testing.T) {
	// Arrange
	loc := time.FixedZone("CST", 8*60*60)
	// 2026-03-01 is a Sunday (日)
	now := time.Date(2026, 3, 1, 8, 0, 0, 0, loc)

	times := []cwa.ForecastTime{
		{
			StartTime: "2026-03-01T06:00:00+08:00",
			EndTime:   "2026-03-01T18:00:00+08:00",
			ElementValue: []map[string]string{
				{"Temperature": "23"},
			},
		},
		{
			StartTime: "2026-03-02T06:00:00+08:00",
			EndTime:   "2026-03-02T18:00:00+08:00",
			ElementValue: []map[string]string{
				{"Temperature": "21"},
			},
		},
	}

	// Act
	output := captureStdout(func() {
		printForecastSummary("溫度", times, now)
	})

	// Assert
	assert.Contains(t, output, "03/01", "output should contain first date")
	assert.Contains(t, output, "03/02", "output should contain second date")
	assert.Contains(t, output, "23", "output should contain first day temperature")
	assert.Contains(t, output, "21", "output should contain second day temperature")
	// Weekday headers: 2026-03-01 = Sunday (日), 2026-03-02 = Monday (一)
	assert.Contains(t, output, "日", "output should contain Sunday weekday label")
	assert.Contains(t, output, "一", "output should contain Monday weekday label")
}

func TestPrintForecastSummary_SkipsBadTimestamps(t *testing.T) {
	// Arrange
	loc := time.FixedZone("CST", 8*60*60)
	now := time.Date(2026, 3, 1, 8, 0, 0, 0, loc)

	times := []cwa.ForecastTime{
		{
			StartTime:    "not-valid",
			ElementValue: []map[string]string{{"Temperature": "99"}},
		},
		{
			StartTime: "2026-03-01T06:00:00+08:00",
			EndTime:   "2026-03-01T12:00:00+08:00",
			ElementValue: []map[string]string{
				{"Temperature": "23"},
			},
		},
	}

	// Act
	output := captureStdout(func() {
		printForecastSummary("溫度", times, now)
	})

	// Assert
	assert.NotContains(t, output, "99", "bad-timestamp entry value should not appear")
	assert.Contains(t, output, "23", "valid entry should appear")
	assert.Contains(t, output, "03/01", "valid date should appear")
}

func TestPrintForecastSummary_GroupsByDate(t *testing.T) {
	// Arrange
	loc := time.FixedZone("CST", 8*60*60)
	now := time.Date(2026, 3, 1, 8, 0, 0, 0, loc)

	// Three entries on the same day — should be grouped under one date header
	times := []cwa.ForecastTime{
		{
			StartTime: "2026-03-01T00:00:00+08:00",
			EndTime:   "2026-03-01T06:00:00+08:00",
			ElementValue: []map[string]string{{"Temperature": "18"}},
		},
		{
			StartTime: "2026-03-01T06:00:00+08:00",
			EndTime:   "2026-03-01T12:00:00+08:00",
			ElementValue: []map[string]string{{"Temperature": "23"}},
		},
		{
			StartTime: "2026-03-01T12:00:00+08:00",
			EndTime:   "2026-03-01T18:00:00+08:00",
			ElementValue: []map[string]string{{"Temperature": "25"}},
		},
	}

	// Act
	output := captureStdout(func() {
		printForecastSummary("溫度", times, now)
	})

	// Assert — "03/01" appears as a date header exactly once
	count := strings.Count(output, "03/01（")
	assert.Equal(t, 1, count, "date header should appear exactly once for same-day entries")
	assert.Contains(t, output, "18")
	assert.Contains(t, output, "23")
	assert.Contains(t, output, "25")
}

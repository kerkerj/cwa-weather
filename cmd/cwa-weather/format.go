package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/kerkerj/cwa-weather/cwa"
)

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func printHeader(title string) {
	fmt.Println(title)
}

// formatElementValue formats ElementValue entries as values only, joined by " / ".
// Keys ending in "Code" or "Scale" (numeric codes) are excluded.
// Keys are sorted for deterministic output order.
func formatElementValue(ev []map[string]string) string {
	if len(ev) == 0 {
		return ""
	}
	var parts []string
	for _, m := range ev {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if strings.HasSuffix(k, "Code") || strings.HasSuffix(k, "Scale") {
				continue
			}
			parts = append(parts, m[k])
		}
	}
	return strings.Join(parts, " / ")
}

// filterForecastByDays filters forecast time entries to only include the specified number of days.
func filterForecastByDays(times []cwa.ForecastTime, now time.Time, days int) []cwa.ForecastTime {
	if days <= 0 {
		days = 1
	}
	if days > 3 {
		days = 3
	}

	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	cutoff := todayStart.AddDate(0, 0, days)

	var result []cwa.ForecastTime
	for _, t := range times {
		ts := t.StartTime
		if ts == "" {
			ts = t.DataTime
		}
		parsed, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			continue
		}
		if parsed.Before(cutoff) {
			result = append(result, t)
		}
	}
	return result
}

// printForecastDetailed prints each time entry as a separate line.
func printForecastDetailed(elemName string, times []cwa.ForecastTime) {
	if len(times) == 0 {
		return
	}
	fmt.Printf("  %s\n", elemName)
	for _, t := range times {
		timeStr := t.DataTime
		if t.StartTime != "" {
			timeStr = formatTimeRange(t.StartTime, t.EndTime)
		}
		val := formatElementValue(t.ElementValue)
		fmt.Printf("    %-28s  %s\n", timeStr, val)
	}
}

// printForecastSummary groups time entries by day and prints day/night summaries.
func printForecastSummary(elemName string, times []cwa.ForecastTime, now time.Time) {
	if len(times) == 0 {
		return
	}

	// Group time entries by date
	type dayEntry struct {
		date  string
		times []cwa.ForecastTime
	}

	days := make(map[string]*dayEntry)
	var order []string

	for _, t := range times {
		ts := t.StartTime
		if ts == "" {
			ts = t.DataTime
		}
		parsed, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			continue
		}

		dateKey := parsed.Format("01/02")
		if _, exists := days[dateKey]; !exists {
			days[dateKey] = &dayEntry{date: dateKey}
			order = append(order, dateKey)
		}
		days[dateKey].times = append(days[dateKey].times, t)
	}

	weekdays := []string{"日", "一", "二", "三", "四", "五", "六"}

	for _, dateKey := range order {
		d := days[dateKey]
		parsed, err := time.Parse("01/02", d.date)
		weekday := ""
		if err == nil {
			full := time.Date(now.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, now.Location())
			weekday = weekdays[full.Weekday()]
		}

		fmt.Printf("\n%s（%s）\n", d.date, weekday)
		for _, t := range d.times {
			timeStr := t.DataTime
			if t.StartTime != "" {
				timeStr = formatTimeRange(t.StartTime, t.EndTime)
			}
			val := formatElementValue(t.ElementValue)
			fmt.Printf("  %-28s  %s\n", timeStr, val)
		}
	}
}

// formatTimeRange formats ISO timestamps to a shorter display format.
func formatTimeRange(start, end string) string {
	s, err1 := time.Parse(time.RFC3339, start)
	e, err2 := time.Parse(time.RFC3339, end)
	if err1 != nil || err2 != nil {
		return start + " ~ " + end
	}
	if s.Day() == e.Day() {
		return fmt.Sprintf("%s %02d:%02d~%02d:%02d",
			s.Format("01/02"), s.Hour(), s.Minute(), e.Hour(), e.Minute())
	}
	return fmt.Sprintf("%s %02d:%02d ~ %s %02d:%02d",
		s.Format("01/02"), s.Hour(), s.Minute(),
		e.Format("01/02"), e.Hour(), e.Minute())
}


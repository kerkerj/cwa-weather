package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/spf13/cobra"
)

var (
	forecastCity    string
	forecastTown    string
	forecastElement string
	forecastFrom    string
	forecastTo      string
	forecastDays    int
	forecastSummary bool
)

var forecastCmd = &cobra.Command{
	Use:   "forecast",
	Short: "Get township-level weather forecast",
	Long:  "Get township-level weather forecast for a city, optionally filtered by town.",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, err := getAPIKey()
		if err != nil {
			return err
		}

		client := cwa.NewClient(apiKey)

		var opts []cwa.ForecastOption
		if forecastElement != "" || forecastFrom != "" || forecastTo != "" {
			opts = append(opts, cwa.ForecastOption{
				Element:  forecastElement,
				TimeFrom: forecastFrom,
				TimeTo:   forecastTo,
			})
		}

		resp, err := client.Forecast(context.Background(), forecastCity, forecastTown, opts...)
		if err != nil {
			return fmt.Errorf("failed to get forecast: %w", err)
		}

		if jsonOutput {
			return printJSON(resp)
		}

		rec, err := resp.ParseForecastRecords()
		if err != nil {
			return err
		}

		now := time.Now()
		for _, locs := range rec.Locations {
			for _, loc := range locs.Location {
				printHeader(fmt.Sprintf("%s（%s）", loc.LocationName, locs.LocationsName))
				for _, elem := range loc.WeatherElement {
					if forecastElement == "" && elem.ElementName != "天氣預報綜合描述" {
						continue
					}

					// Filter time entries by --days
					filtered := filterForecastByDays(elem.Time, now, forecastDays)

					if forecastSummary {
						printForecastSummary(elem.ElementName, filtered, now)
					} else {
						printForecastDetailed(elem.ElementName, filtered)
					}
				}
			}
		}
		return nil
	},
}

func init() {
	forecastCmd.Flags().StringVar(&forecastCity, "city", "", "city name (required)")
	forecastCmd.Flags().StringVar(&forecastTown, "town", "", "town name (optional)")
	forecastCmd.Flags().StringVar(&forecastElement, "element", "", "filter weather elements (comma-separated, e.g. 溫度,天氣現象). Run without this flag to see all available names")
	forecastCmd.Flags().StringVar(&forecastFrom, "time-from", "", "start time filter (yyyy-MM-ddThh:mm:ss)")
	forecastCmd.Flags().StringVar(&forecastTo, "time-to", "", "end time filter (yyyy-MM-ddThh:mm:ss)")
	forecastCmd.Flags().IntVar(&forecastDays, "days", 1, "number of days to show in text output (1-3, default 1)")
	forecastCmd.Flags().BoolVar(&forecastSummary, "summary", false, "show day/night summary instead of 3-hour detail")
	_ = forecastCmd.MarkFlagRequired("city")
	rootCmd.AddCommand(forecastCmd)
}

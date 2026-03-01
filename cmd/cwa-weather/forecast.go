package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/spf13/cobra"
)

var (
	forecastCity    string
	forecastTown    string
	forecastElement string
	forecastFrom    string
	forecastTo      string
)

var forecastCmd = &cobra.Command{
	Use:   "forecast",
	Short: "Get township-level weather forecast",
	Long:  "Get township-level weather forecast for a city, optionally filtered by town.",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := os.Getenv("CWA_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("CWA_API_KEY environment variable is required")
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

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")

		return enc.Encode(resp)
	},
}

func init() {
	forecastCmd.Flags().StringVar(&forecastCity, "city", "", "city name (required)")
	forecastCmd.Flags().StringVar(&forecastTown, "town", "", "town name (optional)")
	forecastCmd.Flags().StringVar(&forecastElement, "element", "", "filter weather elements (comma-separated, e.g. 溫度,天氣現象). Run without this flag to see all available names")
	forecastCmd.Flags().StringVar(&forecastFrom, "time-from", "", "start time filter (yyyy-MM-ddThh:mm:ss)")
	forecastCmd.Flags().StringVar(&forecastTo, "time-to", "", "end time filter (yyyy-MM-ddThh:mm:ss)")
	_ = forecastCmd.MarkFlagRequired("city")
	rootCmd.AddCommand(forecastCmd)
}

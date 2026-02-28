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
	forecastCity string
	forecastTown string
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
		resp, err := client.Forecast(context.Background(), forecastCity, forecastTown)
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
	_ = forecastCmd.MarkFlagRequired("city")
	rootCmd.AddCommand(forecastCmd)
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/spf13/cobra"
)

var alertCity string

var alertCmd = &cobra.Command{
	Use:   "alert",
	Short: "Get current weather alerts (天氣特報)",
	Long:  "Get current weather alerts for all cities or a specific city.",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := os.Getenv("CWA_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("CWA_API_KEY environment variable is required")
		}

		client := cwa.NewClient(apiKey)

		resp, err := client.Alert(context.Background(), alertCity)
		if err != nil {
			return fmt.Errorf("failed to get alerts: %w", err)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")

		return enc.Encode(resp)
	},
}

func init() {
	alertCmd.Flags().StringVar(&alertCity, "city", "", "city name (optional, omit for all cities)")
	rootCmd.AddCommand(alertCmd)
}

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
	observeCity    string
	observeStation string
	observeElement string
)

var observeCmd = &cobra.Command{
	Use:   "observe",
	Short: "Get real-time weather observation",
	Long:  "Get real-time observation data filtered by city or station name (mutually exclusive).",
	RunE: func(cmd *cobra.Command, args []string) error {
		if observeCity == "" && observeStation == "" {
			return fmt.Errorf("either --city or --station is required")
		}
		if observeCity != "" && observeStation != "" {
			return fmt.Errorf("--city and --station are mutually exclusive")
		}

		apiKey := os.Getenv("CWA_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("CWA_API_KEY environment variable is required")
		}

		var opts []cwa.ObserveOption
		if observeCity != "" {
			opts = append(opts, cwa.ObserveByCity(observeCity))
		} else {
			opts = append(opts, cwa.ObserveByStation(observeStation))
		}
		if observeElement != "" {
			opts = append(opts, cwa.ObserveWithElement(observeElement))
		}

		client := cwa.NewClient(apiKey)
		resp, err := client.Observe(context.Background(), opts...)
		if err != nil {
			return fmt.Errorf("failed to get observation: %w", err)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")

		return enc.Encode(resp)
	},
}

func init() {
	observeCmd.Flags().StringVar(&observeCity, "city", "", "city name")
	observeCmd.Flags().StringVar(&observeStation, "station", "", "station name")
	observeCmd.Flags().StringVar(&observeElement, "element", "", "filter weather elements (comma-separated, e.g. AirTemperature,Weather). Run without this flag to see all available names")
	rootCmd.AddCommand(observeCmd)
}

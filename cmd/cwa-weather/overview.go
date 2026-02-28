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
	overviewCity    string
	overviewElement string
	overviewFrom    string
	overviewTo      string
)

var overviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Get 36-hour city-level weather forecast",
	Long:  "Get 36-hour city-level weather forecast for all cities or a specific city.",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := os.Getenv("CWA_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("CWA_API_KEY environment variable is required")
		}

		client := cwa.NewClient(apiKey)

		var opts []cwa.OverviewOption
		if overviewElement != "" || overviewFrom != "" || overviewTo != "" {
			opts = append(opts, cwa.OverviewOption{
				Element:  overviewElement,
				TimeFrom: overviewFrom,
				TimeTo:   overviewTo,
			})
		}

		resp, err := client.Overview(context.Background(), overviewCity, opts...)
		if err != nil {
			return fmt.Errorf("failed to get overview: %w", err)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")

		return enc.Encode(resp)
	},
}

func init() {
	overviewCmd.Flags().StringVar(&overviewCity, "city", "", "city name (optional, omit for all cities)")
	overviewCmd.Flags().StringVar(&overviewElement, "element", "", "filter weather elements (comma-separated, e.g. Wx,PoP,CI,MinT,MaxT)")
	overviewCmd.Flags().StringVar(&overviewFrom, "time-from", "", "start time filter (yyyy-MM-ddThh:mm:ss)")
	overviewCmd.Flags().StringVar(&overviewTo, "time-to", "", "end time filter (yyyy-MM-ddThh:mm:ss)")
	rootCmd.AddCommand(overviewCmd)
}

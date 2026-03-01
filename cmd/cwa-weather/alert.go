package main

import (
	"context"
	"fmt"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/spf13/cobra"
)

var alertCity string

var alertCmd = &cobra.Command{
	Use:   "alert",
	Short: "Get current weather alerts (天氣特報)",
	Long:  "Get current weather alerts for all cities or a specific city.",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, err := getAPIKey()
		if err != nil {
			return err
		}

		client := cwa.NewClient(apiKey)

		resp, err := client.Alert(context.Background(), alertCity)
		if err != nil {
			return fmt.Errorf("failed to get alerts: %w", err)
		}

		if jsonOutput {
			return printJSON(resp)
		}

		rec, err := resp.ParseAlertRecords()
		if err != nil {
			return err
		}

		hasAlert := false
		for _, loc := range rec.Location {
			if len(loc.HazardConditions.Hazards) > 0 {
				hasAlert = true
				printHeader(loc.LocationName)
				for _, h := range loc.HazardConditions.Hazards {
					fmt.Printf("  %s（%s）\n", h.Info.Phenomena, h.Info.Significance)
				}
			}
		}
		if !hasAlert {
			fmt.Println("目前無天氣特報。")
		}
		return nil
	},
}

func init() {
	alertCmd.Flags().StringVar(&alertCity, "city", "", "city name (optional, omit for all cities)")
	rootCmd.AddCommand(alertCmd)
}

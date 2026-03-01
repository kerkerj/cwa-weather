package main

import (
	"context"
	"fmt"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/spf13/cobra"
)

var seaStation string

var seaCmd = &cobra.Command{
	Use:   "sea",
	Short: "Get marine observation data (海象)",
	Long:  "Get marine observation data for all stations or a specific station.",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, err := getAPIKey()
		if err != nil {
			return err
		}

		client := cwa.NewClient(apiKey)

		resp, err := client.Sea(context.Background(), seaStation)
		if err != nil {
			return fmt.Errorf("failed to get marine data: %w", err)
		}

		if jsonOutput {
			return printJSON(resp)
		}

		rec, err := resp.ParseSeaRecords()
		if err != nil {
			return err
		}

		for _, loc := range rec.SeaSurfaceObs.Location {
			name := loc.Station.StationName
			if name == "" {
				name = loc.Station.StationID
			}
			printHeader(fmt.Sprintf("Station %s", name))
			maxObs := 3
			obs := loc.StationObsTimes.StationObsTime
			if len(obs) < maxObs {
				maxObs = len(obs)
			}
			for _, o := range obs[:maxObs] {
				we := o.WeatherElements
				seaTemp := we.SeaTemperature
				if seaTemp == "None" || seaTemp == "" {
					seaTemp = "N/A"
				}
				fmt.Printf("  %s  潮高: %sm  %s  海溫: %s\n",
					o.DateTime, we.TideHeight, we.TideLevel, seaTemp)
			}
			fmt.Println()
		}
		return nil
	},
}

func init() {
	seaCmd.Flags().StringVar(&seaStation, "station", "", "station name (optional, omit for all stations)")
	rootCmd.AddCommand(seaCmd)
}

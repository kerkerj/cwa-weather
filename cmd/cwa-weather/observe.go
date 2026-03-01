package main

import (
	"context"
	"fmt"

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

		apiKey, err := getAPIKey()
		if err != nil {
			return err
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

		if jsonOutput {
			return printJSON(resp)
		}

		rec, err := resp.ParseObserveRecords()
		if err != nil {
			return err
		}

		for _, stn := range rec.Station {
			printHeader(fmt.Sprintf("%s（%s, %s）  %s",
				stn.StationName, stn.GeoInfo.TownName, stn.GeoInfo.CountyName,
				stn.ObsTime.DateTime))
			we := stn.WeatherElement
			fmt.Printf("  天氣: %s    氣溫: %s°C    濕度: %s%%\n", dash(we.Weather), dash(we.AirTemperature), dash(we.RelativeHumidity))
			fmt.Printf("  風速: %s m/s    風向: %s°    氣壓: %s hPa\n", dash(we.WindSpeed), dash(we.WindDirection), dash(we.AirPressure))
			fmt.Printf("  今日降雨: %s mm    紫外線: %s\n", dash(we.Now.Precipitation), dash(we.UVIndex))
			fmt.Println()
		}
		return nil
	},
}

func init() {
	observeCmd.Flags().StringVar(&observeCity, "city", "", "city name")
	observeCmd.Flags().StringVar(&observeStation, "station", "", "station name")
	observeCmd.Flags().StringVar(&observeElement, "element", "", "filter weather elements (comma-separated, e.g. AirTemperature,Weather). Run without this flag to see all available names")
	rootCmd.AddCommand(observeCmd)
}

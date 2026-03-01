package main

import (
	"context"
	"fmt"
	"strings"

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
		apiKey, err := getAPIKey()
		if err != nil {
			return err
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

		if jsonOutput {
			return printJSON(resp)
		}

		rec, err := resp.ParseOverviewRecords()
		if err != nil {
			return err
		}

		for _, loc := range rec.Location {
			printHeader(loc.LocationName)
			// Build a map of time period -> element values
			type timeKey struct{ start, end string }
			type elemVal struct{ name, value, unit string }
			periods := make(map[timeKey][]elemVal)
			var order []timeKey

			for _, elem := range loc.WeatherElement {
				for _, t := range elem.Time {
					tk := timeKey{t.StartTime, t.EndTime}
					if _, exists := periods[tk]; !exists {
						order = append(order, tk)
					}
					periods[tk] = append(periods[tk], elemVal{
						name:  elem.ElementName,
						value: t.Parameter.ParameterName,
						unit:  t.Parameter.ParameterUnit,
					})
				}
			}

			for _, tk := range order {
				fmt.Printf("  %s ~ %s\n", tk.start, tk.end)
				elems := periods[tk]
				var parts []string
				for _, e := range elems {
					label := e.name
					val := e.value
					switch e.name {
					case "Wx":
						label = "天氣"
					case "PoP":
						label = "降雨"
						val += "%"
					case "MinT":
						label = "最低溫"
						val += "°C"
					case "MaxT":
						label = "最高溫"
						val += "°C"
					case "CI":
						label = "舒適度"
					}
					parts = append(parts, fmt.Sprintf("%s: %s", label, val))
				}
				fmt.Printf("    %s\n", strings.Join(parts, "    "))
			}
			fmt.Println()
		}
		return nil
	},
}

func init() {
	overviewCmd.Flags().StringVar(&overviewCity, "city", "", "city name (optional, omit for all cities)")
	overviewCmd.Flags().StringVar(&overviewElement, "element", "", "filter weather elements (comma-separated, e.g. Wx,PoP,CI,MinT,MaxT). Run without this flag to see all available names")
	overviewCmd.Flags().StringVar(&overviewFrom, "time-from", "", "start time filter (yyyy-MM-ddThh:mm:ss)")
	overviewCmd.Flags().StringVar(&overviewTo, "time-to", "", "end time filter (yyyy-MM-ddThh:mm:ss)")
	rootCmd.AddCommand(overviewCmd)
}

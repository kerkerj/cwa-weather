package main

import (
	"context"
	"fmt"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/spf13/cobra"
)

var (
	typhoonTdNo    string
	typhoonDataset string
)

var typhoonCmd = &cobra.Command{
	Use:   "typhoon",
	Short: "Get tropical cyclone tracking data (颱風)",
	Long:  "Get tropical cyclone tracking data for all active typhoons or filter by TD number / dataset.",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, err := getAPIKey()
		if err != nil {
			return err
		}

		client := cwa.NewClient(apiKey)

		var opts []cwa.TyphoonOption
		if typhoonTdNo != "" || typhoonDataset != "" {
			opts = append(opts, cwa.TyphoonOption{
				CwaTdNo: typhoonTdNo,
				Dataset: typhoonDataset,
			})
		}

		resp, err := client.Typhoon(context.Background(), opts...)
		if err != nil {
			return fmt.Errorf("failed to get typhoon data: %w", err)
		}

		if jsonOutput {
			return printJSON(resp)
		}

		rec, err := resp.ParseTyphoonRecords()
		if err != nil {
			return err
		}

		tcs := rec.TropicalCyclones.TropicalCyclone
		if len(tcs) == 0 {
			fmt.Println("目前無颱風資訊。")
			return nil
		}
		for _, tc := range tcs {
			printHeader(fmt.Sprintf("%s (%s) TD-%s", tc.CwaTyphoonName, tc.TyphoonName, tc.CwaTdNo))
			if tc.AnalysisData != nil && len(tc.AnalysisData.Fix) > 0 {
				fix := tc.AnalysisData.Fix[len(tc.AnalysisData.Fix)-1]
				fmt.Printf("  位置: %s°E %s°N    氣壓: %s hPa\n", fix.CoordinateLongitude, fix.CoordinateLatitude, fix.Pressure)
				fmt.Printf("  最大風速: %s m/s    陣風: %s m/s\n", fix.MaxWindSpeed, fix.MaxGustSpeed)
				fmt.Printf("  移動: %s km/h 往 %s\n", fix.MovingSpeed, fix.MovingDirection)
			}
			fmt.Println()
		}
		return nil
	},
}

func init() {
	typhoonCmd.Flags().StringVar(&typhoonTdNo, "td-no", "", "tropical depression number (optional, e.g. 03)")
	typhoonCmd.Flags().StringVar(&typhoonDataset, "dataset", "", "dataset filter (optional, e.g. AnalysisData or ForecastData)")
	rootCmd.AddCommand(typhoonCmd)
}

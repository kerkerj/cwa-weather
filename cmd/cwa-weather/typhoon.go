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
	typhoonTdNo    string
	typhoonDataset string
)

var typhoonCmd = &cobra.Command{
	Use:   "typhoon",
	Short: "Get tropical cyclone tracking data (颱風)",
	Long:  "Get tropical cyclone tracking data for all active typhoons or filter by TD number / dataset.",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := os.Getenv("CWA_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("CWA_API_KEY environment variable is required")
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

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")

		return enc.Encode(resp)
	},
}

func init() {
	typhoonCmd.Flags().StringVar(&typhoonTdNo, "td-no", "", "tropical depression number (optional, e.g. 03)")
	typhoonCmd.Flags().StringVar(&typhoonDataset, "dataset", "", "dataset filter (optional, e.g. AnalysisData or ForecastData)")
	rootCmd.AddCommand(typhoonCmd)
}

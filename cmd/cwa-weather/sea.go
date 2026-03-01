package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/spf13/cobra"
)

var seaStation string

var seaCmd = &cobra.Command{
	Use:   "sea",
	Short: "Get marine observation data (海象)",
	Long:  "Get marine observation data for all stations or a specific station.",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := os.Getenv("CWA_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("CWA_API_KEY environment variable is required")
		}

		client := cwa.NewClient(apiKey)

		resp, err := client.Sea(context.Background(), seaStation)
		if err != nil {
			return fmt.Errorf("failed to get marine data: %w", err)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")

		return enc.Encode(resp)
	},
}

func init() {
	seaCmd.Flags().StringVar(&seaStation, "station", "", "station name (optional, omit for all stations)")
	rootCmd.AddCommand(seaCmd)
}

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/spf13/cobra"
)

var citiesCity string

var citiesCmd = &cobra.Command{
	Use:   "cities",
	Short: "List supported cities or towns",
	Long:  "List all 22 supported cities. If --city is provided, list towns for that city.",
	RunE: func(cmd *cobra.Command, args []string) error {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")

		if citiesCity != "" {
			towns, err := cwa.Towns(citiesCity)
			if err != nil {
				return fmt.Errorf("failed to get towns: %w", err)
			}

			return enc.Encode(towns)
		}

		return enc.Encode(cwa.Cities())
	},
}

func init() {
	citiesCmd.Flags().StringVar(&citiesCity, "city", "", "show towns for this city")
	rootCmd.AddCommand(citiesCmd)
}

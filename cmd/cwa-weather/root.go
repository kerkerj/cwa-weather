package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version    = "dev"
	jsonOutput bool
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output raw JSON instead of human-readable text")
}

func getAPIKey() (string, error) {
	key := os.Getenv("CWA_API_KEY")
	if key == "" {
		return "", fmt.Errorf("CWA_API_KEY environment variable is not set — get a free key at https://opendata.cwa.gov.tw/userLogin")
	}
	return key, nil
}

var rootCmd = &cobra.Command{
	Use:     "cwa-weather",
	Short:   "CWA Open Data API CLI",
	Long:    "CLI tool for Taiwan's Central Weather Administration Open Data API.",
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

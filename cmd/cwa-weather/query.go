package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/kerkerj/cwa-weather/cwa"
	"github.com/spf13/cobra"
)

var queryParams []string

var queryCmd = &cobra.Command{
	Use:   "query [dataid]",
	Short: "Query any CWA dataset by ID",
	Long:  "Send a generic query to the CWA Open Data API for the specified dataset ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, err := getAPIKey()
		if err != nil {
			return err
		}

		dataID := args[0]

		params := make(map[string]string)
		for _, p := range queryParams {
			parts := strings.SplitN(p, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid param format %q, expected key=value", p)
			}
			params[parts[0]] = parts[1]
		}

		client := cwa.NewClient(apiKey)
		resp, err := client.Query(context.Background(), dataID, params)
		if err != nil {
			return fmt.Errorf("failed to query: %w", err)
		}

		return printJSON(resp)
	},
}

func init() {
	queryCmd.Flags().StringArrayVarP(&queryParams, "param", "p", nil, "query parameter in key=value format (repeatable)")
	rootCmd.AddCommand(queryCmd)
}

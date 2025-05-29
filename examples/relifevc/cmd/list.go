/*
Copyright © 2025 sottey

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/sottey/gognitive"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	dateFlag     string
	startFlag    string
	endFlag      string
	limitFlag    int
	timezoneFlag string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent lifelog entries",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := viper.GetString("api_key")
		if apiKey == "" {
			return fmt.Errorf("API key missing. Set it in config.json or pass via --config")
		}

		client := gognitive.NewClient(apiKey)

		// Allow CLI flag to override config value
		timezone := timezoneFlag
		if timezone == "" {
			timezone = viper.GetString("timezone")
		}

		logs, _, err := client.ListLifelogs(
			limitFlag,
			"", // cursor
			dateFlag,
			startFlag,
			endFlag,
			timezone,
		)
		if err != nil {
			return err
		}

		if len(logs) == 0 {
			fmt.Println("No lifelogs found.")
			return nil
		}

		for _, log := range logs {
			fmt.Printf("[%s] %s (ID: %s)\n", log.StartTime, log.Title, log.ID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&dateFlag, "date", "", "Date to fetch entries (YYYY-MM-DD)")
	listCmd.Flags().StringVar(&startFlag, "start", "", "Start datetime (YYYY-MM-DDTHH:MM:SS)")
	listCmd.Flags().StringVar(&endFlag, "end", "", "End datetime (YYYY-MM-DDTHH:MM:SS)")
	listCmd.Flags().IntVar(&limitFlag, "limit", 5, "Max number of entries to return")
	listCmd.Flags().StringVar(&timezoneFlag, "timezone", "", "Override timezone from config")
}

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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sottey/gognitive"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	saveFlag   bool
	allFlag    bool
	formatFlag string
	dirFlag    string
)

var exportCmd = &cobra.Command{
	Use:   "export [id]",
	Short: "Export lifelogs to local disk",
	Args: func(cmd *cobra.Command, args []string) error {
		if !saveFlag {
			return errors.New("you must pass --save")
		}
		if !allFlag && len(args) < 1 {
			return errors.New("you must provide an ID if not using --all")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := viper.GetString("api_key")
		if apiKey == "" {
			return fmt.Errorf("API key missing in config")
		}

		client := gognitive.NewClient(apiKey)

		// Fallback to config values if flags not set
		if !cmd.Flags().Changed("format") {
			if val := viper.GetString("export.format"); val != "" {
				formatFlag = val
			}
		}
		if !cmd.Flags().Changed("dir") {
			if val := viper.GetString("export.dir"); val != "" {
				dirFlag = val
			}
		}
		if dirFlag == "" {
			dirFlag = "./"
		}

		if allFlag {
			return exportAllLifelogs(client)
		}
		return exportOneLifelog(client, args[0])
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().BoolVar(&saveFlag, "save", false, "Save lifelog(s) to disk")
	exportCmd.Flags().BoolVar(&allFlag, "all", false, "Save all lifelogs")
	exportCmd.Flags().StringVar(&formatFlag, "format", "json", "Format to save (json|markdown)")
	exportCmd.Flags().StringVar(&dirFlag, "dir", "./", "Directory to save files to")
}

func exportAllLifelogs(client *gognitive.Client) error {
	cursor := ""
	page := 1

	for {
		logs, nextCursor, err := client.ListLifelogs(50, cursor, "", "", "", viper.GetString("timezone"))
		if err != nil {
			return err
		}

		for _, log := range logs {
			if err := saveLifelogToDisk(log); err != nil {
				fmt.Printf("Failed to save %s: %v\n", log.ID, err)
			}
		}

		if nextCursor == "" {
			break
		}
		cursor = nextCursor
		page++
	}

	fmt.Println("📦 All lifelogs exported successfully.")
	return nil
}

func exportOneLifelog(client *gognitive.Client, id string) error {
	log, err := client.GetLifelog(id)
	if err != nil {
		return err
	}
	return saveLifelogToDisk(*log)
}

func saveLifelogToDisk(log gognitive.Lifelog) error {
	t, err := time.Parse(time.RFC3339, log.StartTime)
	if err != nil {
		return fmt.Errorf("invalid start time for log %s: %w", log.ID, err)
	}
	dateFolder := t.Format("2006-01-02")
	outputDir := filepath.Join(dirFlag, dateFolder)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	dataFile := filepath.Join(outputDir, log.ID)
	metaFile := filepath.Join(outputDir, log.ID+".meta.json")
	var dataPath string

	if formatFlag == "markdown" {
		dataPath = dataFile + ".md"
		if _, err := os.Stat(dataPath); err == nil {
			fmt.Printf("Skipping %s (already exists)\n", log.ID)
			return nil
		}
		if err := os.WriteFile(dataPath, []byte(log.Markdown), 0644); err != nil {
			return err
		}
	} else {
		dataPath = dataFile + ".json"
		if _, err := os.Stat(dataPath); err == nil {
			fmt.Printf("Skipping %s (already exists)\n", log.ID)
			return nil
		}
		data, err := json.MarshalIndent(log, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(dataPath, data, 0644); err != nil {
			return err
		}
	}

	// Write meta.json
	meta := map[string]string{
		"id":         log.ID,
		"title":      log.Title,
		"start_time": log.StartTime,
		"end_time":   log.EndTime,
		"data_file":  filepath.Base(dataPath),
	}
	metaJSON, _ := json.MarshalIndent(meta, "", "  ")
	if err := os.WriteFile(metaFile, metaJSON, 0644); err != nil {
		return err
	}

	fmt.Printf("Saved %s\n", log.ID)
	return nil
}

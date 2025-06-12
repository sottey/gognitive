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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sottey/gognitive"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	lifelogID string
	exportAll bool
	repull    bool
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export lifelog entries to JSON files",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := viper.GetString("api_key")
		if token == "" {
			return fmt.Errorf("token is required")
		}

		client := gognitive.NewClient(token)
		outputDir := viper.GetString("export.dir")
		if outputDir == "" {
			return fmt.Errorf("output_dir is required")
		}
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("creating output dir: %w", err)
		}

		if lifelogID != "" {
			return exportSingleLifeLog(client, outputDir, lifelogID)
		}

		if exportAll {
			return exportAllFromAPI(client, outputDir)
		}

		return fmt.Errorf("please provide --id or --all")
	},
}

func exportSingleLifeLog(client *gognitive.Client, outputDir string, id string) error {
	entry, err := client.GetEnrichedLifelog(id)
	if err != nil {
		return fmt.Errorf("fetching lifelog %s: %w", id, err)
	}
	return saveEntry(outputDir, *entry)
}

func exportAllFromAPI(client *gognitive.Client, outputDir string) error {
	existing := map[string]bool{}
	files, err := os.ReadDir(outputDir)
	if err == nil {
		for _, f := range files {
			if strings.HasPrefix(f.Name(), "lifelog_") && strings.HasSuffix(f.Name(), ".json") {
				id := strings.TrimSuffix(strings.TrimPrefix(f.Name(), "lifelog_"), ".json")
				existing[id] = true
			}
		}
	}

	cursor := ""
	pageCount := 0
	for {
		if pageCount > 100 {
			return fmt.Errorf("too many pagination loops — possible infinite cursor")
		}
		pageCount++

		lifelogs, nextCursor, err := client.ListLifelogs(100, cursor, "", "", "", "")
		if err != nil {
			return fmt.Errorf("listing lifelogs: %w", err)
		}

		fmt.Printf("DEBUG: Got %d lifelogs, nextCursor = %q\n", len(lifelogs), nextCursor)

		for _, l := range lifelogs {
			if !repull && existing[l.ID] {
				continue
			}

			entry, err := client.GetEnrichedLifelog(l.ID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to fetch lifelog %s: %v\n", l.ID, err)
				continue
			}

			fmt.Printf("Exporting: %s %s\n", entry.StartTime, entry.ID)
			if err := saveEntry(outputDir, *entry); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to save lifelog %s: %v\n", l.ID, err)
				continue
			}

			// ✅ Add to existing to prevent re-fetch if seen again in paginated results
			existing[entry.ID] = true
		}

		if nextCursor == "" {
			break
		}
		cursor = nextCursor
		time.Sleep(1 * time.Second)
	}

	return nil
}

func saveEntry(dir string, entry gognitive.EnrichedLifelog) error {
	filename := filepath.Join(dir, fmt.Sprintf("lifelog_%s.json", entry.ID))
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(entry)
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&lifelogID, "id", "i", "", "ID of the lifelog to export")
	exportCmd.Flags().BoolVarP(&exportAll, "all", "a", false, "Export all lifelogs not already saved")
	exportCmd.Flags().BoolVar(&repull, "repull", false, "Re-pull and overwrite even if file already exists")
}

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
			return exportAllLifeLogs(client, outputDir)
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

func exportAllLifeLogs(client *gognitive.Client, outputDir string) error {
	files, err := os.ReadDir(outputDir)
	if err != nil {
		return fmt.Errorf("reading output dir: %w", err)
	}

	for _, file := range files {
		name := file.Name()
		if !strings.HasPrefix(name, "lifelog_") || !strings.HasSuffix(name, ".json") {
			continue
		}

		id := strings.TrimSuffix(strings.TrimPrefix(name, "lifelog_"), ".json")
		fullPath := filepath.Join(outputDir, fmt.Sprintf("lifelog_%s.json", id))

		if !repull {
			if _, err := os.Stat(fullPath); err == nil {
				continue // skip if file exists and not repulling
			}
		}

		entry, err := client.GetEnrichedLifelog(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to fetch lifelog %s: %v\n", id, err)
			continue
		}

		fmt.Printf("Exporting entry: %s %s\n", entry.StartTime, entry.ID)

		if err := saveEntry(outputDir, *entry); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to save lifelog %s: %v\n", id, err)
		}
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
	exportCmd.Flags().BoolVarP(&exportAll, "all", "a", false, "Export all lifelogs found in output directory")
	exportCmd.Flags().BoolVar(&repull, "repull", false, "Re-pull and overwrite even if file already exists")
}

/*package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sottey/gognitive"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export lifelogs to JSON files",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := gognitive.NewClient(viper.GetString("api_key"))
		exportDir := viper.GetString("export_path")
		if exportDir == "" {
			exportDir = "./export"
		}

		if exportAll {
			return exportAllLifelogs(client, exportDir)
		}

		if lifelogID == "" {
			return fmt.Errorf("must provide a lifelog ID or use --all")
		}

		return exportSingleLifelog(client, lifelogID, exportDir)
	},
}

var lifelogID string
var exportAll bool

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&lifelogID, "id", "i", "", "ID of the lifelog to export")
	exportCmd.Flags().BoolVarP(&exportAll, "all", "a", false, "Export all lifelogs")
	exportCmd.Flags().BoolP("repull", "r", false, "Re-pull all entries and overwrite existing files")
}

func exportAllLifelogs(client *gognitive.Client, baseDir string) error {
	cursor := ""
	for {
		lifelogs, nextCursor, err := client.ListLifelogs(50, cursor, "", "", "", "")
		if err != nil {
			return err
		}

		for _, l := range lifelogs {
			enriched, err := client.GetEnrichedLifelog(l.ID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error enriching lifelog %s: %v\n", l.ID, err)
				continue
			}

			dateFolder := filepath.Join(baseDir, l.StartTime[:10])
			os.MkdirAll(dateFolder, 0755)

			outPath := filepath.Join(dateFolder, l.ID+".json")
			outFile, err := os.Create(outPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error creating file %s: %v\n", outPath, err)
				continue
			}
			defer outFile.Close()

			enc := json.NewEncoder(outFile)
			enc.SetIndent("", "  ")
			if err := enc.Encode(enriched); err != nil {
				fmt.Fprintf(os.Stderr, "error writing file %s: %v\n", outPath, err)
				continue
			}
		}

		if nextCursor == "" {
			break
		}
		cursor = nextCursor
		time.Sleep(1 * time.Second)
	}
	return nil
}

func exportSingleLifelog(client *gognitive.Client, id, baseDir string) error {
	enriched, err := client.GetEnrichedLifelog(id)
	if err != nil {
		return err
	}

	dateFolder := filepath.Join(baseDir, enriched.StartTime[:10])
	os.MkdirAll(dateFolder, 0755)

	outPath := filepath.Join(dateFolder, enriched.ID+".json")
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	enc := json.NewEncoder(outFile)
	enc.SetIndent("", "  ")
	return enc.Encode(enriched)
}
*/

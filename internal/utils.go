package internal

import (
	"os"
	"path/filepath"
	"strings"
)

// generateTags is a placeholder that adds basic keyword-based tags
func GenerateTags(text string) []string {
	var tags []string
	keywords := []string{"meeting", "retailer", "support", "form", "demo", "onboarding", "call"}

	for _, kw := range keywords {
		if strings.Contains(strings.ToLower(text), kw) {
			tags = append(tags, kw)
		}
	}

	return tags
}

func GetAlreadyExportedIDs(dir string) (map[string]bool, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	for _, f := range files {
		if filepath.Ext(f.Name()) != ".json" {
			continue
		}
		base := strings.TrimSuffix(f.Name(), ".json")
		parts := strings.Split(base, "_")
		if len(parts) == 2 {
			seen[parts[1]] = true
		}
	}
	return seen, nil
}

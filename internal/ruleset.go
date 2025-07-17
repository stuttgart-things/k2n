package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadRulesets(dir string) ([]string, error) {
	var rulesets []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			rulesets = append(rulesets, fmt.Sprintf("Filename: %s\n%s", filepath.Base(path), content))
		}
		return nil
	})
	return rulesets, err
}

func LoadRulesetsIfExists(dir string) ([]string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, nil // Folder doesn't exist: skip
	}
	return LoadRulesets(dir)
}

func SplitAndTrimPaths(csv string) []string {
	rawPaths := strings.Split(csv, ",")
	paths := make([]string, 0, len(rawPaths))
	for _, p := range rawPaths {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			paths = append(paths, trimmed)
		}
	}
	return paths
}

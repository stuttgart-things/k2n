package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

// LoadCodeExamples loads all  files from the given directory and returns them as strings.
func LoadCodeExamples(dir string) ([]string, error) {
	var examples []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			examples = append(examples, string(content))
		}
		return nil
	})
	return examples, err
}

func LoadExampleFiles(paths []string) ([]string, error) {
	var examples []string
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", path, err)
		}
		examples = append(examples, string(content))
	}
	return examples, nil
}

func DeduplicateStrings(input []string) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, s := range input {
		if _, exists := seen[s]; !exists {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}

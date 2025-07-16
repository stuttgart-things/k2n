package internal

import (
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

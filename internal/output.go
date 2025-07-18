package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SaveOutput writes the content to the given destination or stdout if destination is empty.
// If destination is a directory, it saves each parsed file separately.
// If destination is a file, it combines all parsed files into one with filename comments.
func SaveOutput(destination, content string) error {
	files := ParseGeneratedFiles(content)

	if destination == "" {
		// Output to stdout
		for name, fileContent := range files {
			fmt.Printf("# %s\n%s\n\n", name, fileContent)
		}
		return nil
	}

	info, err := os.Stat(destination)
	if err == nil && info.IsDir() {
		// Destination is a directory: write each file separately
		for name, fileContent := range files {
			destPath := filepath.Join(destination, name)
			if err := os.WriteFile(destPath, []byte(fileContent), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %w", destPath, err)
			}
			fmt.Printf("Saved: %s\n", destPath)
		}
		return nil
	}

	// Destination is a single file: combine all content with filename comments
	var combined strings.Builder
	for name, fileContent := range files {
		combined.WriteString(fmt.Sprintf("# %s\n%s\n\n", name, fileContent))
	}

	if err := os.WriteFile(destination, []byte(combined.String()), 0644); err != nil {
		return fmt.Errorf("failed to write to %s: %w", destination, err)
	}

	fmt.Printf("Result written to %s\n", destination)
	return nil
}

// ParseGeneratedFiles splits AI-generated output into a map of filename -> content
func ParseGeneratedFiles(output string) map[string]string {
	parts := strings.Split(output, "---")
	files := make(map[string]string)

	for _, part := range parts {
		lines := strings.SplitN(strings.TrimSpace(part), "\n", 2)
		if len(lines) < 2 {
			continue
		}
		filename := strings.TrimSpace(lines[0])
		content := strings.TrimSpace(lines[1])
		files[filename] = content
	}

	return files
}

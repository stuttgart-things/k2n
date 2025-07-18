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
	// If destination is empty, print to stdout
	if destination == "" {
		fmt.Println(content)
		return nil
	}

	// Check if destination is a directory
	info, err := os.Stat(destination)
	if err == nil && info.IsDir() {
		// Write as separate files in the directory
		parsedFiles := ParseGeneratedFiles(content)
		for filename, fileContent := range parsedFiles {
			fullPath := filepath.Join(destination, filename)
			if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
				return fmt.Errorf("failed to create directory for %s: %w", fullPath, err)
			}
			if err := os.WriteFile(fullPath, []byte(fileContent), 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", fullPath, err)
			}
			fmt.Printf("Written %s\n", fullPath)
		}
		return nil
	}

	// If destination does not exist, check if it ends with path separator (means directory)
	if os.IsNotExist(err) && os.DirFS(destination) != nil && os.IsPathSeparator(destination[len(destination)-1]) {
		// Try to create directory recursively
		if err := os.MkdirAll(destination, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", destination, err)
		}
		// Then treat as directory again
		return SaveOutput(destination, content)
	}

	// Otherwise, treat as single file
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", destination, err)
	}

	err = os.WriteFile(destination, []byte(content), 0644)
	if err != nil {
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

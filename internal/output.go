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
	if destination == "" {
		fmt.Println(content)
		return nil
	}

	info, err := os.Stat(destination)
	if err == nil && info.IsDir() {
		return writeParsedFilesToDir(destination, content)
	}

	if os.IsNotExist(err) {
		// If destination ends with path separator or parsed files > 1, treat as dir
		parsedFiles := ParseGeneratedFiles(content)
		if strings.HasSuffix(destination, string(os.PathSeparator)) || len(parsedFiles) > 1 {
			if err := os.MkdirAll(destination, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destination, err)
			}
			return writeParsedFilesToDir(destination, content)
		}
	}

	// Otherwise treat as single file
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", destination, err)
	}

	if err := os.WriteFile(destination, []byte(content), 0644); err != nil {
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
		filename := sanitizeFilename(lines[0])
		content := strings.TrimSpace(lines[1])
		files[filename] = content
	}

	return files
}

func writeParsedFilesToDir(destination, content string) error {
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

func sanitizeFilename(name string) string {
	// Remove leading '#' and trim spaces
	name = strings.TrimSpace(strings.TrimPrefix(name, "#"))

	var b strings.Builder
	prevUnderscore := false

	for _, r := range name {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '.' || r == '_' || r == '-' {
			b.WriteRune(r)
			prevUnderscore = false
		} else {
			if !prevUnderscore {
				b.WriteRune('_')
				prevUnderscore = true
			}
		}
	}

	// Trim leading underscores
	result := strings.TrimLeft(b.String(), "_")

	return result
}

package internal

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveOutput(t *testing.T) {
	// Test 1: Save to stdout (destination = "")
	t.Run("stdout", func(t *testing.T) {
		// Capture stdout
		var buf bytes.Buffer
		stdout := os.Stdout
		defer func() { os.Stdout = stdout }()
		r, w, _ := os.Pipe()
		os.Stdout = w

		content := "test stdout content"
		err := SaveOutput("", content)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		w.Close()
		buf.ReadFrom(r)

		if !strings.Contains(buf.String(), content) {
			t.Fatalf("Expected stdout to contain %q, got %q", content, buf.String())
		}
	})

	// Setup temp dir for file tests
	tempDir := t.TempDir()

	// Test 2: Save to single file
	t.Run("single file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "nested", "file.txt")
		content := "file content"

		err := SaveOutput(filePath, content)
		if err != nil {
			t.Fatalf("Expected no error writing file, got %v", err)
		}

		// Verify file content
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read written file: %v", err)
		}
		if string(data) != content {
			t.Fatalf("Expected file content %q, got %q", content, data)
		}
	})

	// Test 3: Save to directory with multiple files
	t.Run("directory with parsed files", func(t *testing.T) {
		content := `
file1.txt
Hello World
---
subdir/file2.txt
Another file
`
		dirPath := filepath.Join(tempDir, "multifiles")

		err := SaveOutput(dirPath, content)
		if err != nil {
			t.Fatalf("Expected no error writing directory output, got %v", err)
		}

		tests := map[string]string{
			"file1.txt":        "Hello World",
			"subdir/file2.txt": "Another file",
		}

		for relPath, expectedContent := range tests {
			fullPath := filepath.Join(dirPath, relPath)
			data, err := os.ReadFile(fullPath)
			if err != nil {
				t.Fatalf("Expected file %s to exist, got error: %v", relPath, err)
			}
			if strings.TrimSpace(string(data)) != expectedContent {
				t.Errorf("Expected content of %s to be %q, got %q", relPath, expectedContent, data)
			}
		}
	})
}

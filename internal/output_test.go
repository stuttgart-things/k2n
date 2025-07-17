package internal

import (
	"os"
	"strings"
	"testing"
)

func TestSaveOutput(t *testing.T) {
	content := "test content"

	// Test stdout (destination empty)
	err := SaveOutput("", content)
	if err != nil {
		t.Errorf("Expected no error when writing to stdout, got: %v", err)
	}

	// Test file writing
	tmpFile := "test_output.txt"
	defer os.Remove(tmpFile) // clean up

	err = SaveOutput(tmpFile, content)
	if err != nil {
		t.Errorf("Expected no error when writing to file, got: %v", err)
	}

	// Validate content
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read test output file: %v", err)
	}
	if strings.TrimSpace(string(data)) != content {
		t.Errorf("Expected file content '%s', got '%s'", content, string(data))
	}
}

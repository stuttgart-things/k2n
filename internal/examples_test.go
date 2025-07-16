package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCodeExamples(t *testing.T) {
	// Setup temp dir
	dir := t.TempDir()

	// Create mock .tf files
	files := []struct {
		name    string
		content string
	}{
		{"main.tf", "resource \"aws_instance\" \"example\" {}"},
		{"variables.tf", "variable \"region\" {}"},
		{"not_tf.txt", "this should be ignored"},
	}

	for _, f := range files {
		path := filepath.Join(dir, f.name)
		if err := os.WriteFile(path, []byte(f.content), 0644); err != nil {
			t.Fatalf("failed to write file %s: %v", path, err)
		}
	}

	// Run LoadExamples
	examples, err := LoadCodeExamples(dir)
	if err != nil {
		t.Fatalf("LoadExamples returned error: %v", err)
	}

	// Expect only the .tf files to be loaded
	expectedCount := 2
	if len(examples) != expectedCount {
		t.Errorf("expected %d examples, got %d", expectedCount, len(examples))
	}

	// Optional: Validate contents
	for _, ex := range examples {
		if ex == "this should be ignored" {
			t.Errorf("unexpected content loaded: %s", ex)
		}
	}
}

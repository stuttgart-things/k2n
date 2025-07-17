package internal

import (
	"os"
	"path/filepath"
	"reflect"
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

func TestLoadExampleFiles(t *testing.T) {
	// Create temporary test files
	tmpDir := t.TempDir()

	file1 := filepath.Join(tmpDir, "example1.yaml")
	content1 := "kind: Example1"
	err := os.WriteFile(file1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	file2 := filepath.Join(tmpDir, "example2.yaml")
	content2 := "kind: Example2"
	err = os.WriteFile(file2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Test
	paths := []string{file1, file2}
	results, err := LoadExampleFiles(paths)
	if err != nil {
		t.Fatalf("LoadExampleFiles returned error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	if results[0] != content1 {
		t.Errorf("Expected first content '%s', got '%s'", content1, results[0])
	}

	if results[1] != content2 {
		t.Errorf("Expected second content '%s', got '%s'", content2, results[1])
	}
}

func TestDeduplicateStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "EmptyInput",
			input:    []string{},
			expected: nil, // Accept nil slice for empty input
		},
		{
			name:     "AllUnique",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "AllDuplicates",
			input:    []string{"x", "x", "x"},
			expected: []string{"x"},
		},
		{
			name:     "MixedDuplicates",
			input:    []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "CaseSensitivity",
			input:    []string{"A", "a", "A", "a"},
			expected: []string{"A", "a"},
		},
		{
			name:     "OrderPreservation",
			input:    []string{"1", "2", "1", "3", "2"},
			expected: []string{"1", "2", "3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DeduplicateStrings(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DeduplicateStrings(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

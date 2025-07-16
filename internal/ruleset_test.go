package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRulesets(t *testing.T) {
	dir := t.TempDir()

	files := []struct {
		name    string
		content string
	}{
		{"env.yaml", "env: prod"},
		{"usecase.yaml", "usecase: db"},
		{"ignore.txt", "this is ignored"},
	}

	for _, f := range files {
		path := filepath.Join(dir, f.name)
		if err := os.WriteFile(path, []byte(f.content), 0644); err != nil {
			t.Fatalf("failed to write file %s: %v", path, err)
		}
	}

	rulesets, err := LoadRulesets(dir)
	if err != nil {
		t.Fatalf("LoadRulesets returned error: %v", err)
	}

	expectedCount := 3 // all files are loaded, even non-yaml in current implementation
	if len(rulesets) != expectedCount {
		t.Errorf("expected %d rulesets, got %d", expectedCount, len(rulesets))
	}
}

func TestLoadRulesetsIfExists(t *testing.T) {
	// Non-existent directory
	rulesets, err := LoadRulesetsIfExists("non-existent-dir")
	if err != nil {
		t.Fatalf("unexpected error for non-existent dir: %v", err)
	}
	if rulesets != nil {
		t.Errorf("expected nil rulesets for non-existent dir, got: %v", rulesets)
	}

	// Existing directory
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "foo.yaml"), []byte("key: value"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	rulesets, err = LoadRulesetsIfExists(dir)
	if err != nil {
		t.Fatalf("unexpected error loading existing dir: %v", err)
	}
	if len(rulesets) != 1 {
		t.Errorf("expected 1 ruleset, got %d", len(rulesets))
	}
}

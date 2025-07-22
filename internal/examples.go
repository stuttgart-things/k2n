// Package internal provides utilities for loading, filtering, and deduplicating code example files.
package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadCodeExamples loads all files from the given directory and returns them as strings.
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

// LoadExampleFiles loads the content of provided file paths as strings.
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

// DeduplicateStrings removes duplicates from a slice of strings.
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

// FilterFilesByExtension filters file paths by allowed extensions.
func FilterFilesByExtension(files []string, allowedExts []string) []string {
	normalizedExts := normalizeExtensions(allowedExts)

	var filtered []string
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file))
		if _, ok := normalizedExts[ext]; ok {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

// LoadCodeExamplesWithExtensions loads files from a directory that match the allowed extensions.
func LoadCodeExamplesWithExtensions(dir string, allowedExts []string) ([]string, error) {
	var examples []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && hasAllowedExtension(path, allowedExts) {
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

// LoadExampleFilesWithExtensions loads files from provided paths that match the allowed extensions.
func LoadExampleFilesWithExtensions(paths []string, allowedExts []string) ([]string, error) {
	var examples []string
	for _, path := range paths {
		if hasAllowedExtension(path, allowedExts) {
			content, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", path, err)
			}
			examples = append(examples, string(content))
		}
	}
	return examples, nil
}

// hasAllowedExtension checks if a file path has one of the allowed extensions.
func hasAllowedExtension(filename string, allowedExts []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	normalizedExts := normalizeExtensions(allowedExts)
	_, ok := normalizedExts[ext]
	return ok
}

// normalizeExtensions cleans and standardizes extension formats.
func normalizeExtensions(exts []string) map[string]struct{} {
	normalized := make(map[string]struct{})
	for _, ext := range exts {
		ext = strings.ToLower(strings.TrimSpace(ext))
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		normalized[ext] = struct{}{}
	}
	return normalized
}

func SplitAndTrimExts(input string) []string {
	parts := strings.Split(input, ",")
	var exts []string
	for _, p := range parts {
		ext := strings.TrimSpace(p)
		if ext != "" {
			exts = append(exts, ext)
		}
	}
	return exts
}

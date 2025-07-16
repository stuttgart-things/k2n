package internal

import (
	"strings"
	"testing"
)

func TestBuildPrompt(t *testing.T) {
	examples := []string{
		"resource \"aws_instance\" \"example\" {}",
		"resource \"google_compute_instance\" \"example\" {}",
	}
	envRules := []string{
		"Filename: env1.yaml\nCPU: 4, RAM: 8GB",
		"Filename: env2.yaml\nCPU: 8, RAM: 16GB",
	}
	usecaseRules := []string{
		"Filename: usecase1.yaml\nOptimized for database workloads",
	}

	technology := "Terraform"
	instruction := "Generate a config for a high-memory VM."

	prompt := BuildPrompt(examples, envRules, usecaseRules, technology, instruction)

	if !strings.Contains(prompt, "Terraform expert") {
		t.Error("Prompt does not mention technology expert")
	}
	if !strings.Contains(prompt, "Environment Rules:") {
		t.Error("Prompt missing Environment Rules section")
	}
	if !strings.Contains(prompt, "Use Case Rules:") {
		t.Error("Prompt missing Use Case Rules section")
	}
	if !strings.Contains(prompt, examples[0]) || !strings.Contains(prompt, examples[1]) {
		t.Error("Prompt missing examples")
	}
	if !strings.Contains(prompt, instruction) {
		t.Error("Prompt missing instruction")
	}
}

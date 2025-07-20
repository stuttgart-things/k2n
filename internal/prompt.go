package internal

import (
	"fmt"
	"strings"
)

func BuildPrompt(
	examples []string,
	envRules []string,
	usecaseRules []string,
	technology,
	instruction string) string {

	var builder strings.Builder

	tech := technology
	if tech == "" {
		tech = "technology"
	}
	builder.WriteString("You are a " + tech + " expert.\n\n")

	builder.WriteString("You are a " + technology + " expert.\n\n")
	builder.WriteString("General Output Formatting Rules:\n")
	builder.WriteString("- add the marker three dashes.\n")
	builder.WriteString("- add a potential file name (not a file path) e.g. playbook.yaml\n")
	builder.WriteString("- Use '.yaml' as the extension for YAML files.\n")
	builder.WriteString("- Do NOT include syntax highlighting or markdown code fences.\n\n")

	if len(envRules) > 0 {
		builder.WriteString("Environment Rules:\n")
		for _, rule := range envRules {
			builder.WriteString(rule + "\n---\n")
		}
		builder.WriteString("\n")
	}

	if len(usecaseRules) > 0 {
		builder.WriteString("Use Case Rules:\n")
		for _, rule := range usecaseRules {
			builder.WriteString(rule + "\n---\n")
		}
		builder.WriteString("\n")
	}

	builder.WriteString("Examples:\n")
	for i, ex := range examples {
		builder.WriteString(fmt.Sprintf("Example %d:\n%s\n\n", i+1, ex))
	}

	builder.WriteString(fmt.Sprintf("Instruction:\n%s\n", instruction))

	return builder.String()
}

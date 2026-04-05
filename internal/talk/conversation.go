package talk

import (
	"encoding/json"
	"fmt"
	"strings"
)

// AIResponse represents the structured response parsed from the AI output.
type AIResponse struct {
	TemplateName string                 `json:"templateName"`
	Parameters   map[string]interface{} `json:"parameters"`
	Explanation  string                 `json:"explanation"`
}

// BuildSystemPrompt constructs a system prompt that includes all available
// claim templates and their parameter schemas so the AI can match user intent
// to a specific template and fill in the parameters.
func BuildSystemPrompt(templates []ClaimTemplate) string {
	var b strings.Builder

	b.WriteString("You are an infrastructure assistant that helps users provision Crossplane claims.\n")
	b.WriteString("You have access to the following claim templates.\n\n")
	b.WriteString("AVAILABLE TEMPLATES:\n")
	b.WriteString(strings.Repeat("=", 60) + "\n\n")

	for _, t := range templates {
		b.WriteString(fmt.Sprintf("Template: %s\n", t.Metadata.Name))
		if t.Metadata.Title != "" {
			b.WriteString(fmt.Sprintf("  Title: %s\n", t.Metadata.Title))
		}
		if t.Metadata.Description != "" {
			b.WriteString(fmt.Sprintf("  Description: %s\n", t.Metadata.Description))
		}
		b.WriteString(fmt.Sprintf("  Type: %s\n", t.Spec.Type))
		if len(t.Metadata.Tags) > 0 {
			b.WriteString(fmt.Sprintf("  Tags: %s\n", strings.Join(t.Metadata.Tags, ", ")))
		}
		b.WriteString("  Parameters:\n")
		for _, p := range t.Spec.Parameters {
			if p.Hidden {
				continue
			}
			req := ""
			if p.Required {
				req = " (REQUIRED)"
			}
			b.WriteString(fmt.Sprintf("    - %s (%s)%s: %s\n", p.Name, p.Type, req, p.Title))
			if p.Description != "" {
				b.WriteString(fmt.Sprintf("      Description: %s\n", p.Description))
			}
			if p.Default != nil {
				b.WriteString(fmt.Sprintf("      Default: %v\n", p.Default))
			}
			if len(p.Enum) > 0 {
				b.WriteString(fmt.Sprintf("      Allowed values: %s\n", strings.Join(p.Enum, ", ")))
			}
		}
		b.WriteString("\n")
	}

	b.WriteString(strings.Repeat("=", 60) + "\n\n")
	b.WriteString("INSTRUCTIONS:\n")
	b.WriteString("Based on the user's request, select the most appropriate template and fill in the parameters.\n")
	b.WriteString("Respond with ONLY a JSON block in the following format (no markdown fences, no extra text):\n\n")
	b.WriteString(`{
  "templateName": "<name of the selected template>",
  "parameters": {
    "<param1>": "<value1>",
    "<param2>": "<value2>"
  },
  "explanation": "<brief explanation of why this template was chosen and what values were set>"
}`)
	b.WriteString("\n\n")
	b.WriteString("Rules:\n")
	b.WriteString("- Always include all required parameters.\n")
	b.WriteString("- Use default values for optional parameters the user did not mention.\n")
	b.WriteString("- If the user's request does not match any template, set templateName to \"\" and explain why in the explanation field.\n")
	b.WriteString("- For enum parameters, only use allowed values.\n")
	b.WriteString("- For array parameters, provide a JSON array.\n")

	return b.String()
}

// BuildUserPrompt wraps the user's natural language instruction into the conversation.
func BuildUserPrompt(systemPrompt, userInstruction string) string {
	return systemPrompt + "\n\nUser request:\n" + userInstruction + "\n"
}

// ParseAIResponse extracts the structured AIResponse from the AI's text output.
func ParseAIResponse(aiOutput string) (*AIResponse, error) {
	// Try to find JSON in the output
	trimmed := strings.TrimSpace(aiOutput)

	// Strip markdown code fences if present
	if strings.HasPrefix(trimmed, "```") {
		lines := strings.Split(trimmed, "\n")
		// Remove first and last lines (fences)
		if len(lines) >= 3 {
			trimmed = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}

	// Try direct parse first
	var resp AIResponse
	if err := json.Unmarshal([]byte(trimmed), &resp); err == nil {
		return &resp, nil
	}

	// Try to extract JSON object from surrounding text
	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start >= 0 && end > start {
		jsonStr := trimmed[start : end+1]
		if err := json.Unmarshal([]byte(jsonStr), &resp); err == nil {
			return &resp, nil
		}
	}

	return nil, fmt.Errorf("could not parse AI response as JSON: %s", truncateForError(trimmed, 200))
}

func truncateForError(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

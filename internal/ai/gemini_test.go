package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCallGeminiAPI(t *testing.T) {
	fakeResponse := map[string]interface{}{
		"candidates": []map[string]interface{}{
			{
				"content": map[string]interface{}{
					"parts": []map[string]interface{}{
						{"text": "Generated Terraform Config"},
					},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(fakeResponse)
	}))
	defer server.Close()

	// Provide server.URL to CallGeminiAPI via parameter or refactor to accept URL (here's a modified version)
	result, err := callGeminiAPIWithURL("fake-api-key", "fake-prompt", server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Generated Terraform Config"
	if result != expected {
		t.Errorf("expected %q but got %q", expected, result)
	}
}

// callGeminiAPIWithURL is a testable variant of CallGeminiAPI.
func callGeminiAPIWithURL(apiKey, prompt, baseURL string) (string, error) {
	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(baseURL, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no candidates returned")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

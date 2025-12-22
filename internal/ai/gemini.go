package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const (
	GeminiURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-3-pro-preview:generateContent"
)

func CallGeminiAPI(apiKey, prompt string) (string, error) {
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

	resp, err := http.Post(
		GeminiURL+"?key="+apiKey,
		"application/json",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Raw response:", string(respBody))

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no candidates returned")
	}

	text := geminiResp.Candidates[0].Content.Parts[0].Text
	return cleanCodeBlock(text), nil
}

// cleanCodeBlock removes surrounding triple backticks and optional language hints.
func cleanCodeBlock(text string) string {
	// Match patterns like ```yaml\ncontent\n```
	re := regexp.MustCompile("^```[a-zA-Z]*\\n([\\s\\S]*?)\\n```$")
	matches := re.FindStringSubmatch(strings.TrimSpace(text))
	if len(matches) > 1 {
		return matches[1]
	}
	return text
}

package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CallOpenRouterWithURL(apiKey, prompt, baseURL, model string) (string, error) {
	reqBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var orResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(respBody, &orResp); err != nil {
		return "", err
	}

	if orResp.Error != nil {
		return "", fmt.Errorf("openrouter error: %s", orResp.Error.Message)
	}

	if len(orResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	return cleanCodeBlock(orResp.Choices[0].Message.Content), nil
}

package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func CallOpenRouterApi(apiKey, prompt, baseURL, model string) (string, error) {
	log.Printf("[OpenRouter] Starting API call with model: %s", model)
	log.Printf("[OpenRouter] Base URL: %s", baseURL)
	log.Printf("[OpenRouter] Prompt length: %d characters", len(prompt))

	reqBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("[OpenRouter] ERROR marshaling request body: %v", err)
		return "", err
	}
	log.Printf("[OpenRouter] Request body marshaled successfully (%d bytes)", len(bodyBytes))

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(bodyBytes))
	if err != nil {
		log.Printf("[OpenRouter] ERROR creating HTTP request: %v", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	log.Printf("[OpenRouter] HTTP request created and headers set")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[OpenRouter] ERROR executing HTTP request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	log.Printf("[OpenRouter] HTTP response received with status: %d %s", resp.StatusCode, resp.Status)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[OpenRouter] ERROR reading response body: %v", err)
		return "", err
	}
	log.Printf("[OpenRouter] Response body read successfully (%d bytes)", len(respBody))
	log.Printf("[OpenRouter] Response body content: %s", string(respBody))

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
		log.Printf("[OpenRouter] ERROR unmarshaling response: %v", err)
		return "", err
	}
	log.Printf("[OpenRouter] Response unmarshaled successfully")

	if orResp.Error != nil {
		log.Printf("[OpenRouter] API returned error: %s", orResp.Error.Message)
		return "", fmt.Errorf("openrouter error: %s", orResp.Error.Message)
	}

	if len(orResp.Choices) == 0 {
		log.Printf("[OpenRouter] ERROR: no choices returned in response")
		return "", fmt.Errorf("no choices returned")
	}

	result := cleanCodeBlock(orResp.Choices[0].Message.Content)
	log.Printf("[OpenRouter] Response processed successfully. Result length: %d characters", len(result))
	return result, nil
}

package ai

import (
	"fmt"
	"os"
	"strings"
)

// AIProvider defines the interface for AI providers
type AIProvider interface {
	Call(apiKey, prompt string) (string, error)
}

// ProviderType represents the type of AI provider
type ProviderType string

const (
	ProviderOpenRouter ProviderType = "openrouter"
	ProviderGemini     ProviderType = "gemini"
)

// ProviderConfig holds configuration for the AI provider
type ProviderConfig struct {
	Type    ProviderType
	APIKey  string
	Model   string
	BaseURL string
}

// GetProviderFromEnv creates a provider configuration from environment variables
// Environment variables:
//   - AI_PROVIDER: "openrouter" or "gemini" (default: "openrouter")
//   - AI_API_KEY: API key for the provider
//   - AI_MODEL: Model name (for OpenRouter, e.g., "openai/gpt-4")
//   - AI_BASE_URL: Base URL for OpenRouter API (optional)
func GetProviderFromEnv() (*ProviderConfig, error) {
	provider := strings.ToLower(os.Getenv("AI_PROVIDER"))
	if provider == "" {
		provider = "openrouter" // default
	}

	apiKey := os.Getenv("AI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("AI_API_KEY environment variable is required")
	}

	config := &ProviderConfig{
		APIKey: apiKey,
	}

	switch provider {
	case "openrouter":
		config.Type = ProviderOpenRouter
		config.Model = os.Getenv("AI_MODEL")
		if config.Model == "" {
			config.Model = "openai/gpt-3.5-turbo" // default model
		}
		config.BaseURL = os.Getenv("AI_BASE_URL")
		if config.BaseURL == "" {
			config.BaseURL = "https://openrouter.ai/api/v1/chat/completions"
		}
	case "gemini":
		config.Type = ProviderGemini
	default:
		return nil, fmt.Errorf("unknown AI_PROVIDER: %s (supported: openrouter, gemini)", provider)
	}

	return config, nil
}

// NewProvider creates a new AI provider instance based on the configuration
func NewProvider(config *ProviderConfig) (AIProvider, error) {
	switch config.Type {
	case ProviderOpenRouter:
		return &OpenRouterProvider{
			APIKey:  config.APIKey,
			Model:   config.Model,
			BaseURL: config.BaseURL,
		}, nil
	case ProviderGemini:
		return &GeminiProvider{
			APIKey: config.APIKey,
		}, nil
	default:
		return nil, fmt.Errorf("unknown provider type: %v", config.Type)
	}
}

// CallAI calls the configured AI provider with the given prompt
func CallAI(config *ProviderConfig, prompt string) (string, error) {
	provider, err := NewProvider(config)
	if err != nil {
		return "", err
	}
	return provider.Call(config.APIKey, prompt)
}

// CallAIWithProvider calls the configured AI provider using environment variables
func CallAIWithProvider(prompt string) (string, error) {
	config, err := GetProviderFromEnv()
	if err != nil {
		return "", err
	}
	return CallAI(config, prompt)
}

// OpenRouterProvider implements AIProvider for OpenRouter
type OpenRouterProvider struct {
	APIKey  string
	Model   string
	BaseURL string
}

// Call implements AIProvider.Call for OpenRouter
func (p *OpenRouterProvider) Call(apiKey, prompt string) (string, error) {
	return CallOpenRouterWithURL(apiKey, prompt, p.BaseURL, p.Model)
}

// GeminiProvider implements AIProvider for Gemini
type GeminiProvider struct {
	APIKey string
}

// Call implements AIProvider.Call for Gemini
func (p *GeminiProvider) Call(apiKey, prompt string) (string, error) {
	return CallGeminiAPI(apiKey, prompt)
}

// Package talk provides a client for the claim-machinery-api and AI-powered
// conversation logic for rendering Crossplane claims.
package talk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ClaimTemplate represents a claim template from the claim-machinery-api.
type ClaimTemplate struct {
	APIVersion string                `json:"apiVersion"`
	Kind       string                `json:"kind"`
	Metadata   ClaimTemplateMetadata `json:"metadata"`
	Spec       ClaimTemplateSpec     `json:"spec"`
}

// ClaimTemplateMetadata holds template metadata.
type ClaimTemplateMetadata struct {
	Name        string   `json:"name"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Profile     string   `json:"profile,omitempty"`
}

// ClaimTemplateSpec holds the template specification.
type ClaimTemplateSpec struct {
	Type       string      `json:"type"`
	Source     string      `json:"source"`
	Tag        string      `json:"tag,omitempty"`
	Parameters []Parameter `json:"parameters"`
}

// Parameter describes a single template parameter.
type Parameter struct {
	Name        string      `json:"name"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Type        string      `json:"type"`
	Default     interface{} `json:"default,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
	Hidden      bool        `json:"hidden,omitempty"`
	Pattern     string      `json:"pattern,omitempty"`
	MinLength   *int        `json:"minLength,omitempty"`
	MaxLength   *int        `json:"maxLength,omitempty"`
}

// ClaimTemplateListResponse is the response from GET /api/v1/claim-templates.
type ClaimTemplateListResponse struct {
	APIVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	Items      []ClaimTemplate `json:"items"`
}

// OrderRequest is the request body for POST /api/v1/claim-templates/{name}/order.
type OrderRequest struct {
	Parameters map[string]interface{} `json:"parameters"`
	Author     string                 `json:"author,omitempty"`
}

// OrderResponse is the response from POST /api/v1/claim-templates/{name}/order.
type OrderResponse struct {
	APIVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   map[string]interface{} `json:"metadata"`
	Rendered   string                 `json:"rendered"`
}

// Client communicates with the claim-machinery-api.
type Client struct {
	BaseURL    string
	AuthToken  string
	HTTPClient *http.Client
}

// NewClient creates a new claim-machinery-api client.
func NewClient(baseURL, authToken string) *Client {
	return &Client{
		BaseURL:   baseURL,
		AuthToken: authToken,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) doRequest(method, path string, body io.Reader) ([]byte, error) {
	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(data))
	}

	return data, nil
}

// ListTemplates fetches all available claim templates.
func (c *Client) ListTemplates() ([]ClaimTemplate, error) {
	data, err := c.doRequest(http.MethodGet, "/api/v1/claim-templates", nil)
	if err != nil {
		return nil, err
	}

	var listResp ClaimTemplateListResponse
	if err := json.Unmarshal(data, &listResp); err != nil {
		return nil, fmt.Errorf("decoding template list: %w", err)
	}

	return listResp.Items, nil
}

// GetTemplate fetches a specific claim template by name.
func (c *Client) GetTemplate(name string) (*ClaimTemplate, error) {
	data, err := c.doRequest(http.MethodGet, "/api/v1/claim-templates/"+name, nil)
	if err != nil {
		return nil, err
	}

	var tmpl ClaimTemplate
	if err := json.Unmarshal(data, &tmpl); err != nil {
		return nil, fmt.Errorf("decoding template: %w", err)
	}

	return &tmpl, nil
}

// OrderClaim renders a claim by posting parameters to the order endpoint.
func (c *Client) OrderClaim(templateName string, params map[string]interface{}, author string) (*OrderResponse, error) {
	reqBody := OrderRequest{
		Parameters: params,
		Author:     author,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("encoding order request: %w", err)
	}

	data, err := c.doRequest(http.MethodPost, "/api/v1/claim-templates/"+templateName+"/order", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	var orderResp OrderResponse
	if err := json.Unmarshal(data, &orderResp); err != nil {
		return nil, fmt.Errorf("decoding order response: %w", err)
	}

	return &orderResp, nil
}

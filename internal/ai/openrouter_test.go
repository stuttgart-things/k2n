package ai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCallOpenRouterWithURL(t *testing.T) {
	fakeResponse := map[string]interface{}{
		"choices": []map[string]interface{}{
			{
				"message": map[string]interface{}{
					"content": "```yaml\nGenerated Terraform Config\n```",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure the Authorization header is passed
		if r.Header.Get("Authorization") != "Bearer fake-api-key" {
			t.Errorf("expected Authorization header but got %q", r.Header.Get("Authorization"))
		}
		_ = json.NewEncoder(w).Encode(fakeResponse)
	}))
	defer server.Close()

	result, err := CallOpenRouterWithURL("fake-api-key", "fake-prompt", server.URL, "deepseek/deepseek-r1-0528:free")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Generated Terraform Config"
	if result != expected {
		t.Errorf("expected %q but got %q", expected, result)
	}
}

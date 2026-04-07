package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultClaudeEndpoint = "https://api.anthropic.com/v1/messages"
	claudeModel           = "claude-sonnet-4-20250514"
	classificationSystem  = `You are a financial classification assistant. Given a transaction description and a list of available categories, return the best matching category and relevant tags. Return ONLY a JSON object: {"category": "name", "tags": ["tag1"]}`
)

// ClaudeProvider implements Provider using the Anthropic Messages API.
type ClaudeProvider struct {
	apiKey   string
	endpoint string
	client   *http.Client
}

// NewClaudeProvider creates a new ClaudeProvider. If endpoint is empty, the default
// Anthropic API endpoint is used.
func NewClaudeProvider(apiKey, endpoint string) *ClaudeProvider {
	if endpoint == "" {
		endpoint = defaultClaudeEndpoint
	}
	return &ClaudeProvider{
		apiKey:   apiKey,
		endpoint: endpoint,
		client:   &http.Client{Timeout: 5 * time.Second},
	}
}

// Classify sends a transaction description to Claude and returns the predicted
// category and tags.
func (p *ClaudeProvider) Classify(ctx context.Context, description string, categories []string) (*Classification, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("claude API key not configured")
	}

	catJSON, _ := json.Marshal(categories)
	userMsg := fmt.Sprintf("Transaction description: %q\nAvailable categories: %s", description, string(catJSON))

	body := map[string]any{
		"model":      claudeModel,
		"max_tokens": 256,
		"system":     classificationSystem,
		"messages": []map[string]string{
			{"role": "user", "content": userMsg},
		},
	}

	bodyBytes, _ := json.Marshal(body)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("claude API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("claude API error %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse claude response: %w", err)
	}
	if len(result.Content) == 0 {
		return nil, fmt.Errorf("empty claude response")
	}

	return extractClassification(result.Content[0].Text)
}

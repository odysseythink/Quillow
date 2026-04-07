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
	defaultOpenAIEndpoint = "https://api.openai.com/v1/chat/completions"
	openAIModel           = "gpt-4o-mini"
)

// OpenAIProvider implements Provider using the OpenAI Chat Completions API.
type OpenAIProvider struct {
	apiKey   string
	endpoint string
	client   *http.Client
}

// NewOpenAIProvider creates a new OpenAIProvider. If endpoint is empty, the default
// OpenAI API endpoint is used.
func NewOpenAIProvider(apiKey, endpoint string) *OpenAIProvider {
	if endpoint == "" {
		endpoint = defaultOpenAIEndpoint
	}
	return &OpenAIProvider{
		apiKey:   apiKey,
		endpoint: endpoint,
		client:   &http.Client{Timeout: 5 * time.Second},
	}
}

// Classify sends a transaction description to OpenAI and returns the predicted
// category and tags.
func (p *OpenAIProvider) Classify(ctx context.Context, description string, categories []string) (*Classification, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("openai API key not configured")
	}

	catJSON, _ := json.Marshal(categories)
	userMsg := fmt.Sprintf("Transaction description: %q\nAvailable categories: %s", description, string(catJSON))

	body := map[string]any{
		"model":      openAIModel,
		"max_tokens": 256,
		"messages": []map[string]string{
			{"role": "system", "content": classificationSystem},
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
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openai API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openai API error %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse openai response: %w", err)
	}
	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("empty openai response")
	}

	return extractClassification(result.Choices[0].Message.Content)
}

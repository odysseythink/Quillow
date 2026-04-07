package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

// extractClassification parses a Classification from raw AI response text,
// handling potential markdown code block wrappers.
func extractClassification(text string) (*Classification, error) {
	text = strings.TrimSpace(text)

	// Extract JSON from potential markdown code blocks
	if idx := strings.Index(text, "{"); idx >= 0 {
		if end := strings.LastIndex(text, "}"); end >= idx {
			text = text[idx : end+1]
		}
	}

	var classification Classification
	if err := json.Unmarshal([]byte(text), &classification); err != nil {
		return nil, fmt.Errorf("failed to parse classification JSON: %w", err)
	}

	return &classification, nil
}

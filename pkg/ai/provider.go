package ai

import "context"

// Classification represents the result of an AI-powered transaction classification.
type Classification struct {
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

// Provider defines the interface for AI-powered transaction classification.
type Provider interface {
	Classify(ctx context.Context, description string, categories []string) (*Classification, error)
}

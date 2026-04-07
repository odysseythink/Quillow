package ai

import (
	"context"
	"log"
)

// SuggestResult holds the result of an AI classification suggestion.
type SuggestResult struct {
	CategoryID   uint     `json:"category_id"`
	CategoryName string   `json:"category_name"`
	Tags         []string `json:"tags"`
	Source       string   `json:"source"` // "local", "llm", "none"
}

// CategoryInfo is a minimal category representation for matching.
type CategoryInfo struct {
	ID   uint
	Name string
}

// Service orchestrates the 3-level fallback classification.
type Service struct {
	matcher  *LocalMatcher
	provider Provider // LLM provider (may be nil)
}

// NewService creates a new AI classification service.
func NewService(provider Provider) *Service {
	return &Service{
		matcher:  NewLocalMatcher(),
		provider: provider,
	}
}

// Matcher returns the local matcher for external cache loading.
func (s *Service) Matcher() *LocalMatcher {
	return s.matcher
}

// Provider returns the LLM provider (may be nil).
func (s *Service) Provider() Provider {
	return s.provider
}

// Suggest classifies a transaction description using the 3-level fallback:
//  1. Local pattern matching
//  2. LLM API
//  3. No suggestion
func (s *Service) Suggest(ctx context.Context, userID uint, description string, categories []CategoryInfo) *SuggestResult {
	// Level 1: Local pattern match
	if match := s.matcher.Match(userID, description); match != nil {
		name := ""
		for _, c := range categories {
			if c.ID == match.CategoryID {
				name = c.Name
				break
			}
		}
		return &SuggestResult{
			CategoryID:   match.CategoryID,
			CategoryName: name,
			Tags:         nil,
			Source:       "local",
		}
	}

	// Level 2: LLM API
	if s.provider != nil {
		catNames := make([]string, len(categories))
		for i, c := range categories {
			catNames[i] = c.Name
		}

		result, err := s.provider.Classify(ctx, description, catNames)
		if err != nil {
			log.Printf("AI provider classify failed: %v", err)
		} else if result != nil && result.Category != "" {
			// Find category ID by name
			var catID uint
			for _, c := range categories {
				if c.Name == result.Category {
					catID = c.ID
					break
				}
			}
			return &SuggestResult{
				CategoryID:   catID,
				CategoryName: result.Category,
				Tags:         result.Tags,
				Source:       "llm",
			}
		}
	}

	// Level 3: No suggestion
	return &SuggestResult{Source: "none"}
}

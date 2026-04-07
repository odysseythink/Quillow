package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anthropics/quillow/internal/adapter/repository"
	repomodel "github.com/anthropics/quillow/internal/adapter/repository/model"
	"github.com/anthropics/quillow/internal/port"
	"github.com/anthropics/quillow/pkg/ai"
)

// ChatResponse is the unified response for the /ai/chat endpoint.
type ChatResponse struct {
	Intent        string                `json:"intent"` // "record" or "query"
	Parsed        *ai.ParsedTransaction `json:"parsed,omitempty"`
	Confidence    string                `json:"confidence,omitempty"`
	Created       bool                  `json:"created,omitempty"`
	TransactionID uint                  `json:"transaction_id,omitempty"`
	Answer        string                `json:"answer,omitempty"`
	Data          any                   `json:"data,omitempty"`
}

type UseCase struct {
	aiSvc         *ai.Service
	patternRepo   *repository.ClassificationPatternRepository
	categoryRepo  port.CategoryRepository
	queryRegistry *ai.QueryRegistry
}

func NewUseCase(
	aiSvc *ai.Service,
	patternRepo *repository.ClassificationPatternRepository,
	categoryRepo port.CategoryRepository,
	queryRegistry *ai.QueryRegistry,
) *UseCase {
	return &UseCase{
		aiSvc:         aiSvc,
		patternRepo:   patternRepo,
		categoryRepo:  categoryRepo,
		queryRegistry: queryRegistry,
	}
}

// LoadUserPatterns loads patterns from DB into the in-memory cache.
func (uc *UseCase) LoadUserPatterns(ctx context.Context, userID uint) error {
	patterns, err := uc.patternRepo.FindByUser(ctx, userID)
	if err != nil {
		return err
	}
	entries := make([]ai.PatternEntry, len(patterns))
	for i, p := range patterns {
		entries[i] = ai.PatternEntry{
			ID:         p.ID,
			Pattern:    p.Pattern,
			CategoryID: p.CategoryID,
			TagIDs:     p.TagIDs,
			HitCount:   p.HitCount,
		}
	}
	uc.aiSvc.Matcher().LoadPatterns(userID, entries)
	return nil
}

// Suggest returns a classification suggestion for a description.
func (uc *UseCase) Suggest(ctx context.Context, userID uint, description string) (*ai.SuggestResult, error) {
	// Ensure patterns are loaded
	_ = uc.LoadUserPatterns(ctx, userID)

	// Get user's categories
	categories, _, err := uc.categoryRepo.List(ctx, 0, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	catInfos := make([]ai.CategoryInfo, len(categories))
	for i, c := range categories {
		catInfos[i] = ai.CategoryInfo{ID: c.ID, Name: c.Name}
	}

	result := uc.aiSvc.Suggest(ctx, userID, description, catInfos)
	return result, nil
}

// LearnPattern records a user's classification choice for future matching.
func (uc *UseCase) LearnPattern(ctx context.Context, userID uint, pattern string, categoryID uint, tagIDs []uint) error {
	existing, err := uc.patternRepo.FindByUserAndPattern(ctx, userID, pattern, categoryID)
	if err == nil && existing != nil {
		// Increment hit count
		if err := uc.patternRepo.IncrementHitCount(ctx, existing.ID); err != nil {
			return err
		}
		existing.HitCount++
		uc.aiSvc.Matcher().AddPattern(userID, ai.PatternEntry{
			ID: existing.ID, Pattern: existing.Pattern,
			CategoryID: existing.CategoryID, TagIDs: existing.TagIDs,
			HitCount: existing.HitCount,
		})
		return nil
	}

	tagJSON, _ := json.Marshal(tagIDs)
	m := &repomodel.ClassificationPatternModel{
		UserID:     userID,
		Pattern:    pattern,
		CategoryID: categoryID,
		TagIDs:     string(tagJSON),
		HitCount:   1,
	}
	if err := uc.patternRepo.Create(ctx, m); err != nil {
		return err
	}
	uc.aiSvc.Matcher().AddPattern(userID, ai.PatternEntry{
		ID: m.ID, Pattern: m.Pattern,
		CategoryID: m.CategoryID, TagIDs: m.TagIDs,
		HitCount: 1,
	})
	return nil
}

// Chat is the unified entry point for the chat bubble.
// It detects intent and dispatches to record or query flow.
func (uc *UseCase) Chat(ctx context.Context, userID uint, message string) (*ChatResponse, error) {
	intent := ai.DetectIntent(message)

	switch intent {
	case "query":
		answer, data, err := uc.Insight(ctx, userID, message)
		if err != nil {
			return nil, err
		}
		return &ChatResponse{Intent: "query", Answer: answer, Data: data}, nil
	case "record":
		parsed, confidence := ai.ParseLocal(message, time.Now())
		// Try to get category suggestion
		if parsed.Description != "" {
			if suggestion, err := uc.Suggest(ctx, userID, parsed.Description); err == nil && suggestion.CategoryID != 0 {
				parsed.CategoryID = suggestion.CategoryID
				parsed.Category = suggestion.CategoryName
			}
		}
		resp := &ChatResponse{
			Intent:     "record",
			Parsed:     parsed,
			Confidence: confidence,
		}
		// High confidence: auto-create would happen here when transaction creation is wired
		// For now, return parsed data for frontend to confirm
		return resp, nil
	default:
		// Unknown intent — try record first, fall back to query
		parsed, confidence := ai.ParseLocal(message, time.Now())
		if parsed.Amount != "" {
			if parsed.Description != "" {
				if suggestion, err := uc.Suggest(ctx, userID, parsed.Description); err == nil && suggestion.CategoryID != 0 {
					parsed.CategoryID = suggestion.CategoryID
					parsed.Category = suggestion.CategoryName
				}
			}
			return &ChatResponse{Intent: "record", Parsed: parsed, Confidence: confidence}, nil
		}
		// Fall back to query
		answer, data, err := uc.Insight(ctx, userID, message)
		if err != nil {
			return &ChatResponse{Intent: "query", Answer: "暂不支持该查询"}, nil
		}
		return &ChatResponse{Intent: "query", Answer: answer, Data: data}, nil
	}
}

// Insight handles query-type messages by executing predefined query functions.
func (uc *UseCase) Insight(ctx context.Context, userID uint, message string) (string, any, error) {
	if uc.queryRegistry == nil {
		return "查询功能未启用", nil, nil
	}

	// Use LLM to determine which function to call
	if uc.aiSvc != nil && uc.aiSvc.Provider() != nil {
		categories, _, _ := uc.categoryRepo.List(ctx, 0, 1000, 0)
		catNames := make([]string, len(categories))
		for i, c := range categories {
			catNames[i] = c.Name
		}

		prompt := fmt.Sprintf(
			"用户问题: \"%s\"\n今天日期: %s\n可用查询函数:\n%s\n可用分类: %v\n\n返回 JSON: {\"function\": \"函数名\", \"params\": {\"key\": \"value\"}}\n仅返回 JSON。",
			message, time.Now().Format("2006-01-02"), uc.queryRegistry.Describe(), catNames,
		)

		result, err := uc.aiSvc.Provider().Classify(ctx, prompt, nil)
		if err == nil && result != nil {
			// Parse the function call from LLM response
			var call struct {
				Function string            `json:"function"`
				Params   map[string]string `json:"params"`
			}
			if err := json.Unmarshal([]byte(result.Category), &call); err == nil && call.Function != "" {
				data, err := uc.queryRegistry.Execute(ctx, userID, call.Function, call.Params)
				if err == nil {
					// Second LLM call to generate natural language answer
					dataJSON, _ := json.Marshal(data)
					answerPrompt := fmt.Sprintf("用户问题: \"%s\"\n查询结果: %s\n\n用简洁友好的中文回答。", message, string(dataJSON))
					answerResult, err := uc.aiSvc.Provider().Classify(ctx, answerPrompt, nil)
					if err == nil && answerResult != nil {
						return answerResult.Category, data, nil
					}
					// If second call fails, return raw data
					return fmt.Sprintf("查询结果: %s", string(dataJSON)), data, nil
				}
			}
		}
	}

	return "暂不支持该查询，请确认 AI 服务已配置", nil, nil
}

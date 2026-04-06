package i18n

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type contextKey string

const LocaleKey contextKey = "locale"
const DefaultLocale = "en_US"

type Service struct {
	basePath string
	locales  map[string]map[string]string
	mu       sync.RWMutex
}

func NewService(basePath string) (*Service, error) {
	return &Service{
		basePath: basePath,
		locales:  make(map[string]map[string]string),
	}, nil
}

func (s *Service) LoadLocale(locale string) error {
	path := filepath.Join(s.basePath, locale, "messages.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to load locale %s: %w", locale, err)
	}

	messages := make(map[string]string)
	if err := json.Unmarshal(data, &messages); err != nil {
		return fmt.Errorf("failed to parse locale %s: %w", locale, err)
	}

	s.mu.Lock()
	s.locales[locale] = messages
	s.mu.Unlock()

	return nil
}

func (s *Service) T(ctx context.Context, key string, params ...any) string {
	locale := DefaultLocale
	if v := ctx.Value(LocaleKey); v != nil {
		if l, ok := v.(string); ok {
			locale = l
		}
	}

	s.mu.RLock()
	messages, ok := s.locales[locale]
	if !ok {
		messages = s.locales[DefaultLocale]
	}
	s.mu.RUnlock()

	if messages == nil {
		return key
	}

	msg, ok := messages[key]
	if !ok {
		return key
	}

	if len(params) > 0 {
		return fmt.Sprintf(msg, params...)
	}
	return msg
}

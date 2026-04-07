package ai

import (
	"strings"
	"sync"
)

// PatternEntry represents a learned classification pattern.
type PatternEntry struct {
	ID         uint
	Pattern    string
	CategoryID uint
	TagIDs     string
	HitCount   uint
}

// LocalMatcher provides in-memory pattern matching for transaction classification.
type LocalMatcher struct {
	mu       sync.RWMutex
	patterns map[uint][]PatternEntry // userID -> patterns (sorted by hit_count desc)
}

func NewLocalMatcher() *LocalMatcher {
	return &LocalMatcher{patterns: make(map[uint][]PatternEntry)}
}

// LoadPatterns replaces all patterns for a user.
func (m *LocalMatcher) LoadPatterns(userID uint, patterns []PatternEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.patterns[userID] = patterns
}

// Match finds the best matching pattern for a description.
// Returns nil if no pattern matches.
func (m *LocalMatcher) Match(userID uint, description string) *PatternEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	desc := strings.ToLower(description)
	for i := range m.patterns[userID] {
		if strings.Contains(desc, strings.ToLower(m.patterns[userID][i].Pattern)) {
			return &m.patterns[userID][i]
		}
	}
	return nil
}

// AddPattern adds or updates a pattern in the cache.
func (m *LocalMatcher) AddPattern(userID uint, entry PatternEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()

	patterns := m.patterns[userID]
	for i := range patterns {
		if patterns[i].ID == entry.ID {
			patterns[i] = entry
			m.patterns[userID] = patterns
			return
		}
	}
	m.patterns[userID] = append(patterns, entry)
}

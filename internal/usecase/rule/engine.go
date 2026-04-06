package rule

import (
	"context"
	"strings"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/internal/port"
)

// TriggerFunc checks if a trigger condition matches.
type TriggerFunc func(description, value string) bool

// ActionFunc executes an action (placeholder for SP5).
type ActionFunc func(journalID uint, value string) error

var triggerRegistry = map[string]TriggerFunc{
	"description_is": func(desc, val string) bool {
		return strings.EqualFold(desc, val)
	},
	"description_contains": func(desc, val string) bool {
		return strings.Contains(strings.ToLower(desc), strings.ToLower(val))
	},
	"description_starts": func(desc, val string) bool {
		return strings.HasPrefix(strings.ToLower(desc), strings.ToLower(val))
	},
	"description_ends": func(desc, val string) bool {
		return strings.HasSuffix(strings.ToLower(desc), strings.ToLower(val))
	},
}

// Engine evaluates rules against transaction data.
type Engine struct {
	ruleGroupRepo port.RuleGroupRepository
	ruleRepo      port.RuleRepository
}

// NewEngine creates a new rule evaluation engine.
func NewEngine(rgRepo port.RuleGroupRepository, rRepo port.RuleRepository) *Engine {
	return &Engine{ruleGroupRepo: rgRepo, ruleRepo: rRepo}
}

// EvaluateRules checks all active rules against a transaction description.
// It returns the list of rules whose triggers all match the given description.
func (e *Engine) EvaluateRules(ctx context.Context, userGroupID uint, description string) ([]entity.Rule, error) {
	// List all rule groups ordered by order
	groups, _, err := e.ruleGroupRepo.List(ctx, userGroupID, 1000, 0)
	if err != nil {
		return nil, err
	}

	var matched []entity.Rule

	for _, group := range groups {
		if !group.Active {
			continue
		}

		// For each group, get rules
		rules, err := e.ruleGroupRepo.ListRules(ctx, group.ID)
		if err != nil {
			return nil, err
		}

		for _, rule := range rules {
			if !rule.Active {
				continue
			}

			// For each rule, get triggers
			triggers, err := e.ruleRepo.GetTriggers(ctx, rule.ID)
			if err != nil {
				return nil, err
			}

			if len(triggers) == 0 {
				continue
			}

			// Check if triggers match
			if e.matchesTriggers(rule.Strict, triggers, description) {
				matched = append(matched, rule)
			}

			if rule.StopProcessing {
				break
			}
		}

		if group.StopProcessing {
			break
		}
	}

	return matched, nil
}

// matchesTriggers checks whether triggers match the description.
// If strict is true, ALL triggers must match. Otherwise, ANY trigger matching is sufficient.
func (e *Engine) matchesTriggers(strict bool, triggers []entity.RuleTrigger, description string) bool {
	if strict {
		for _, t := range triggers {
			if !t.Active {
				continue
			}
			fn, ok := triggerRegistry[t.TriggerType]
			if !ok {
				return false
			}
			if !fn(description, t.TriggerValue) {
				return false
			}
		}
		return true
	}

	// Non-strict: any trigger match is sufficient
	for _, t := range triggers {
		if !t.Active {
			continue
		}
		fn, ok := triggerRegistry[t.TriggerType]
		if !ok {
			continue
		}
		if fn(description, t.TriggerValue) {
			return true
		}
	}
	return false
}

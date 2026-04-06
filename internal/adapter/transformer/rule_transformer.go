package transformer

import (
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/pkg/response"
)

func TransformRuleGroup(rg *entity.RuleGroup) response.Resource {
	return response.Resource{
		Type: "rule_groups",
		ID:   fmt.Sprintf("%d", rg.ID),
		Attributes: map[string]any{
			"created_at":      rg.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at":      rg.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"title":           rg.Title,
			"description":     rg.Description,
			"order":           rg.Order,
			"active":          rg.Active,
			"stop_processing": rg.StopProcessing,
		},
	}
}

func TransformRule(r *entity.Rule, triggers []entity.RuleTrigger, actions []entity.RuleAction) response.Resource {
	triggerData := make([]map[string]any, len(triggers))
	for i, t := range triggers {
		triggerData[i] = map[string]any{
			"id":              fmt.Sprintf("%d", t.ID),
			"type":            t.TriggerType,
			"value":           t.TriggerValue,
			"order":           t.Order,
			"active":          t.Active,
			"stop_processing": t.StopProcessing,
		}
	}

	actionData := make([]map[string]any, len(actions))
	for i, a := range actions {
		actionData[i] = map[string]any{
			"id":              fmt.Sprintf("%d", a.ID),
			"type":            a.ActionType,
			"value":           a.ActionValue,
			"order":           a.Order,
			"active":          a.Active,
			"stop_processing": a.StopProcessing,
		}
	}

	return response.Resource{
		Type: "rules",
		ID:   fmt.Sprintf("%d", r.ID),
		Attributes: map[string]any{
			"created_at":      r.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at":      r.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"title":           r.Title,
			"description":     r.Description,
			"rule_group_id":   fmt.Sprintf("%d", r.RuleGroupID),
			"order":           r.Order,
			"active":          r.Active,
			"strict":          r.Strict,
			"stop_processing": r.StopProcessing,
			"triggers":        triggerData,
			"actions":         actionData,
		},
	}
}

package port

import (
	"context"

	"github.com/anthropics/quillow/internal/entity"
)

type RuleGroupRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.RuleGroup, error)
	List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.RuleGroup, int64, error)
	Create(ctx context.Context, rg *entity.RuleGroup) error
	Update(ctx context.Context, rg *entity.RuleGroup) error
	Delete(ctx context.Context, id uint) error
	ListRules(ctx context.Context, ruleGroupID uint) ([]entity.Rule, error)
}

type RuleRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.Rule, error)
	List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Rule, int64, error)
	Create(ctx context.Context, rule *entity.Rule) error
	Update(ctx context.Context, rule *entity.Rule) error
	Delete(ctx context.Context, id uint) error
	GetTriggers(ctx context.Context, ruleID uint) ([]entity.RuleTrigger, error)
	SetTriggers(ctx context.Context, ruleID uint, triggers []entity.RuleTrigger) error
	GetActions(ctx context.Context, ruleID uint) ([]entity.RuleAction, error)
	SetActions(ctx context.Context, ruleID uint, actions []entity.RuleAction) error
}

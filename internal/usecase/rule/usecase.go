package rule

import (
	"context"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/internal/port"
)

type UseCase struct {
	ruleGroupRepo port.RuleGroupRepository
	ruleRepo      port.RuleRepository
}

func NewUseCase(rgRepo port.RuleGroupRepository, rRepo port.RuleRepository) *UseCase {
	return &UseCase{ruleGroupRepo: rgRepo, ruleRepo: rRepo}
}

// RuleGroup methods

func (uc *UseCase) GetGroupByID(ctx context.Context, id uint) (*entity.RuleGroup, error) {
	return uc.ruleGroupRepo.FindByID(ctx, id)
}

func (uc *UseCase) ListGroups(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.RuleGroup, int64, error) {
	return uc.ruleGroupRepo.List(ctx, userGroupID, limit, offset)
}

func (uc *UseCase) CreateGroup(ctx context.Context, rg *entity.RuleGroup) error {
	return uc.ruleGroupRepo.Create(ctx, rg)
}

func (uc *UseCase) UpdateGroup(ctx context.Context, rg *entity.RuleGroup) error {
	return uc.ruleGroupRepo.Update(ctx, rg)
}

func (uc *UseCase) DeleteGroup(ctx context.Context, id uint) error {
	return uc.ruleGroupRepo.Delete(ctx, id)
}

func (uc *UseCase) ListGroupRules(ctx context.Context, ruleGroupID uint) ([]entity.Rule, error) {
	return uc.ruleGroupRepo.ListRules(ctx, ruleGroupID)
}

// Rule methods

func (uc *UseCase) GetRuleByID(ctx context.Context, id uint) (*entity.Rule, error) {
	return uc.ruleRepo.FindByID(ctx, id)
}

func (uc *UseCase) ListRules(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Rule, int64, error) {
	return uc.ruleRepo.List(ctx, userGroupID, limit, offset)
}

func (uc *UseCase) CreateRule(ctx context.Context, rule *entity.Rule) error {
	return uc.ruleRepo.Create(ctx, rule)
}

func (uc *UseCase) UpdateRule(ctx context.Context, rule *entity.Rule) error {
	return uc.ruleRepo.Update(ctx, rule)
}

func (uc *UseCase) DeleteRule(ctx context.Context, id uint) error {
	return uc.ruleRepo.Delete(ctx, id)
}

func (uc *UseCase) GetTriggers(ctx context.Context, ruleID uint) ([]entity.RuleTrigger, error) {
	return uc.ruleRepo.GetTriggers(ctx, ruleID)
}

func (uc *UseCase) SetTriggers(ctx context.Context, ruleID uint, triggers []entity.RuleTrigger) error {
	return uc.ruleRepo.SetTriggers(ctx, ruleID, triggers)
}

func (uc *UseCase) GetActions(ctx context.Context, ruleID uint) ([]entity.RuleAction, error) {
	return uc.ruleRepo.GetActions(ctx, ruleID)
}

func (uc *UseCase) SetActions(ctx context.Context, ruleID uint, actions []entity.RuleAction) error {
	return uc.ruleRepo.SetActions(ctx, ruleID, actions)
}

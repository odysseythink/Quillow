package repository

import (
	"context"
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/adapter/repository/model"
	"github.com/anthropics/firefly-iii-go/internal/entity"
	"gorm.io/gorm"
)

// RuleGroupRepo implements port.RuleGroupRepository

type RuleGroupRepo struct {
	db *gorm.DB
}

func NewRuleGroupRepository(db *gorm.DB) *RuleGroupRepo {
	return &RuleGroupRepo{db: db}
}

func (r *RuleGroupRepo) FindByID(ctx context.Context, id uint) (*entity.RuleGroup, error) {
	var m model.RuleGroupModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("rule group not found: %w", err)
	}
	return ruleGroupModelToEntity(&m), nil
}

func (r *RuleGroupRepo) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.RuleGroup, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.RuleGroupModel{}).Where("deleted_at IS NULL")
	if userGroupID > 0 {
		query = query.Where("user_group_id = ?", userGroupID)
	}

	var total int64
	query.Count(&total)

	var models []model.RuleGroupModel
	if err := query.Order("`order` ASC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	groups := make([]entity.RuleGroup, len(models))
	for i, m := range models {
		groups[i] = *ruleGroupModelToEntity(&m)
	}
	return groups, total, nil
}

func (r *RuleGroupRepo) Create(ctx context.Context, rg *entity.RuleGroup) error {
	m := ruleGroupEntityToModel(rg)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	rg.ID = m.ID
	rg.CreatedAt = m.CreatedAt
	rg.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *RuleGroupRepo) Update(ctx context.Context, rg *entity.RuleGroup) error {
	m := ruleGroupEntityToModel(rg)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *RuleGroupRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.RuleGroupModel{}, id).Error
}

func (r *RuleGroupRepo) ListRules(ctx context.Context, ruleGroupID uint) ([]entity.Rule, error) {
	var models []model.RuleModel
	if err := r.db.WithContext(ctx).Where("rule_group_id = ? AND deleted_at IS NULL", ruleGroupID).Order("`order` ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	rules := make([]entity.Rule, len(models))
	for i, m := range models {
		rules[i] = *ruleModelToEntity(&m)
	}
	return rules, nil
}

// RuleRepo implements port.RuleRepository

type RuleRepo struct {
	db *gorm.DB
}

func NewRuleRepository(db *gorm.DB) *RuleRepo {
	return &RuleRepo{db: db}
}

func (r *RuleRepo) FindByID(ctx context.Context, id uint) (*entity.Rule, error) {
	var m model.RuleModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("rule not found: %w", err)
	}
	return ruleModelToEntity(&m), nil
}

func (r *RuleRepo) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Rule, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.RuleModel{}).Where("deleted_at IS NULL")
	if userGroupID > 0 {
		query = query.Where("user_group_id = ?", userGroupID)
	}

	var total int64
	query.Count(&total)

	var models []model.RuleModel
	if err := query.Order("`order` ASC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	rules := make([]entity.Rule, len(models))
	for i, m := range models {
		rules[i] = *ruleModelToEntity(&m)
	}
	return rules, total, nil
}

func (r *RuleRepo) Create(ctx context.Context, rule *entity.Rule) error {
	m := ruleEntityToModel(rule)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	rule.ID = m.ID
	rule.CreatedAt = m.CreatedAt
	rule.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *RuleRepo) Update(ctx context.Context, rule *entity.Rule) error {
	m := ruleEntityToModel(rule)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *RuleRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.RuleModel{}, id).Error
}

func (r *RuleRepo) GetTriggers(ctx context.Context, ruleID uint) ([]entity.RuleTrigger, error) {
	var models []model.RuleTriggerModel
	if err := r.db.WithContext(ctx).Where("rule_id = ?", ruleID).Order("`order` ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	triggers := make([]entity.RuleTrigger, len(models))
	for i, m := range models {
		triggers[i] = *ruleTriggerModelToEntity(&m)
	}
	return triggers, nil
}

func (r *RuleRepo) SetTriggers(ctx context.Context, ruleID uint, triggers []entity.RuleTrigger) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("rule_id = ?", ruleID).Delete(&model.RuleTriggerModel{}).Error; err != nil {
			return err
		}
		if len(triggers) == 0 {
			return nil
		}
		models := make([]model.RuleTriggerModel, len(triggers))
		for i, t := range triggers {
			models[i] = *ruleTriggerEntityToModel(&t)
			models[i].RuleID = ruleID
		}
		return tx.Create(&models).Error
	})
}

func (r *RuleRepo) GetActions(ctx context.Context, ruleID uint) ([]entity.RuleAction, error) {
	var models []model.RuleActionModel
	if err := r.db.WithContext(ctx).Where("rule_id = ?", ruleID).Order("`order` ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	actions := make([]entity.RuleAction, len(models))
	for i, m := range models {
		actions[i] = *ruleActionModelToEntity(&m)
	}
	return actions, nil
}

func (r *RuleRepo) SetActions(ctx context.Context, ruleID uint, actions []entity.RuleAction) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("rule_id = ?", ruleID).Delete(&model.RuleActionModel{}).Error; err != nil {
			return err
		}
		if len(actions) == 0 {
			return nil
		}
		models := make([]model.RuleActionModel, len(actions))
		for i, a := range actions {
			models[i] = *ruleActionEntityToModel(&a)
			models[i].RuleID = ruleID
		}
		return tx.Create(&models).Error
	})
}

// Model-to-entity conversion functions

func ruleGroupModelToEntity(m *model.RuleGroupModel) *entity.RuleGroup {
	desc := ""
	if m.Description != nil {
		desc = *m.Description
	}
	return &entity.RuleGroup{
		ID:             m.ID,
		UserID:         m.UserID,
		UserGroupID:    m.UserGroupID,
		Title:          m.Title,
		Description:    desc,
		Order:          m.Order,
		Active:         m.Active,
		StopProcessing: m.StopProcessing,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      m.DeletedAt,
	}
}

func ruleGroupEntityToModel(rg *entity.RuleGroup) *model.RuleGroupModel {
	var desc *string
	if rg.Description != "" {
		desc = &rg.Description
	}
	return &model.RuleGroupModel{
		ID:             rg.ID,
		UserID:         rg.UserID,
		UserGroupID:    rg.UserGroupID,
		Title:          rg.Title,
		Description:    desc,
		Order:          rg.Order,
		Active:         rg.Active,
		StopProcessing: rg.StopProcessing,
	}
}

func ruleModelToEntity(m *model.RuleModel) *entity.Rule {
	desc := ""
	if m.Description != nil {
		desc = *m.Description
	}
	return &entity.Rule{
		ID:             m.ID,
		UserID:         m.UserID,
		UserGroupID:    m.UserGroupID,
		RuleGroupID:    m.RuleGroupID,
		Title:          m.Title,
		Description:    desc,
		Order:          m.Order,
		Active:         m.Active,
		StopProcessing: m.StopProcessing,
		Strict:         m.Strict,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      m.DeletedAt,
	}
}

func ruleEntityToModel(r *entity.Rule) *model.RuleModel {
	var desc *string
	if r.Description != "" {
		desc = &r.Description
	}
	return &model.RuleModel{
		ID:             r.ID,
		UserID:         r.UserID,
		UserGroupID:    r.UserGroupID,
		RuleGroupID:    r.RuleGroupID,
		Title:          r.Title,
		Description:    desc,
		Order:          r.Order,
		Active:         r.Active,
		StopProcessing: r.StopProcessing,
		Strict:         r.Strict,
	}
}

func ruleTriggerModelToEntity(m *model.RuleTriggerModel) *entity.RuleTrigger {
	return &entity.RuleTrigger{
		ID:             m.ID,
		RuleID:         m.RuleID,
		TriggerType:    m.TriggerType,
		TriggerValue:   m.TriggerValue,
		Order:          m.Order,
		Active:         m.Active,
		StopProcessing: m.StopProcessing,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func ruleTriggerEntityToModel(t *entity.RuleTrigger) *model.RuleTriggerModel {
	return &model.RuleTriggerModel{
		ID:             t.ID,
		RuleID:         t.RuleID,
		TriggerType:    t.TriggerType,
		TriggerValue:   t.TriggerValue,
		Order:          t.Order,
		Active:         t.Active,
		StopProcessing: t.StopProcessing,
	}
}

func ruleActionModelToEntity(m *model.RuleActionModel) *entity.RuleAction {
	return &entity.RuleAction{
		ID:             m.ID,
		RuleID:         m.RuleID,
		ActionType:     m.ActionType,
		ActionValue:    m.ActionValue,
		Order:          m.Order,
		Active:         m.Active,
		StopProcessing: m.StopProcessing,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func ruleActionEntityToModel(a *entity.RuleAction) *model.RuleActionModel {
	return &model.RuleActionModel{
		ID:             a.ID,
		RuleID:         a.RuleID,
		ActionType:     a.ActionType,
		ActionValue:    a.ActionValue,
		Order:          a.Order,
		Active:         a.Active,
		StopProcessing: a.StopProcessing,
	}
}

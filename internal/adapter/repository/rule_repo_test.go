package repository

import (
	"context"
	"testing"

	"github.com/anthropics/firefly-iii-go/internal/adapter/repository/model"
	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRuleDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&model.RuleGroupModel{},
		&model.RuleModel{},
		&model.RuleTriggerModel{},
		&model.RuleActionModel{},
	))
	return db
}

func TestRuleGroupRepo_CreateAndFind(t *testing.T) {
	db := setupRuleDB(t)
	repo := NewRuleGroupRepository(db)
	ctx := context.Background()

	rg := &entity.RuleGroup{
		UserID:      1,
		UserGroupID: 1,
		Title:       "Default Rules",
		Description: "Auto-generated rule group",
		Order:       1,
		Active:      true,
	}
	require.NoError(t, repo.Create(ctx, rg))
	assert.NotZero(t, rg.ID)

	found, err := repo.FindByID(ctx, rg.ID)
	require.NoError(t, err)
	assert.Equal(t, "Default Rules", found.Title)
	assert.Equal(t, "Auto-generated rule group", found.Description)
	assert.True(t, found.Active)
}

func TestRuleGroupRepo_List(t *testing.T) {
	db := setupRuleDB(t)
	repo := NewRuleGroupRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &entity.RuleGroup{UserID: 1, UserGroupID: 1, Title: "Group A", Order: 1, Active: true})
	repo.Create(ctx, &entity.RuleGroup{UserID: 1, UserGroupID: 1, Title: "Group B", Order: 2, Active: true})
	repo.Create(ctx, &entity.RuleGroup{UserID: 1, UserGroupID: 1, Title: "Group C", Order: 3, Active: true})

	groups, total, err := repo.List(ctx, 1, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, groups, 3)

	// Test pagination
	groups, total, err = repo.List(ctx, 1, 2, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, groups, 2)
}

func TestRuleGroupRepo_ListRules(t *testing.T) {
	db := setupRuleDB(t)
	rgRepo := NewRuleGroupRepository(db)
	rRepo := NewRuleRepository(db)
	ctx := context.Background()

	rg := &entity.RuleGroup{UserID: 1, UserGroupID: 1, Title: "Group", Order: 1, Active: true}
	require.NoError(t, rgRepo.Create(ctx, rg))

	rRepo.Create(ctx, &entity.Rule{UserID: 1, UserGroupID: 1, RuleGroupID: rg.ID, Title: "Rule 1", Order: 1, Active: true})
	rRepo.Create(ctx, &entity.Rule{UserID: 1, UserGroupID: 1, RuleGroupID: rg.ID, Title: "Rule 2", Order: 2, Active: true})

	rules, err := rgRepo.ListRules(ctx, rg.ID)
	require.NoError(t, err)
	assert.Len(t, rules, 2)
	assert.Equal(t, "Rule 1", rules[0].Title)
	assert.Equal(t, "Rule 2", rules[1].Title)
}

func TestRuleGroupRepo_Delete(t *testing.T) {
	db := setupRuleDB(t)
	repo := NewRuleGroupRepository(db)
	ctx := context.Background()

	rg := &entity.RuleGroup{UserID: 1, UserGroupID: 1, Title: "To Delete", Order: 1, Active: true}
	require.NoError(t, repo.Create(ctx, rg))
	require.NoError(t, repo.Delete(ctx, rg.ID))

	_, err := repo.FindByID(ctx, rg.ID)
	assert.Error(t, err)
}

func TestRuleRepo_CreateAndFind(t *testing.T) {
	db := setupRuleDB(t)
	repo := NewRuleRepository(db)
	ctx := context.Background()

	rule := &entity.Rule{
		UserID:      1,
		UserGroupID: 1,
		RuleGroupID: 1,
		Title:       "Auto-categorize groceries",
		Description: "Matches grocery store transactions",
		Order:       1,
		Active:      true,
		Strict:      true,
	}
	require.NoError(t, repo.Create(ctx, rule))
	assert.NotZero(t, rule.ID)

	found, err := repo.FindByID(ctx, rule.ID)
	require.NoError(t, err)
	assert.Equal(t, "Auto-categorize groceries", found.Title)
	assert.True(t, found.Strict)
}

func TestRuleRepo_TriggersCRUD(t *testing.T) {
	db := setupRuleDB(t)
	repo := NewRuleRepository(db)
	ctx := context.Background()

	rule := &entity.Rule{UserID: 1, UserGroupID: 1, RuleGroupID: 1, Title: "Test", Order: 1, Active: true}
	require.NoError(t, repo.Create(ctx, rule))

	triggers := []entity.RuleTrigger{
		{TriggerType: "description_contains", TriggerValue: "grocery", Order: 1, Active: true},
		{TriggerType: "description_starts", TriggerValue: "walmart", Order: 2, Active: true},
	}
	require.NoError(t, repo.SetTriggers(ctx, rule.ID, triggers))

	got, err := repo.GetTriggers(ctx, rule.ID)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "description_contains", got[0].TriggerType)
	assert.Equal(t, "grocery", got[0].TriggerValue)
	assert.Equal(t, rule.ID, got[0].RuleID)

	// Replace triggers
	newTriggers := []entity.RuleTrigger{
		{TriggerType: "description_is", TriggerValue: "exact match", Order: 1, Active: true},
	}
	require.NoError(t, repo.SetTriggers(ctx, rule.ID, newTriggers))

	got, err = repo.GetTriggers(ctx, rule.ID)
	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, "description_is", got[0].TriggerType)
}

func TestRuleRepo_ActionsCRUD(t *testing.T) {
	db := setupRuleDB(t)
	repo := NewRuleRepository(db)
	ctx := context.Background()

	rule := &entity.Rule{UserID: 1, UserGroupID: 1, RuleGroupID: 1, Title: "Test", Order: 1, Active: true}
	require.NoError(t, repo.Create(ctx, rule))

	actions := []entity.RuleAction{
		{ActionType: "set_category", ActionValue: "Groceries", Order: 1, Active: true},
		{ActionType: "add_tag", ActionValue: "food", Order: 2, Active: true},
		{ActionType: "set_description", ActionValue: "Grocery purchase", Order: 3, Active: true},
	}
	require.NoError(t, repo.SetActions(ctx, rule.ID, actions))

	got, err := repo.GetActions(ctx, rule.ID)
	require.NoError(t, err)
	assert.Len(t, got, 3)
	assert.Equal(t, "set_category", got[0].ActionType)
	assert.Equal(t, "Groceries", got[0].ActionValue)
	assert.Equal(t, rule.ID, got[0].RuleID)

	// Replace with fewer actions
	newActions := []entity.RuleAction{
		{ActionType: "set_category", ActionValue: "Food", Order: 1, Active: true},
	}
	require.NoError(t, repo.SetActions(ctx, rule.ID, newActions))

	got, err = repo.GetActions(ctx, rule.ID)
	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, "Food", got[0].ActionValue)
}

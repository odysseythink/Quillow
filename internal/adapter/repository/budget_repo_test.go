package repository

import (
	"context"
	"testing"
	"time"

	"github.com/anthropics/quillow/internal/adapter/repository/model"
	"github.com/anthropics/quillow/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupBudgetDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.BudgetModel{}, &model.BudgetLimitModel{}, &model.NoteModel{}))
	return db
}

func TestBudgetRepo_CreateAndFind(t *testing.T) {
	db := setupBudgetDB(t)
	repo := NewBudgetRepository(db)
	ctx := context.Background()

	budget := &entity.Budget{UserID: 1, UserGroupID: 1, Name: "Groceries", Active: true}
	require.NoError(t, repo.Create(ctx, budget))
	assert.NotZero(t, budget.ID)

	found, err := repo.FindByID(ctx, budget.ID)
	require.NoError(t, err)
	assert.Equal(t, "Groceries", found.Name)
	assert.True(t, found.Active)
}

func TestBudgetRepo_List(t *testing.T) {
	db := setupBudgetDB(t)
	repo := NewBudgetRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &entity.Budget{UserID: 1, UserGroupID: 1, Name: "A", Active: true})
	repo.Create(ctx, &entity.Budget{UserID: 1, UserGroupID: 1, Name: "B", Active: true})
	repo.Create(ctx, &entity.Budget{UserID: 1, UserGroupID: 2, Name: "C", Active: true})

	budgets, total, err := repo.List(ctx, 1, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, budgets, 2)

	// List all
	budgets, total, err = repo.List(ctx, 0, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, budgets, 3)
}

func TestBudgetRepo_Delete(t *testing.T) {
	db := setupBudgetDB(t)
	repo := NewBudgetRepository(db)
	ctx := context.Background()

	budget := &entity.Budget{UserID: 1, UserGroupID: 1, Name: "Del", Active: true}
	require.NoError(t, repo.Create(ctx, budget))
	require.NoError(t, repo.Delete(ctx, budget.ID))

	_, err := repo.FindByID(ctx, budget.ID)
	assert.Error(t, err)
}

func TestBudgetLimitRepo_CRUD(t *testing.T) {
	db := setupBudgetDB(t)
	budgetRepo := NewBudgetRepository(db)
	limitRepo := NewBudgetLimitRepository(db)
	ctx := context.Background()

	budget := &entity.Budget{UserID: 1, UserGroupID: 1, Name: "Food", Active: true}
	require.NoError(t, budgetRepo.Create(ctx, budget))

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	limit := &entity.BudgetLimit{
		BudgetID:              budget.ID,
		TransactionCurrencyID: 1,
		StartDate:             start,
		EndDate:               end,
		Amount:                "500.00",
		Period:                "monthly",
	}
	require.NoError(t, limitRepo.Create(ctx, limit))
	assert.NotZero(t, limit.ID)

	// Find
	found, err := limitRepo.FindByID(ctx, limit.ID)
	require.NoError(t, err)
	assert.Equal(t, "500.00", found.Amount)
	assert.Equal(t, "monthly", found.Period)

	// Update
	limit.Amount = "600.00"
	require.NoError(t, limitRepo.Update(ctx, limit))
	found, _ = limitRepo.FindByID(ctx, limit.ID)
	assert.Equal(t, "600.00", found.Amount)

	// ListByBudget
	limits, err := limitRepo.ListByBudget(ctx, budget.ID)
	require.NoError(t, err)
	assert.Len(t, limits, 1)

	// ListByPeriod
	limits, err = limitRepo.ListByPeriod(ctx, budget.ID, start, end)
	require.NoError(t, err)
	assert.Len(t, limits, 1)

	// Delete
	require.NoError(t, limitRepo.Delete(ctx, limit.ID))
	_, err = limitRepo.FindByID(ctx, limit.ID)
	assert.Error(t, err)
}

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/anthropics/firefly-iii-go/internal/adapter/repository/model"
	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRecurrenceDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&model.RecurrenceModel{},
		&model.RecurrenceRepetitionModel{},
		&model.RecurrenceTransactionModel{},
	))
	return db
}

func TestRecurrenceRepo_CreateAndFind(t *testing.T) {
	db := setupRecurrenceDB(t)
	repo := NewRecurrenceRepository(db)
	ctx := context.Background()

	rec := &entity.Recurrence{
		UserID:                1,
		UserGroupID:           1,
		TransactionTypeID:     1,
		TransactionCurrencyID: 1,
		Title:                 "Monthly Salary",
		Description:           "Recurring salary deposit",
		FirstDate:             time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Repetitions:           0,
		ApplyRules:            true,
		Active:                true,
	}
	require.NoError(t, repo.Create(ctx, rec))
	assert.NotZero(t, rec.ID)

	found, err := repo.FindByID(ctx, rec.ID)
	require.NoError(t, err)
	assert.Equal(t, "Monthly Salary", found.Title)
	assert.Equal(t, "Recurring salary deposit", found.Description)
	assert.True(t, found.Active)
	assert.True(t, found.ApplyRules)
}

func TestRecurrenceRepo_List(t *testing.T) {
	db := setupRecurrenceDB(t)
	repo := NewRecurrenceRepository(db)
	ctx := context.Background()

	firstDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	repo.Create(ctx, &entity.Recurrence{UserID: 1, UserGroupID: 1, TransactionTypeID: 1, TransactionCurrencyID: 1, Title: "A", FirstDate: firstDate, Active: true})
	repo.Create(ctx, &entity.Recurrence{UserID: 1, UserGroupID: 1, TransactionTypeID: 1, TransactionCurrencyID: 1, Title: "B", FirstDate: firstDate, Active: true})
	repo.Create(ctx, &entity.Recurrence{UserID: 1, UserGroupID: 2, TransactionTypeID: 1, TransactionCurrencyID: 1, Title: "C", FirstDate: firstDate, Active: true})

	recs, total, err := repo.List(ctx, 1, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, recs, 2)

	// List all
	recs, total, err = repo.List(ctx, 0, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, recs, 3)
}

func TestRecurrenceRepo_Delete(t *testing.T) {
	db := setupRecurrenceDB(t)
	repo := NewRecurrenceRepository(db)
	ctx := context.Background()

	rec := &entity.Recurrence{
		UserID:                1,
		UserGroupID:           1,
		TransactionTypeID:     1,
		TransactionCurrencyID: 1,
		Title:                 "ToDelete",
		FirstDate:             time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Active:                true,
	}
	require.NoError(t, repo.Create(ctx, rec))
	require.NoError(t, repo.Delete(ctx, rec.ID))

	_, err := repo.FindByID(ctx, rec.ID)
	assert.Error(t, err)
}

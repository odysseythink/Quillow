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

func setupPiggyBankDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.PiggyBankModel{}, &model.PiggyBankEventModel{}, &model.NoteModel{}))
	return db
}

func TestPiggyBankRepo_CreateAndFind(t *testing.T) {
	db := setupPiggyBankDB(t)
	repo := NewPiggyBankRepository(db)
	ctx := context.Background()

	piggy := &entity.PiggyBank{
		AccountID:    1,
		Name:         "Vacation Fund",
		TargetAmount: "1000.00",
		Active:       true,
	}
	require.NoError(t, repo.Create(ctx, piggy))
	assert.NotZero(t, piggy.ID)

	found, err := repo.FindByID(ctx, piggy.ID)
	require.NoError(t, err)
	assert.Equal(t, "Vacation Fund", found.Name)
	assert.Equal(t, "1000.00", found.TargetAmount)
	assert.Equal(t, uint(1), found.AccountID)
}

func TestPiggyBankRepo_List(t *testing.T) {
	db := setupPiggyBankDB(t)
	repo := NewPiggyBankRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &entity.PiggyBank{AccountID: 1, Name: "Fund A", Active: true, Order: 2})
	repo.Create(ctx, &entity.PiggyBank{AccountID: 1, Name: "Fund B", Active: true, Order: 1})
	repo.Create(ctx, &entity.PiggyBank{AccountID: 2, Name: "Fund C", Active: true, Order: 3})

	items, total, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, items, 3)
	// Should be ordered by order ASC
	assert.Equal(t, "Fund B", items[0].Name)

	byAccount, err := repo.ListByAccount(ctx, 1)
	require.NoError(t, err)
	assert.Len(t, byAccount, 2)
}

func TestPiggyBankRepo_Events(t *testing.T) {
	db := setupPiggyBankDB(t)
	repo := NewPiggyBankRepository(db)
	ctx := context.Background()

	piggy := &entity.PiggyBank{AccountID: 1, Name: "Test", Active: true}
	require.NoError(t, repo.Create(ctx, piggy))

	now := time.Now()
	evt := &entity.PiggyBankEvent{
		PiggyBankID: piggy.ID,
		Amount:      "50.00",
		Date:        now,
	}
	require.NoError(t, repo.CreateEvent(ctx, evt))
	assert.NotZero(t, evt.ID)

	events, err := repo.ListEvents(ctx, piggy.ID)
	require.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, "50.00", events[0].Amount)
}

func TestPiggyBankRepo_Delete(t *testing.T) {
	db := setupPiggyBankDB(t)
	repo := NewPiggyBankRepository(db)
	ctx := context.Background()

	piggy := &entity.PiggyBank{AccountID: 1, Name: "To Delete", Active: true}
	require.NoError(t, repo.Create(ctx, piggy))

	require.NoError(t, repo.Delete(ctx, piggy.ID))

	_, err := repo.FindByID(ctx, piggy.ID)
	assert.Error(t, err)
}

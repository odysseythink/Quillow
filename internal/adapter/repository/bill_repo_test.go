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

func setupBillDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.BillModel{}, &model.NoteModel{}))
	return db
}

func TestBillRepo_CreateAndFind(t *testing.T) {
	db := setupBillDB(t)
	repo := NewBillRepository(db)
	ctx := context.Background()

	bill := &entity.Bill{
		UserID:                1,
		UserGroupID:           1,
		TransactionCurrencyID: 1,
		Name:                  "Rent",
		AmountMin:             "1000.00",
		AmountMax:             "1200.00",
		Date:                  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		RepeatFreq:            "monthly",
		Active:                true,
		Automatch:             true,
	}
	require.NoError(t, repo.Create(ctx, bill))
	assert.NotZero(t, bill.ID)

	found, err := repo.FindByID(ctx, bill.ID)
	require.NoError(t, err)
	assert.Equal(t, "Rent", found.Name)
	assert.Equal(t, "1000.00", found.AmountMin)
	assert.True(t, found.Active)
}

func TestBillRepo_List(t *testing.T) {
	db := setupBillDB(t)
	repo := NewBillRepository(db)
	ctx := context.Background()

	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	repo.Create(ctx, &entity.Bill{UserID: 1, UserGroupID: 1, TransactionCurrencyID: 1, Name: "A", AmountMin: "10", AmountMax: "20", Date: date, RepeatFreq: "monthly", Active: true})
	repo.Create(ctx, &entity.Bill{UserID: 1, UserGroupID: 1, TransactionCurrencyID: 1, Name: "B", AmountMin: "10", AmountMax: "20", Date: date, RepeatFreq: "monthly", Active: true})
	repo.Create(ctx, &entity.Bill{UserID: 1, UserGroupID: 2, TransactionCurrencyID: 1, Name: "C", AmountMin: "10", AmountMax: "20", Date: date, RepeatFreq: "monthly", Active: true})

	bills, total, err := repo.List(ctx, 1, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, bills, 2)

	// List all
	bills, total, err = repo.List(ctx, 0, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, bills, 3)
}

func TestBillRepo_Delete(t *testing.T) {
	db := setupBillDB(t)
	repo := NewBillRepository(db)
	ctx := context.Background()

	bill := &entity.Bill{
		UserID:                1,
		UserGroupID:           1,
		TransactionCurrencyID: 1,
		Name:                  "Del",
		AmountMin:             "10",
		AmountMax:             "20",
		Date:                  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		RepeatFreq:            "monthly",
		Active:                true,
	}
	require.NoError(t, repo.Create(ctx, bill))
	require.NoError(t, repo.Delete(ctx, bill.ID))

	_, err := repo.FindByID(ctx, bill.ID)
	assert.Error(t, err)
}

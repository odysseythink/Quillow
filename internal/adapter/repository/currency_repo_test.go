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

func setupCurrencyDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.TransactionCurrencyModel{}, &model.CurrencyExchangeRateModel{}))
	return db
}

func TestCurrencyRepo_CreateAndFind(t *testing.T) {
	db := setupCurrencyDB(t)
	repo := NewCurrencyRepository(db)
	ctx := context.Background()

	curr := &entity.TransactionCurrency{Code: "EUR", Name: "Euro", Symbol: "€", DecimalPlaces: 2, Enabled: true}
	require.NoError(t, repo.Create(ctx, curr))
	assert.NotZero(t, curr.ID)

	found, err := repo.FindByCode(ctx, "EUR")
	require.NoError(t, err)
	assert.Equal(t, "Euro", found.Name)
	assert.Equal(t, "€", found.Symbol)
}

func TestCurrencyRepo_List(t *testing.T) {
	db := setupCurrencyDB(t)
	repo := NewCurrencyRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &entity.TransactionCurrency{Code: "EUR", Name: "Euro", Symbol: "€", DecimalPlaces: 2, Enabled: true})
	repo.Create(ctx, &entity.TransactionCurrency{Code: "USD", Name: "US Dollar", Symbol: "$", DecimalPlaces: 2, Enabled: true})

	currencies, total, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, currencies, 2)
}

func TestCurrencyRepo_GetPrimary(t *testing.T) {
	db := setupCurrencyDB(t)
	repo := NewCurrencyRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &entity.TransactionCurrency{Code: "EUR", Name: "Euro", Symbol: "€", DecimalPlaces: 2, Enabled: true})

	primary, err := repo.GetPrimary(ctx)
	require.NoError(t, err)
	assert.Equal(t, "EUR", primary.Code)
}

func TestCurrencyRepo_Delete(t *testing.T) {
	db := setupCurrencyDB(t)
	repo := NewCurrencyRepository(db)
	ctx := context.Background()

	curr := &entity.TransactionCurrency{Code: "GBP", Name: "Pound", Symbol: "£", DecimalPlaces: 2, Enabled: true}
	require.NoError(t, repo.Create(ctx, curr))
	require.NoError(t, repo.Delete(ctx, curr.ID))

	_, err := repo.FindByID(ctx, curr.ID)
	assert.Error(t, err)
}

func TestExchangeRateRepo_CreateAndFind(t *testing.T) {
	db := setupCurrencyDB(t)
	repo := NewExchangeRateRepository(db)
	ctx := context.Background()

	rate := &entity.CurrencyExchangeRate{
		UserGroupID:    1,
		FromCurrencyID: 1,
		ToCurrencyID:   2,
		Rate:           "1.0850",
	}
	require.NoError(t, repo.Create(ctx, rate))
	assert.NotZero(t, rate.ID)

	found, err := repo.FindByID(ctx, rate.ID)
	require.NoError(t, err)
	assert.Equal(t, "1.0850", found.Rate)
}

func TestExchangeRateRepo_List(t *testing.T) {
	db := setupCurrencyDB(t)
	repo := NewExchangeRateRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &entity.CurrencyExchangeRate{UserGroupID: 1, FromCurrencyID: 1, ToCurrencyID: 2, Rate: "1.08"})
	repo.Create(ctx, &entity.CurrencyExchangeRate{UserGroupID: 1, FromCurrencyID: 2, ToCurrencyID: 1, Rate: "0.92"})

	rates, total, err := repo.List(ctx, 1, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, rates, 2)
}

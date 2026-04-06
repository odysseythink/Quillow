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

func setupAccountDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.AccountModel{}, &model.AccountTypeModel{}, &model.AccountMetaModel{}, &model.NoteModel{}))
	// Seed account types
	db.Create(&model.AccountTypeModel{ID: 1, Type: "Asset account"})
	db.Create(&model.AccountTypeModel{ID: 2, Type: "Expense account"})
	db.Create(&model.AccountTypeModel{ID: 3, Type: "Revenue account"})
	return db
}

func TestAccountRepo_CreateAndFind(t *testing.T) {
	db := setupAccountDB(t)
	repo := NewAccountRepository(db)
	ctx := context.Background()

	acct := &entity.Account{UserID: 1, UserGroupID: 1, AccountTypeID: 1, Name: "Checking", Active: true}
	require.NoError(t, repo.Create(ctx, acct))
	assert.NotZero(t, acct.ID)

	found, err := repo.FindByID(ctx, acct.ID)
	require.NoError(t, err)
	assert.Equal(t, "Checking", found.Name)
	assert.True(t, found.Active)
}

func TestAccountRepo_List(t *testing.T) {
	db := setupAccountDB(t)
	repo := NewAccountRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &entity.Account{UserID: 1, UserGroupID: 1, AccountTypeID: 1, Name: "A", Active: true})
	repo.Create(ctx, &entity.Account{UserID: 1, UserGroupID: 1, AccountTypeID: 2, Name: "B", Active: true})
	repo.Create(ctx, &entity.Account{UserID: 1, UserGroupID: 1, AccountTypeID: 1, Name: "C", Active: true})

	// List all
	accounts, total, err := repo.List(ctx, 1, nil, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, accounts, 3)

	// Filter by type
	accounts, total, err = repo.List(ctx, 1, []string{"Asset account"}, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, accounts, 2)
}

func TestAccountRepo_Search(t *testing.T) {
	db := setupAccountDB(t)
	repo := NewAccountRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &entity.Account{UserID: 1, UserGroupID: 1, AccountTypeID: 1, Name: "Checking Account", Active: true})
	repo.Create(ctx, &entity.Account{UserID: 1, UserGroupID: 1, AccountTypeID: 1, Name: "Savings Account", Active: true})
	repo.Create(ctx, &entity.Account{UserID: 1, UserGroupID: 1, AccountTypeID: 2, Name: "Groceries", Active: true})

	results, err := repo.Search(ctx, 1, "Account", nil, 10)
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestAccountRepo_Meta(t *testing.T) {
	db := setupAccountDB(t)
	repo := NewAccountRepository(db)
	ctx := context.Background()

	acct := &entity.Account{UserID: 1, UserGroupID: 1, AccountTypeID: 1, Name: "Test", Active: true}
	require.NoError(t, repo.Create(ctx, acct))

	require.NoError(t, repo.SetMeta(ctx, acct.ID, "account_role", "defaultAsset"))
	val, err := repo.GetMeta(ctx, acct.ID, "account_role")
	require.NoError(t, err)
	assert.Equal(t, "defaultAsset", val)

	// Update existing meta
	require.NoError(t, repo.SetMeta(ctx, acct.ID, "account_role", "savingAsset"))
	val, _ = repo.GetMeta(ctx, acct.ID, "account_role")
	assert.Equal(t, "savingAsset", val)
}

func TestAccountRepo_Delete(t *testing.T) {
	db := setupAccountDB(t)
	repo := NewAccountRepository(db)
	ctx := context.Background()

	acct := &entity.Account{UserID: 1, UserGroupID: 1, AccountTypeID: 1, Name: "Del", Active: true}
	require.NoError(t, repo.Create(ctx, acct))
	require.NoError(t, repo.Delete(ctx, acct.ID))

	_, err := repo.FindByID(ctx, acct.ID)
	assert.Error(t, err)
}

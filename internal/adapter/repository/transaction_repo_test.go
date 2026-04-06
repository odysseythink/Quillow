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

func setupTransactionDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&model.TransactionGroupModel{},
		&model.TransactionJournalModel{},
		&model.TransactionModel{},
		&model.TransactionTypeModel{},
		&model.TransactionJournalMetaModel{},
		&model.NoteModel{},
	))
	// Seed transaction types
	db.Create(&model.TransactionTypeModel{ID: 1, Type: "Withdrawal"})
	db.Create(&model.TransactionTypeModel{ID: 2, Type: "Deposit"})
	db.Create(&model.TransactionTypeModel{ID: 3, Type: "Transfer"})
	return db
}

func TestTxRepo_CreateGroupAndJournal(t *testing.T) {
	db := setupTransactionDB(t)
	repo := NewTransactionRepository(db)
	ctx := context.Background()

	group := &entity.TransactionGroup{UserID: 1, UserGroupID: 1, Title: "Test Group"}
	require.NoError(t, repo.CreateGroup(ctx, group))
	assert.NotZero(t, group.ID)

	journal := &entity.TransactionJournal{
		UserID: 1, UserGroupID: 1,
		TransactionTypeID:     1,
		TransactionCurrencyID: 1,
		Description:           "Groceries",
		Date:                  time.Now(),
		Completed:             true,
		TransactionGroupID:    group.ID,
	}
	require.NoError(t, repo.CreateJournal(ctx, journal))
	assert.NotZero(t, journal.ID)

	// Create debit and credit transactions
	source := &entity.Transaction{
		TransactionJournalID:  journal.ID,
		AccountID:             1,
		TransactionCurrencyID: 1,
		Amount:                "-50.00",
	}
	dest := &entity.Transaction{
		TransactionJournalID:  journal.ID,
		AccountID:             2,
		TransactionCurrencyID: 1,
		Amount:                "50.00",
	}
	require.NoError(t, repo.CreateTransaction(ctx, source))
	require.NoError(t, repo.CreateTransaction(ctx, dest))

	// Verify
	txns, err := repo.ListByJournalID(ctx, journal.ID)
	require.NoError(t, err)
	assert.Len(t, txns, 2)
}

func TestTxRepo_FindGroupByID(t *testing.T) {
	db := setupTransactionDB(t)
	repo := NewTransactionRepository(db)
	ctx := context.Background()

	group := &entity.TransactionGroup{UserID: 1, UserGroupID: 1, Title: "Find me"}
	require.NoError(t, repo.CreateGroup(ctx, group))

	found, err := repo.FindGroupByID(ctx, group.ID)
	require.NoError(t, err)
	assert.Equal(t, "Find me", found.Title)
}

func TestTxRepo_ListGroups(t *testing.T) {
	db := setupTransactionDB(t)
	repo := NewTransactionRepository(db)
	ctx := context.Background()

	g1 := &entity.TransactionGroup{UserID: 1, UserGroupID: 1}
	g2 := &entity.TransactionGroup{UserID: 1, UserGroupID: 1}
	repo.CreateGroup(ctx, g1)
	repo.CreateGroup(ctx, g2)

	groups, total, err := repo.ListGroups(ctx, 0, "", nil, nil, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, groups, 2)
}

func TestTxRepo_DeleteGroup(t *testing.T) {
	db := setupTransactionDB(t)
	repo := NewTransactionRepository(db)
	ctx := context.Background()

	group := &entity.TransactionGroup{UserID: 1, UserGroupID: 1}
	repo.CreateGroup(ctx, group)
	require.NoError(t, repo.DeleteGroup(ctx, group.ID))

	_, err := repo.FindGroupByID(ctx, group.ID)
	assert.Error(t, err)
}

func TestTxRepo_JournalMeta(t *testing.T) {
	db := setupTransactionDB(t)
	repo := NewTransactionRepository(db)
	ctx := context.Background()

	group := &entity.TransactionGroup{UserID: 1, UserGroupID: 1}
	repo.CreateGroup(ctx, group)
	journal := &entity.TransactionJournal{
		UserID: 1, UserGroupID: 1,
		TransactionTypeID:     1,
		TransactionCurrencyID: 1,
		Description:           "Test",
		Date:                  time.Now(),
		TransactionGroupID:    group.ID,
	}
	repo.CreateJournal(ctx, journal)

	require.NoError(t, repo.SetJournalMeta(ctx, journal.ID, "external_id", "EXT-123"))
	val, _ := repo.GetJournalMeta(ctx, journal.ID, "external_id")
	assert.Equal(t, "EXT-123", val)
}

func TestTxRepo_SearchGroups(t *testing.T) {
	db := setupTransactionDB(t)
	repo := NewTransactionRepository(db)
	ctx := context.Background()

	g1 := &entity.TransactionGroup{UserID: 1, UserGroupID: 1}
	repo.CreateGroup(ctx, g1)
	j1 := &entity.TransactionJournal{UserID: 1, UserGroupID: 1, TransactionTypeID: 1, TransactionCurrencyID: 1, Description: "Groceries at Store", Date: time.Now(), TransactionGroupID: g1.ID}
	repo.CreateJournal(ctx, j1)

	g2 := &entity.TransactionGroup{UserID: 1, UserGroupID: 1}
	repo.CreateGroup(ctx, g2)
	j2 := &entity.TransactionJournal{UserID: 1, UserGroupID: 1, TransactionTypeID: 1, TransactionCurrencyID: 1, Description: "Coffee Shop", Date: time.Now(), TransactionGroupID: g2.ID}
	repo.CreateJournal(ctx, j2)

	results, total, err := repo.SearchGroups(ctx, 0, "Groceries", 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, results, 1)
}

func TestTxRepo_GetTransactionType(t *testing.T) {
	db := setupTransactionDB(t)
	repo := NewTransactionRepository(db)
	ctx := context.Background()

	tt, err := repo.GetTransactionType(ctx, "Withdrawal")
	require.NoError(t, err)
	assert.Equal(t, "Withdrawal", tt.Type)
	assert.Equal(t, uint(1), tt.ID)
}

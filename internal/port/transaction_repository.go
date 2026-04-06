package port

import (
	"context"
	"time"

	"github.com/anthropics/firefly-iii-go/internal/entity"
)

type TransactionRepository interface {
	// TransactionGroup operations
	FindGroupByID(ctx context.Context, id uint) (*entity.TransactionGroup, error)
	ListGroups(ctx context.Context, userGroupID uint, transactionType string, start, end *time.Time, limit, offset int) ([]entity.TransactionGroup, int64, error)
	CreateGroup(ctx context.Context, group *entity.TransactionGroup) error
	UpdateGroup(ctx context.Context, group *entity.TransactionGroup) error
	DeleteGroup(ctx context.Context, id uint) error

	// TransactionJournal operations
	FindJournalByID(ctx context.Context, id uint) (*entity.TransactionJournal, error)
	ListJournalsByGroupID(ctx context.Context, groupID uint) ([]entity.TransactionJournal, error)
	CreateJournal(ctx context.Context, journal *entity.TransactionJournal) error
	UpdateJournal(ctx context.Context, journal *entity.TransactionJournal) error
	DeleteJournal(ctx context.Context, id uint) error

	// Transaction (line item) operations
	ListByJournalID(ctx context.Context, journalID uint) ([]entity.Transaction, error)
	CreateTransaction(ctx context.Context, txn *entity.Transaction) error
	UpdateTransaction(ctx context.Context, txn *entity.Transaction) error

	// Meta operations
	GetJournalMeta(ctx context.Context, journalID uint, name string) (string, error)
	SetJournalMeta(ctx context.Context, journalID uint, name, value string) error
	GetJournalNotes(ctx context.Context, journalID uint) (string, error)
	SetJournalNotes(ctx context.Context, journalID uint, text string) error
	GetJournalTags(ctx context.Context, journalID uint) ([]string, error)
	SetJournalTags(ctx context.Context, journalID uint, tags []string) error

	// Transaction type
	GetTransactionType(ctx context.Context, typeName string) (*entity.TransactionType, error)
	GetTransactionTypeByID(ctx context.Context, id uint) (*entity.TransactionType, error)

	// Search
	SearchGroups(ctx context.Context, userGroupID uint, query string, limit, offset int) ([]entity.TransactionGroup, int64, error)
	CountByQuery(ctx context.Context, userGroupID uint, query string) (int64, error)

	// Insight aggregations
	SumByType(ctx context.Context, userGroupID uint, transactionType string, start, end time.Time) ([]InsightEntry, error)
	SumByAccount(ctx context.Context, userGroupID uint, accountID uint, transactionType string, start, end time.Time) ([]InsightEntry, error)
}

type InsightEntry struct {
	ID              uint
	Name            string
	Difference      string
	DifferenceFloat float64
	CurrencyID      uint
	CurrencyCode    string
}

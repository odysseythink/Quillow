package entity

import "time"

type TransactionType struct {
	ID        uint
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type TransactionGroup struct {
	ID          uint
	UserID      uint
	UserGroupID uint
	Title       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type TransactionJournal struct {
	ID                    uint
	UserID                uint
	UserGroupID           uint
	TransactionTypeID     uint
	BillID                *uint
	TransactionCurrencyID uint
	Description           string
	Date                  time.Time
	Order                 uint
	TagCount              int
	Encrypted             bool
	Completed             bool
	TransactionGroupID    uint
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time
}

type Transaction struct {
	ID                    uint
	TransactionJournalID  uint
	AccountID             uint
	TransactionCurrencyID uint
	ForeignCurrencyID     *uint
	Amount                string
	ForeignAmount         string
	Description           string
	Reconciled            bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time
}

type TransactionJournalMeta struct {
	ID                   uint
	TransactionJournalID uint
	Name                 string
	Data                 string
	Hash                 string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type TransactionJournalLink struct {
	ID            uint
	LinkTypeID    uint
	SourceID      uint
	DestinationID uint
	Comment       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

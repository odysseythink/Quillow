package model

import "time"

type TransactionTypeModel struct {
	ID        uint       `gorm:"primaryKey;column:id"`
	Type      string     `gorm:"column:type"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (TransactionTypeModel) TableName() string { return "transaction_types" }

type TransactionGroupModel struct {
	ID          uint       `gorm:"primaryKey;column:id"`
	UserID      uint       `gorm:"column:user_id"`
	UserGroupID uint       `gorm:"column:user_group_id"`
	Title       *string    `gorm:"column:title"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (TransactionGroupModel) TableName() string { return "transaction_groups" }

type TransactionJournalModel struct {
	ID                    uint       `gorm:"primaryKey;column:id"`
	UserID                uint       `gorm:"column:user_id"`
	UserGroupID           uint       `gorm:"column:user_group_id"`
	TransactionTypeID     uint       `gorm:"column:transaction_type_id"`
	BillID                *uint      `gorm:"column:bill_id"`
	TransactionCurrencyID uint       `gorm:"column:transaction_currency_id"`
	Description           string     `gorm:"column:description"`
	Date                  time.Time  `gorm:"column:date"`
	Order                 uint       `gorm:"column:order"`
	TagCount              int        `gorm:"column:tag_count"`
	Encrypted             bool       `gorm:"column:encrypted"`
	Completed             bool       `gorm:"column:completed"`
	TransactionGroupID    uint       `gorm:"column:transaction_group_id"`
	CreatedAt             time.Time  `gorm:"column:created_at"`
	UpdatedAt             time.Time  `gorm:"column:updated_at"`
	DeletedAt             *time.Time `gorm:"column:deleted_at"`
}

func (TransactionJournalModel) TableName() string { return "transaction_journals" }

type TransactionModel struct {
	ID                    uint       `gorm:"primaryKey;column:id"`
	TransactionJournalID  uint       `gorm:"column:transaction_journal_id"`
	AccountID             uint       `gorm:"column:account_id"`
	TransactionCurrencyID uint       `gorm:"column:transaction_currency_id"`
	ForeignCurrencyID     *uint      `gorm:"column:foreign_currency_id"`
	Amount                string     `gorm:"column:amount"`
	ForeignAmount         *string    `gorm:"column:foreign_amount"`
	Description           *string    `gorm:"column:description"`
	Reconciled            bool       `gorm:"column:reconciled"`
	CreatedAt             time.Time  `gorm:"column:created_at"`
	UpdatedAt             time.Time  `gorm:"column:updated_at"`
	DeletedAt             *time.Time `gorm:"column:deleted_at"`
}

func (TransactionModel) TableName() string { return "transactions" }

type TransactionJournalMetaModel struct {
	ID                   uint      `gorm:"primaryKey;column:id"`
	TransactionJournalID uint      `gorm:"column:transaction_journal_id"`
	Name                 string    `gorm:"column:name"`
	Data                 string    `gorm:"column:data"`
	Hash                 string    `gorm:"column:hash"`
	CreatedAt            time.Time `gorm:"column:created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at"`
}

func (TransactionJournalMetaModel) TableName() string { return "journal_meta" }

type TransactionJournalLinkModel struct {
	ID            uint      `gorm:"primaryKey;column:id"`
	LinkTypeID    uint      `gorm:"column:link_type_id"`
	SourceID      uint      `gorm:"column:source_id"`
	DestinationID uint      `gorm:"column:destination_id"`
	Comment       *string   `gorm:"column:comment"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (TransactionJournalLinkModel) TableName() string { return "journal_links" }

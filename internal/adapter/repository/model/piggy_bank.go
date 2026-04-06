package model

import "time"

type PiggyBankModel struct {
	ID           uint       `gorm:"primaryKey;column:id"`
	AccountID    uint       `gorm:"column:account_id"`
	Name         string     `gorm:"column:name"`
	TargetAmount *string    `gorm:"column:targetamount"`
	StartDate    *time.Time `gorm:"column:startdate"`
	TargetDate   *time.Time `gorm:"column:targetdate"`
	Order        uint       `gorm:"column:order"`
	Active       bool       `gorm:"column:active"`
	Encrypted    bool       `gorm:"column:encrypted"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}

func (PiggyBankModel) TableName() string { return "piggy_banks" }

type PiggyBankEventModel struct {
	ID                   uint      `gorm:"primaryKey;column:id"`
	PiggyBankID          uint      `gorm:"column:piggy_bank_id"`
	TransactionJournalID *uint     `gorm:"column:transaction_journal_id"`
	Amount               string    `gorm:"column:amount"`
	Date                 time.Time `gorm:"column:date"`
	CreatedAt            time.Time `gorm:"column:created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at"`
}

func (PiggyBankEventModel) TableName() string { return "piggy_bank_events" }

type PiggyBankRepetitionModel struct {
	ID            uint       `gorm:"primaryKey;column:id"`
	PiggyBankID   uint       `gorm:"column:piggy_bank_id"`
	StartDate     *time.Time `gorm:"column:startdate"`
	TargetDate    *time.Time `gorm:"column:targetdate"`
	CurrentAmount string     `gorm:"column:currentamount"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
}

func (PiggyBankRepetitionModel) TableName() string { return "piggy_bank_repetitions" }

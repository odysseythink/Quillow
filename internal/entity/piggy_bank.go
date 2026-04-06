package entity

import "time"

type PiggyBank struct {
	ID           uint
	AccountID    uint
	Name         string
	TargetAmount string
	StartDate    *time.Time
	TargetDate   *time.Time
	Order        uint
	Active       bool
	Encrypted    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type PiggyBankEvent struct {
	ID                   uint
	PiggyBankID          uint
	TransactionJournalID *uint
	Amount               string
	Date                 time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type PiggyBankRepetition struct {
	ID            uint
	PiggyBankID   uint
	StartDate     *time.Time
	TargetDate    *time.Time
	CurrentAmount string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

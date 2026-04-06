package entity

import "time"

type Budget struct {
	ID          uint
	UserID      uint
	UserGroupID uint
	Name        string
	Active      bool
	Encrypted   bool
	Order       uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type BudgetLimit struct {
	ID                    uint
	BudgetID              uint
	TransactionCurrencyID uint
	StartDate             time.Time
	EndDate               time.Time
	Amount                string
	Period                string
	Generated             bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type AutoBudget struct {
	ID                    uint
	BudgetID              uint
	TransactionCurrencyID uint
	AutoBudgetType        int
	Amount                string
	Period                string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type AvailableBudget struct {
	ID                    uint
	UserID                uint
	UserGroupID           uint
	TransactionCurrencyID uint
	Amount                string
	StartDate             time.Time
	EndDate               time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

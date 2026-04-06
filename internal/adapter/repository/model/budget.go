package model

import "time"

type BudgetModel struct {
	ID          uint       `gorm:"primaryKey;column:id"`
	UserID      uint       `gorm:"column:user_id"`
	UserGroupID uint       `gorm:"column:user_group_id"`
	Name        string     `gorm:"column:name"`
	Active      bool       `gorm:"column:active"`
	Encrypted   bool       `gorm:"column:encrypted"`
	Order       uint       `gorm:"column:order"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (BudgetModel) TableName() string { return "budgets" }

type BudgetLimitModel struct {
	ID                    uint      `gorm:"primaryKey;column:id"`
	BudgetID              uint      `gorm:"column:budget_id"`
	TransactionCurrencyID uint      `gorm:"column:transaction_currency_id"`
	StartDate             time.Time `gorm:"column:start_date"`
	EndDate               time.Time `gorm:"column:end_date"`
	Amount                string    `gorm:"column:amount"`
	Period                *string   `gorm:"column:period"`
	Generated             bool      `gorm:"column:generated"`
	CreatedAt             time.Time `gorm:"column:created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at"`
}

func (BudgetLimitModel) TableName() string { return "budget_limits" }

type AutoBudgetModel struct {
	ID                    uint      `gorm:"primaryKey;column:id"`
	BudgetID              uint      `gorm:"column:budget_id"`
	TransactionCurrencyID uint      `gorm:"column:transaction_currency_id"`
	AutoBudgetType        int       `gorm:"column:auto_budget_type"`
	Amount                string    `gorm:"column:amount"`
	Period                string    `gorm:"column:period"`
	CreatedAt             time.Time `gorm:"column:created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at"`
}

func (AutoBudgetModel) TableName() string { return "auto_budgets" }

type AvailableBudgetModel struct {
	ID                    uint      `gorm:"primaryKey;column:id"`
	UserID                uint      `gorm:"column:user_id"`
	UserGroupID           uint      `gorm:"column:user_group_id"`
	TransactionCurrencyID uint      `gorm:"column:transaction_currency_id"`
	Amount                string    `gorm:"column:amount"`
	StartDate             time.Time `gorm:"column:start_date"`
	EndDate               time.Time `gorm:"column:end_date"`
	CreatedAt             time.Time `gorm:"column:created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at"`
}

func (AvailableBudgetModel) TableName() string { return "available_budgets" }

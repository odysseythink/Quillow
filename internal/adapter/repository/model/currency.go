package model

import "time"

type TransactionCurrencyModel struct {
	ID            uint       `gorm:"primaryKey;column:id"`
	Code          string     `gorm:"column:code"`
	Name          string     `gorm:"column:name"`
	Symbol        string     `gorm:"column:symbol"`
	DecimalPlaces int        `gorm:"column:decimal_places"`
	Enabled       bool       `gorm:"column:enabled"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
}

func (TransactionCurrencyModel) TableName() string { return "transaction_currencies" }

type CurrencyExchangeRateModel struct {
	ID             uint       `gorm:"primaryKey;column:id"`
	UserGroupID    uint       `gorm:"column:user_group_id"`
	FromCurrencyID uint       `gorm:"column:from_currency_id"`
	ToCurrencyID   uint       `gorm:"column:to_currency_id"`
	Date           time.Time  `gorm:"column:date"`
	Rate           string     `gorm:"column:rate"`
	UserRate       *string    `gorm:"column:user_rate"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at"`
}

func (CurrencyExchangeRateModel) TableName() string { return "currency_exchange_rates" }

package entity

import "time"

type TransactionCurrency struct {
	ID            uint
	Code          string
	Name          string
	Symbol        string
	DecimalPlaces int
	Enabled       bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

type CurrencyExchangeRate struct {
	ID             uint
	UserGroupID    uint
	FromCurrencyID uint
	ToCurrencyID   uint
	Date           time.Time
	Rate           string
	UserRate       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

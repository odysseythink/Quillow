package entity

import "time"

type Recurrence struct {
	ID                    uint
	UserID                uint
	UserGroupID           uint
	TransactionTypeID     uint
	TransactionCurrencyID uint
	Title                 string
	Description           string
	FirstDate             time.Time
	RepeatUntil           *time.Time
	LatestDate            *time.Time
	Repetitions           uint
	ApplyRules            bool
	Active                bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time
}

type RecurrenceRepetition struct {
	ID               uint
	RecurrenceID     uint
	RepetitionType   string
	RepetitionMoment string
	RepetitionSkip   uint
	Weekend          uint
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type RecurrenceTransaction struct {
	ID                    uint
	RecurrenceID          uint
	TransactionCurrencyID uint
	ForeignCurrencyID     *uint
	SourceID              uint
	DestinationID         uint
	Amount                string
	ForeignAmount         string
	Description           string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type RecurrenceMeta struct {
	ID           uint
	RecurrenceID uint
	Name         string
	Value        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

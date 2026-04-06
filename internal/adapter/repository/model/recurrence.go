package model

import "time"

type RecurrenceModel struct {
	ID                    uint       `gorm:"primaryKey;column:id"`
	UserID                uint       `gorm:"column:user_id"`
	UserGroupID           uint       `gorm:"column:user_group_id"`
	TransactionTypeID     uint       `gorm:"column:transaction_type_id"`
	TransactionCurrencyID uint       `gorm:"column:transaction_currency_id"`
	Title                 string     `gorm:"column:title"`
	Description           string     `gorm:"column:description"`
	FirstDate             time.Time  `gorm:"column:first_date"`
	RepeatUntil           *time.Time `gorm:"column:repeat_until"`
	LatestDate            *time.Time `gorm:"column:latest_date"`
	Repetitions           uint       `gorm:"column:repetitions"`
	ApplyRules            bool       `gorm:"column:apply_rules"`
	Active                bool       `gorm:"column:active"`
	CreatedAt             time.Time  `gorm:"column:created_at"`
	UpdatedAt             time.Time  `gorm:"column:updated_at"`
	DeletedAt             *time.Time `gorm:"column:deleted_at"`
}

func (RecurrenceModel) TableName() string { return "recurrences" }

type RecurrenceRepetitionModel struct {
	ID               uint      `gorm:"primaryKey;column:id"`
	RecurrenceID     uint      `gorm:"column:recurrence_id"`
	RepetitionType   string    `gorm:"column:repetition_type"`
	RepetitionMoment string    `gorm:"column:repetition_moment"`
	RepetitionSkip   uint      `gorm:"column:repetition_skip"`
	Weekend          uint      `gorm:"column:weekend"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (RecurrenceRepetitionModel) TableName() string { return "recurrences_repetitions" }

type RecurrenceTransactionModel struct {
	ID                    uint      `gorm:"primaryKey;column:id"`
	RecurrenceID          uint      `gorm:"column:recurrence_id"`
	TransactionCurrencyID uint      `gorm:"column:transaction_currency_id"`
	ForeignCurrencyID     *uint     `gorm:"column:foreign_currency_id"`
	SourceID              uint      `gorm:"column:source_id"`
	DestinationID         uint      `gorm:"column:destination_id"`
	Amount                string    `gorm:"column:amount"`
	ForeignAmount         *string   `gorm:"column:foreign_amount"`
	Description           string    `gorm:"column:description"`
	CreatedAt             time.Time `gorm:"column:created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at"`
}

func (RecurrenceTransactionModel) TableName() string { return "recurrences_transactions" }

type RecurrenceMetaModel struct {
	ID           uint      `gorm:"primaryKey;column:id"`
	RecurrenceID uint      `gorm:"column:recurrence_id"`
	Name         string    `gorm:"column:name"`
	Value        string    `gorm:"column:value"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (RecurrenceMetaModel) TableName() string { return "recurrences_meta" }

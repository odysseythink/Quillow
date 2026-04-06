package model

import "time"

type AuditLogEntryModel struct {
	ID            uint      `gorm:"primaryKey;column:id"`
	AuditableID   uint      `gorm:"column:auditable_id"`
	AuditableType string    `gorm:"column:auditable_type"`
	ChangerID     uint      `gorm:"column:changer_id"`
	ChangerType   string    `gorm:"column:changer_type"`
	Action        string    `gorm:"column:action"`
	Before        *string   `gorm:"column:before"`
	After         *string   `gorm:"column:after"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (AuditLogEntryModel) TableName() string { return "audit_log_entries" }

type NoteModel struct {
	ID           uint       `gorm:"primaryKey;column:id"`
	NoteableID   uint       `gorm:"column:noteable_id"`
	NoteableType string     `gorm:"column:noteable_type"`
	Title        *string    `gorm:"column:title"`
	Text         *string    `gorm:"column:text"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}

func (NoteModel) TableName() string { return "notes" }

type LocationModel struct {
	ID            uint       `gorm:"primaryKey;column:id"`
	LocatableID   uint       `gorm:"column:locatable_id"`
	LocatableType string     `gorm:"column:locatable_type"`
	Latitude      *float64   `gorm:"column:latitude"`
	Longitude     *float64   `gorm:"column:longitude"`
	ZoomLevel     *int       `gorm:"column:zoom_level"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
}

func (LocationModel) TableName() string { return "locations" }

type PeriodStatisticModel struct {
	ID                    uint      `gorm:"primaryKey;column:id"`
	StatisticalID         uint      `gorm:"column:statistical_id"`
	StatisticalType       string    `gorm:"column:statistical_type"`
	Period                string    `gorm:"column:period"`
	TransactionCurrencyID uint      `gorm:"column:transaction_currency_id"`
	StartDate             time.Time `gorm:"column:start_date"`
	EndDate               time.Time `gorm:"column:end_date"`
	Amount                string    `gorm:"column:amount"`
	CreatedAt             time.Time `gorm:"column:created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at"`
}

func (PeriodStatisticModel) TableName() string { return "period_statistics" }

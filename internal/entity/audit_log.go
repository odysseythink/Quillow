package entity

import "time"

type AuditLogEntry struct {
	ID            uint
	AuditableID   uint
	AuditableType string
	ChangerID     uint
	ChangerType   string
	Action        string
	Before        string
	After         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Note struct {
	ID           uint
	NoteableID   uint
	NoteableType string
	Title        string
	Text         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type Location struct {
	ID            uint
	LocatableID   uint
	LocatableType string
	Latitude      *float64
	Longitude     *float64
	ZoomLevel     *int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

type PeriodStatistic struct {
	ID                    uint
	StatisticalID         uint
	StatisticalType       string
	Period                string
	TransactionCurrencyID uint
	StartDate             time.Time
	EndDate               time.Time
	Amount                string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

package entity

import "time"

type Bill struct {
	ID                    uint
	UserID                uint
	UserGroupID           uint
	TransactionCurrencyID uint
	Name                  string
	AmountMin             string
	AmountMax             string
	Date                  time.Time
	EndDate               *time.Time
	ExtensionDate         *time.Time
	RepeatFreq            string
	Skip                  uint
	Automatch             bool
	Active                bool
	NameEncrypted         bool
	MatchEncrypted        bool
	Order                 uint
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time
}

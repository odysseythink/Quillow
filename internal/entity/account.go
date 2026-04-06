package entity

import "time"

type AccountType struct {
	ID        uint
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Account struct {
	ID             uint
	UserID         uint
	UserGroupID    uint
	AccountTypeID  uint
	Name           string
	VirtualBalance string
	IBAN           string
	Active         bool
	Encrypted      bool
	Order          int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

type AccountMeta struct {
	ID        uint
	AccountID uint
	Name      string
	Data      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

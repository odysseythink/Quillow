package model

import "time"

type AccountTypeModel struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	Type      string    `gorm:"column:type"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (AccountTypeModel) TableName() string { return "account_types" }

type AccountModel struct {
	ID             uint       `gorm:"primaryKey;column:id"`
	UserID         uint       `gorm:"column:user_id"`
	UserGroupID    uint       `gorm:"column:user_group_id"`
	AccountTypeID  uint       `gorm:"column:account_type_id"`
	Name           string     `gorm:"column:name"`
	VirtualBalance *string    `gorm:"column:virtual_balance"`
	IBAN           *string    `gorm:"column:iban"`
	Active         bool       `gorm:"column:active"`
	Encrypted      bool       `gorm:"column:encrypted"`
	Order          int        `gorm:"column:order"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at"`
}

func (AccountModel) TableName() string { return "accounts" }

type AccountMetaModel struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	AccountID uint      `gorm:"column:account_id"`
	Name      string    `gorm:"column:name"`
	Data      string    `gorm:"column:data"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (AccountMetaModel) TableName() string { return "account_meta" }

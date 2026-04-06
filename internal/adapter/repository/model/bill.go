package model

import "time"

type BillModel struct {
	ID                    uint       `gorm:"primaryKey;column:id"`
	UserID                uint       `gorm:"column:user_id"`
	UserGroupID           uint       `gorm:"column:user_group_id"`
	TransactionCurrencyID uint       `gorm:"column:transaction_currency_id"`
	Name                  string     `gorm:"column:name"`
	AmountMin             string     `gorm:"column:amount_min"`
	AmountMax             string     `gorm:"column:amount_max"`
	Date                  time.Time  `gorm:"column:date"`
	EndDate               *time.Time `gorm:"column:end_date"`
	ExtensionDate         *time.Time `gorm:"column:extension_date"`
	RepeatFreq            string     `gorm:"column:repeat_freq"`
	Skip                  uint       `gorm:"column:skip"`
	Automatch             bool       `gorm:"column:automatch"`
	Active                bool       `gorm:"column:active"`
	NameEncrypted         bool       `gorm:"column:name_encrypted"`
	MatchEncrypted        bool       `gorm:"column:match_encrypted"`
	Order                 uint       `gorm:"column:order"`
	CreatedAt             time.Time  `gorm:"column:created_at"`
	UpdatedAt             time.Time  `gorm:"column:updated_at"`
	DeletedAt             *time.Time `gorm:"column:deleted_at"`
}

func (BillModel) TableName() string { return "bills" }

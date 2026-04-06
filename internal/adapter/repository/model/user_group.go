package model

import "time"

type UserGroupModel struct {
	ID        uint       `gorm:"primaryKey;column:id"`
	Title     string     `gorm:"column:title"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (UserGroupModel) TableName() string { return "user_groups" }

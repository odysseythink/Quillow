package model

import "time"

type CategoryModel struct {
	ID          uint       `gorm:"primaryKey;column:id"`
	UserID      uint       `gorm:"column:user_id"`
	UserGroupID uint       `gorm:"column:user_group_id"`
	Name        string     `gorm:"column:name"`
	Encrypted   bool       `gorm:"column:encrypted"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (CategoryModel) TableName() string { return "categories" }

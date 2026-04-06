package model

import "time"

type ObjectGroupModel struct {
	ID          uint       `gorm:"primaryKey;column:id"`
	UserID      uint       `gorm:"column:user_id"`
	UserGroupID uint       `gorm:"column:user_group_id"`
	Title       string     `gorm:"column:title"`
	Order       uint       `gorm:"column:order"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (ObjectGroupModel) TableName() string { return "object_groups" }

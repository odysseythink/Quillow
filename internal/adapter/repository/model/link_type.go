package model

import "time"

type LinkTypeModel struct {
	ID        uint       `gorm:"primaryKey;column:id"`
	Name      string     `gorm:"column:name"`
	Outward   string     `gorm:"column:outward"`
	Inward    string     `gorm:"column:inward"`
	Editable  bool       `gorm:"column:editable"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (LinkTypeModel) TableName() string { return "link_types" }

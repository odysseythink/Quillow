package model

import "time"

type ConfigurationModel struct {
	ID        uint       `gorm:"primaryKey;column:id"`
	Name      string     `gorm:"column:name"`
	Data      string     `gorm:"column:data"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (ConfigurationModel) TableName() string { return "configuration" }

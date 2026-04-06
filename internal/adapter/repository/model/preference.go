package model

import "time"

type PreferenceModel struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	UserID    uint      `gorm:"column:user_id"`
	Name      string    `gorm:"column:name"`
	Data      string    `gorm:"column:data"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (PreferenceModel) TableName() string { return "preferences" }

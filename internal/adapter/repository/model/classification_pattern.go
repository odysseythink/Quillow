package model

import "time"

type ClassificationPatternModel struct {
	ID         uint      `gorm:"primaryKey;column:id"`
	UserID     uint      `gorm:"column:user_id;index"`
	Pattern    string    `gorm:"column:pattern"`
	CategoryID uint      `gorm:"column:category_id"`
	TagIDs     string    `gorm:"column:tag_ids"`
	HitCount   uint      `gorm:"column:hit_count;default:1"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (ClassificationPatternModel) TableName() string { return "classification_patterns" }

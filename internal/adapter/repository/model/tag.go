package model

import "time"

type TagModel struct {
	ID          uint       `gorm:"primaryKey;column:id"`
	UserID      uint       `gorm:"column:user_id"`
	UserGroupID uint       `gorm:"column:user_group_id"`
	Tag         string     `gorm:"column:tag"`
	TagMode     string     `gorm:"column:tagMode"`
	Date        *time.Time `gorm:"column:date"`
	Description *string    `gorm:"column:description"`
	Latitude    *float64   `gorm:"column:latitude"`
	Longitude   *float64   `gorm:"column:longitude"`
	ZoomLevel   *int       `gorm:"column:zoomLevel"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (TagModel) TableName() string { return "tags" }

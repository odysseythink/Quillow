package entity

import "time"

type Tag struct {
	ID          uint
	UserID      uint
	UserGroupID uint
	Tag         string
	TagMode     string
	Date        *time.Time
	Description string
	Latitude    *float64
	Longitude   *float64
	ZoomLevel   *int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

package entity

import "time"

type ObjectGroup struct {
	ID          uint
	UserID      uint
	UserGroupID uint
	Title       string
	Order       uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

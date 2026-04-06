package entity

import "time"

type UserGroup struct {
	ID        uint
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

package entity

import "time"

type Category struct {
	ID          uint
	UserID      uint
	UserGroupID uint
	Name        string
	Encrypted   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

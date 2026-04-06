package entity

import "time"

type Preference struct {
	ID        uint
	UserID    uint
	Name      string
	Data      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

package entity

import "time"

type LinkType struct {
	ID        uint
	Name      string
	Outward   string
	Inward    string
	Editable  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

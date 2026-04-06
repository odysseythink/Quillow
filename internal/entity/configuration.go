package entity

import (
	"errors"
	"time"
)

var ErrNotEditable = errors.New("this configuration key is not editable")

type Configuration struct {
	ID        uint
	Name      string
	Data      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

package entity

import "time"

type Attachment struct {
	ID             uint
	UserID         uint
	UserGroupID    uint
	AttachableID   uint
	AttachableType string
	MD5            string
	Filename       string
	Title          string
	Description    string
	Mime           string
	Size           uint
	Uploaded       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

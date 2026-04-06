package model

import "time"

type AttachmentModel struct {
	ID             uint       `gorm:"primaryKey;column:id"`
	UserID         uint       `gorm:"column:user_id"`
	UserGroupID    uint       `gorm:"column:user_group_id"`
	AttachableID   uint       `gorm:"column:attachable_id"`
	AttachableType string     `gorm:"column:attachable_type"`
	MD5            string     `gorm:"column:md5"`
	Filename       string     `gorm:"column:filename"`
	Title          *string    `gorm:"column:title"`
	Description    *string    `gorm:"column:description"`
	Mime           string     `gorm:"column:mime"`
	Size           uint       `gorm:"column:size"`
	Uploaded       bool       `gorm:"column:uploaded"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at"`
}

func (AttachmentModel) TableName() string { return "attachments" }

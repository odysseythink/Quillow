package model

import "time"

type WebhookModel struct {
	ID          uint       `gorm:"primaryKey;column:id"`
	UserID      uint       `gorm:"column:user_id"`
	UserGroupID uint       `gorm:"column:user_group_id"`
	Active      bool       `gorm:"column:active"`
	Title       string     `gorm:"column:title"`
	Secret      string     `gorm:"column:secret"`
	Trigger     int        `gorm:"column:trigger"`
	Response    int        `gorm:"column:response"`
	Delivery    int        `gorm:"column:delivery"`
	URL         string     `gorm:"column:url"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (WebhookModel) TableName() string { return "webhooks" }

type WebhookMessageModel struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	WebhookID uint      `gorm:"column:webhook_id"`
	Sent      bool      `gorm:"column:sent"`
	Errored   bool      `gorm:"column:errored"`
	UUID      string    `gorm:"column:uuid"`
	Message   *string   `gorm:"column:message"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (WebhookMessageModel) TableName() string { return "webhook_messages" }

type WebhookAttemptModel struct {
	ID               uint      `gorm:"primaryKey;column:id"`
	WebhookMessageID uint      `gorm:"column:webhook_message_id"`
	StatusCode       int       `gorm:"column:status_code"`
	Logs             *string   `gorm:"column:logs"`
	Response         *string   `gorm:"column:response"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (WebhookAttemptModel) TableName() string { return "webhook_attempts" }

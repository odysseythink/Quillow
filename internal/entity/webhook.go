package entity

import "time"

type Webhook struct {
	ID          uint
	UserID      uint
	UserGroupID uint
	Active      bool
	Title       string
	Secret      string
	Trigger     int
	Response    int
	Delivery    int
	URL         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type WebhookMessage struct {
	ID        uint
	WebhookID uint
	Sent      bool
	Errored   bool
	UUID      string
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WebhookAttempt struct {
	ID               uint
	WebhookMessageID uint
	StatusCode       int
	Logs             string
	Response         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

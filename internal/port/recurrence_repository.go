package port

import (
	"context"

	"github.com/anthropics/quillow/internal/entity"
)

type RecurrenceRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.Recurrence, error)
	List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Recurrence, int64, error)
	Create(ctx context.Context, rec *entity.Recurrence) error
	Update(ctx context.Context, rec *entity.Recurrence) error
	Delete(ctx context.Context, id uint) error
	GetRepetitions(ctx context.Context, recurrenceID uint) ([]entity.RecurrenceRepetition, error)
	GetTransactions(ctx context.Context, recurrenceID uint) ([]entity.RecurrenceTransaction, error)
}

type WebhookRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.Webhook, error)
	List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Webhook, int64, error)
	Create(ctx context.Context, wh *entity.Webhook) error
	Update(ctx context.Context, wh *entity.Webhook) error
	Delete(ctx context.Context, id uint) error
	ListMessages(ctx context.Context, webhookID uint) ([]entity.WebhookMessage, error)
	CreateMessage(ctx context.Context, msg *entity.WebhookMessage) error
	ListAttempts(ctx context.Context, messageID uint) ([]entity.WebhookAttempt, error)
	CreateAttempt(ctx context.Context, attempt *entity.WebhookAttempt) error
}

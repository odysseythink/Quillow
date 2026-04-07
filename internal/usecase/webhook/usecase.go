package webhook

import (
	"context"

	"github.com/anthropics/quillow/internal/entity"
	"github.com/anthropics/quillow/internal/port"
)

type UseCase struct {
	webhookRepo port.WebhookRepository
}

func NewUseCase(webhookRepo port.WebhookRepository) *UseCase {
	return &UseCase{webhookRepo: webhookRepo}
}

func (uc *UseCase) GetByID(ctx context.Context, id uint) (*entity.Webhook, error) {
	return uc.webhookRepo.FindByID(ctx, id)
}

func (uc *UseCase) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Webhook, int64, error) {
	return uc.webhookRepo.List(ctx, userGroupID, limit, offset)
}

func (uc *UseCase) Create(ctx context.Context, wh *entity.Webhook) error {
	return uc.webhookRepo.Create(ctx, wh)
}

func (uc *UseCase) Update(ctx context.Context, wh *entity.Webhook) error {
	return uc.webhookRepo.Update(ctx, wh)
}

func (uc *UseCase) Delete(ctx context.Context, id uint) error {
	return uc.webhookRepo.Delete(ctx, id)
}

func (uc *UseCase) ListMessages(ctx context.Context, webhookID uint) ([]entity.WebhookMessage, error) {
	return uc.webhookRepo.ListMessages(ctx, webhookID)
}

func (uc *UseCase) ListAttempts(ctx context.Context, messageID uint) ([]entity.WebhookAttempt, error) {
	return uc.webhookRepo.ListAttempts(ctx, messageID)
}

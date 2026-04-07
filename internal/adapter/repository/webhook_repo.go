package repository

import (
	"context"
	"fmt"

	"github.com/anthropics/quillow/internal/adapter/repository/model"
	"github.com/anthropics/quillow/internal/entity"
	"gorm.io/gorm"
)

type WebhookRepository struct {
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

func (r *WebhookRepository) FindByID(ctx context.Context, id uint) (*entity.Webhook, error) {
	var m model.WebhookModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("webhook not found: %w", err)
	}
	return webhookModelToEntity(&m), nil
}

func (r *WebhookRepository) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Webhook, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.WebhookModel{}).Where("deleted_at IS NULL")
	if userGroupID > 0 {
		query = query.Where("user_group_id = ?", userGroupID)
	}

	var total int64
	query.Count(&total)

	var models []model.WebhookModel
	if err := query.Order("id ASC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]entity.Webhook, len(models))
	for i, m := range models {
		result[i] = *webhookModelToEntity(&m)
	}
	return result, total, nil
}

func (r *WebhookRepository) Create(ctx context.Context, wh *entity.Webhook) error {
	m := webhookEntityToModel(wh)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	wh.ID = m.ID
	wh.CreatedAt = m.CreatedAt
	wh.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *WebhookRepository) Update(ctx context.Context, wh *entity.Webhook) error {
	m := webhookEntityToModel(wh)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *WebhookRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.WebhookModel{}, id).Error
}

func (r *WebhookRepository) ListMessages(ctx context.Context, webhookID uint) ([]entity.WebhookMessage, error) {
	var models []model.WebhookMessageModel
	if err := r.db.WithContext(ctx).Where("webhook_id = ?", webhookID).Order("id ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	result := make([]entity.WebhookMessage, len(models))
	for i, m := range models {
		msg := ""
		if m.Message != nil {
			msg = *m.Message
		}
		result[i] = entity.WebhookMessage{
			ID:        m.ID,
			WebhookID: m.WebhookID,
			Sent:      m.Sent,
			Errored:   m.Errored,
			UUID:      m.UUID,
			Message:   msg,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}
	return result, nil
}

func (r *WebhookRepository) CreateMessage(ctx context.Context, msg *entity.WebhookMessage) error {
	m := &model.WebhookMessageModel{
		WebhookID: msg.WebhookID,
		Sent:      msg.Sent,
		Errored:   msg.Errored,
		UUID:      msg.UUID,
		Message:   &msg.Message,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	msg.ID = m.ID
	msg.CreatedAt = m.CreatedAt
	msg.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *WebhookRepository) ListAttempts(ctx context.Context, messageID uint) ([]entity.WebhookAttempt, error) {
	var models []model.WebhookAttemptModel
	if err := r.db.WithContext(ctx).Where("webhook_message_id = ?", messageID).Order("id ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	result := make([]entity.WebhookAttempt, len(models))
	for i, m := range models {
		logs := ""
		if m.Logs != nil {
			logs = *m.Logs
		}
		resp := ""
		if m.Response != nil {
			resp = *m.Response
		}
		result[i] = entity.WebhookAttempt{
			ID:               m.ID,
			WebhookMessageID: m.WebhookMessageID,
			StatusCode:       m.StatusCode,
			Logs:             logs,
			Response:         resp,
			CreatedAt:        m.CreatedAt,
			UpdatedAt:        m.UpdatedAt,
		}
	}
	return result, nil
}

func (r *WebhookRepository) CreateAttempt(ctx context.Context, attempt *entity.WebhookAttempt) error {
	m := &model.WebhookAttemptModel{
		WebhookMessageID: attempt.WebhookMessageID,
		StatusCode:       attempt.StatusCode,
		Logs:             &attempt.Logs,
		Response:         &attempt.Response,
		CreatedAt:        attempt.CreatedAt,
		UpdatedAt:        attempt.UpdatedAt,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	attempt.ID = m.ID
	attempt.CreatedAt = m.CreatedAt
	attempt.UpdatedAt = m.UpdatedAt
	return nil
}

// ---------------------------------------------------------------------------
// Conversion helpers
// ---------------------------------------------------------------------------

func webhookModelToEntity(m *model.WebhookModel) *entity.Webhook {
	return &entity.Webhook{
		ID:          m.ID,
		UserID:      m.UserID,
		UserGroupID: m.UserGroupID,
		Active:      m.Active,
		Title:       m.Title,
		Secret:      m.Secret,
		Trigger:     m.Trigger,
		Response:    m.Response,
		Delivery:    m.Delivery,
		URL:         m.URL,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}
}

func webhookEntityToModel(e *entity.Webhook) *model.WebhookModel {
	return &model.WebhookModel{
		ID:          e.ID,
		UserID:      e.UserID,
		UserGroupID: e.UserGroupID,
		Active:      e.Active,
		Title:       e.Title,
		Secret:      e.Secret,
		Trigger:     e.Trigger,
		Response:    e.Response,
		Delivery:    e.Delivery,
		URL:         e.URL,
	}
}

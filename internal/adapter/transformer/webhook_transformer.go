package transformer

import (
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/pkg/response"
)

func TransformWebhook(w *entity.Webhook) response.Resource {
	return response.Resource{
		Type: "webhooks",
		ID:   fmt.Sprintf("%d", w.ID),
		Attributes: map[string]any{
			"created_at": w.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at": w.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"active":     w.Active,
			"title":      w.Title,
			"trigger":    w.Trigger,
			"response":   w.Response,
			"delivery":   w.Delivery,
			"url":        w.URL,
		},
	}
}

func TransformWebhookMessage(m *entity.WebhookMessage) response.Resource {
	return response.Resource{
		Type: "webhook_messages",
		ID:   fmt.Sprintf("%d", m.ID),
		Attributes: map[string]any{
			"created_at": m.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at": m.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"webhook_id": fmt.Sprintf("%d", m.WebhookID),
			"sent":       m.Sent,
			"errored":    m.Errored,
			"uuid":       m.UUID,
			"message":    m.Message,
		},
	}
}

func TransformWebhookAttempt(a *entity.WebhookAttempt) response.Resource {
	return response.Resource{
		Type: "webhook_attempts",
		ID:   fmt.Sprintf("%d", a.ID),
		Attributes: map[string]any{
			"created_at":         a.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at":         a.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"webhook_message_id": fmt.Sprintf("%d", a.WebhookMessageID),
			"status_code":        a.StatusCode,
			"response":           a.Response,
		},
	}
}

package transformer

import (
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/pkg/response"
)

func TransformTag(t *entity.Tag) response.Resource {
	attrs := map[string]any{
		"created_at":  t.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
		"updated_at":  t.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
		"tag":         t.Tag,
		"description": t.Description,
		"latitude":    t.Latitude,
		"longitude":   t.Longitude,
		"zoom_level":  t.ZoomLevel,
	}
	if t.Date != nil {
		attrs["date"] = t.Date.Format("2006-01-02")
	}
	return response.Resource{
		Type:       "tags",
		ID:         fmt.Sprintf("%d", t.ID),
		Attributes: attrs,
	}
}

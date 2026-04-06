package transformer

import (
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/pkg/response"
)

func TransformObjectGroup(og *entity.ObjectGroup) response.Resource {
	return response.Resource{
		Type: "object_groups",
		ID:   fmt.Sprintf("%d", og.ID),
		Attributes: map[string]any{
			"created_at": og.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at": og.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"title":      og.Title,
			"order":      og.Order,
		},
	}
}

package transformer

import (
	"fmt"

	"github.com/anthropics/quillow/internal/entity"
	"github.com/anthropics/quillow/pkg/response"
)

func TransformCategory(c *entity.Category, notes string) response.Resource {
	return response.Resource{
		Type: "categories",
		ID:   fmt.Sprintf("%d", c.ID),
		Attributes: map[string]any{
			"created_at": c.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at": c.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"name":       c.Name,
			"notes":      notes,
		},
	}
}

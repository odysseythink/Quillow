package transformer

import (
	"fmt"

	"github.com/anthropics/quillow/internal/entity"
	"github.com/anthropics/quillow/pkg/response"
)

func TransformUser(user *entity.User, role string) response.Resource {
	return response.Resource{
		Type: "users",
		ID:   fmt.Sprintf("%d", user.ID),
		Attributes: map[string]any{
			"created_at":   user.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at":   user.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"email":        user.Email,
			"blocked":      user.Blocked,
			"blocked_code": user.BlockedCode,
			"role":         role,
		},
	}
}

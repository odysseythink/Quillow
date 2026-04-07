package transformer

import (
	"encoding/json"
	"fmt"

	"github.com/anthropics/quillow/internal/entity"
	"github.com/anthropics/quillow/pkg/response"
)

func TransformPreference(pref *entity.Preference) response.Resource {
	var data any
	if err := json.Unmarshal([]byte(pref.Data), &data); err != nil {
		data = pref.Data
	}

	return response.Resource{
		Type: "preferences",
		ID:   fmt.Sprintf("%d", pref.ID),
		Attributes: map[string]any{
			"created_at": pref.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at": pref.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"name":       pref.Name,
			"data":       data,
		},
	}
}

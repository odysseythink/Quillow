package transformer

import (
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/pkg/response"
)

func TransformAttachment(att *entity.Attachment) response.Resource {
	attachableType := att.AttachableType
	// Simplify the type name for API response
	switch {
	case len(attachableType) > 20:
		// Extract last part after backslash
		for i := len(attachableType) - 1; i >= 0; i-- {
			if attachableType[i] == '\\' {
				attachableType = attachableType[i+1:]
				break
			}
		}
	}

	return response.Resource{
		Type: "attachments",
		ID:   fmt.Sprintf("%d", att.ID),
		Attributes: map[string]any{
			"created_at":      att.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at":      att.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"attachable_id":   fmt.Sprintf("%d", att.AttachableID),
			"attachable_type": attachableType,
			"hash":            att.MD5,
			"filename":        att.Filename,
			"title":           att.Title,
			"notes":           att.Description,
			"mime":            att.Mime,
			"size":            att.Size,
			"download_url":    fmt.Sprintf("/api/v1/attachments/%d/download", att.ID),
			"upload_url":      fmt.Sprintf("/api/v1/attachments/%d/upload", att.ID),
		},
	}
}

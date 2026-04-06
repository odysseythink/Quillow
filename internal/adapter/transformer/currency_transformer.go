package transformer

import (
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/pkg/response"
)

func TransformCurrency(curr *entity.TransactionCurrency, isPrimary bool) response.Resource {
	return response.Resource{
		Type: "currencies",
		ID:   fmt.Sprintf("%d", curr.ID),
		Attributes: map[string]any{
			"created_at":     curr.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at":     curr.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"enabled":        curr.Enabled,
			"primary":        isPrimary,
			"default":        isPrimary,
			"native":         isPrimary,
			"name":           curr.Name,
			"code":           curr.Code,
			"symbol":         curr.Symbol,
			"decimal_places": curr.DecimalPlaces,
		},
	}
}

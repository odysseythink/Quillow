package transformer

import (
	"fmt"

	"github.com/anthropics/quillow/internal/entity"
	"github.com/anthropics/quillow/pkg/response"
)

func TransformBill(b *entity.Bill, notes string) response.Resource {
	attrs := map[string]any{
		"created_at":              b.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
		"updated_at":              b.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
		"name":                    b.Name,
		"amount_min":              b.AmountMin,
		"amount_max":              b.AmountMax,
		"date":                    b.Date.Format("2006-01-02"),
		"repeat_freq":             b.RepeatFreq,
		"skip":                    b.Skip,
		"automatch":               b.Automatch,
		"active":                  b.Active,
		"order":                   b.Order,
		"transaction_currency_id": fmt.Sprintf("%d", b.TransactionCurrencyID),
		"notes":                   notes,
	}
	if b.EndDate != nil {
		attrs["end_date"] = b.EndDate.Format("2006-01-02")
	}
	if b.ExtensionDate != nil {
		attrs["extension_date"] = b.ExtensionDate.Format("2006-01-02")
	}
	return response.Resource{
		Type:       "bills",
		ID:         fmt.Sprintf("%d", b.ID),
		Attributes: attrs,
	}
}

package transformer

import (
	"fmt"

	"github.com/anthropics/quillow/internal/entity"
	"github.com/anthropics/quillow/pkg/response"
)

func TransformPiggyBank(p *entity.PiggyBank, notes string) response.Resource {
	attrs := map[string]any{
		"created_at":    p.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
		"updated_at":    p.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
		"account_id":    fmt.Sprintf("%d", p.AccountID),
		"name":          p.Name,
		"target_amount": p.TargetAmount,
		"order":         p.Order,
		"active":        p.Active,
		"notes":         notes,
	}
	if p.StartDate != nil {
		attrs["start_date"] = p.StartDate.Format("2006-01-02")
	}
	if p.TargetDate != nil {
		attrs["target_date"] = p.TargetDate.Format("2006-01-02")
	}
	return response.Resource{
		Type:       "piggy_banks",
		ID:         fmt.Sprintf("%d", p.ID),
		Attributes: attrs,
	}
}

func TransformPiggyBankEvent(e *entity.PiggyBankEvent) response.Resource {
	return response.Resource{
		Type: "piggy_bank_events",
		ID:   fmt.Sprintf("%d", e.ID),
		Attributes: map[string]any{
			"created_at":              e.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at":              e.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"piggy_bank_id":           fmt.Sprintf("%d", e.PiggyBankID),
			"transaction_journal_id":  e.TransactionJournalID,
			"amount":                  e.Amount,
			"date":                    e.Date.Format("2006-01-02"),
		},
	}
}

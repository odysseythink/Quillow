package transformer

import (
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/pkg/response"
)

func TransformBudget(b *entity.Budget, notes string) response.Resource {
	return response.Resource{
		Type: "budgets",
		ID:   fmt.Sprintf("%d", b.ID),
		Attributes: map[string]any{
			"created_at": b.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at": b.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"name":       b.Name,
			"active":     b.Active,
			"order":      b.Order,
			"notes":      notes,
		},
	}
}

func TransformBudgetLimit(bl *entity.BudgetLimit) response.Resource {
	return response.Resource{
		Type: "budget_limits",
		ID:   fmt.Sprintf("%d", bl.ID),
		Attributes: map[string]any{
			"created_at":            bl.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at":            bl.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"budget_id":             fmt.Sprintf("%d", bl.BudgetID),
			"transaction_currency_id": fmt.Sprintf("%d", bl.TransactionCurrencyID),
			"amount":                bl.Amount,
			"start":                 bl.StartDate.Format("2006-01-02"),
			"end":                   bl.EndDate.Format("2006-01-02"),
			"period":                bl.Period,
		},
	}
}

package transformer

import (
	"fmt"

	"github.com/anthropics/quillow/internal/entity"
	"github.com/anthropics/quillow/pkg/response"
)

func TransformRecurrence(r *entity.Recurrence, reps []entity.RecurrenceRepetition, txns []entity.RecurrenceTransaction) response.Resource {
	repData := make([]map[string]any, len(reps))
	for i, rep := range reps {
		repData[i] = map[string]any{
			"id":                fmt.Sprintf("%d", rep.ID),
			"repetition_type":   rep.RepetitionType,
			"repetition_moment": rep.RepetitionMoment,
			"repetition_skip":   rep.RepetitionSkip,
			"weekend":           rep.Weekend,
		}
	}

	txnData := make([]map[string]any, len(txns))
	for i, tx := range txns {
		txnData[i] = map[string]any{
			"id":                      fmt.Sprintf("%d", tx.ID),
			"transaction_currency_id": fmt.Sprintf("%d", tx.TransactionCurrencyID),
			"source_id":               fmt.Sprintf("%d", tx.SourceID),
			"destination_id":          fmt.Sprintf("%d", tx.DestinationID),
			"amount":                  tx.Amount,
			"description":             tx.Description,
		}
	}

	attrs := map[string]any{
		"created_at":  r.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
		"updated_at":  r.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
		"title":       r.Title,
		"description": r.Description,
		"first_date":  r.FirstDate.Format("2006-01-02"),
		"active":      r.Active,
		"apply_rules": r.ApplyRules,
		"repetitions": repData,
		"transactions": txnData,
	}
	if r.RepeatUntil != nil {
		attrs["repeat_until"] = r.RepeatUntil.Format("2006-01-02")
	}

	return response.Resource{
		Type:       "recurrences",
		ID:         fmt.Sprintf("%d", r.ID),
		Attributes: attrs,
	}
}

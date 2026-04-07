package transformer

import (
	"fmt"

	"github.com/anthropics/quillow/internal/entity"
	"github.com/anthropics/quillow/pkg/response"
)

type TransactionJournalExtra struct {
	TypeName        string
	CurrencyCode    string
	CurrencySymbol  string
	CurrencyDP      int
	SourceID        uint
	SourceName      string
	SourceType      string
	DestinationID   uint
	DestinationName string
	DestinationType string
	Amount          string
	ForeignAmount   string
	BudgetID        string
	BudgetName      string
	CategoryID      string
	CategoryName    string
	BillID          string
	BillName        string
	Tags            []string
	Notes           string
	Reconciled      bool
	ExternalID      string
	ExternalURL     string
	InternalRef     string
	HasAttachments  bool
}

func TransformTransactionGroup(group *entity.TransactionGroup, journals []TransactionJournalExtra) response.Resource {
	txns := make([]map[string]any, len(journals))
	for i, j := range journals {
		txns[i] = map[string]any{
			"type":                   j.TypeName,
			"description":           "", // set below
			"currency_code":         j.CurrencyCode,
			"currency_symbol":       j.CurrencySymbol,
			"currency_decimal_places": j.CurrencyDP,
			"amount":                j.Amount,
			"foreign_amount":        j.ForeignAmount,
			"source_id":             fmt.Sprintf("%d", j.SourceID),
			"source_name":           j.SourceName,
			"source_type":           j.SourceType,
			"destination_id":        fmt.Sprintf("%d", j.DestinationID),
			"destination_name":      j.DestinationName,
			"destination_type":      j.DestinationType,
			"budget_id":             j.BudgetID,
			"budget_name":           j.BudgetName,
			"category_id":           j.CategoryID,
			"category_name":         j.CategoryName,
			"bill_id":               j.BillID,
			"bill_name":             j.BillName,
			"reconciled":            j.Reconciled,
			"tags":                  j.Tags,
			"notes":                 j.Notes,
			"external_id":           j.ExternalID,
			"external_url":          j.ExternalURL,
			"internal_reference":    j.InternalRef,
			"has_attachments":       j.HasAttachments,
		}
	}

	return response.Resource{
		Type: "transactions",
		ID:   fmt.Sprintf("%d", group.ID),
		Attributes: map[string]any{
			"created_at":   group.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at":   group.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"group_title":  group.Title,
			"transactions": txns,
		},
	}
}

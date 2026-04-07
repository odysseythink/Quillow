package transformer

import (
	"fmt"

	"github.com/anthropics/quillow/internal/entity"
	"github.com/anthropics/quillow/pkg/response"
)

type AccountExtra struct {
	AccountType    string
	AccountRole    string
	CurrencyID     string
	CurrencyName   string
	CurrencyCode   string
	CurrencySymbol string
	CurrencyDP     int
	CurrentBalance string
	Notes          string
	AccountNumber  string
	IBAN           string
}

func TransformAccount(acct *entity.Account, extra AccountExtra) response.Resource {
	return response.Resource{
		Type: "accounts",
		ID:   fmt.Sprintf("%d", acct.ID),
		Attributes: map[string]any{
			"created_at":              acct.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at":             acct.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"active":                 acct.Active,
			"order":                  acct.Order,
			"name":                   acct.Name,
			"type":                   extra.AccountType,
			"account_role":           extra.AccountRole,
			"currency_id":            extra.CurrencyID,
			"currency_name":          extra.CurrencyName,
			"currency_code":          extra.CurrencyCode,
			"currency_symbol":        extra.CurrencySymbol,
			"currency_decimal_places": extra.CurrencyDP,
			"current_balance":        extra.CurrentBalance,
			"virtual_balance":        acct.VirtualBalance,
			"iban":                   acct.IBAN,
			"account_number":         extra.AccountNumber,
			"notes":                  extra.Notes,
		},
	}
}

type AutocompleteAccount struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Active       bool   `json:"active"`
	Type         string `json:"type"`
	CurrencyID   string `json:"currency_id"`
	CurrencyCode string `json:"currency_code"`
	CurrencyName string `json:"currency_name"`
}

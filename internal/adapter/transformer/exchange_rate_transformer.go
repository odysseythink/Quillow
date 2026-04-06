package transformer

import (
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/pkg/response"
)

type ExchangeRateExtra struct {
	FromCurrency *entity.TransactionCurrency
	ToCurrency   *entity.TransactionCurrency
}

func TransformExchangeRate(rate *entity.CurrencyExchangeRate, extra ExchangeRateExtra) response.Resource {
	attrs := map[string]any{
		"created_at": rate.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
		"updated_at": rate.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
		"rate":       rate.Rate,
		"date":       rate.Date.Format("2006-01-02T15:04:05-07:00"),
	}

	if extra.FromCurrency != nil {
		attrs["from_currency_id"] = fmt.Sprintf("%d", extra.FromCurrency.ID)
		attrs["from_currency_name"] = extra.FromCurrency.Name
		attrs["from_currency_code"] = extra.FromCurrency.Code
		attrs["from_currency_symbol"] = extra.FromCurrency.Symbol
		attrs["from_currency_decimal_places"] = extra.FromCurrency.DecimalPlaces
	}
	if extra.ToCurrency != nil {
		attrs["to_currency_id"] = fmt.Sprintf("%d", extra.ToCurrency.ID)
		attrs["to_currency_name"] = extra.ToCurrency.Name
		attrs["to_currency_code"] = extra.ToCurrency.Code
		attrs["to_currency_symbol"] = extra.ToCurrency.Symbol
		attrs["to_currency_decimal_places"] = extra.ToCurrency.DecimalPlaces
	}

	return response.Resource{
		Type: "exchange_rates",
		ID:   fmt.Sprintf("%d", rate.ID),
		Attributes: attrs,
	}
}

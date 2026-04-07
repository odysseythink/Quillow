package port

import (
	"context"
	"time"

	"github.com/anthropics/quillow/internal/entity"
)

type CurrencyRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.TransactionCurrency, error)
	FindByCode(ctx context.Context, code string) (*entity.TransactionCurrency, error)
	List(ctx context.Context, limit, offset int) ([]entity.TransactionCurrency, int64, error)
	Create(ctx context.Context, currency *entity.TransactionCurrency) error
	Update(ctx context.Context, currency *entity.TransactionCurrency) error
	Delete(ctx context.Context, id uint) error
	GetPrimary(ctx context.Context) (*entity.TransactionCurrency, error)
}

type ExchangeRateRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.CurrencyExchangeRate, error)
	FindByPair(ctx context.Context, fromID, toID uint, date time.Time) (*entity.CurrencyExchangeRate, error)
	List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.CurrencyExchangeRate, int64, error)
	ListByPair(ctx context.Context, fromCode, toCode string) ([]entity.CurrencyExchangeRate, error)
	Create(ctx context.Context, rate *entity.CurrencyExchangeRate) error
	Update(ctx context.Context, rate *entity.CurrencyExchangeRate) error
	Delete(ctx context.Context, id uint) error
	DeleteByPair(ctx context.Context, fromID, toID uint) error
}

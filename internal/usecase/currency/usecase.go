package currency

import (
	"context"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/internal/port"
)

type UseCase struct {
	repo port.CurrencyRepository
}

func NewUseCase(repo port.CurrencyRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetByID(ctx context.Context, id uint) (*entity.TransactionCurrency, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *UseCase) GetByCode(ctx context.Context, code string) (*entity.TransactionCurrency, error) {
	return uc.repo.FindByCode(ctx, code)
}

func (uc *UseCase) GetPrimary(ctx context.Context) (*entity.TransactionCurrency, error) {
	return uc.repo.GetPrimary(ctx)
}

func (uc *UseCase) List(ctx context.Context, limit, offset int) ([]entity.TransactionCurrency, int64, error) {
	return uc.repo.List(ctx, limit, offset)
}

func (uc *UseCase) Create(ctx context.Context, name, code, symbol string, decimalPlaces int, enabled bool) (*entity.TransactionCurrency, error) {
	curr := &entity.TransactionCurrency{
		Name:          name,
		Code:          code,
		Symbol:        symbol,
		DecimalPlaces: decimalPlaces,
		Enabled:       enabled,
	}
	if err := uc.repo.Create(ctx, curr); err != nil {
		return nil, err
	}
	return curr, nil
}

func (uc *UseCase) Update(ctx context.Context, id uint, name, code, symbol string, decimalPlaces int, enabled bool) (*entity.TransactionCurrency, error) {
	curr, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	curr.Name = name
	curr.Code = code
	curr.Symbol = symbol
	curr.DecimalPlaces = decimalPlaces
	curr.Enabled = enabled
	if err := uc.repo.Update(ctx, curr); err != nil {
		return nil, err
	}
	return curr, nil
}

func (uc *UseCase) Enable(ctx context.Context, code string) (*entity.TransactionCurrency, error) {
	curr, err := uc.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	curr.Enabled = true
	if err := uc.repo.Update(ctx, curr); err != nil {
		return nil, err
	}
	return curr, nil
}

func (uc *UseCase) Disable(ctx context.Context, code string) (*entity.TransactionCurrency, error) {
	curr, err := uc.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	curr.Enabled = false
	if err := uc.repo.Update(ctx, curr); err != nil {
		return nil, err
	}
	return curr, nil
}

func (uc *UseCase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}

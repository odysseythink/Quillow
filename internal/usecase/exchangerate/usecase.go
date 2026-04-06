package exchangerate

import (
	"context"
	"time"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/internal/port"
)

type UseCase struct {
	repo     port.ExchangeRateRepository
	currRepo port.CurrencyRepository
}

func NewUseCase(repo port.ExchangeRateRepository, currRepo port.CurrencyRepository) *UseCase {
	return &UseCase{repo: repo, currRepo: currRepo}
}

func (uc *UseCase) GetByID(ctx context.Context, id uint) (*entity.CurrencyExchangeRate, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *UseCase) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.CurrencyExchangeRate, int64, error) {
	return uc.repo.List(ctx, userGroupID, limit, offset)
}

func (uc *UseCase) ListByPair(ctx context.Context, fromCode, toCode string) ([]entity.CurrencyExchangeRate, error) {
	return uc.repo.ListByPair(ctx, fromCode, toCode)
}

func (uc *UseCase) GetByPairAndDate(ctx context.Context, fromCode, toCode string, date time.Time) (*entity.CurrencyExchangeRate, error) {
	from, err := uc.currRepo.FindByCode(ctx, fromCode)
	if err != nil {
		return nil, err
	}
	to, err := uc.currRepo.FindByCode(ctx, toCode)
	if err != nil {
		return nil, err
	}
	return uc.repo.FindByPair(ctx, from.ID, to.ID, date)
}

func (uc *UseCase) Create(ctx context.Context, userGroupID, fromCurrencyID, toCurrencyID uint, rate string, date time.Time) (*entity.CurrencyExchangeRate, error) {
	er := &entity.CurrencyExchangeRate{
		UserGroupID:    userGroupID,
		FromCurrencyID: fromCurrencyID,
		ToCurrencyID:   toCurrencyID,
		Rate:           rate,
		Date:           date,
	}
	if err := uc.repo.Create(ctx, er); err != nil {
		return nil, err
	}
	return er, nil
}

func (uc *UseCase) Update(ctx context.Context, id uint, rate string) (*entity.CurrencyExchangeRate, error) {
	er, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	er.Rate = rate
	if err := uc.repo.Update(ctx, er); err != nil {
		return nil, err
	}
	return er, nil
}

func (uc *UseCase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *UseCase) DeleteByPair(ctx context.Context, fromCode, toCode string) error {
	from, err := uc.currRepo.FindByCode(ctx, fromCode)
	if err != nil {
		return err
	}
	to, err := uc.currRepo.FindByCode(ctx, toCode)
	if err != nil {
		return err
	}
	return uc.repo.DeleteByPair(ctx, from.ID, to.ID)
}

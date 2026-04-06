package piggybank

import (
	"context"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/internal/port"
)

type UseCase struct {
	repo port.PiggyBankRepository
}

func NewUseCase(repo port.PiggyBankRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetByID(ctx context.Context, id uint) (*entity.PiggyBank, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *UseCase) List(ctx context.Context, limit, offset int) ([]entity.PiggyBank, int64, error) {
	return uc.repo.List(ctx, limit, offset)
}

func (uc *UseCase) ListByAccount(ctx context.Context, accountID uint) ([]entity.PiggyBank, error) {
	return uc.repo.ListByAccount(ctx, accountID)
}

func (uc *UseCase) Create(ctx context.Context, piggy *entity.PiggyBank) error {
	return uc.repo.Create(ctx, piggy)
}

func (uc *UseCase) Update(ctx context.Context, piggy *entity.PiggyBank) error {
	return uc.repo.Update(ctx, piggy)
}

func (uc *UseCase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *UseCase) ListEvents(ctx context.Context, piggyBankID uint) ([]entity.PiggyBankEvent, error) {
	return uc.repo.ListEvents(ctx, piggyBankID)
}

func (uc *UseCase) AddEvent(ctx context.Context, event *entity.PiggyBankEvent) error {
	return uc.repo.CreateEvent(ctx, event)
}

func (uc *UseCase) GetNotes(ctx context.Context, piggyBankID uint) (string, error) {
	return uc.repo.GetNotes(ctx, piggyBankID)
}

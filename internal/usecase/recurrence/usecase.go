package recurrence

import (
	"context"

	"github.com/anthropics/quillow/internal/entity"
	"github.com/anthropics/quillow/internal/port"
)

type UseCase struct {
	recurrenceRepo port.RecurrenceRepository
}

func NewUseCase(recurrenceRepo port.RecurrenceRepository) *UseCase {
	return &UseCase{recurrenceRepo: recurrenceRepo}
}

func (uc *UseCase) GetByID(ctx context.Context, id uint) (*entity.Recurrence, error) {
	return uc.recurrenceRepo.FindByID(ctx, id)
}

func (uc *UseCase) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Recurrence, int64, error) {
	return uc.recurrenceRepo.List(ctx, userGroupID, limit, offset)
}

func (uc *UseCase) Create(ctx context.Context, rec *entity.Recurrence) error {
	return uc.recurrenceRepo.Create(ctx, rec)
}

func (uc *UseCase) Update(ctx context.Context, rec *entity.Recurrence) error {
	return uc.recurrenceRepo.Update(ctx, rec)
}

func (uc *UseCase) Delete(ctx context.Context, id uint) error {
	return uc.recurrenceRepo.Delete(ctx, id)
}

func (uc *UseCase) GetRepetitions(ctx context.Context, recurrenceID uint) ([]entity.RecurrenceRepetition, error) {
	return uc.recurrenceRepo.GetRepetitions(ctx, recurrenceID)
}

func (uc *UseCase) GetTransactions(ctx context.Context, recurrenceID uint) ([]entity.RecurrenceTransaction, error) {
	return uc.recurrenceRepo.GetTransactions(ctx, recurrenceID)
}

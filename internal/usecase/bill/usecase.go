package bill

import (
	"context"

	"github.com/anthropics/quillow/internal/entity"
	"github.com/anthropics/quillow/internal/port"
)

type UseCase struct {
	billRepo port.BillRepository
}

func NewUseCase(billRepo port.BillRepository) *UseCase {
	return &UseCase{billRepo: billRepo}
}

func (uc *UseCase) GetByID(ctx context.Context, id uint) (*entity.Bill, error) {
	return uc.billRepo.FindByID(ctx, id)
}

func (uc *UseCase) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Bill, int64, error) {
	return uc.billRepo.List(ctx, userGroupID, limit, offset)
}

func (uc *UseCase) ListActive(ctx context.Context, userGroupID uint) ([]entity.Bill, error) {
	return uc.billRepo.ListActive(ctx, userGroupID)
}

func (uc *UseCase) Create(ctx context.Context, bill *entity.Bill) error {
	return uc.billRepo.Create(ctx, bill)
}

func (uc *UseCase) Update(ctx context.Context, bill *entity.Bill) error {
	return uc.billRepo.Update(ctx, bill)
}

func (uc *UseCase) Delete(ctx context.Context, id uint) error {
	return uc.billRepo.Delete(ctx, id)
}

func (uc *UseCase) GetNotes(ctx context.Context, billID uint) (string, error) {
	return uc.billRepo.GetNotes(ctx, billID)
}

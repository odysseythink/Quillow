package budget

import (
	"context"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/internal/port"
)

type UseCase struct {
	budgetRepo port.BudgetRepository
	limitRepo  port.BudgetLimitRepository
}

func NewUseCase(budgetRepo port.BudgetRepository, limitRepo port.BudgetLimitRepository) *UseCase {
	return &UseCase{budgetRepo: budgetRepo, limitRepo: limitRepo}
}

func (uc *UseCase) GetByID(ctx context.Context, id uint) (*entity.Budget, error) {
	return uc.budgetRepo.FindByID(ctx, id)
}

func (uc *UseCase) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Budget, int64, error) {
	return uc.budgetRepo.List(ctx, userGroupID, limit, offset)
}

func (uc *UseCase) Create(ctx context.Context, budget *entity.Budget) error {
	return uc.budgetRepo.Create(ctx, budget)
}

func (uc *UseCase) Update(ctx context.Context, budget *entity.Budget) error {
	return uc.budgetRepo.Update(ctx, budget)
}

func (uc *UseCase) Delete(ctx context.Context, id uint) error {
	return uc.budgetRepo.Delete(ctx, id)
}

func (uc *UseCase) GetNotes(ctx context.Context, budgetID uint) (string, error) {
	return uc.budgetRepo.GetNotes(ctx, budgetID)
}

func (uc *UseCase) ListLimits(ctx context.Context, budgetID uint) ([]entity.BudgetLimit, error) {
	return uc.limitRepo.ListByBudget(ctx, budgetID)
}

func (uc *UseCase) CreateLimit(ctx context.Context, limit *entity.BudgetLimit) error {
	return uc.limitRepo.Create(ctx, limit)
}

func (uc *UseCase) UpdateLimit(ctx context.Context, limit *entity.BudgetLimit) error {
	return uc.limitRepo.Update(ctx, limit)
}

func (uc *UseCase) DeleteLimit(ctx context.Context, id uint) error {
	return uc.limitRepo.Delete(ctx, id)
}

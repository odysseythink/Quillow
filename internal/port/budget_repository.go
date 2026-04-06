package port

import (
	"context"
	"time"

	"github.com/anthropics/firefly-iii-go/internal/entity"
)

type BudgetRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.Budget, error)
	List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Budget, int64, error)
	ListActive(ctx context.Context, userGroupID uint) ([]entity.Budget, error)
	Create(ctx context.Context, budget *entity.Budget) error
	Update(ctx context.Context, budget *entity.Budget) error
	Delete(ctx context.Context, id uint) error
	GetNotes(ctx context.Context, budgetID uint) (string, error)
	SetNotes(ctx context.Context, budgetID uint, text string) error
}

type BudgetLimitRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.BudgetLimit, error)
	ListByBudget(ctx context.Context, budgetID uint) ([]entity.BudgetLimit, error)
	ListByPeriod(ctx context.Context, budgetID uint, start, end time.Time) ([]entity.BudgetLimit, error)
	Create(ctx context.Context, limit *entity.BudgetLimit) error
	Update(ctx context.Context, limit *entity.BudgetLimit) error
	Delete(ctx context.Context, id uint) error
}

type BillRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.Bill, error)
	List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Bill, int64, error)
	ListActive(ctx context.Context, userGroupID uint) ([]entity.Bill, error)
	Create(ctx context.Context, bill *entity.Bill) error
	Update(ctx context.Context, bill *entity.Bill) error
	Delete(ctx context.Context, id uint) error
	GetNotes(ctx context.Context, billID uint) (string, error)
	SetNotes(ctx context.Context, billID uint, text string) error
}

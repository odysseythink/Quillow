package category

import (
	"context"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/internal/port"
)

type UseCase struct {
	repo port.CategoryRepository
}

func NewUseCase(repo port.CategoryRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetByID(ctx context.Context, id uint) (*entity.Category, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *UseCase) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Category, int64, error) {
	return uc.repo.List(ctx, userGroupID, limit, offset)
}

func (uc *UseCase) Search(ctx context.Context, userGroupID uint, query string, limit int) ([]entity.Category, error) {
	return uc.repo.Search(ctx, userGroupID, query, limit)
}

func (uc *UseCase) Create(ctx context.Context, category *entity.Category) error {
	return uc.repo.Create(ctx, category)
}

func (uc *UseCase) Update(ctx context.Context, category *entity.Category) error {
	return uc.repo.Update(ctx, category)
}

func (uc *UseCase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *UseCase) GetNotes(ctx context.Context, categoryID uint) (string, error) {
	return uc.repo.GetNotes(ctx, categoryID)
}

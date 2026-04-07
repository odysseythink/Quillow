package tag

import (
	"context"

	"github.com/anthropics/quillow/internal/entity"
	"github.com/anthropics/quillow/internal/port"
)

type UseCase struct {
	repo port.TagRepository
}

func NewUseCase(repo port.TagRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetByID(ctx context.Context, id uint) (*entity.Tag, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *UseCase) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Tag, int64, error) {
	return uc.repo.List(ctx, userGroupID, limit, offset)
}

func (uc *UseCase) Search(ctx context.Context, userGroupID uint, query string, limit int) ([]entity.Tag, error) {
	return uc.repo.Search(ctx, userGroupID, query, limit)
}

func (uc *UseCase) Create(ctx context.Context, tag *entity.Tag) error {
	return uc.repo.Create(ctx, tag)
}

func (uc *UseCase) Update(ctx context.Context, tag *entity.Tag) error {
	return uc.repo.Update(ctx, tag)
}

func (uc *UseCase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}

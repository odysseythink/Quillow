package objectgroup

import (
	"context"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/internal/port"
)

type UseCase struct {
	repo port.ObjectGroupRepository
}

func NewUseCase(repo port.ObjectGroupRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetByID(ctx context.Context, id uint) (*entity.ObjectGroup, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *UseCase) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.ObjectGroup, int64, error) {
	return uc.repo.List(ctx, userGroupID, limit, offset)
}

func (uc *UseCase) Create(ctx context.Context, og *entity.ObjectGroup) error {
	return uc.repo.Create(ctx, og)
}

func (uc *UseCase) Update(ctx context.Context, og *entity.ObjectGroup) error {
	return uc.repo.Update(ctx, og)
}

func (uc *UseCase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}

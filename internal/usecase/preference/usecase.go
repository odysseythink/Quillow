package preference

import (
	"context"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/internal/port"
)

type UseCase struct {
	repo port.PreferenceRepository
}

func NewUseCase(repo port.PreferenceRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetByName(ctx context.Context, userID uint, name string) (*entity.Preference, error) {
	return uc.repo.FindByUserAndName(ctx, userID, name)
}

func (uc *UseCase) List(ctx context.Context, userID uint, limit, offset int) ([]entity.Preference, int64, error) {
	return uc.repo.ListByUser(ctx, userID, limit, offset)
}

func (uc *UseCase) Store(ctx context.Context, userID uint, name, data string) (*entity.Preference, error) {
	existing, err := uc.repo.FindByUserAndName(ctx, userID, name)
	if err == nil {
		existing.Data = data
		if err := uc.repo.Update(ctx, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	pref := &entity.Preference{
		UserID: userID,
		Name:   name,
		Data:   data,
	}
	if err := uc.repo.Create(ctx, pref); err != nil {
		return nil, err
	}
	return pref, nil
}

func (uc *UseCase) Update(ctx context.Context, userID uint, name, data string) (*entity.Preference, error) {
	pref, err := uc.repo.FindByUserAndName(ctx, userID, name)
	if err != nil {
		return nil, err
	}
	pref.Data = data
	if err := uc.repo.Update(ctx, pref); err != nil {
		return nil, err
	}
	return pref, nil
}

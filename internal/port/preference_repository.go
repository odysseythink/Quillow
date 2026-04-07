package port

import (
	"context"

	"github.com/anthropics/quillow/internal/entity"
)

type PreferenceRepository interface {
	FindByUserAndName(ctx context.Context, userID uint, name string) (*entity.Preference, error)
	ListByUser(ctx context.Context, userID uint, limit, offset int) ([]entity.Preference, int64, error)
	Create(ctx context.Context, pref *entity.Preference) error
	Update(ctx context.Context, pref *entity.Preference) error
}

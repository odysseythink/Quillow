package port

import (
	"context"

	"github.com/anthropics/quillow/internal/entity"
)

type ConfigurationRepository interface {
	FindByName(ctx context.Context, name string) (*entity.Configuration, error)
	List(ctx context.Context) ([]entity.Configuration, error)
	Upsert(ctx context.Context, config *entity.Configuration) error
}

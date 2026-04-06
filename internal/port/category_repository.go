package port

import (
	"context"

	"github.com/anthropics/firefly-iii-go/internal/entity"
)

type CategoryRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.Category, error)
	FindByName(ctx context.Context, userGroupID uint, name string) (*entity.Category, error)
	List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Category, int64, error)
	Search(ctx context.Context, userGroupID uint, query string, limit int) ([]entity.Category, error)
	Create(ctx context.Context, category *entity.Category) error
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uint) error
	GetNotes(ctx context.Context, categoryID uint) (string, error)
	SetNotes(ctx context.Context, categoryID uint, text string) error
}

type TagRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.Tag, error)
	FindByTag(ctx context.Context, userGroupID uint, tag string) (*entity.Tag, error)
	List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Tag, int64, error)
	Search(ctx context.Context, userGroupID uint, query string, limit int) ([]entity.Tag, error)
	Create(ctx context.Context, tag *entity.Tag) error
	Update(ctx context.Context, tag *entity.Tag) error
	Delete(ctx context.Context, id uint) error
}

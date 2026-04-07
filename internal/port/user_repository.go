package port

import (
	"context"

	"github.com/anthropics/quillow/internal/entity"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	List(ctx context.Context, limit, offset int) ([]entity.User, int64, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error
	GetRole(ctx context.Context, userID uint) (string, error)
}

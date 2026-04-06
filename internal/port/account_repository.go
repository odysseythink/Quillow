package port

import (
	"context"

	"github.com/anthropics/firefly-iii-go/internal/entity"
)

type AccountRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.Account, error)
	List(ctx context.Context, userGroupID uint, accountTypes []string, limit, offset int) ([]entity.Account, int64, error)
	Search(ctx context.Context, userGroupID uint, query string, accountTypes []string, limit int) ([]entity.Account, error)
	Create(ctx context.Context, account *entity.Account) error
	Update(ctx context.Context, account *entity.Account) error
	Delete(ctx context.Context, id uint) error
	GetAccountType(ctx context.Context, accountID uint) (*entity.AccountType, error)
	GetAccountTypeByType(ctx context.Context, typeName string) (*entity.AccountType, error)
	GetMeta(ctx context.Context, accountID uint, name string) (string, error)
	SetMeta(ctx context.Context, accountID uint, name, value string) error
	GetNotes(ctx context.Context, accountID uint) (string, error)
	SetNotes(ctx context.Context, accountID uint, text string) error
}

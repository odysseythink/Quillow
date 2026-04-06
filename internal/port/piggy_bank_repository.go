package port

import (
	"context"

	"github.com/anthropics/firefly-iii-go/internal/entity"
)

type PiggyBankRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.PiggyBank, error)
	List(ctx context.Context, limit, offset int) ([]entity.PiggyBank, int64, error)
	ListByAccount(ctx context.Context, accountID uint) ([]entity.PiggyBank, error)
	Create(ctx context.Context, piggy *entity.PiggyBank) error
	Update(ctx context.Context, piggy *entity.PiggyBank) error
	Delete(ctx context.Context, id uint) error
	ListEvents(ctx context.Context, piggyBankID uint) ([]entity.PiggyBankEvent, error)
	CreateEvent(ctx context.Context, event *entity.PiggyBankEvent) error
	GetNotes(ctx context.Context, piggyBankID uint) (string, error)
	SetNotes(ctx context.Context, piggyBankID uint, text string) error
}

type ObjectGroupRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.ObjectGroup, error)
	List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.ObjectGroup, int64, error)
	Create(ctx context.Context, og *entity.ObjectGroup) error
	Update(ctx context.Context, og *entity.ObjectGroup) error
	Delete(ctx context.Context, id uint) error
}

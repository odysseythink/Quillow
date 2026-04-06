package account

import (
	"context"
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/internal/port"
)

type UseCase struct {
	repo     port.AccountRepository
	currRepo port.CurrencyRepository
}

func NewUseCase(repo port.AccountRepository, currRepo port.CurrencyRepository) *UseCase {
	return &UseCase{repo: repo, currRepo: currRepo}
}

func (uc *UseCase) GetByID(ctx context.Context, id uint) (*entity.Account, *entity.AccountType, error) {
	acct, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	at, err := uc.repo.GetAccountType(ctx, id)
	if err != nil {
		return acct, nil, nil
	}
	return acct, at, nil
}

func (uc *UseCase) List(ctx context.Context, userGroupID uint, accountTypes []string, limit, offset int) ([]entity.Account, int64, error) {
	return uc.repo.List(ctx, userGroupID, accountTypes, limit, offset)
}

func (uc *UseCase) Search(ctx context.Context, userGroupID uint, query string, accountTypes []string, limit int) ([]entity.Account, error) {
	return uc.repo.Search(ctx, userGroupID, query, accountTypes, limit)
}

type CreateAccountInput struct {
	UserID         uint
	UserGroupID    uint
	Name           string
	Type           string
	IBAN           string
	AccountNumber  string
	VirtualBalance string
	Active         bool
	Order          int
	AccountRole    string
	CurrencyID     uint
	Notes          string
}

func (uc *UseCase) Create(ctx context.Context, input CreateAccountInput) (*entity.Account, error) {
	at, err := uc.repo.GetAccountTypeByType(ctx, input.Type)
	if err != nil {
		return nil, err
	}

	acct := &entity.Account{
		UserID:         input.UserID,
		UserGroupID:    input.UserGroupID,
		AccountTypeID:  at.ID,
		Name:           input.Name,
		IBAN:           input.IBAN,
		VirtualBalance: input.VirtualBalance,
		Active:         input.Active,
		Order:          input.Order,
	}

	if err := uc.repo.Create(ctx, acct); err != nil {
		return nil, err
	}

	// Store meta
	if input.AccountRole != "" {
		_ = uc.repo.SetMeta(ctx, acct.ID, "account_role", input.AccountRole)
	}
	if input.AccountNumber != "" {
		_ = uc.repo.SetMeta(ctx, acct.ID, "account_number", input.AccountNumber)
	}
	if input.CurrencyID > 0 {
		_ = uc.repo.SetMeta(ctx, acct.ID, "currency_id", fmt.Sprintf("%d", input.CurrencyID))
	}
	if input.Notes != "" {
		_ = uc.repo.SetNotes(ctx, acct.ID, input.Notes)
	}

	return acct, nil
}

type UpdateAccountInput struct {
	Name           string
	IBAN           string
	VirtualBalance string
	Active         bool
	Order          int
	AccountRole    string
	AccountNumber  string
	Notes          string
}

func (uc *UseCase) Update(ctx context.Context, id uint, input UpdateAccountInput) (*entity.Account, error) {
	acct, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	acct.Name = input.Name
	acct.IBAN = input.IBAN
	acct.VirtualBalance = input.VirtualBalance
	acct.Active = input.Active
	acct.Order = input.Order

	if err := uc.repo.Update(ctx, acct); err != nil {
		return nil, err
	}

	if input.AccountRole != "" {
		_ = uc.repo.SetMeta(ctx, id, "account_role", input.AccountRole)
	}
	if input.AccountNumber != "" {
		_ = uc.repo.SetMeta(ctx, id, "account_number", input.AccountNumber)
	}
	if input.Notes != "" {
		_ = uc.repo.SetNotes(ctx, id, input.Notes)
	}

	return acct, nil
}

func (uc *UseCase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *UseCase) GetMeta(ctx context.Context, accountID uint, key string) (string, error) {
	return uc.repo.GetMeta(ctx, accountID, key)
}

func (uc *UseCase) GetNotes(ctx context.Context, accountID uint) (string, error) {
	return uc.repo.GetNotes(ctx, accountID)
}

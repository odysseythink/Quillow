package transaction

import (
	"context"
	"time"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/internal/port"
)

type UseCase struct {
	repo     port.TransactionRepository
	acctRepo port.AccountRepository
}

func NewUseCase(repo port.TransactionRepository, acctRepo port.AccountRepository) *UseCase {
	return &UseCase{repo: repo, acctRepo: acctRepo}
}

func (uc *UseCase) GetGroupByID(ctx context.Context, id uint) (*entity.TransactionGroup, []entity.TransactionJournal, error) {
	group, err := uc.repo.FindGroupByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	journals, err := uc.repo.ListJournalsByGroupID(ctx, id)
	if err != nil {
		return group, nil, err
	}
	return group, journals, nil
}

func (uc *UseCase) ListGroups(ctx context.Context, userGroupID uint, transactionType string, start, end *time.Time, limit, offset int) ([]entity.TransactionGroup, int64, error) {
	return uc.repo.ListGroups(ctx, userGroupID, transactionType, start, end, limit, offset)
}

type TransactionInput struct {
	Type              string
	Description       string
	Date              time.Time
	Amount            string
	ForeignAmount     string
	CurrencyID        uint
	ForeignCurrencyID *uint
	SourceID          uint
	SourceName        string
	DestinationID     uint
	DestinationName   string
	BudgetID          *uint
	CategoryID        *uint
	BillID            *uint
	Tags              []string
	Notes             string
	Reconciled        bool
	ExternalID        string
	ExternalURL       string
	InternalRef       string
}

type CreateGroupInput struct {
	UserID       uint
	UserGroupID  uint
	GroupTitle    string
	Transactions []TransactionInput
}

func (uc *UseCase) CreateGroup(ctx context.Context, input CreateGroupInput) (*entity.TransactionGroup, error) {
	group := &entity.TransactionGroup{
		UserID:      input.UserID,
		UserGroupID: input.UserGroupID,
		Title:       input.GroupTitle,
	}
	if err := uc.repo.CreateGroup(ctx, group); err != nil {
		return nil, err
	}

	for i, txInput := range input.Transactions {
		txType, err := uc.repo.GetTransactionType(ctx, txInput.Type)
		if err != nil {
			return nil, err
		}

		journal := &entity.TransactionJournal{
			UserID:                input.UserID,
			UserGroupID:           input.UserGroupID,
			TransactionTypeID:     txType.ID,
			TransactionCurrencyID: txInput.CurrencyID,
			BillID:                txInput.BillID,
			Description:           txInput.Description,
			Date:                  txInput.Date,
			Order:                 uint(i),
			Completed:             true,
			TransactionGroupID:    group.ID,
		}
		if err := uc.repo.CreateJournal(ctx, journal); err != nil {
			return nil, err
		}

		// Source transaction (negative amount)
		source := &entity.Transaction{
			TransactionJournalID:  journal.ID,
			AccountID:             txInput.SourceID,
			TransactionCurrencyID: txInput.CurrencyID,
			ForeignCurrencyID:     txInput.ForeignCurrencyID,
			Amount:                "-" + txInput.Amount,
			ForeignAmount:         txInput.ForeignAmount,
			Reconciled:            txInput.Reconciled,
		}
		if err := uc.repo.CreateTransaction(ctx, source); err != nil {
			return nil, err
		}

		// Destination transaction (positive amount)
		dest := &entity.Transaction{
			TransactionJournalID:  journal.ID,
			AccountID:             txInput.DestinationID,
			TransactionCurrencyID: txInput.CurrencyID,
			ForeignCurrencyID:     txInput.ForeignCurrencyID,
			Amount:                txInput.Amount,
			ForeignAmount:         txInput.ForeignAmount,
			Reconciled:            txInput.Reconciled,
		}
		if err := uc.repo.CreateTransaction(ctx, dest); err != nil {
			return nil, err
		}

		// Set meta
		if txInput.Notes != "" {
			uc.repo.SetJournalNotes(ctx, journal.ID, txInput.Notes)
		}
		if txInput.ExternalID != "" {
			uc.repo.SetJournalMeta(ctx, journal.ID, "external_id", txInput.ExternalID)
		}
		if txInput.ExternalURL != "" {
			uc.repo.SetJournalMeta(ctx, journal.ID, "external_url", txInput.ExternalURL)
		}
		if txInput.InternalRef != "" {
			uc.repo.SetJournalMeta(ctx, journal.ID, "internal_reference", txInput.InternalRef)
		}
		if len(txInput.Tags) > 0 {
			uc.repo.SetJournalTags(ctx, journal.ID, txInput.Tags)
		}
	}

	return group, nil
}

func (uc *UseCase) DeleteGroup(ctx context.Context, id uint) error {
	journals, err := uc.repo.ListJournalsByGroupID(ctx, id)
	if err != nil {
		return err
	}
	for _, j := range journals {
		uc.repo.DeleteJournal(ctx, j.ID)
	}
	return uc.repo.DeleteGroup(ctx, id)
}

func (uc *UseCase) GetJournalByID(ctx context.Context, id uint) (*entity.TransactionJournal, error) {
	return uc.repo.FindJournalByID(ctx, id)
}

func (uc *UseCase) DeleteJournal(ctx context.Context, id uint) error {
	return uc.repo.DeleteJournal(ctx, id)
}

func (uc *UseCase) GetTransactions(ctx context.Context, journalID uint) ([]entity.Transaction, error) {
	return uc.repo.ListByJournalID(ctx, journalID)
}

func (uc *UseCase) SearchGroups(ctx context.Context, userGroupID uint, query string, limit, offset int) ([]entity.TransactionGroup, int64, error) {
	return uc.repo.SearchGroups(ctx, userGroupID, query, limit, offset)
}

func (uc *UseCase) CountByQuery(ctx context.Context, userGroupID uint, query string) (int64, error) {
	return uc.repo.CountByQuery(ctx, userGroupID, query)
}

func (uc *UseCase) GetJournalMeta(ctx context.Context, journalID uint, key string) (string, error) {
	return uc.repo.GetJournalMeta(ctx, journalID, key)
}

func (uc *UseCase) GetJournalNotes(ctx context.Context, journalID uint) (string, error) {
	return uc.repo.GetJournalNotes(ctx, journalID)
}

func (uc *UseCase) GetJournalTags(ctx context.Context, journalID uint) ([]string, error) {
	return uc.repo.GetJournalTags(ctx, journalID)
}

// Insight operations
func (uc *UseCase) InsightByType(ctx context.Context, userGroupID uint, transactionType string, start, end time.Time) ([]port.InsightEntry, error) {
	return uc.repo.SumByType(ctx, userGroupID, transactionType, start, end)
}

func (uc *UseCase) InsightByAccount(ctx context.Context, userGroupID uint, accountID uint, transactionType string, start, end time.Time) ([]port.InsightEntry, error) {
	return uc.repo.SumByAccount(ctx, userGroupID, accountID, transactionType, start, end)
}

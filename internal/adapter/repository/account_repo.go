package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthropics/quillow/internal/adapter/repository/model"
	"github.com/anthropics/quillow/internal/entity"
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) FindByID(ctx context.Context, id uint) (*entity.Account, error) {
	var m model.AccountModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}
	return accountModelToEntity(&m), nil
}

func (r *AccountRepository) List(ctx context.Context, userGroupID uint, accountTypes []string, limit, offset int) ([]entity.Account, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.AccountModel{}).Where("deleted_at IS NULL")
	if userGroupID > 0 {
		query = query.Where("user_group_id = ?", userGroupID)
	}
	if len(accountTypes) > 0 {
		var typeIDs []uint
		r.db.WithContext(ctx).Model(&model.AccountTypeModel{}).
			Where("type IN ?", accountTypes).Pluck("id", &typeIDs)
		if len(typeIDs) > 0 {
			query = query.Where("account_type_id IN ?", typeIDs)
		}
	}

	var total int64
	query.Count(&total)

	var models []model.AccountModel
	if err := query.Order("`order` ASC, name ASC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	accounts := make([]entity.Account, len(models))
	for i, m := range models {
		accounts[i] = *accountModelToEntity(&m)
	}
	return accounts, total, nil
}

func (r *AccountRepository) Search(ctx context.Context, userGroupID uint, query string, accountTypes []string, limit int) ([]entity.Account, error) {
	q := r.db.WithContext(ctx).Model(&model.AccountModel{}).Where("deleted_at IS NULL")
	if userGroupID > 0 {
		q = q.Where("user_group_id = ?", userGroupID)
	}
	if query != "" {
		q = q.Where("name LIKE ?", "%"+query+"%")
	}
	if len(accountTypes) > 0 {
		var typeIDs []uint
		r.db.WithContext(ctx).Model(&model.AccountTypeModel{}).
			Where("type IN ?", accountTypes).Pluck("id", &typeIDs)
		if len(typeIDs) > 0 {
			q = q.Where("account_type_id IN ?", typeIDs)
		}
	}

	var models []model.AccountModel
	if err := q.Limit(limit).Find(&models).Error; err != nil {
		return nil, err
	}

	accounts := make([]entity.Account, len(models))
	for i, m := range models {
		accounts[i] = *accountModelToEntity(&m)
	}
	return accounts, nil
}

func (r *AccountRepository) Create(ctx context.Context, account *entity.Account) error {
	m := accountEntityToModel(account)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	account.ID = m.ID
	account.CreatedAt = m.CreatedAt
	account.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *AccountRepository) Update(ctx context.Context, account *entity.Account) error {
	m := accountEntityToModel(account)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *AccountRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.AccountModel{}, id).Error
}

func (r *AccountRepository) GetAccountType(ctx context.Context, accountID uint) (*entity.AccountType, error) {
	var account model.AccountModel
	if err := r.db.WithContext(ctx).First(&account, accountID).Error; err != nil {
		return nil, err
	}
	var at model.AccountTypeModel
	if err := r.db.WithContext(ctx).First(&at, account.AccountTypeID).Error; err != nil {
		return nil, err
	}
	return &entity.AccountType{ID: at.ID, Type: at.Type}, nil
}

func (r *AccountRepository) GetAccountTypeByType(ctx context.Context, typeName string) (*entity.AccountType, error) {
	var at model.AccountTypeModel
	if err := r.db.WithContext(ctx).Where("type = ?", typeName).First(&at).Error; err != nil {
		return nil, fmt.Errorf("account type not found: %s", typeName)
	}
	return &entity.AccountType{ID: at.ID, Type: at.Type}, nil
}

func (r *AccountRepository) GetMeta(ctx context.Context, accountID uint, name string) (string, error) {
	var m model.AccountMetaModel
	if err := r.db.WithContext(ctx).Where("account_id = ? AND name = ?", accountID, name).First(&m).Error; err != nil {
		return "", nil
	}
	return m.Data, nil
}

func (r *AccountRepository) SetMeta(ctx context.Context, accountID uint, name, value string) error {
	var existing model.AccountMetaModel
	err := r.db.WithContext(ctx).Where("account_id = ? AND name = ?", accountID, name).First(&existing).Error
	if err != nil {
		return r.db.WithContext(ctx).Create(&model.AccountMetaModel{
			AccountID: accountID,
			Name:      name,
			Data:      value,
		}).Error
	}
	existing.Data = value
	return r.db.WithContext(ctx).Save(&existing).Error
}

func (r *AccountRepository) GetNotes(ctx context.Context, accountID uint) (string, error) {
	var note model.NoteModel
	if err := r.db.WithContext(ctx).Where("noteable_id = ? AND noteable_type LIKE ? AND deleted_at IS NULL", accountID, "%Account%").First(&note).Error; err != nil {
		return "", nil
	}
	if note.Text != nil {
		return *note.Text, nil
	}
	return "", nil
}

func (r *AccountRepository) SetNotes(ctx context.Context, accountID uint, text string) error {
	var existing model.NoteModel
	noteableType := "Quillow\\Models\\Account"
	err := r.db.WithContext(ctx).Where("noteable_id = ? AND noteable_type = ?", accountID, noteableType).First(&existing).Error
	if err != nil {
		return r.db.WithContext(ctx).Create(&model.NoteModel{
			NoteableID:   accountID,
			NoteableType: noteableType,
			Text:         &text,
		}).Error
	}
	existing.Text = &text
	return r.db.WithContext(ctx).Save(&existing).Error
}

func accountModelToEntity(m *model.AccountModel) *entity.Account {
	a := &entity.Account{
		ID:            m.ID,
		UserID:        m.UserID,
		UserGroupID:   m.UserGroupID,
		AccountTypeID: m.AccountTypeID,
		Name:          m.Name,
		Active:        m.Active,
		Encrypted:     m.Encrypted,
		Order:         m.Order,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		DeletedAt:     m.DeletedAt,
	}
	if m.VirtualBalance != nil {
		a.VirtualBalance = *m.VirtualBalance
	}
	if m.IBAN != nil {
		a.IBAN = *m.IBAN
	}
	return a
}

func accountEntityToModel(a *entity.Account) *model.AccountModel {
	m := &model.AccountModel{
		ID:            a.ID,
		UserID:        a.UserID,
		UserGroupID:   a.UserGroupID,
		AccountTypeID: a.AccountTypeID,
		Name:          a.Name,
		Active:        a.Active,
		Encrypted:     a.Encrypted,
		Order:         a.Order,
	}
	if a.VirtualBalance != "" {
		m.VirtualBalance = &a.VirtualBalance
	}
	if a.IBAN != "" {
		iban := strings.ToUpper(a.IBAN)
		m.IBAN = &iban
	}
	return m
}

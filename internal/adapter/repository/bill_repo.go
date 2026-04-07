package repository

import (
	"context"
	"fmt"

	"github.com/anthropics/quillow/internal/adapter/repository/model"
	"github.com/anthropics/quillow/internal/entity"
	"gorm.io/gorm"
)

type BillRepository struct {
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) *BillRepository {
	return &BillRepository{db: db}
}

func (r *BillRepository) FindByID(ctx context.Context, id uint) (*entity.Bill, error) {
	var m model.BillModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("bill not found: %w", err)
	}
	return billModelToEntity(&m), nil
}

func (r *BillRepository) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Bill, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.BillModel{}).Where("deleted_at IS NULL")
	if userGroupID > 0 {
		query = query.Where("user_group_id = ?", userGroupID)
	}

	var total int64
	query.Count(&total)

	var models []model.BillModel
	if err := query.Order("`order` ASC, name ASC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	bills := make([]entity.Bill, len(models))
	for i, m := range models {
		bills[i] = *billModelToEntity(&m)
	}
	return bills, total, nil
}

func (r *BillRepository) ListActive(ctx context.Context, userGroupID uint) ([]entity.Bill, error) {
	query := r.db.WithContext(ctx).Model(&model.BillModel{}).Where("deleted_at IS NULL AND active = ?", true)
	if userGroupID > 0 {
		query = query.Where("user_group_id = ?", userGroupID)
	}

	var models []model.BillModel
	if err := query.Order("`order` ASC, name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	bills := make([]entity.Bill, len(models))
	for i, m := range models {
		bills[i] = *billModelToEntity(&m)
	}
	return bills, nil
}

func (r *BillRepository) Create(ctx context.Context, bill *entity.Bill) error {
	m := billEntityToModel(bill)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	bill.ID = m.ID
	bill.CreatedAt = m.CreatedAt
	bill.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *BillRepository) Update(ctx context.Context, bill *entity.Bill) error {
	m := billEntityToModel(bill)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *BillRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.BillModel{}, id).Error
}

func (r *BillRepository) GetNotes(ctx context.Context, billID uint) (string, error) {
	var note model.NoteModel
	if err := r.db.WithContext(ctx).Where("noteable_id = ? AND noteable_type LIKE ? AND deleted_at IS NULL", billID, "%Bill%").First(&note).Error; err != nil {
		return "", nil
	}
	if note.Text != nil {
		return *note.Text, nil
	}
	return "", nil
}

func (r *BillRepository) SetNotes(ctx context.Context, billID uint, text string) error {
	var existing model.NoteModel
	noteableType := "Quillow\\Models\\Bill"
	err := r.db.WithContext(ctx).Where("noteable_id = ? AND noteable_type = ?", billID, noteableType).First(&existing).Error
	if err != nil {
		return r.db.WithContext(ctx).Create(&model.NoteModel{
			NoteableID:   billID,
			NoteableType: noteableType,
			Text:         &text,
		}).Error
	}
	existing.Text = &text
	return r.db.WithContext(ctx).Save(&existing).Error
}

// ---------------------------------------------------------------------------
// Conversion helpers
// ---------------------------------------------------------------------------

func billModelToEntity(m *model.BillModel) *entity.Bill {
	return &entity.Bill{
		ID:                    m.ID,
		UserID:                m.UserID,
		UserGroupID:           m.UserGroupID,
		TransactionCurrencyID: m.TransactionCurrencyID,
		Name:                  m.Name,
		AmountMin:             m.AmountMin,
		AmountMax:             m.AmountMax,
		Date:                  m.Date,
		EndDate:               m.EndDate,
		ExtensionDate:         m.ExtensionDate,
		RepeatFreq:            m.RepeatFreq,
		Skip:                  m.Skip,
		Automatch:             m.Automatch,
		Active:                m.Active,
		NameEncrypted:         m.NameEncrypted,
		MatchEncrypted:        m.MatchEncrypted,
		Order:                 m.Order,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
		DeletedAt:             m.DeletedAt,
	}
}

func billEntityToModel(b *entity.Bill) *model.BillModel {
	return &model.BillModel{
		ID:                    b.ID,
		UserID:                b.UserID,
		UserGroupID:           b.UserGroupID,
		TransactionCurrencyID: b.TransactionCurrencyID,
		Name:                  b.Name,
		AmountMin:             b.AmountMin,
		AmountMax:             b.AmountMax,
		Date:                  b.Date,
		EndDate:               b.EndDate,
		ExtensionDate:         b.ExtensionDate,
		RepeatFreq:            b.RepeatFreq,
		Skip:                  b.Skip,
		Automatch:             b.Automatch,
		Active:                b.Active,
		NameEncrypted:         b.NameEncrypted,
		MatchEncrypted:        b.MatchEncrypted,
		Order:                 b.Order,
	}
}

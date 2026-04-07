package repository

import (
	"context"
	"fmt"

	"github.com/anthropics/quillow/internal/adapter/repository/model"
	"github.com/anthropics/quillow/internal/entity"
	"gorm.io/gorm"
)

type PiggyBankRepository struct {
	db *gorm.DB
}

func NewPiggyBankRepository(db *gorm.DB) *PiggyBankRepository {
	return &PiggyBankRepository{db: db}
}

func (r *PiggyBankRepository) FindByID(ctx context.Context, id uint) (*entity.PiggyBank, error) {
	var m model.PiggyBankModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("piggy bank not found: %w", err)
	}
	return piggyBankModelToEntity(&m), nil
}

func (r *PiggyBankRepository) List(ctx context.Context, limit, offset int) ([]entity.PiggyBank, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&model.PiggyBankModel{}).Where("deleted_at IS NULL").Count(&total)

	var models []model.PiggyBankModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").
		Order("`order` ASC").Limit(limit).Offset(offset).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	items := make([]entity.PiggyBank, len(models))
	for i, m := range models {
		items[i] = *piggyBankModelToEntity(&m)
	}
	return items, total, nil
}

func (r *PiggyBankRepository) ListByAccount(ctx context.Context, accountID uint) ([]entity.PiggyBank, error) {
	var models []model.PiggyBankModel
	if err := r.db.WithContext(ctx).
		Where("account_id = ? AND deleted_at IS NULL", accountID).
		Order("`order` ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]entity.PiggyBank, len(models))
	for i, m := range models {
		items[i] = *piggyBankModelToEntity(&m)
	}
	return items, nil
}

func (r *PiggyBankRepository) Create(ctx context.Context, piggy *entity.PiggyBank) error {
	m := piggyBankEntityToModel(piggy)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	piggy.ID = m.ID
	piggy.CreatedAt = m.CreatedAt
	piggy.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *PiggyBankRepository) Update(ctx context.Context, piggy *entity.PiggyBank) error {
	m := piggyBankEntityToModel(piggy)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *PiggyBankRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.PiggyBankModel{}, id).Error
}

func (r *PiggyBankRepository) ListEvents(ctx context.Context, piggyBankID uint) ([]entity.PiggyBankEvent, error) {
	var models []model.PiggyBankEventModel
	if err := r.db.WithContext(ctx).
		Where("piggy_bank_id = ?", piggyBankID).
		Order("date DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]entity.PiggyBankEvent, len(models))
	for i, m := range models {
		items[i] = piggyBankEventModelToEntity(&m)
	}
	return items, nil
}

func (r *PiggyBankRepository) CreateEvent(ctx context.Context, event *entity.PiggyBankEvent) error {
	m := piggyBankEventEntityToModel(event)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	event.ID = m.ID
	event.CreatedAt = m.CreatedAt
	event.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *PiggyBankRepository) GetNotes(ctx context.Context, piggyBankID uint) (string, error) {
	var note model.NoteModel
	if err := r.db.WithContext(ctx).Where("noteable_id = ? AND noteable_type LIKE ? AND deleted_at IS NULL", piggyBankID, "%PiggyBank%").First(&note).Error; err != nil {
		return "", nil
	}
	if note.Text != nil {
		return *note.Text, nil
	}
	return "", nil
}

func (r *PiggyBankRepository) SetNotes(ctx context.Context, piggyBankID uint, text string) error {
	var existing model.NoteModel
	noteableType := "Quillow\\Models\\PiggyBank"
	err := r.db.WithContext(ctx).Where("noteable_id = ? AND noteable_type = ?", piggyBankID, noteableType).First(&existing).Error
	if err != nil {
		return r.db.WithContext(ctx).Create(&model.NoteModel{
			NoteableID:   piggyBankID,
			NoteableType: noteableType,
			Text:         &text,
		}).Error
	}
	existing.Text = &text
	return r.db.WithContext(ctx).Save(&existing).Error
}

func piggyBankModelToEntity(m *model.PiggyBankModel) *entity.PiggyBank {
	p := &entity.PiggyBank{
		ID:        m.ID,
		AccountID: m.AccountID,
		Name:      m.Name,
		StartDate: m.StartDate,
		TargetDate: m.TargetDate,
		Order:     m.Order,
		Active:    m.Active,
		Encrypted: m.Encrypted,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
	if m.TargetAmount != nil {
		p.TargetAmount = *m.TargetAmount
	}
	return p
}

func piggyBankEntityToModel(p *entity.PiggyBank) *model.PiggyBankModel {
	m := &model.PiggyBankModel{
		ID:         p.ID,
		AccountID:  p.AccountID,
		Name:       p.Name,
		StartDate:  p.StartDate,
		TargetDate: p.TargetDate,
		Order:      p.Order,
		Active:     p.Active,
		Encrypted:  p.Encrypted,
	}
	if p.TargetAmount != "" {
		m.TargetAmount = &p.TargetAmount
	}
	return m
}

func piggyBankEventModelToEntity(m *model.PiggyBankEventModel) entity.PiggyBankEvent {
	return entity.PiggyBankEvent{
		ID:                   m.ID,
		PiggyBankID:          m.PiggyBankID,
		TransactionJournalID: m.TransactionJournalID,
		Amount:               m.Amount,
		Date:                 m.Date,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func piggyBankEventEntityToModel(e *entity.PiggyBankEvent) *model.PiggyBankEventModel {
	return &model.PiggyBankEventModel{
		ID:                   e.ID,
		PiggyBankID:          e.PiggyBankID,
		TransactionJournalID: e.TransactionJournalID,
		Amount:               e.Amount,
		Date:                 e.Date,
	}
}

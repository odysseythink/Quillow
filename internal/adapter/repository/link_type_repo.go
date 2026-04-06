package repository

import (
	"context"
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/adapter/repository/model"
	"github.com/anthropics/firefly-iii-go/internal/entity"
	"gorm.io/gorm"
)

type LinkTypeRepository struct {
	db *gorm.DB
}

func NewLinkTypeRepository(db *gorm.DB) *LinkTypeRepository {
	return &LinkTypeRepository{db: db}
}

func (r *LinkTypeRepository) FindByID(ctx context.Context, id uint) (*entity.LinkType, error) {
	var m model.LinkTypeModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("link type not found: %w", err)
	}
	return linkTypeModelToEntity(&m), nil
}

func (r *LinkTypeRepository) List(ctx context.Context, limit, offset int) ([]entity.LinkType, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&model.LinkTypeModel{}).Where("deleted_at IS NULL").Count(&total)

	var models []model.LinkTypeModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	items := make([]entity.LinkType, len(models))
	for i, m := range models {
		items[i] = *linkTypeModelToEntity(&m)
	}
	return items, total, nil
}

func (r *LinkTypeRepository) Create(ctx context.Context, lt *entity.LinkType) error {
	m := linkTypeEntityToModel(lt)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	lt.ID = m.ID
	return nil
}

func (r *LinkTypeRepository) Update(ctx context.Context, lt *entity.LinkType) error {
	m := linkTypeEntityToModel(lt)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *LinkTypeRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.LinkTypeModel{}, id).Error
}

func linkTypeModelToEntity(m *model.LinkTypeModel) *entity.LinkType {
	return &entity.LinkType{
		ID:        m.ID,
		Name:      m.Name,
		Outward:   m.Outward,
		Inward:    m.Inward,
		Editable:  m.Editable,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}

func linkTypeEntityToModel(lt *entity.LinkType) *model.LinkTypeModel {
	return &model.LinkTypeModel{
		ID:       lt.ID,
		Name:     lt.Name,
		Outward:  lt.Outward,
		Inward:   lt.Inward,
		Editable: lt.Editable,
	}
}

// TransactionLink Repository

type TransactionLinkRepository struct {
	db *gorm.DB
}

func NewTransactionLinkRepository(db *gorm.DB) *TransactionLinkRepository {
	return &TransactionLinkRepository{db: db}
}

func (r *TransactionLinkRepository) FindByID(ctx context.Context, id uint) (*entity.TransactionJournalLink, error) {
	var m model.TransactionJournalLinkModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("transaction link not found: %w", err)
	}
	return txLinkModelToEntity(&m), nil
}

func (r *TransactionLinkRepository) List(ctx context.Context, limit, offset int) ([]entity.TransactionJournalLink, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&model.TransactionJournalLinkModel{}).Count(&total)

	var models []model.TransactionJournalLinkModel
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	items := make([]entity.TransactionJournalLink, len(models))
	for i, m := range models {
		items[i] = *txLinkModelToEntity(&m)
	}
	return items, total, nil
}

func (r *TransactionLinkRepository) ListByJournalID(ctx context.Context, journalID uint) ([]entity.TransactionJournalLink, error) {
	var models []model.TransactionJournalLinkModel
	if err := r.db.WithContext(ctx).
		Where("source_id = ? OR destination_id = ?", journalID, journalID).
		Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]entity.TransactionJournalLink, len(models))
	for i, m := range models {
		items[i] = *txLinkModelToEntity(&m)
	}
	return items, nil
}

func (r *TransactionLinkRepository) Create(ctx context.Context, link *entity.TransactionJournalLink) error {
	m := txLinkEntityToModel(link)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	link.ID = m.ID
	return nil
}

func (r *TransactionLinkRepository) Update(ctx context.Context, link *entity.TransactionJournalLink) error {
	m := txLinkEntityToModel(link)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *TransactionLinkRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.TransactionJournalLinkModel{}, id).Error
}

func txLinkModelToEntity(m *model.TransactionJournalLinkModel) *entity.TransactionJournalLink {
	link := &entity.TransactionJournalLink{
		ID:            m.ID,
		LinkTypeID:    m.LinkTypeID,
		SourceID:      m.SourceID,
		DestinationID: m.DestinationID,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
	if m.Comment != nil {
		link.Comment = *m.Comment
	}
	return link
}

func txLinkEntityToModel(l *entity.TransactionJournalLink) *model.TransactionJournalLinkModel {
	m := &model.TransactionJournalLinkModel{
		ID:            l.ID,
		LinkTypeID:    l.LinkTypeID,
		SourceID:      l.SourceID,
		DestinationID: l.DestinationID,
	}
	if l.Comment != "" {
		m.Comment = &l.Comment
	}
	return m
}

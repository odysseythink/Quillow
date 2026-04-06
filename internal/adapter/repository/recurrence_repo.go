package repository

import (
	"context"
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/adapter/repository/model"
	"github.com/anthropics/firefly-iii-go/internal/entity"
	"gorm.io/gorm"
)

type RecurrenceRepository struct {
	db *gorm.DB
}

func NewRecurrenceRepository(db *gorm.DB) *RecurrenceRepository {
	return &RecurrenceRepository{db: db}
}

func (r *RecurrenceRepository) FindByID(ctx context.Context, id uint) (*entity.Recurrence, error) {
	var m model.RecurrenceModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("recurrence not found: %w", err)
	}
	return recurrenceModelToEntity(&m), nil
}

func (r *RecurrenceRepository) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Recurrence, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.RecurrenceModel{}).Where("deleted_at IS NULL")
	if userGroupID > 0 {
		query = query.Where("user_group_id = ?", userGroupID)
	}

	var total int64
	query.Count(&total)

	var models []model.RecurrenceModel
	if err := query.Order("id ASC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]entity.Recurrence, len(models))
	for i, m := range models {
		result[i] = *recurrenceModelToEntity(&m)
	}
	return result, total, nil
}

func (r *RecurrenceRepository) Create(ctx context.Context, rec *entity.Recurrence) error {
	m := recurrenceEntityToModel(rec)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	rec.ID = m.ID
	rec.CreatedAt = m.CreatedAt
	rec.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *RecurrenceRepository) Update(ctx context.Context, rec *entity.Recurrence) error {
	m := recurrenceEntityToModel(rec)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *RecurrenceRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.RecurrenceModel{}, id).Error
}

func (r *RecurrenceRepository) GetRepetitions(ctx context.Context, recurrenceID uint) ([]entity.RecurrenceRepetition, error) {
	var models []model.RecurrenceRepetitionModel
	if err := r.db.WithContext(ctx).Where("recurrence_id = ?", recurrenceID).Find(&models).Error; err != nil {
		return nil, err
	}

	result := make([]entity.RecurrenceRepetition, len(models))
	for i, m := range models {
		result[i] = entity.RecurrenceRepetition{
			ID:               m.ID,
			RecurrenceID:     m.RecurrenceID,
			RepetitionType:   m.RepetitionType,
			RepetitionMoment: m.RepetitionMoment,
			RepetitionSkip:   m.RepetitionSkip,
			Weekend:          m.Weekend,
			CreatedAt:        m.CreatedAt,
			UpdatedAt:        m.UpdatedAt,
		}
	}
	return result, nil
}

func (r *RecurrenceRepository) GetTransactions(ctx context.Context, recurrenceID uint) ([]entity.RecurrenceTransaction, error) {
	var models []model.RecurrenceTransactionModel
	if err := r.db.WithContext(ctx).Where("recurrence_id = ?", recurrenceID).Find(&models).Error; err != nil {
		return nil, err
	}

	result := make([]entity.RecurrenceTransaction, len(models))
	for i, m := range models {
		fa := ""
		if m.ForeignAmount != nil {
			fa = *m.ForeignAmount
		}
		result[i] = entity.RecurrenceTransaction{
			ID:                    m.ID,
			RecurrenceID:          m.RecurrenceID,
			TransactionCurrencyID: m.TransactionCurrencyID,
			ForeignCurrencyID:     m.ForeignCurrencyID,
			SourceID:              m.SourceID,
			DestinationID:         m.DestinationID,
			Amount:                m.Amount,
			ForeignAmount:         fa,
			Description:           m.Description,
			CreatedAt:             m.CreatedAt,
			UpdatedAt:             m.UpdatedAt,
		}
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// Conversion helpers
// ---------------------------------------------------------------------------

func recurrenceModelToEntity(m *model.RecurrenceModel) *entity.Recurrence {
	return &entity.Recurrence{
		ID:                    m.ID,
		UserID:                m.UserID,
		UserGroupID:           m.UserGroupID,
		TransactionTypeID:     m.TransactionTypeID,
		TransactionCurrencyID: m.TransactionCurrencyID,
		Title:                 m.Title,
		Description:           m.Description,
		FirstDate:             m.FirstDate,
		RepeatUntil:           m.RepeatUntil,
		LatestDate:            m.LatestDate,
		Repetitions:           m.Repetitions,
		ApplyRules:            m.ApplyRules,
		Active:                m.Active,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
		DeletedAt:             m.DeletedAt,
	}
}

func recurrenceEntityToModel(e *entity.Recurrence) *model.RecurrenceModel {
	return &model.RecurrenceModel{
		ID:                    e.ID,
		UserID:                e.UserID,
		UserGroupID:           e.UserGroupID,
		TransactionTypeID:     e.TransactionTypeID,
		TransactionCurrencyID: e.TransactionCurrencyID,
		Title:                 e.Title,
		Description:           e.Description,
		FirstDate:             e.FirstDate,
		RepeatUntil:           e.RepeatUntil,
		LatestDate:            e.LatestDate,
		Repetitions:           e.Repetitions,
		ApplyRules:            e.ApplyRules,
		Active:                e.Active,
	}
}

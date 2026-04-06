package repository

import (
	"context"
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/adapter/repository/model"
	"github.com/anthropics/firefly-iii-go/internal/entity"
	"gorm.io/gorm"
)

type PreferenceRepository struct {
	db *gorm.DB
}

func NewPreferenceRepository(db *gorm.DB) *PreferenceRepository {
	return &PreferenceRepository{db: db}
}

func (r *PreferenceRepository) FindByUserAndName(ctx context.Context, userID uint, name string) (*entity.Preference, error) {
	var m model.PreferenceModel
	if err := r.db.WithContext(ctx).Where("user_id = ? AND name = ?", userID, name).First(&m).Error; err != nil {
		return nil, fmt.Errorf("preference not found: %w", err)
	}
	return prefModelToEntity(&m), nil
}

func (r *PreferenceRepository) ListByUser(ctx context.Context, userID uint, limit, offset int) ([]entity.Preference, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&model.PreferenceModel{}).Where("user_id = ?", userID).Count(&total)

	var models []model.PreferenceModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	prefs := make([]entity.Preference, len(models))
	for i, m := range models {
		prefs[i] = *prefModelToEntity(&m)
	}
	return prefs, total, nil
}

func (r *PreferenceRepository) Create(ctx context.Context, pref *entity.Preference) error {
	m := prefEntityToModel(pref)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	pref.ID = m.ID
	return nil
}

func (r *PreferenceRepository) Update(ctx context.Context, pref *entity.Preference) error {
	m := prefEntityToModel(pref)
	return r.db.WithContext(ctx).Save(m).Error
}

func prefModelToEntity(m *model.PreferenceModel) *entity.Preference {
	return &entity.Preference{
		ID:        m.ID,
		UserID:    m.UserID,
		Name:      m.Name,
		Data:      m.Data,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func prefEntityToModel(p *entity.Preference) *model.PreferenceModel {
	return &model.PreferenceModel{
		ID:     p.ID,
		UserID: p.UserID,
		Name:   p.Name,
		Data:   p.Data,
	}
}

package repository

import (
	"context"

	"github.com/anthropics/quillow/internal/adapter/repository/model"
	"gorm.io/gorm"
)

type ClassificationPatternRepository struct {
	db *gorm.DB
}

func NewClassificationPatternRepository(db *gorm.DB) *ClassificationPatternRepository {
	return &ClassificationPatternRepository{db: db}
}

func (r *ClassificationPatternRepository) FindByUser(ctx context.Context, userID uint) ([]model.ClassificationPatternModel, error) {
	var patterns []model.ClassificationPatternModel
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("hit_count DESC").Find(&patterns).Error
	return patterns, err
}

func (r *ClassificationPatternRepository) FindByUserAndPattern(ctx context.Context, userID uint, pattern string, categoryID uint) (*model.ClassificationPatternModel, error) {
	var m model.ClassificationPatternModel
	err := r.db.WithContext(ctx).Where("user_id = ? AND pattern = ? AND category_id = ?", userID, pattern, categoryID).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *ClassificationPatternRepository) Create(ctx context.Context, m *model.ClassificationPatternModel) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *ClassificationPatternRepository) IncrementHitCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.ClassificationPatternModel{}).Where("id = ?", id).
		UpdateColumn("hit_count", gorm.Expr("hit_count + 1")).Error
}

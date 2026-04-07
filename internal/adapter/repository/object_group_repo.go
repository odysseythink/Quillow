package repository

import (
	"context"
	"fmt"

	"github.com/anthropics/quillow/internal/adapter/repository/model"
	"github.com/anthropics/quillow/internal/entity"
	"gorm.io/gorm"
)

type ObjectGroupRepository struct {
	db *gorm.DB
}

func NewObjectGroupRepository(db *gorm.DB) *ObjectGroupRepository {
	return &ObjectGroupRepository{db: db}
}

func (r *ObjectGroupRepository) FindByID(ctx context.Context, id uint) (*entity.ObjectGroup, error) {
	var m model.ObjectGroupModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("object group not found: %w", err)
	}
	return objectGroupModelToEntity(&m), nil
}

func (r *ObjectGroupRepository) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.ObjectGroup, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.ObjectGroupModel{}).Where("deleted_at IS NULL")
	if userGroupID > 0 {
		query = query.Where("user_group_id = ?", userGroupID)
	}

	var total int64
	query.Count(&total)

	var models []model.ObjectGroupModel
	if err := query.Order("`order` ASC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	items := make([]entity.ObjectGroup, len(models))
	for i, m := range models {
		items[i] = *objectGroupModelToEntity(&m)
	}
	return items, total, nil
}

func (r *ObjectGroupRepository) Create(ctx context.Context, og *entity.ObjectGroup) error {
	m := objectGroupEntityToModel(og)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	og.ID = m.ID
	og.CreatedAt = m.CreatedAt
	og.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *ObjectGroupRepository) Update(ctx context.Context, og *entity.ObjectGroup) error {
	m := objectGroupEntityToModel(og)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *ObjectGroupRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.ObjectGroupModel{}, id).Error
}

func objectGroupModelToEntity(m *model.ObjectGroupModel) *entity.ObjectGroup {
	return &entity.ObjectGroup{
		ID:          m.ID,
		UserID:      m.UserID,
		UserGroupID: m.UserGroupID,
		Title:       m.Title,
		Order:       m.Order,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}
}

func objectGroupEntityToModel(og *entity.ObjectGroup) *model.ObjectGroupModel {
	return &model.ObjectGroupModel{
		ID:          og.ID,
		UserID:      og.UserID,
		UserGroupID: og.UserGroupID,
		Title:       og.Title,
		Order:       og.Order,
	}
}

package repository

import (
	"context"
	"fmt"

	"github.com/anthropics/quillow/internal/adapter/repository/model"
	"github.com/anthropics/quillow/internal/entity"
	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) FindByID(ctx context.Context, id uint) (*entity.Tag, error) {
	var m model.TagModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}
	return tagModelToEntity(&m), nil
}

func (r *TagRepository) FindByTag(ctx context.Context, userGroupID uint, tag string) (*entity.Tag, error) {
	var m model.TagModel
	if err := r.db.WithContext(ctx).Where("user_group_id = ? AND tag = ? AND deleted_at IS NULL", userGroupID, tag).First(&m).Error; err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}
	return tagModelToEntity(&m), nil
}

func (r *TagRepository) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Tag, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.TagModel{}).Where("deleted_at IS NULL")
	if userGroupID > 0 {
		query = query.Where("user_group_id = ?", userGroupID)
	}

	var total int64
	query.Count(&total)

	var models []model.TagModel
	if err := query.Order("tag ASC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	tags := make([]entity.Tag, len(models))
	for i, m := range models {
		tags[i] = *tagModelToEntity(&m)
	}
	return tags, total, nil
}

func (r *TagRepository) Search(ctx context.Context, userGroupID uint, query string, limit int) ([]entity.Tag, error) {
	q := r.db.WithContext(ctx).Model(&model.TagModel{}).Where("deleted_at IS NULL")
	if userGroupID > 0 {
		q = q.Where("user_group_id = ?", userGroupID)
	}
	if query != "" {
		q = q.Where("tag LIKE ?", "%"+query+"%")
	}

	var models []model.TagModel
	if err := q.Limit(limit).Find(&models).Error; err != nil {
		return nil, err
	}

	tags := make([]entity.Tag, len(models))
	for i, m := range models {
		tags[i] = *tagModelToEntity(&m)
	}
	return tags, nil
}

func (r *TagRepository) Create(ctx context.Context, tag *entity.Tag) error {
	m := tagEntityToModel(tag)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	tag.ID = m.ID
	tag.CreatedAt = m.CreatedAt
	tag.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *TagRepository) Update(ctx context.Context, tag *entity.Tag) error {
	m := tagEntityToModel(tag)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *TagRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.TagModel{}, id).Error
}

func tagModelToEntity(m *model.TagModel) *entity.Tag {
	t := &entity.Tag{
		ID:          m.ID,
		UserID:      m.UserID,
		UserGroupID: m.UserGroupID,
		Tag:         m.Tag,
		TagMode:     m.TagMode,
		Date:        m.Date,
		Latitude:    m.Latitude,
		Longitude:   m.Longitude,
		ZoomLevel:   m.ZoomLevel,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}
	if m.Description != nil {
		t.Description = *m.Description
	}
	return t
}

func tagEntityToModel(t *entity.Tag) *model.TagModel {
	m := &model.TagModel{
		ID:          t.ID,
		UserID:      t.UserID,
		UserGroupID: t.UserGroupID,
		Tag:         t.Tag,
		TagMode:     t.TagMode,
		Date:        t.Date,
		Latitude:    t.Latitude,
		Longitude:   t.Longitude,
		ZoomLevel:   t.ZoomLevel,
	}
	if t.Description != "" {
		m.Description = &t.Description
	}
	return m
}

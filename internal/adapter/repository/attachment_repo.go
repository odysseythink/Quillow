package repository

import (
	"context"
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/adapter/repository/model"
	"github.com/anthropics/firefly-iii-go/internal/entity"
	"gorm.io/gorm"
)

type AttachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{db: db}
}

func (r *AttachmentRepository) FindByID(ctx context.Context, id uint) (*entity.Attachment, error) {
	var m model.AttachmentModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("attachment not found: %w", err)
	}
	return attachModelToEntity(&m), nil
}

func (r *AttachmentRepository) List(ctx context.Context, limit, offset int) ([]entity.Attachment, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&model.AttachmentModel{}).Where("deleted_at IS NULL").Count(&total)

	var models []model.AttachmentModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	items := make([]entity.Attachment, len(models))
	for i, m := range models {
		items[i] = *attachModelToEntity(&m)
	}
	return items, total, nil
}

func (r *AttachmentRepository) ListByAttachable(ctx context.Context, attachableType string, attachableID uint) ([]entity.Attachment, error) {
	var models []model.AttachmentModel
	if err := r.db.WithContext(ctx).
		Where("attachable_type LIKE ? AND attachable_id = ? AND deleted_at IS NULL", "%"+attachableType+"%", attachableID).
		Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]entity.Attachment, len(models))
	for i, m := range models {
		items[i] = *attachModelToEntity(&m)
	}
	return items, nil
}

func (r *AttachmentRepository) Create(ctx context.Context, attachment *entity.Attachment) error {
	m := attachEntityToModel(attachment)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	attachment.ID = m.ID
	attachment.CreatedAt = m.CreatedAt
	attachment.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *AttachmentRepository) Update(ctx context.Context, attachment *entity.Attachment) error {
	m := attachEntityToModel(attachment)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *AttachmentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.AttachmentModel{}, id).Error
}

func attachModelToEntity(m *model.AttachmentModel) *entity.Attachment {
	a := &entity.Attachment{
		ID:             m.ID,
		UserID:         m.UserID,
		UserGroupID:    m.UserGroupID,
		AttachableID:   m.AttachableID,
		AttachableType: m.AttachableType,
		MD5:            m.MD5,
		Filename:       m.Filename,
		Mime:           m.Mime,
		Size:           m.Size,
		Uploaded:       m.Uploaded,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      m.DeletedAt,
	}
	if m.Title != nil {
		a.Title = *m.Title
	}
	if m.Description != nil {
		a.Description = *m.Description
	}
	return a
}

func attachEntityToModel(a *entity.Attachment) *model.AttachmentModel {
	m := &model.AttachmentModel{
		ID:             a.ID,
		UserID:         a.UserID,
		UserGroupID:    a.UserGroupID,
		AttachableID:   a.AttachableID,
		AttachableType: a.AttachableType,
		MD5:            a.MD5,
		Filename:       a.Filename,
		Mime:           a.Mime,
		Size:           a.Size,
		Uploaded:       a.Uploaded,
	}
	if a.Title != "" {
		m.Title = &a.Title
	}
	if a.Description != "" {
		m.Description = &a.Description
	}
	return m
}

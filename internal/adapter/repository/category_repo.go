package repository

import (
	"context"
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/adapter/repository/model"
	"github.com/anthropics/firefly-iii-go/internal/entity"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindByID(ctx context.Context, id uint) (*entity.Category, error) {
	var m model.CategoryModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	return categoryModelToEntity(&m), nil
}

func (r *CategoryRepository) FindByName(ctx context.Context, userGroupID uint, name string) (*entity.Category, error) {
	var m model.CategoryModel
	if err := r.db.WithContext(ctx).Where("user_group_id = ? AND name = ? AND deleted_at IS NULL", userGroupID, name).First(&m).Error; err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	return categoryModelToEntity(&m), nil
}

func (r *CategoryRepository) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Category, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.CategoryModel{}).Where("deleted_at IS NULL")
	if userGroupID > 0 {
		query = query.Where("user_group_id = ?", userGroupID)
	}

	var total int64
	query.Count(&total)

	var models []model.CategoryModel
	if err := query.Order("name ASC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	categories := make([]entity.Category, len(models))
	for i, m := range models {
		categories[i] = *categoryModelToEntity(&m)
	}
	return categories, total, nil
}

func (r *CategoryRepository) Search(ctx context.Context, userGroupID uint, query string, limit int) ([]entity.Category, error) {
	q := r.db.WithContext(ctx).Model(&model.CategoryModel{}).Where("deleted_at IS NULL")
	if userGroupID > 0 {
		q = q.Where("user_group_id = ?", userGroupID)
	}
	if query != "" {
		q = q.Where("name LIKE ?", "%"+query+"%")
	}

	var models []model.CategoryModel
	if err := q.Limit(limit).Find(&models).Error; err != nil {
		return nil, err
	}

	categories := make([]entity.Category, len(models))
	for i, m := range models {
		categories[i] = *categoryModelToEntity(&m)
	}
	return categories, nil
}

func (r *CategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	m := categoryEntityToModel(category)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	category.ID = m.ID
	category.CreatedAt = m.CreatedAt
	category.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	m := categoryEntityToModel(category)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *CategoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.CategoryModel{}, id).Error
}

func (r *CategoryRepository) GetNotes(ctx context.Context, categoryID uint) (string, error) {
	var note model.NoteModel
	if err := r.db.WithContext(ctx).Where("noteable_id = ? AND noteable_type LIKE ? AND deleted_at IS NULL", categoryID, "%Category%").First(&note).Error; err != nil {
		return "", nil
	}
	if note.Text != nil {
		return *note.Text, nil
	}
	return "", nil
}

func (r *CategoryRepository) SetNotes(ctx context.Context, categoryID uint, text string) error {
	var existing model.NoteModel
	noteableType := "FireflyIII\\Models\\Category"
	err := r.db.WithContext(ctx).Where("noteable_id = ? AND noteable_type = ?", categoryID, noteableType).First(&existing).Error
	if err != nil {
		return r.db.WithContext(ctx).Create(&model.NoteModel{
			NoteableID:   categoryID,
			NoteableType: noteableType,
			Text:         &text,
		}).Error
	}
	existing.Text = &text
	return r.db.WithContext(ctx).Save(&existing).Error
}

func categoryModelToEntity(m *model.CategoryModel) *entity.Category {
	return &entity.Category{
		ID:          m.ID,
		UserID:      m.UserID,
		UserGroupID: m.UserGroupID,
		Name:        m.Name,
		Encrypted:   m.Encrypted,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}
}

func categoryEntityToModel(c *entity.Category) *model.CategoryModel {
	return &model.CategoryModel{
		ID:          c.ID,
		UserID:      c.UserID,
		UserGroupID: c.UserGroupID,
		Name:        c.Name,
		Encrypted:   c.Encrypted,
	}
}

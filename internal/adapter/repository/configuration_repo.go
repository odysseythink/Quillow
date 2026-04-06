package repository

import (
	"context"
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/adapter/repository/model"
	"github.com/anthropics/firefly-iii-go/internal/entity"
	"gorm.io/gorm"
)

type ConfigurationRepository struct {
	db *gorm.DB
}

func NewConfigurationRepository(db *gorm.DB) *ConfigurationRepository {
	return &ConfigurationRepository{db: db}
}

func (r *ConfigurationRepository) FindByName(ctx context.Context, name string) (*entity.Configuration, error) {
	var m model.ConfigurationModel
	if err := r.db.WithContext(ctx).Where("name = ? AND deleted_at IS NULL", name).First(&m).Error; err != nil {
		return nil, fmt.Errorf("configuration not found: %w", err)
	}
	return configModelToEntity(&m), nil
}

func (r *ConfigurationRepository) List(ctx context.Context) ([]entity.Configuration, error) {
	var models []model.ConfigurationModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&models).Error; err != nil {
		return nil, err
	}
	configs := make([]entity.Configuration, len(models))
	for i, m := range models {
		configs[i] = *configModelToEntity(&m)
	}
	return configs, nil
}

func (r *ConfigurationRepository) Upsert(ctx context.Context, cfg *entity.Configuration) error {
	var existing model.ConfigurationModel
	err := r.db.WithContext(ctx).Where("name = ?", cfg.Name).First(&existing).Error
	if err != nil {
		m := configEntityToModel(cfg)
		if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
			return fmt.Errorf("failed to create config: %w", err)
		}
		cfg.ID = m.ID
		return nil
	}
	existing.Data = cfg.Data
	return r.db.WithContext(ctx).Save(&existing).Error
}

func configModelToEntity(m *model.ConfigurationModel) *entity.Configuration {
	return &entity.Configuration{
		ID:        m.ID,
		Name:      m.Name,
		Data:      m.Data,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}

func configEntityToModel(c *entity.Configuration) *model.ConfigurationModel {
	return &model.ConfigurationModel{
		ID:   c.ID,
		Name: c.Name,
		Data: c.Data,
	}
}

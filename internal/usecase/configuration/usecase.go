package configuration

import (
	"context"

	"github.com/anthropics/quillow/internal/entity"
	"github.com/anthropics/quillow/internal/port"
)

var editableKeys = map[string]bool{
	"single_user_mode":        true,
	"is_demo_site":            true,
	"permission_update_check": true,
	"last_update_check":       true,
	"latest_version":          true,
}

type ConfigItem struct {
	Title    string `json:"title"`
	Value    string `json:"value"`
	Editable bool   `json:"editable"`
}

type UseCase struct {
	repo port.ConfigurationRepository
}

func NewUseCase(repo port.ConfigurationRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) List(ctx context.Context) ([]ConfigItem, error) {
	configs, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]ConfigItem, len(configs))
	for i, c := range configs {
		items[i] = ConfigItem{
			Title:    c.Name,
			Value:    c.Data,
			Editable: editableKeys[c.Name],
		}
	}
	return items, nil
}

func (uc *UseCase) GetByName(ctx context.Context, name string) (*ConfigItem, error) {
	cfg, err := uc.repo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return &ConfigItem{
		Title:    cfg.Name,
		Value:    cfg.Data,
		Editable: editableKeys[cfg.Name],
	}, nil
}

func (uc *UseCase) Update(ctx context.Context, name, value string) (*ConfigItem, error) {
	if !editableKeys[name] {
		return nil, entity.ErrNotEditable
	}
	cfg := &entity.Configuration{Name: name, Data: value}
	if err := uc.repo.Upsert(ctx, cfg); err != nil {
		return nil, err
	}
	return &ConfigItem{
		Title:    name,
		Value:    value,
		Editable: true,
	}, nil
}

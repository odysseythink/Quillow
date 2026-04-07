package repository

import (
	"context"
	"testing"

	"github.com/anthropics/quillow/internal/adapter/repository/model"
	"github.com/anthropics/quillow/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupConfigDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.ConfigurationModel{}))
	return db
}

func TestConfigRepo_UpsertAndFind(t *testing.T) {
	db := setupConfigDB(t)
	repo := NewConfigurationRepository(db)
	ctx := context.Background()

	cfg := &entity.Configuration{Name: "single_user_mode", Data: "true"}
	require.NoError(t, repo.Upsert(ctx, cfg))

	found, err := repo.FindByName(ctx, "single_user_mode")
	require.NoError(t, err)
	assert.Equal(t, "true", found.Data)

	cfg.Data = "false"
	require.NoError(t, repo.Upsert(ctx, cfg))

	found, err = repo.FindByName(ctx, "single_user_mode")
	require.NoError(t, err)
	assert.Equal(t, "false", found.Data)
}

func TestConfigRepo_List(t *testing.T) {
	db := setupConfigDB(t)
	repo := NewConfigurationRepository(db)
	ctx := context.Background()

	require.NoError(t, repo.Upsert(ctx, &entity.Configuration{Name: "key1", Data: "\"val1\""}))
	require.NoError(t, repo.Upsert(ctx, &entity.Configuration{Name: "key2", Data: "\"val2\""}))

	configs, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, configs, 2)
}

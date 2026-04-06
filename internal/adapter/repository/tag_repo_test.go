package repository

import (
	"context"
	"testing"

	"github.com/anthropics/firefly-iii-go/internal/adapter/repository/model"
	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTagDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.TagModel{}))
	return db
}

func TestTagRepo_CreateAndFind(t *testing.T) {
	db := setupTagDB(t)
	repo := NewTagRepository(db)
	ctx := context.Background()

	tag := &entity.Tag{UserID: 1, UserGroupID: 1, Tag: "vacation", TagMode: "nothing"}
	require.NoError(t, repo.Create(ctx, tag))
	assert.NotZero(t, tag.ID)

	found, err := repo.FindByID(ctx, tag.ID)
	require.NoError(t, err)
	assert.Equal(t, "vacation", found.Tag)

	foundByTag, err := repo.FindByTag(ctx, 1, "vacation")
	require.NoError(t, err)
	assert.Equal(t, tag.ID, foundByTag.ID)
}

func TestTagRepo_List(t *testing.T) {
	db := setupTagDB(t)
	repo := NewTagRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &entity.Tag{UserID: 1, UserGroupID: 1, Tag: "alpha", TagMode: "nothing"})
	repo.Create(ctx, &entity.Tag{UserID: 1, UserGroupID: 1, Tag: "beta", TagMode: "nothing"})
	repo.Create(ctx, &entity.Tag{UserID: 1, UserGroupID: 1, Tag: "gamma", TagMode: "nothing"})

	tags, total, err := repo.List(ctx, 1, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, tags, 3)

	// Test pagination
	tags, total, err = repo.List(ctx, 1, 2, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, tags, 2)
}

func TestTagRepo_Search(t *testing.T) {
	db := setupTagDB(t)
	repo := NewTagRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &entity.Tag{UserID: 1, UserGroupID: 1, Tag: "vacation", TagMode: "nothing"})
	repo.Create(ctx, &entity.Tag{UserID: 1, UserGroupID: 1, Tag: "work", TagMode: "nothing"})
	repo.Create(ctx, &entity.Tag{UserID: 1, UserGroupID: 1, Tag: "vacation-2024", TagMode: "nothing"})

	results, err := repo.Search(ctx, 1, "vacat", 10)
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestTagRepo_Delete(t *testing.T) {
	db := setupTagDB(t)
	repo := NewTagRepository(db)
	ctx := context.Background()

	tag := &entity.Tag{UserID: 1, UserGroupID: 1, Tag: "del", TagMode: "nothing"}
	require.NoError(t, repo.Create(ctx, tag))
	require.NoError(t, repo.Delete(ctx, tag.ID))

	_, err := repo.FindByID(ctx, tag.ID)
	assert.Error(t, err)
}

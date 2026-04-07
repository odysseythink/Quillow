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

func setupCategoryDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.CategoryModel{}, &model.NoteModel{}))
	return db
}

func TestCategoryRepo_CreateAndFind(t *testing.T) {
	db := setupCategoryDB(t)
	repo := NewCategoryRepository(db)
	ctx := context.Background()

	cat := &entity.Category{UserID: 1, UserGroupID: 1, Name: "Groceries"}
	require.NoError(t, repo.Create(ctx, cat))
	assert.NotZero(t, cat.ID)

	found, err := repo.FindByID(ctx, cat.ID)
	require.NoError(t, err)
	assert.Equal(t, "Groceries", found.Name)

	foundByName, err := repo.FindByName(ctx, 1, "Groceries")
	require.NoError(t, err)
	assert.Equal(t, cat.ID, foundByName.ID)
}

func TestCategoryRepo_List(t *testing.T) {
	db := setupCategoryDB(t)
	repo := NewCategoryRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &entity.Category{UserID: 1, UserGroupID: 1, Name: "A"})
	repo.Create(ctx, &entity.Category{UserID: 1, UserGroupID: 1, Name: "B"})
	repo.Create(ctx, &entity.Category{UserID: 1, UserGroupID: 1, Name: "C"})

	categories, total, err := repo.List(ctx, 1, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, categories, 3)

	// Test pagination
	categories, total, err = repo.List(ctx, 1, 2, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, categories, 2)
}

func TestCategoryRepo_Search(t *testing.T) {
	db := setupCategoryDB(t)
	repo := NewCategoryRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &entity.Category{UserID: 1, UserGroupID: 1, Name: "Groceries"})
	repo.Create(ctx, &entity.Category{UserID: 1, UserGroupID: 1, Name: "Transport"})
	repo.Create(ctx, &entity.Category{UserID: 1, UserGroupID: 1, Name: "Grocery Store"})

	results, err := repo.Search(ctx, 1, "Grocer", 10)
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestCategoryRepo_Delete(t *testing.T) {
	db := setupCategoryDB(t)
	repo := NewCategoryRepository(db)
	ctx := context.Background()

	cat := &entity.Category{UserID: 1, UserGroupID: 1, Name: "Del"}
	require.NoError(t, repo.Create(ctx, cat))
	require.NoError(t, repo.Delete(ctx, cat.ID))

	_, err := repo.FindByID(ctx, cat.ID)
	assert.Error(t, err)
}

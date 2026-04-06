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

func setupObjectGroupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.ObjectGroupModel{}))
	return db
}

func TestObjectGroupRepo_CreateAndFind(t *testing.T) {
	db := setupObjectGroupDB(t)
	repo := NewObjectGroupRepository(db)
	ctx := context.Background()

	og := &entity.ObjectGroup{
		UserID:      1,
		UserGroupID: 1,
		Title:       "Bills",
		Order:       1,
	}
	require.NoError(t, repo.Create(ctx, og))
	assert.NotZero(t, og.ID)

	found, err := repo.FindByID(ctx, og.ID)
	require.NoError(t, err)
	assert.Equal(t, "Bills", found.Title)
	assert.Equal(t, uint(1), found.UserGroupID)
}

func TestObjectGroupRepo_List(t *testing.T) {
	db := setupObjectGroupDB(t)
	repo := NewObjectGroupRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &entity.ObjectGroup{UserID: 1, UserGroupID: 1, Title: "Group B", Order: 2})
	repo.Create(ctx, &entity.ObjectGroup{UserID: 1, UserGroupID: 1, Title: "Group A", Order: 1})
	repo.Create(ctx, &entity.ObjectGroup{UserID: 1, UserGroupID: 2, Title: "Group C", Order: 3})

	items, total, err := repo.List(ctx, 1, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, items, 2)
	// Should be ordered by order ASC
	assert.Equal(t, "Group A", items[0].Title)
}

func TestObjectGroupRepo_Delete(t *testing.T) {
	db := setupObjectGroupDB(t)
	repo := NewObjectGroupRepository(db)
	ctx := context.Background()

	og := &entity.ObjectGroup{UserID: 1, UserGroupID: 1, Title: "To Delete", Order: 1}
	require.NoError(t, repo.Create(ctx, og))

	require.NoError(t, repo.Delete(ctx, og.ID))

	_, err := repo.FindByID(ctx, og.ID)
	assert.Error(t, err)
}

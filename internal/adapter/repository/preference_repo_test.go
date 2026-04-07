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

func setupPrefDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.PreferenceModel{}))
	return db
}

func TestPrefRepo_CreateAndFindByName(t *testing.T) {
	db := setupPrefDB(t)
	repo := NewPreferenceRepository(db)
	ctx := context.Background()

	pref := &entity.Preference{UserID: 1, Name: "language", Data: "\"en\""}
	require.NoError(t, repo.Create(ctx, pref))
	assert.NotZero(t, pref.ID)

	found, err := repo.FindByUserAndName(ctx, 1, "language")
	require.NoError(t, err)
	assert.Equal(t, "\"en\"", found.Data)
}

func TestPrefRepo_ListByUser(t *testing.T) {
	db := setupPrefDB(t)
	repo := NewPreferenceRepository(db)
	ctx := context.Background()

	require.NoError(t, repo.Create(ctx, &entity.Preference{UserID: 1, Name: "lang", Data: "\"en\""}))
	require.NoError(t, repo.Create(ctx, &entity.Preference{UserID: 1, Name: "theme", Data: "\"dark\""}))
	require.NoError(t, repo.Create(ctx, &entity.Preference{UserID: 2, Name: "lang", Data: "\"de\""}))

	prefs, total, err := repo.ListByUser(ctx, 1, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, prefs, 2)
}

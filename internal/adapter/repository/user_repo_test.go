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

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.UserModel{}, &model.RoleModel{}, &model.RoleUserModel{}))
	return db
}

func TestUserRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &entity.User{Email: "test@example.com", Password: "hashed"}
	err := repo.Create(ctx, user)
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserRepo_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &entity.User{Email: "find@example.com", Password: "hashed"}
	require.NoError(t, repo.Create(ctx, user))

	found, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "find@example.com", found.Email)
}

func TestUserRepo_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &entity.User{Email: "email@example.com", Password: "hashed"}
	require.NoError(t, repo.Create(ctx, user))

	found, err := repo.FindByEmail(ctx, "email@example.com")
	require.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
}

func TestUserRepo_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	require.NoError(t, repo.Create(ctx, &entity.User{Email: "a@b.com", Password: "h"}))
	require.NoError(t, repo.Create(ctx, &entity.User{Email: "c@d.com", Password: "h"}))

	users, total, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, users, 2)
}

func TestUserRepo_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &entity.User{Email: "old@example.com", Password: "hashed"}
	require.NoError(t, repo.Create(ctx, user))

	user.Email = "new@example.com"
	require.NoError(t, repo.Update(ctx, user))

	found, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "new@example.com", found.Email)
}

func TestUserRepo_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &entity.User{Email: "del@example.com", Password: "hashed"}
	require.NoError(t, repo.Create(ctx, user))
	require.NoError(t, repo.Delete(ctx, user.ID))

	_, err := repo.FindByID(ctx, user.ID)
	assert.Error(t, err)
}

func TestUserRepo_GetRole(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &entity.User{Email: "role@example.com", Password: "hashed"}
	require.NoError(t, repo.Create(ctx, user))

	db.Create(&model.RoleModel{ID: 1, Name: "owner"})
	db.Create(&model.RoleUserModel{UserID: user.ID, RoleID: 1})

	role, err := repo.GetRole(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "owner", role)
}

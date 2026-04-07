package repository

import (
	"context"
	"fmt"

	"github.com/anthropics/quillow/internal/adapter/repository/model"
	"github.com/anthropics/quillow/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var m model.UserModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return userModelToEntity(&m), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var m model.UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return userModelToEntity(&m), nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]entity.User, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&model.UserModel{}).Count(&total)

	var models []model.UserModel
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	users := make([]entity.User, len(models))
	for i, m := range models {
		users[i] = *userModelToEntity(&m)
	}
	return users, total, nil
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	m := userEntityToModel(user)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	user.ID = m.ID
	user.CreatedAt = m.CreatedAt
	user.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	m := userEntityToModel(user)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.UserModel{}, id).Error
}

func (r *UserRepository) GetRole(ctx context.Context, userID uint) (string, error) {
	var role model.RoleModel
	err := r.db.WithContext(ctx).
		Joins("JOIN role_user ON role_user.role_id = roles.id").
		Where("role_user.user_id = ?", userID).
		First(&role).Error
	if err != nil {
		return "", nil
	}
	return role.Name, nil
}

func userModelToEntity(m *model.UserModel) *entity.User {
	u := &entity.User{
		ID:        m.ID,
		Email:     m.Email,
		Password:  m.Password,
		Blocked:   m.Blocked,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.UserGroupID != nil {
		u.UserGroupID = *m.UserGroupID
	}
	if m.BlockedCode != nil {
		u.BlockedCode = *m.BlockedCode
	}
	if m.RememberToken != nil {
		u.RememberToken = *m.RememberToken
	}
	if m.Reset != nil {
		u.Reset = *m.Reset
	}
	if m.Domain != nil {
		u.Domain = *m.Domain
	}
	return u
}

func userEntityToModel(u *entity.User) *model.UserModel {
	m := &model.UserModel{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		Blocked:  u.Blocked,
	}
	if u.UserGroupID != 0 {
		m.UserGroupID = &u.UserGroupID
	}
	if u.BlockedCode != "" {
		m.BlockedCode = &u.BlockedCode
	}
	if u.RememberToken != "" {
		m.RememberToken = &u.RememberToken
	}
	if u.Reset != "" {
		m.Reset = &u.Reset
	}
	if u.Domain != "" {
		m.Domain = &u.Domain
	}
	return m
}

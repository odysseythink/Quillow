package user

import (
	"context"
	"fmt"

	"github.com/anthropics/quillow/internal/entity"
	"github.com/anthropics/quillow/internal/port"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repo port.UserRepository
}

func NewUseCase(repo port.UserRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetByID(ctx context.Context, id uint) (*entity.User, string, error) {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, "", err
	}
	role, _ := uc.repo.GetRole(ctx, id)
	return user, role, nil
}

func (uc *UseCase) List(ctx context.Context, limit, offset int) ([]entity.User, int64, error) {
	return uc.repo.List(ctx, limit, offset)
}

func (uc *UseCase) Create(ctx context.Context, email, password string, blocked bool, blockedCode, role string) (*entity.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user := &entity.User{
		Email:       email,
		Password:    string(hashed),
		Blocked:     blocked,
		BlockedCode: blockedCode,
	}
	if err := uc.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UseCase) Update(ctx context.Context, id uint, email string, blocked bool, blockedCode string) (*entity.User, error) {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Email = email
	user.Blocked = blocked
	user.BlockedCode = blockedCode
	if err := uc.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UseCase) ChangePassword(ctx context.Context, id uint, currentPassword, newPassword string) error {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return fmt.Errorf("current password is incorrect")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashed)
	return uc.repo.Update(ctx, user)
}

func (uc *UseCase) ChangeEmail(ctx context.Context, id uint, password, newEmail string) (*entity.User, error) {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("password is incorrect")
	}
	user.Email = newEmail
	if err := uc.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UseCase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}

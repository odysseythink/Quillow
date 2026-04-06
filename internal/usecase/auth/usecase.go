package auth

import (
	"context"
	"errors"

	"github.com/anthropics/firefly-iii-go/internal/port"
	"github.com/anthropics/firefly-iii-go/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type LoginResult struct {
	AccessToken  string
	RefreshToken string
}

type UseCase struct {
	userRepo port.UserRepository
	jwt      *jwt.Service
}

func NewUseCase(userRepo port.UserRepository, jwt *jwt.Service) *UseCase {
	return &UseCase{userRepo: userRepo, jwt: jwt}
}

func (uc *UseCase) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user.Blocked {
		return nil, errors.New("account is blocked")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := uc.jwt.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *UseCase) Refresh(ctx context.Context, refreshToken string) (*LoginResult, error) {
	claims, err := uc.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	user, err := uc.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := uc.jwt.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := uc.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

package jwt

import (
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret             []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

type AccessClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	gojwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID uint `json:"user_id"`
	gojwt.RegisteredClaims
}

func NewService(secret string, accessExpiryHours, refreshExpiryDays int) *Service {
	return &Service{
		secret:             []byte(secret),
		accessTokenExpiry:  time.Duration(accessExpiryHours) * time.Hour,
		refreshTokenExpiry: time.Duration(refreshExpiryDays) * 24 * time.Hour,
	}
}

func (s *Service) GenerateAccessToken(userID uint, email string) (string, error) {
	claims := AccessClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(s.accessTokenExpiry)),
			IssuedAt:  gojwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) ValidateAccessToken(tokenStr string) (*AccessClaims, error) {
	token, err := gojwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(t *gojwt.Token) (any, error) {
		if _, ok := t.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

func (s *Service) GenerateRefreshToken(userID uint) (string, error) {
	claims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(s.refreshTokenExpiry)),
			IssuedAt:  gojwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) ValidateRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := gojwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(t *gojwt.Token) (any, error) {
		if _, ok := t.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAccessToken(t *testing.T) {
	svc := NewService("test-secret-key-that-is-long-enough", 24, 14)
	token, err := svc.GenerateAccessToken(1, "test@example.com")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateAccessToken(t *testing.T) {
	svc := NewService("test-secret-key-that-is-long-enough", 24, 14)
	token, err := svc.GenerateAccessToken(42, "user@example.com")
	require.NoError(t, err)

	claims, err := svc.ValidateAccessToken(token)
	require.NoError(t, err)
	assert.Equal(t, uint(42), claims.UserID)
	assert.Equal(t, "user@example.com", claims.Email)
}

func TestExpiredAccessToken(t *testing.T) {
	svc := NewService("test-secret-key-that-is-long-enough", 0, 14)
	token, err := svc.GenerateAccessToken(1, "test@example.com")
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)
	_, err = svc.ValidateAccessToken(token)
	assert.Error(t, err)
}

func TestInvalidSecret(t *testing.T) {
	svc1 := NewService("secret-one-long-enough-here-yes", 24, 14)
	svc2 := NewService("secret-two-long-enough-here-yes", 24, 14)

	token, err := svc1.GenerateAccessToken(1, "test@example.com")
	require.NoError(t, err)

	_, err = svc2.ValidateAccessToken(token)
	assert.Error(t, err)
}

func TestGenerateRefreshToken(t *testing.T) {
	svc := NewService("test-secret-key-that-is-long-enough", 24, 14)
	token, err := svc.GenerateRefreshToken(1)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateRefreshToken(t *testing.T) {
	svc := NewService("test-secret-key-that-is-long-enough", 24, 14)
	token, err := svc.GenerateRefreshToken(99)
	require.NoError(t, err)

	claims, err := svc.ValidateRefreshToken(token)
	require.NoError(t, err)
	assert.Equal(t, uint(99), claims.UserID)
}

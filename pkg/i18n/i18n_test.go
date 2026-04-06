package i18n

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadLocale(t *testing.T) {
	svc, err := NewService("../../locales")
	require.NoError(t, err)

	err = svc.LoadLocale("en_US")
	require.NoError(t, err)
}

func TestTranslate(t *testing.T) {
	svc, err := NewService("../../locales")
	require.NoError(t, err)
	require.NoError(t, svc.LoadLocale("en_US"))

	ctx := context.WithValue(context.Background(), LocaleKey, "en_US")
	result := svc.T(ctx, "auth.login_success")
	assert.Equal(t, "Login successful", result)
}

func TestTranslateWithParams(t *testing.T) {
	svc, err := NewService("../../locales")
	require.NoError(t, err)
	require.NoError(t, svc.LoadLocale("en_US"))

	ctx := context.WithValue(context.Background(), LocaleKey, "en_US")
	result := svc.T(ctx, "validation.required", "email")
	assert.Equal(t, "The email field is required", result)
}

func TestFallbackToKey(t *testing.T) {
	svc, err := NewService("../../locales")
	require.NoError(t, err)
	require.NoError(t, svc.LoadLocale("en_US"))

	ctx := context.WithValue(context.Background(), LocaleKey, "en_US")
	result := svc.T(ctx, "nonexistent.key")
	assert.Equal(t, "nonexistent.key", result)
}

func TestFallbackLocale(t *testing.T) {
	svc, err := NewService("../../locales")
	require.NoError(t, err)
	require.NoError(t, svc.LoadLocale("en_US"))

	ctx := context.WithValue(context.Background(), LocaleKey, "fr_FR")
	result := svc.T(ctx, "auth.login_success")
	assert.Equal(t, "Login successful", result)
}

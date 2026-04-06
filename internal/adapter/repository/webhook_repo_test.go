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

func setupWebhookDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&model.WebhookModel{},
		&model.WebhookMessageModel{},
		&model.WebhookAttemptModel{},
	))
	return db
}

func TestWebhookRepo_CreateAndFind(t *testing.T) {
	db := setupWebhookDB(t)
	repo := NewWebhookRepository(db)
	ctx := context.Background()

	wh := &entity.Webhook{
		UserID:      1,
		UserGroupID: 1,
		Active:      true,
		Title:       "Deploy Hook",
		Secret:      "s3cret",
		Trigger:     1,
		Response:    1,
		Delivery:    1,
		URL:         "https://example.com/hook",
	}
	require.NoError(t, repo.Create(ctx, wh))
	assert.NotZero(t, wh.ID)

	found, err := repo.FindByID(ctx, wh.ID)
	require.NoError(t, err)
	assert.Equal(t, "Deploy Hook", found.Title)
	assert.Equal(t, "https://example.com/hook", found.URL)
	assert.True(t, found.Active)
}

func TestWebhookRepo_List(t *testing.T) {
	db := setupWebhookDB(t)
	repo := NewWebhookRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &entity.Webhook{UserID: 1, UserGroupID: 1, Title: "A", URL: "https://a.com", Active: true})
	repo.Create(ctx, &entity.Webhook{UserID: 1, UserGroupID: 1, Title: "B", URL: "https://b.com", Active: true})
	repo.Create(ctx, &entity.Webhook{UserID: 1, UserGroupID: 2, Title: "C", URL: "https://c.com", Active: true})

	hooks, total, err := repo.List(ctx, 1, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, hooks, 2)

	// List all
	hooks, total, err = repo.List(ctx, 0, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, hooks, 3)
}

func TestWebhookRepo_Messages(t *testing.T) {
	db := setupWebhookDB(t)
	repo := NewWebhookRepository(db)
	ctx := context.Background()

	wh := &entity.Webhook{UserID: 1, UserGroupID: 1, Title: "Hook", URL: "https://example.com", Active: true}
	require.NoError(t, repo.Create(ctx, wh))

	msg := &entity.WebhookMessage{
		WebhookID: wh.ID,
		Sent:      false,
		Errored:   false,
		UUID:      "uuid-1234",
		Message:   `{"event":"test"}`,
	}
	require.NoError(t, repo.CreateMessage(ctx, msg))
	assert.NotZero(t, msg.ID)

	msgs, err := repo.ListMessages(ctx, wh.ID)
	require.NoError(t, err)
	assert.Len(t, msgs, 1)
	assert.Equal(t, "uuid-1234", msgs[0].UUID)

	attempt := &entity.WebhookAttempt{
		WebhookMessageID: msg.ID,
		StatusCode:       200,
		Logs:             "OK",
		Response:         `{"status":"ok"}`,
	}
	require.NoError(t, repo.CreateAttempt(ctx, attempt))
	assert.NotZero(t, attempt.ID)

	attempts, err := repo.ListAttempts(ctx, msg.ID)
	require.NoError(t, err)
	assert.Len(t, attempts, 1)
	assert.Equal(t, 200, attempts[0].StatusCode)
}

func TestWebhookRepo_Delete(t *testing.T) {
	db := setupWebhookDB(t)
	repo := NewWebhookRepository(db)
	ctx := context.Background()

	wh := &entity.Webhook{
		UserID:      1,
		UserGroupID: 1,
		Title:       "ToDelete",
		URL:         "https://example.com/del",
		Active:      true,
	}
	require.NoError(t, repo.Create(ctx, wh))
	require.NoError(t, repo.Delete(ctx, wh.ID))

	_, err := repo.FindByID(ctx, wh.ID)
	assert.Error(t, err)
}

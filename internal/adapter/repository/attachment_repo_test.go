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

func setupAttachmentDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.AttachmentModel{}, &model.LinkTypeModel{}, &model.TransactionJournalLinkModel{}))
	return db
}

func TestAttachmentRepo_CreateAndFind(t *testing.T) {
	db := setupAttachmentDB(t)
	repo := NewAttachmentRepository(db)
	ctx := context.Background()

	att := &entity.Attachment{
		UserID: 1, UserGroupID: 1,
		AttachableID:   100,
		AttachableType: "FireflyIII\\Models\\TransactionJournal",
		Filename:       "receipt.pdf",
		Mime:           "application/pdf",
		Size:           1024,
		Uploaded:       true,
	}
	require.NoError(t, repo.Create(ctx, att))
	assert.NotZero(t, att.ID)

	found, err := repo.FindByID(ctx, att.ID)
	require.NoError(t, err)
	assert.Equal(t, "receipt.pdf", found.Filename)
}

func TestAttachmentRepo_ListByAttachable(t *testing.T) {
	db := setupAttachmentDB(t)
	repo := NewAttachmentRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &entity.Attachment{UserID: 1, UserGroupID: 1, AttachableID: 1, AttachableType: "FireflyIII\\Models\\TransactionJournal", Filename: "a.pdf", Mime: "application/pdf", Uploaded: true})
	repo.Create(ctx, &entity.Attachment{UserID: 1, UserGroupID: 1, AttachableID: 1, AttachableType: "FireflyIII\\Models\\TransactionJournal", Filename: "b.pdf", Mime: "application/pdf", Uploaded: true})
	repo.Create(ctx, &entity.Attachment{UserID: 1, UserGroupID: 1, AttachableID: 2, AttachableType: "FireflyIII\\Models\\Bill", Filename: "c.pdf", Mime: "application/pdf", Uploaded: true})

	atts, err := repo.ListByAttachable(ctx, "TransactionJournal", 1)
	require.NoError(t, err)
	assert.Len(t, atts, 2)
}

func TestLinkTypeRepo_CreateAndList(t *testing.T) {
	db := setupAttachmentDB(t)
	repo := NewLinkTypeRepository(db)
	ctx := context.Background()

	lt := &entity.LinkType{Name: "Paid", Outward: "is paid by", Inward: "pays for", Editable: false}
	require.NoError(t, repo.Create(ctx, lt))
	assert.NotZero(t, lt.ID)

	items, total, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, items, 1)
	assert.Equal(t, "Paid", items[0].Name)
}

func TestTransactionLinkRepo_CRUD(t *testing.T) {
	db := setupAttachmentDB(t)
	repo := NewTransactionLinkRepository(db)
	ctx := context.Background()

	// Seed a link type
	db.Create(&model.LinkTypeModel{ID: 1, Name: "Paid", Outward: "is paid by", Inward: "pays for"})

	link := &entity.TransactionJournalLink{LinkTypeID: 1, SourceID: 10, DestinationID: 20, Comment: "test"}
	require.NoError(t, repo.Create(ctx, link))
	assert.NotZero(t, link.ID)

	found, err := repo.FindByID(ctx, link.ID)
	require.NoError(t, err)
	assert.Equal(t, uint(10), found.SourceID)

	links, err := repo.ListByJournalID(ctx, 10)
	require.NoError(t, err)
	assert.Len(t, links, 1)
}

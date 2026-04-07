package port

import (
	"context"

	"github.com/anthropics/quillow/internal/entity"
)

type AttachmentRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.Attachment, error)
	List(ctx context.Context, limit, offset int) ([]entity.Attachment, int64, error)
	ListByAttachable(ctx context.Context, attachableType string, attachableID uint) ([]entity.Attachment, error)
	Create(ctx context.Context, attachment *entity.Attachment) error
	Update(ctx context.Context, attachment *entity.Attachment) error
	Delete(ctx context.Context, id uint) error
}

type LinkTypeRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.LinkType, error)
	List(ctx context.Context, limit, offset int) ([]entity.LinkType, int64, error)
	Create(ctx context.Context, lt *entity.LinkType) error
	Update(ctx context.Context, lt *entity.LinkType) error
	Delete(ctx context.Context, id uint) error
}

type TransactionLinkRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.TransactionJournalLink, error)
	List(ctx context.Context, limit, offset int) ([]entity.TransactionJournalLink, int64, error)
	ListByJournalID(ctx context.Context, journalID uint) ([]entity.TransactionJournalLink, error)
	Create(ctx context.Context, link *entity.TransactionJournalLink) error
	Update(ctx context.Context, link *entity.TransactionJournalLink) error
	Delete(ctx context.Context, id uint) error
}

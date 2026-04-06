package transformer

import (
	"fmt"

	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/pkg/response"
)

func TransformLinkType(lt *entity.LinkType) response.Resource {
	return response.Resource{
		Type: "link_types",
		ID:   fmt.Sprintf("%d", lt.ID),
		Attributes: map[string]any{
			"created_at": lt.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at": lt.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"name":       lt.Name,
			"inward":     lt.Inward,
			"outward":    lt.Outward,
			"editable":   lt.Editable,
		},
	}
}

func TransformTransactionLink(link *entity.TransactionJournalLink) response.Resource {
	return response.Resource{
		Type: "transaction_links",
		ID:   fmt.Sprintf("%d", link.ID),
		Attributes: map[string]any{
			"created_at":   link.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			"updated_at":   link.UpdatedAt.Format("2006-01-02T15:04:05-07:00"),
			"link_type_id": fmt.Sprintf("%d", link.LinkTypeID),
			"inward_id":    fmt.Sprintf("%d", link.SourceID),
			"outward_id":   fmt.Sprintf("%d", link.DestinationID),
			"notes":        link.Comment,
		},
	}
}

package importer

import (
	"context"
	"io"

	"github.com/anthropics/quillow/pkg/importer"
)

type UseCase struct{}

func NewUseCase() *UseCase {
	return &UseCase{}
}

// PreviewResult is the response for the preview endpoint.
type PreviewResult struct {
	Source       string              `json:"source"`
	Total        int                 `json:"total"`
	Transactions []PreviewTransaction `json:"transactions"`
}

// PreviewTransaction extends ImportedTransaction with selection and category info.
type PreviewTransaction struct {
	importer.ImportedTransaction
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	Selected     bool   `json:"selected"`
}

// Preview parses a CSV file and returns a preview of transactions.
func (uc *UseCase) Preview(_ context.Context, reader io.Reader) (*PreviewResult, error) {
	source, txs, err := importer.ParseFrom(reader)
	if err != nil {
		return nil, err
	}

	previews := make([]PreviewTransaction, len(txs))
	for i, tx := range txs {
		selected := true
		// Deselect non-successful transactions
		if tx.Status != "" && tx.Status != "支付成功" && tx.Status != "交易成功" {
			selected = false
		}
		previews[i] = PreviewTransaction{
			ImportedTransaction: tx,
			Selected:            selected,
		}
	}

	return &PreviewResult{
		Source:       string(source),
		Total:        len(previews),
		Transactions: previews,
	}, nil
}

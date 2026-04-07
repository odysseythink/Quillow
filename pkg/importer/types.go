package importer

// ImportedTransaction represents a single parsed transaction from a CSV file.
type ImportedTransaction struct {
	Index        int    `json:"index"`
	Date         string `json:"date"`
	Type         string `json:"type"` // withdrawal, deposit, transfer
	Description  string `json:"description"`
	Counterparty string `json:"counterparty"`
	Amount       string `json:"amount"`
	Status       string `json:"status"`
	ExternalID   string `json:"external_id"`
}

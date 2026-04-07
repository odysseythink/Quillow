package importer

import (
	"encoding/csv"
	"io"
	"strings"
)

// ParseWechat parses a WeChat Pay CSV bill export.
func ParseWechat(reader io.Reader) ([]ImportedTransaction, error) {
	r := csv.NewReader(reader)
	r.LazyQuotes = true
	r.FieldsPerRecord = -1

	var results []ImportedTransaction
	lineNum := 0
	headerIdx := map[string]int{}
	index := 0

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			lineNum++
			continue
		}
		lineNum++

		// Skip first 16 lines (overview info)
		if lineNum <= 16 {
			// Line 17 is the header
			if lineNum == 17 {
				for i, col := range record {
					headerIdx[strings.TrimSpace(col)] = i
				}
			}
			continue
		}

		if len(headerIdx) == 0 {
			continue
		}

		getCol := func(name string) string {
			if idx, ok := headerIdx[name]; ok && idx < len(record) {
				return strings.TrimSpace(record[idx])
			}
			return ""
		}

		direction := getCol("收/支")
		if direction != "支出" && direction != "收入" {
			continue // skip non-transaction rows (e.g., "/")
		}

		txType := "withdrawal"
		if direction == "收入" {
			txType = "deposit"
		}

		amount := getCol("金额(元)")
		amount = strings.TrimPrefix(amount, "¥")
		amount = strings.TrimSpace(amount)

		dateStr := getCol("交易时间")
		if len(dateStr) >= 10 {
			dateStr = dateStr[:10] // "2026-04-07 12:34:56" -> "2026-04-07"
		}

		results = append(results, ImportedTransaction{
			Index:        index,
			Date:         dateStr,
			Type:         txType,
			Description:  getCol("商品"),
			Counterparty: getCol("交易对方"),
			Amount:       amount,
			Status:       getCol("当前状态"),
			ExternalID:   getCol("交易单号"),
		})
		index++
	}

	return results, nil
}

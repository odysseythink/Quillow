package importer

import (
	"encoding/csv"
	"io"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// ParseAlipay parses an Alipay CSV bill export (GBK encoded).
func ParseAlipay(reader io.Reader) ([]ImportedTransaction, error) {
	// Convert GBK to UTF-8
	utf8Reader := transform.NewReader(reader, simplifiedchinese.GBK.NewDecoder())

	r := csv.NewReader(utf8Reader)
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

		// Skip first 4 lines (title info), line 5 is header
		if lineNum <= 4 {
			continue
		}
		if lineNum == 5 {
			for i, col := range record {
				headerIdx[strings.TrimSpace(col)] = i
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
			continue
		}

		status := getCol("交易状态")
		if status != "交易成功" && status != "支付成功" {
			continue
		}

		txType := "withdrawal"
		if direction == "收入" {
			txType = "deposit"
		}

		amount := strings.TrimSpace(getCol("金额（元）"))
		if amount == "" {
			amount = strings.TrimSpace(getCol("金额(元)"))
		}

		dateStr := getCol("交易创建时间")
		if dateStr == "" {
			dateStr = getCol("付款时间")
		}
		if len(dateStr) >= 10 {
			dateStr = dateStr[:10]
		}

		results = append(results, ImportedTransaction{
			Index:        index,
			Date:         dateStr,
			Type:         txType,
			Description:  getCol("商品名称"),
			Counterparty: getCol("交易对方"),
			Amount:       amount,
			Status:       status,
			ExternalID:   strings.TrimSpace(getCol("交易号")),
		})
		index++
	}

	return results, nil
}

package importer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Source represents the detected CSV source.
type Source string

const (
	SourceWechat  Source = "wechat"
	SourceAlipay  Source = "alipay"
	SourceUnknown Source = "unknown"
)

// Detect reads the first few lines to determine the CSV source.
func Detect(data []byte) Source {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	lines := 0
	for scanner.Scan() && lines < 20 {
		line := scanner.Text()
		if strings.Contains(line, "微信支付账单") {
			return SourceWechat
		}
		if strings.Contains(line, "支付宝") || strings.Contains(line, "交易号") {
			return SourceAlipay
		}
		lines++
	}
	return SourceUnknown
}

// Parse auto-detects the source and parses the CSV data.
func Parse(data []byte) (Source, []ImportedTransaction, error) {
	source := Detect(data)
	reader := bytes.NewReader(data)

	switch source {
	case SourceWechat:
		txs, err := ParseWechat(reader)
		return source, txs, err
	case SourceAlipay:
		txs, err := ParseAlipay(reader)
		return source, txs, err
	default:
		return source, nil, fmt.Errorf("unrecognized CSV format")
	}
}

// ParseFrom parses from a reader, reading all data first for detection.
func ParseFrom(r io.Reader) (Source, []ImportedTransaction, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return SourceUnknown, nil, fmt.Errorf("failed to read file: %w", err)
	}
	return Parse(data)
}

package ai

import (
	"regexp"
	"strings"
	"time"
)

// ParsedTransaction holds the result of natural language transaction parsing.
type ParsedTransaction struct {
	Type            string `json:"type"`
	Description     string `json:"description"`
	Amount          string `json:"amount"`
	Date            string `json:"date"`
	Category        string `json:"category"`
	CategoryID      uint   `json:"category_id"`
	SourceName      string `json:"source_name"`
	DestinationName string `json:"destination_name"`
}

var (
	amountPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(\d+\.?\d*)\s*元`),
		regexp.MustCompile(`¥\s*(\d+\.?\d*)`),
		regexp.MustCompile(`(\d+\.?\d*)\s*块`),
		regexp.MustCompile(`(\d+\.?\d*)\s*刀`),
	}

	dateKeywords = map[string]int{
		"今天": 0, "今日": 0,
		"昨天": -1, "昨日": -1,
		"前天": -2, "前日": -2,
	}

	datePattern  = regexp.MustCompile(`(\d{1,2})\s*月\s*(\d{1,2})\s*[日号]`)
	isoPattern   = regexp.MustCompile(`(\d{4})-(\d{1,2})-(\d{1,2})`)

	depositKeywords  = []string{"收入", "工资", "奖金", "转入", "红包收", "退款"}
	transferKeywords = []string{"转账", "转给", "还款"}

	queryKeywords = []string{
		"多少", "花了", "支出", "收入了", "趋势", "最大", "最小", "平均",
		"比较", "哪个", "几笔", "统计", "分析", "总共", "合计", "余额",
		"预算", "对比", "变化", "增长", "减少", "排名", "top",
	}
)

// ParseLocal extracts transaction fields from natural language using regex.
// Returns the parsed result and confidence ("high" or "low").
func ParseLocal(message string, today time.Time) (*ParsedTransaction, string) {
	result := &ParsedTransaction{
		Type: "withdrawal",
		Date: today.Format("2006-01-02"),
	}

	remaining := message

	// Extract amount
	for _, p := range amountPatterns {
		if m := p.FindStringSubmatch(message); len(m) > 1 {
			result.Amount = m[1]
			remaining = p.ReplaceAllString(remaining, "")
			break
		}
	}

	// Extract date
	for kw, offset := range dateKeywords {
		if strings.Contains(message, kw) {
			result.Date = today.AddDate(0, 0, offset).Format("2006-01-02")
			remaining = strings.ReplaceAll(remaining, kw, "")
			break
		}
	}
	if m := isoPattern.FindStringSubmatch(message); len(m) == 4 {
		result.Date = m[0]
		remaining = isoPattern.ReplaceAllString(remaining, "")
	} else if m := datePattern.FindStringSubmatch(message); len(m) == 3 {
		month := m[1]
		day := m[2]
		if len(month) == 1 {
			month = "0" + month
		}
		if len(day) == 1 {
			day = "0" + day
		}
		result.Date = today.Format("2006") + "-" + month + "-" + day
		remaining = datePattern.ReplaceAllString(remaining, "")
	}

	// Detect type
	for _, kw := range depositKeywords {
		if strings.Contains(message, kw) {
			result.Type = "deposit"
			remaining = strings.ReplaceAll(remaining, kw, "")
			break
		}
	}
	if result.Type == "withdrawal" {
		for _, kw := range transferKeywords {
			if strings.Contains(message, kw) {
				result.Type = "transfer"
				remaining = strings.ReplaceAll(remaining, kw, "")
				break
			}
		}
	}

	// Description is the remaining text
	remaining = strings.TrimSpace(remaining)
	remaining = strings.Trim(remaining, "，。,. ")
	if remaining != "" {
		result.Description = remaining
	}

	// Confidence: high if we have amount + description
	confidence := "low"
	if result.Amount != "" && result.Description != "" {
		confidence = "high"
	}

	return result, confidence
}

// DetectIntent determines whether a message is a record or query intent.
func DetectIntent(message string) string {
	msg := strings.ToLower(message)

	// Check query keywords first
	for _, kw := range queryKeywords {
		if strings.Contains(msg, kw) {
			return "query"
		}
	}

	// Check if contains amount pattern (likely record)
	for _, p := range amountPatterns {
		if p.MatchString(message) {
			return "record"
		}
	}

	return "unknown"
}

package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// QueryFunction defines a safe, predefined query that LLM can invoke.
type QueryFunction struct {
	Name        string
	Description string
	Execute     func(ctx context.Context, db *gorm.DB, userID uint, params map[string]string) (any, error)
}

// QueryRegistry holds all available query functions.
type QueryRegistry struct {
	functions map[string]QueryFunction
	db        *gorm.DB
}

// NewQueryRegistry creates a registry with all predefined query functions.
func NewQueryRegistry(db *gorm.DB) *QueryRegistry {
	r := &QueryRegistry{
		functions: make(map[string]QueryFunction),
		db:        db,
	}
	r.register()
	return r
}

func (r *QueryRegistry) register() {
	r.functions["total_spending"] = QueryFunction{
		Name:        "total_spending",
		Description: "total_spending(start, end): Total spending in a date range",
		Execute: func(ctx context.Context, db *gorm.DB, userID uint, params map[string]string) (any, error) {
			var total float64
			err := db.WithContext(ctx).Raw(`
				SELECT COALESCE(SUM(ABS(t.amount)), 0) FROM transactions t
				JOIN transaction_journals tj ON t.transaction_journal_id = tj.id
				JOIN transaction_types tt ON tj.transaction_type_id = tt.id
				WHERE tj.user_id = ? AND tt.type = 'Withdrawal'
				AND tj.date >= ? AND tj.date <= ? AND t.amount < 0
			`, userID, params["start"], params["end"]).Scan(&total).Error
			return map[string]any{"total": fmt.Sprintf("%.2f", total), "start": params["start"], "end": params["end"]}, err
		},
	}

	r.functions["total_income"] = QueryFunction{
		Name:        "total_income",
		Description: "total_income(start, end): Total income in a date range",
		Execute: func(ctx context.Context, db *gorm.DB, userID uint, params map[string]string) (any, error) {
			var total float64
			err := db.WithContext(ctx).Raw(`
				SELECT COALESCE(SUM(t.amount), 0) FROM transactions t
				JOIN transaction_journals tj ON t.transaction_journal_id = tj.id
				JOIN transaction_types tt ON tj.transaction_type_id = tt.id
				WHERE tj.user_id = ? AND tt.type = 'Deposit'
				AND tj.date >= ? AND tj.date <= ? AND t.amount > 0
			`, userID, params["start"], params["end"]).Scan(&total).Error
			return map[string]any{"total": fmt.Sprintf("%.2f", total), "start": params["start"], "end": params["end"]}, err
		},
	}

	r.functions["sum_by_category"] = QueryFunction{
		Name:        "sum_by_category",
		Description: "sum_by_category(start, end, type): Sum grouped by category",
		Execute: func(ctx context.Context, db *gorm.DB, userID uint, params map[string]string) (any, error) {
			type row struct {
				Category string  `json:"category"`
				Total    float64 `json:"total"`
				Count    int     `json:"count"`
			}
			var rows []row
			err := db.WithContext(ctx).Raw(`
				SELECT COALESCE(c.name, 'Uncategorized') as category,
					SUM(ABS(t.amount)) as total, COUNT(*) as count
				FROM transactions t
				JOIN transaction_journals tj ON t.transaction_journal_id = tj.id
				LEFT JOIN categories c ON tj.category_id = c.id
				JOIN transaction_types tt ON tj.transaction_type_id = tt.id
				WHERE tj.user_id = ? AND tj.date >= ? AND tj.date <= ?
				AND t.amount < 0
				GROUP BY c.name ORDER BY total DESC
			`, userID, params["start"], params["end"]).Scan(&rows).Error
			return rows, err
		},
	}

	r.functions["top_transactions"] = QueryFunction{
		Name:        "top_transactions",
		Description: "top_transactions(start, end, type, order, limit): Top N transactions",
		Execute: func(ctx context.Context, db *gorm.DB, userID uint, params map[string]string) (any, error) {
			type row struct {
				Description string  `json:"description"`
				Amount      float64 `json:"amount"`
				Date        string  `json:"date"`
			}
			order := "DESC"
			if params["order"] == "asc" {
				order = "ASC"
			}
			limit := "5"
			if params["limit"] != "" {
				limit = params["limit"]
			}
			var rows []row
			err := db.WithContext(ctx).Raw(fmt.Sprintf(`
				SELECT tj.description, ABS(t.amount) as amount,
					DATE(tj.date) as date
				FROM transactions t
				JOIN transaction_journals tj ON t.transaction_journal_id = tj.id
				WHERE tj.user_id = ? AND tj.date >= ? AND tj.date <= ?
				AND t.amount < 0
				ORDER BY ABS(t.amount) %s LIMIT %s
			`, order, limit), userID, params["start"], params["end"]).Scan(&rows).Error
			return rows, err
		},
	}

	r.functions["daily_trend"] = QueryFunction{
		Name:        "daily_trend",
		Description: "daily_trend(start, end, type): Daily spending/income trend",
		Execute: func(ctx context.Context, db *gorm.DB, userID uint, params map[string]string) (any, error) {
			type row struct {
				Date  string  `json:"date"`
				Total float64 `json:"total"`
			}
			var rows []row
			err := db.WithContext(ctx).Raw(`
				SELECT DATE(tj.date) as date, SUM(ABS(t.amount)) as total
				FROM transactions t
				JOIN transaction_journals tj ON t.transaction_journal_id = tj.id
				WHERE tj.user_id = ? AND tj.date >= ? AND tj.date <= ?
				AND t.amount < 0
				GROUP BY DATE(tj.date) ORDER BY date
			`, userID, params["start"], params["end"]).Scan(&rows).Error
			return rows, err
		},
	}

	r.functions["monthly_comparison"] = QueryFunction{
		Name:        "monthly_comparison",
		Description: "monthly_comparison(month1, month2, type): Compare two months",
		Execute: func(ctx context.Context, db *gorm.DB, userID uint, params map[string]string) (any, error) {
			getMonthTotal := func(month string) (float64, error) {
				t, err := time.Parse("2006-01", month)
				if err != nil {
					return 0, err
				}
				start := t.Format("2006-01-02")
				end := t.AddDate(0, 1, -1).Format("2006-01-02")
				var total float64
				err = db.WithContext(ctx).Raw(`
					SELECT COALESCE(SUM(ABS(t.amount)), 0) FROM transactions t
					JOIN transaction_journals tj ON t.transaction_journal_id = tj.id
					WHERE tj.user_id = ? AND tj.date >= ? AND tj.date <= ? AND t.amount < 0
				`, userID, start, end).Scan(&total).Error
				return total, err
			}
			t1, _ := getMonthTotal(params["month1"])
			t2, _ := getMonthTotal(params["month2"])
			diff := t1 - t2
			pct := 0.0
			if t2 > 0 {
				pct = (diff / t2) * 100
			}
			return map[string]any{
				"month1": params["month1"], "month1_total": fmt.Sprintf("%.2f", t1),
				"month2": params["month2"], "month2_total": fmt.Sprintf("%.2f", t2),
				"diff": fmt.Sprintf("%.2f", diff), "pct": fmt.Sprintf("%.1f%%", pct),
			}, nil
		},
	}

	r.functions["category_comparison"] = QueryFunction{
		Name:        "category_comparison",
		Description: "category_comparison(category, month1, month2): Compare a category across months",
		Execute: func(ctx context.Context, db *gorm.DB, userID uint, params map[string]string) (any, error) {
			getCatTotal := func(month, category string) (float64, error) {
				t, err := time.Parse("2006-01", month)
				if err != nil {
					return 0, err
				}
				start := t.Format("2006-01-02")
				end := t.AddDate(0, 1, -1).Format("2006-01-02")
				var total float64
				err = db.WithContext(ctx).Raw(`
					SELECT COALESCE(SUM(ABS(t.amount)), 0) FROM transactions t
					JOIN transaction_journals tj ON t.transaction_journal_id = tj.id
					LEFT JOIN categories c ON tj.category_id = c.id
					WHERE tj.user_id = ? AND tj.date >= ? AND tj.date <= ?
					AND t.amount < 0 AND c.name = ?
				`, userID, start, end, category).Scan(&total).Error
				return total, err
			}
			t1, _ := getCatTotal(params["month1"], params["category"])
			t2, _ := getCatTotal(params["month2"], params["category"])
			return map[string]any{
				"category": params["category"],
				"month1": params["month1"], "month1_total": fmt.Sprintf("%.2f", t1),
				"month2": params["month2"], "month2_total": fmt.Sprintf("%.2f", t2),
			}, nil
		},
	}

	r.functions["account_balance"] = QueryFunction{
		Name:        "account_balance",
		Description: "account_balance(account_name): Current balance of an account",
		Execute: func(ctx context.Context, db *gorm.DB, userID uint, params map[string]string) (any, error) {
			var balance float64
			err := db.WithContext(ctx).Raw(`
				SELECT COALESCE(SUM(t.amount), 0) FROM transactions t
				JOIN accounts a ON t.account_id = a.id
				WHERE a.user_id = ? AND a.name = ?
			`, userID, params["account_name"]).Scan(&balance).Error
			return map[string]any{"account": params["account_name"], "balance": fmt.Sprintf("%.2f", balance)}, err
		},
	}

	r.functions["budget_status"] = QueryFunction{
		Name:        "budget_status",
		Description: "budget_status(budget_name, month): Budget usage status",
		Execute: func(ctx context.Context, db *gorm.DB, userID uint, params map[string]string) (any, error) {
			type result struct {
				BudgetName string  `json:"budget_name"`
				LimitAmt   float64 `json:"limit_amount"`
				SpentAmt   float64 `json:"spent_amount"`
			}
			var r result
			r.BudgetName = params["budget_name"]
			// Get budget limit
			db.WithContext(ctx).Raw(`
				SELECT COALESCE(bl.amount, 0) FROM budget_limits bl
				JOIN budgets b ON bl.budget_id = b.id
				WHERE b.user_id = ? AND b.name = ?
				ORDER BY bl.start_date DESC LIMIT 1
			`, userID, params["budget_name"]).Scan(&r.LimitAmt)
			return r, nil
		},
	}
}

// Execute runs a query function by name.
func (r *QueryRegistry) Execute(ctx context.Context, userID uint, name string, params map[string]string) (any, error) {
	fn, ok := r.functions[name]
	if !ok {
		return nil, fmt.Errorf("unknown query function: %s", name)
	}
	return fn.Execute(ctx, r.db, userID, params)
}

// Describe returns a description of all available functions for LLM prompting.
func (r *QueryRegistry) Describe() string {
	var lines []string
	for _, fn := range r.functions {
		lines = append(lines, "- "+fn.Description)
	}
	return strings.Join(lines, "\n")
}

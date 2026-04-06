package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/anthropics/firefly-iii-go/internal/adapter/transformer"
	accountuc "github.com/anthropics/firefly-iii-go/internal/usecase/account"
	txuc "github.com/anthropics/firefly-iii-go/internal/usecase/transaction"
	"github.com/anthropics/firefly-iii-go/pkg/pagination"
	"github.com/anthropics/firefly-iii-go/pkg/response"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	uc     *txuc.UseCase
	acctUC *accountuc.UseCase
}

func NewTransactionHandler(uc *txuc.UseCase, acctUC *accountuc.UseCase) *TransactionHandler {
	return &TransactionHandler{uc: uc, acctUC: acctUC}
}

func (h *TransactionHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	txType := c.Query("type")

	var start, end *time.Time
	if s := c.Query("start"); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			start = &t
		}
	}
	if e := c.Query("end"); e != "" {
		if t, err := time.Parse("2006-01-02", e); err == nil {
			end = &t
		}
	}

	groups, total, err := h.uc.ListGroups(c.Request.Context(), 0, txType, start, end, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(groups))
	for i, group := range groups {
		extras := h.enrichGroup(c, group.ID)
		items[i] = transformer.TransformTransactionGroup(&group, extras)
	}

	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *TransactionHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid transaction ID")
		return
	}

	group, _, err := h.uc.GetGroupByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Transaction not found")
		return
	}

	extras := h.enrichGroup(c, group.ID)
	resource := transformer.TransformTransactionGroup(group, extras)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type storeTransactionRequest struct {
	GroupTitle    string                    `json:"group_title"`
	Transactions []storeTransactionJournal `json:"transactions" binding:"required,min=1"`
}

type storeTransactionJournal struct {
	Type              string   `json:"type" binding:"required"`
	Description       string   `json:"description" binding:"required"`
	Date              string   `json:"date" binding:"required"`
	Amount            string   `json:"amount" binding:"required"`
	ForeignAmount     string   `json:"foreign_amount"`
	CurrencyID        uint     `json:"currency_id"`
	CurrencyCode      string   `json:"currency_code"`
	ForeignCurrencyID *uint    `json:"foreign_currency_id"`
	SourceID          uint     `json:"source_id"`
	SourceName        string   `json:"source_name"`
	DestinationID     uint     `json:"destination_id"`
	DestinationName   string   `json:"destination_name"`
	BudgetID          *uint    `json:"budget_id"`
	CategoryID        *uint    `json:"category_id"`
	BillID            *uint    `json:"bill_id"`
	Tags              []string `json:"tags"`
	Notes             string   `json:"notes"`
	Reconciled        bool     `json:"reconciled"`
	ExternalID        string   `json:"external_id"`
	ExternalURL       string   `json:"external_url"`
	InternalReference string   `json:"internal_reference"`
}

func (h *TransactionHandler) Store(c *gin.Context) {
	var req storeTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	txInputs := make([]txuc.TransactionInput, len(req.Transactions))
	for i, t := range req.Transactions {
		date, _ := time.Parse("2006-01-02", t.Date)
		txInputs[i] = txuc.TransactionInput{
			Type:              t.Type,
			Description:       t.Description,
			Date:              date,
			Amount:            t.Amount,
			ForeignAmount:     t.ForeignAmount,
			CurrencyID:        t.CurrencyID,
			ForeignCurrencyID: t.ForeignCurrencyID,
			SourceID:          t.SourceID,
			SourceName:        t.SourceName,
			DestinationID:     t.DestinationID,
			DestinationName:   t.DestinationName,
			BudgetID:          t.BudgetID,
			CategoryID:        t.CategoryID,
			BillID:            t.BillID,
			Tags:              t.Tags,
			Notes:             t.Notes,
			Reconciled:        t.Reconciled,
			ExternalID:        t.ExternalID,
			ExternalURL:       t.ExternalURL,
			InternalRef:       t.InternalReference,
		}
	}

	input := txuc.CreateGroupInput{
		UserID:       userID,
		GroupTitle:    req.GroupTitle,
		Transactions: txInputs,
	}

	group, err := h.uc.CreateGroup(c.Request.Context(), input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	extras := h.enrichGroup(c, group.ID)
	resource := transformer.TransformTransactionGroup(group, extras)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *TransactionHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid transaction ID")
		return
	}
	if err := h.uc.DeleteGroup(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "Transaction not found")
		return
	}
	response.NoContent(c)
}

// Search endpoints
func (h *TransactionHandler) Search(c *gin.Context) {
	query := c.Query("query")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	groups, total, err := h.uc.SearchGroups(c.Request.Context(), 0, query, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(groups))
	for i, group := range groups {
		extras := h.enrichGroup(c, group.ID)
		items[i] = transformer.TransformTransactionGroup(&group, extras)
	}

	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *TransactionHandler) SearchCount(c *gin.Context) {
	query := c.Query("query")
	count, err := h.uc.CountByQuery(c.Request.Context(), 0, query)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

// Insight endpoints
func (h *TransactionHandler) InsightExpense(c *gin.Context) {
	h.insightByType(c, "Withdrawal")
}

func (h *TransactionHandler) InsightIncome(c *gin.Context) {
	h.insightByType(c, "Deposit")
}

func (h *TransactionHandler) InsightTransfer(c *gin.Context) {
	h.insightByType(c, "Transfer")
}

func (h *TransactionHandler) insightByType(c *gin.Context, txType string) {
	start, end := parseDateRange(c)
	entries, err := h.uc.InsightByType(c.Request.Context(), 0, txType, start, end)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result := make([]map[string]any, len(entries))
	for i, e := range entries {
		result[i] = map[string]any{
			"id":               fmt.Sprintf("%d", e.ID),
			"name":             e.Name,
			"difference":       e.Difference,
			"difference_float": e.DifferenceFloat,
			"currency_id":      fmt.Sprintf("%d", e.CurrencyID),
			"currency_code":    e.CurrencyCode,
		}
	}
	c.JSON(http.StatusOK, result)
}

// Summary
func (h *TransactionHandler) Summary(c *gin.Context) {
	start, end := parseDateRange(c)

	result := make(map[string]any)

	// Spent (withdrawals)
	spent, _ := h.uc.InsightByType(c.Request.Context(), 0, "Withdrawal", start, end)
	for _, s := range spent {
		key := fmt.Sprintf("spent-in-%s", s.CurrencyCode)
		result[key] = map[string]any{
			"key":                     key,
			"title":                   "Spent",
			"monetary_value":          s.Difference,
			"currency_id":             fmt.Sprintf("%d", s.CurrencyID),
			"currency_code":           s.CurrencyCode,
			"value_parsed":            s.Difference,
			"local_icon":              "balance-scale",
			"sub_title":               "",
		}
	}

	// Earned (deposits)
	earned, _ := h.uc.InsightByType(c.Request.Context(), 0, "Deposit", start, end)
	for _, e := range earned {
		key := fmt.Sprintf("earned-in-%s", e.CurrencyCode)
		result[key] = map[string]any{
			"key":                     key,
			"title":                   "Earned",
			"monetary_value":          e.Difference,
			"currency_id":             fmt.Sprintf("%d", e.CurrencyID),
			"currency_code":           e.CurrencyCode,
			"value_parsed":            e.Difference,
			"local_icon":              "balance-scale",
			"sub_title":               "",
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *TransactionHandler) enrichGroup(c *gin.Context, groupID uint) []transformer.TransactionJournalExtra {
	ctx := c.Request.Context()
	_, journals, _ := h.uc.GetGroupByID(ctx, groupID)

	extras := make([]transformer.TransactionJournalExtra, len(journals))
	for i, j := range journals {
		extra := transformer.TransactionJournalExtra{
			Tags: []string{},
		}

		txns, _ := h.uc.GetTransactions(ctx, j.ID)
		for _, t := range txns {
			if t.Amount != "" && t.Amount[0] == '-' {
				extra.SourceID = t.AccountID
				extra.Amount = strings.TrimPrefix(t.Amount, "-")
				extra.ForeignAmount = t.ForeignAmount
			} else {
				extra.DestinationID = t.AccountID
			}
		}

		// Enrich account names
		if acct, at, err := h.acctUC.GetByID(ctx, extra.SourceID); err == nil {
			extra.SourceName = acct.Name
			if at != nil {
				extra.SourceType = strings.ToLower(at.Type)
			}
		}
		if acct, at, err := h.acctUC.GetByID(ctx, extra.DestinationID); err == nil {
			extra.DestinationName = acct.Name
			if at != nil {
				extra.DestinationType = strings.ToLower(at.Type)
			}
		}

		extra.Notes, _ = h.uc.GetJournalNotes(ctx, j.ID)
		extra.Reconciled = false
		if tags, err := h.uc.GetJournalTags(ctx, j.ID); err == nil && len(tags) > 0 {
			extra.Tags = tags
		}
		if eid, _ := h.uc.GetJournalMeta(ctx, j.ID, "external_id"); eid != "" {
			extra.ExternalID = eid
		}
		if eurl, _ := h.uc.GetJournalMeta(ctx, j.ID, "external_url"); eurl != "" {
			extra.ExternalURL = eurl
		}
		if iref, _ := h.uc.GetJournalMeta(ctx, j.ID, "internal_reference"); iref != "" {
			extra.InternalRef = iref
		}

		extras[i] = extra
	}

	return extras
}

func parseDateRange(c *gin.Context) (time.Time, time.Time) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, -1)

	if s := c.Query("start"); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			start = t
		}
	}
	if e := c.Query("end"); e != "" {
		if t, err := time.Parse("2006-01-02", e); err == nil {
			end = t
		}
	}
	return start, end
}

package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/anthropics/quillow/internal/adapter/transformer"
	accountuc "github.com/anthropics/quillow/internal/usecase/account"
	currencyuc "github.com/anthropics/quillow/internal/usecase/currency"
	"github.com/anthropics/quillow/pkg/pagination"
	"github.com/anthropics/quillow/pkg/response"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	uc     *accountuc.UseCase
	currUC *currencyuc.UseCase
}

func NewAccountHandler(uc *accountuc.UseCase, currUC *currencyuc.UseCase) *AccountHandler {
	return &AccountHandler{uc: uc, currUC: currUC}
}

func (h *AccountHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 1000 {
		limit = 50
	}
	offset := (page - 1) * limit

	var accountTypes []string
	if t := c.Query("type"); t != "" {
		accountTypes = strings.Split(t, ",")
	}

	userGroupID := uint(0) // TODO: get from user context in later SP

	accounts, total, err := h.uc.List(c.Request.Context(), userGroupID, accountTypes, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(accounts))
	for i, acct := range accounts {
		extra := h.enrichAccount(c, acct.ID)
		items[i] = transformer.TransformAccount(&acct, extra)
	}

	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *AccountHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("account"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid account ID")
		return
	}

	acct, at, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Account not found")
		return
	}

	extra := h.enrichAccount(c, acct.ID)
	if at != nil {
		extra.AccountType = strings.ToLower(at.Type)
	}
	resource := transformer.TransformAccount(acct, extra)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type flexUint uint

func (f *flexUint) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	if s == "" || s == "null" {
		*f = 0
		return nil
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	*f = flexUint(v)
	return nil
}

type storeAccountRequest struct {
	Name           string   `json:"name" binding:"required"`
	Type           string   `json:"type" binding:"required"`
	IBAN           string   `json:"iban"`
	AccountNumber  string   `json:"account_number"`
	VirtualBalance string   `json:"virtual_balance"`
	Active         *bool    `json:"active"`
	Order          int      `json:"order"`
	AccountRole    string   `json:"account_role"`
	CurrencyID     flexUint `json:"currency_id"`
	CurrencyCode   string   `json:"currency_code"`
	Notes          string   `json:"notes"`
}

func (h *AccountHandler) Store(c *gin.Context) {
	var req storeAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	// Resolve currency
	currencyID := uint(req.CurrencyID)
	if currencyID == 0 && req.CurrencyCode != "" {
		curr, err := h.currUC.GetByCode(c.Request.Context(), req.CurrencyCode)
		if err == nil {
			currencyID = curr.ID
		}
	}

	userID := c.GetUint("user_id")
	input := accountuc.CreateAccountInput{
		UserID:         userID,
		UserGroupID:    0, // TODO: resolve from user
		Name:           req.Name,
		Type:           req.Type,
		IBAN:           req.IBAN,
		AccountNumber:  req.AccountNumber,
		VirtualBalance: req.VirtualBalance,
		Active:         active,
		Order:          req.Order,
		AccountRole:    req.AccountRole,
		CurrencyID:     currencyID,
		Notes:          req.Notes,
	}

	acct, err := h.uc.Create(c.Request.Context(), input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	extra := h.enrichAccount(c, acct.ID)
	resource := transformer.TransformAccount(acct, extra)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type updateAccountRequest struct {
	Name           string `json:"name" binding:"required"`
	IBAN           string `json:"iban"`
	AccountNumber  string `json:"account_number"`
	VirtualBalance string `json:"virtual_balance"`
	Active         *bool  `json:"active"`
	Order          int    `json:"order"`
	AccountRole    string `json:"account_role"`
	Notes          string `json:"notes"`
}

func (h *AccountHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("account"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid account ID")
		return
	}

	var req updateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	input := accountuc.UpdateAccountInput{
		Name:           req.Name,
		IBAN:           req.IBAN,
		VirtualBalance: req.VirtualBalance,
		Active:         active,
		Order:          req.Order,
		AccountRole:    req.AccountRole,
		AccountNumber:  req.AccountNumber,
		Notes:          req.Notes,
	}

	acct, err := h.uc.Update(c.Request.Context(), uint(id), input)
	if err != nil {
		response.NotFound(c, "Account not found")
		return
	}

	extra := h.enrichAccount(c, acct.ID)
	resource := transformer.TransformAccount(acct, extra)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *AccountHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("account"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid account ID")
		return
	}
	if err := h.uc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "Account not found")
		return
	}
	response.NoContent(c)
}

func (h *AccountHandler) enrichAccount(c *gin.Context, accountID uint) transformer.AccountExtra {
	ctx := c.Request.Context()
	extra := transformer.AccountExtra{
		CurrentBalance: "0",
	}

	if _, accountType, err := h.uc.GetByID(ctx, accountID); err == nil && accountType != nil {
		extra.AccountType = strings.ToLower(accountType.Type)
	}

	if role, err := h.uc.GetMeta(ctx, accountID, "account_role"); err == nil {
		extra.AccountRole = role
	}
	if acctNum, err := h.uc.GetMeta(ctx, accountID, "account_number"); err == nil {
		extra.AccountNumber = acctNum
	}
	if notes, err := h.uc.GetNotes(ctx, accountID); err == nil {
		extra.Notes = notes
	}
	if currIDStr, err := h.uc.GetMeta(ctx, accountID, "currency_id"); err == nil && currIDStr != "" {
		if currID, parseErr := strconv.ParseUint(currIDStr, 10, 32); parseErr == nil {
			if curr, currErr := h.currUC.GetByID(ctx, uint(currID)); currErr == nil {
				extra.CurrencyID = fmt.Sprintf("%d", curr.ID)
				extra.CurrencyCode = curr.Code
				extra.CurrencyName = curr.Name
				extra.CurrencySymbol = curr.Symbol
				extra.CurrencyDP = curr.DecimalPlaces
			}
		}
	}

	return extra
}

// Autocomplete endpoint
func (h *AccountHandler) Autocomplete(c *gin.Context) {
	query := c.Query("query")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var accountTypes []string
	if t := c.Query("types"); t != "" {
		accountTypes = strings.Split(t, ",")
	}

	accounts, err := h.uc.Search(c.Request.Context(), 0, query, accountTypes, limit)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result := make([]transformer.AutocompleteAccount, len(accounts))
	for i, acct := range accounts {
		_, at, _ := h.uc.GetByID(c.Request.Context(), acct.ID)
		typeName := ""
		if at != nil {
			typeName = strings.ToLower(at.Type)
		}
		result[i] = transformer.AutocompleteAccount{
			ID:     fmt.Sprintf("%d", acct.ID),
			Name:   acct.Name,
			Active: acct.Active,
			Type:   typeName,
		}
	}

	c.JSON(http.StatusOK, result)
}

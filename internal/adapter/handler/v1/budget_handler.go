package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anthropics/quillow/internal/adapter/transformer"
	"github.com/anthropics/quillow/internal/entity"
	budgetuc "github.com/anthropics/quillow/internal/usecase/budget"
	"github.com/anthropics/quillow/pkg/pagination"
	"github.com/anthropics/quillow/pkg/response"
	"github.com/gin-gonic/gin"
)

type BudgetHandler struct {
	uc *budgetuc.UseCase
}

func NewBudgetHandler(uc *budgetuc.UseCase) *BudgetHandler {
	return &BudgetHandler{uc: uc}
}

func (h *BudgetHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 1000 {
		limit = 50
	}
	offset := (page - 1) * limit

	userGroupID := uint(0) // TODO: get from user context in later SP

	budgets, total, err := h.uc.List(c.Request.Context(), userGroupID, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(budgets))
	for i, b := range budgets {
		notes, _ := h.uc.GetNotes(c.Request.Context(), b.ID)
		items[i] = transformer.TransformBudget(&b, notes)
	}

	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *BudgetHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("budget"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid budget ID")
		return
	}

	b, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Budget not found")
		return
	}

	notes, _ := h.uc.GetNotes(c.Request.Context(), b.ID)
	resource := transformer.TransformBudget(b, notes)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type storeBudgetRequest struct {
	Name   string `json:"name" binding:"required"`
	Active *bool  `json:"active"`
	Order  uint   `json:"order"`
	Notes  string `json:"notes"`
}

func (h *BudgetHandler) Store(c *gin.Context) {
	var req storeBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	b := &entity.Budget{
		Name:   req.Name,
		Active: active,
		Order:  req.Order,
	}

	if err := h.uc.Create(c.Request.Context(), b); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	notes, _ := h.uc.GetNotes(c.Request.Context(), b.ID)
	resource := transformer.TransformBudget(b, notes)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type updateBudgetRequest struct {
	Name   string `json:"name" binding:"required"`
	Active *bool  `json:"active"`
	Order  uint   `json:"order"`
	Notes  string `json:"notes"`
}

func (h *BudgetHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("budget"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid budget ID")
		return
	}

	b, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Budget not found")
		return
	}

	var req updateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	b.Name = req.Name
	b.Order = req.Order
	if req.Active != nil {
		b.Active = *req.Active
	}

	if err := h.uc.Update(c.Request.Context(), b); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	notes, _ := h.uc.GetNotes(c.Request.Context(), b.ID)
	resource := transformer.TransformBudget(b, notes)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *BudgetHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("budget"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid budget ID")
		return
	}
	if err := h.uc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "Budget not found")
		return
	}
	response.NoContent(c)
}

func (h *BudgetHandler) ListLimits(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("budget"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid budget ID")
		return
	}

	limits, err := h.uc.ListLimits(c.Request.Context(), uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(limits))
	for i, bl := range limits {
		items[i] = transformer.TransformBudgetLimit(&bl)
	}

	pg := pagination.NewMeta(len(limits), len(limits), 1)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

type storeBudgetLimitRequest struct {
	CurrencyID uint   `json:"currency_id" binding:"required"`
	Start      string `json:"start" binding:"required"`
	End        string `json:"end" binding:"required"`
	Amount     string `json:"amount" binding:"required"`
}

func (h *BudgetHandler) StoreLimits(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("budget"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid budget ID")
		return
	}

	var req storeBudgetLimitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	startDate, err := time.Parse("2006-01-02", req.Start)
	if err != nil {
		response.BadRequest(c, "Invalid start date format, expected YYYY-MM-DD")
		return
	}
	endDate, err := time.Parse("2006-01-02", req.End)
	if err != nil {
		response.BadRequest(c, "Invalid end date format, expected YYYY-MM-DD")
		return
	}

	bl := &entity.BudgetLimit{
		BudgetID:              uint(id),
		TransactionCurrencyID: req.CurrencyID,
		StartDate:             startDate,
		EndDate:               endDate,
		Amount:                req.Amount,
	}

	if err := h.uc.CreateLimit(c.Request.Context(), bl); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resource := transformer.TransformBudgetLimit(bl)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

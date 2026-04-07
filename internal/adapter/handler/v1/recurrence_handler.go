package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anthropics/quillow/internal/adapter/transformer"
	"github.com/anthropics/quillow/internal/entity"
	recurrenceuc "github.com/anthropics/quillow/internal/usecase/recurrence"
	"github.com/anthropics/quillow/pkg/pagination"
	"github.com/anthropics/quillow/pkg/response"
	"github.com/gin-gonic/gin"
)

type RecurrenceHandler struct {
	uc *recurrenceuc.UseCase
}

func NewRecurrenceHandler(uc *recurrenceuc.UseCase) *RecurrenceHandler {
	return &RecurrenceHandler{uc: uc}
}

func (h *RecurrenceHandler) Index(c *gin.Context) {
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

	recurrences, total, err := h.uc.List(c.Request.Context(), userGroupID, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(recurrences))
	for i, r := range recurrences {
		reps, _ := h.uc.GetRepetitions(c.Request.Context(), r.ID)
		txns, _ := h.uc.GetTransactions(c.Request.Context(), r.ID)
		items[i] = transformer.TransformRecurrence(&r, reps, txns)
	}

	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *RecurrenceHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid recurrence ID")
		return
	}

	rec, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Recurrence not found")
		return
	}

	reps, _ := h.uc.GetRepetitions(c.Request.Context(), rec.ID)
	txns, _ := h.uc.GetTransactions(c.Request.Context(), rec.ID)
	resource := transformer.TransformRecurrence(rec, reps, txns)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type storeRecurrenceRequest struct {
	Title                 string `json:"title" binding:"required"`
	Description           string `json:"description"`
	FirstDate             string `json:"first_date" binding:"required"`
	RepeatUntil           string `json:"repeat_until"`
	Repetitions           uint   `json:"repetitions"`
	TransactionTypeID     uint   `json:"transaction_type_id"`
	TransactionCurrencyID uint   `json:"transaction_currency_id"`
	ApplyRules            *bool  `json:"apply_rules"`
	Active                *bool  `json:"active"`
}

func (h *RecurrenceHandler) Store(c *gin.Context) {
	var req storeRecurrenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	firstDate, err := time.Parse("2006-01-02", req.FirstDate)
	if err != nil {
		response.BadRequest(c, "Invalid first_date format, expected YYYY-MM-DD")
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}
	applyRules := true
	if req.ApplyRules != nil {
		applyRules = *req.ApplyRules
	}

	rec := &entity.Recurrence{
		Title:                 req.Title,
		Description:           req.Description,
		FirstDate:             firstDate,
		Repetitions:           req.Repetitions,
		TransactionTypeID:     req.TransactionTypeID,
		TransactionCurrencyID: req.TransactionCurrencyID,
		ApplyRules:            applyRules,
		Active:                active,
	}

	if req.RepeatUntil != "" {
		t, err := time.Parse("2006-01-02", req.RepeatUntil)
		if err != nil {
			response.BadRequest(c, "Invalid repeat_until format, expected YYYY-MM-DD")
			return
		}
		rec.RepeatUntil = &t
	}

	if err := h.uc.Create(c.Request.Context(), rec); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	reps, _ := h.uc.GetRepetitions(c.Request.Context(), rec.ID)
	txns, _ := h.uc.GetTransactions(c.Request.Context(), rec.ID)
	resource := transformer.TransformRecurrence(rec, reps, txns)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type updateRecurrenceRequest struct {
	Title                 string `json:"title" binding:"required"`
	Description           string `json:"description"`
	FirstDate             string `json:"first_date"`
	RepeatUntil           string `json:"repeat_until"`
	Repetitions           uint   `json:"repetitions"`
	TransactionTypeID     uint   `json:"transaction_type_id"`
	TransactionCurrencyID uint   `json:"transaction_currency_id"`
	ApplyRules            *bool  `json:"apply_rules"`
	Active                *bool  `json:"active"`
}

func (h *RecurrenceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid recurrence ID")
		return
	}

	rec, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Recurrence not found")
		return
	}

	var req updateRecurrenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rec.Title = req.Title
	rec.Description = req.Description
	rec.Repetitions = req.Repetitions
	rec.TransactionTypeID = req.TransactionTypeID
	rec.TransactionCurrencyID = req.TransactionCurrencyID

	if req.FirstDate != "" {
		t, err := time.Parse("2006-01-02", req.FirstDate)
		if err != nil {
			response.BadRequest(c, "Invalid first_date format, expected YYYY-MM-DD")
			return
		}
		rec.FirstDate = t
	}
	if req.RepeatUntil != "" {
		t, err := time.Parse("2006-01-02", req.RepeatUntil)
		if err != nil {
			response.BadRequest(c, "Invalid repeat_until format, expected YYYY-MM-DD")
			return
		}
		rec.RepeatUntil = &t
	}
	if req.Active != nil {
		rec.Active = *req.Active
	}
	if req.ApplyRules != nil {
		rec.ApplyRules = *req.ApplyRules
	}

	if err := h.uc.Update(c.Request.Context(), rec); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	reps, _ := h.uc.GetRepetitions(c.Request.Context(), rec.ID)
	txns, _ := h.uc.GetTransactions(c.Request.Context(), rec.ID)
	resource := transformer.TransformRecurrence(rec, reps, txns)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *RecurrenceHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid recurrence ID")
		return
	}
	if err := h.uc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "Recurrence not found")
		return
	}
	response.NoContent(c)
}

package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anthropics/firefly-iii-go/internal/adapter/transformer"
	"github.com/anthropics/firefly-iii-go/internal/entity"
	piggybankuc "github.com/anthropics/firefly-iii-go/internal/usecase/piggybank"
	"github.com/anthropics/firefly-iii-go/pkg/pagination"
	"github.com/anthropics/firefly-iii-go/pkg/response"
	"github.com/gin-gonic/gin"
)

type PiggyBankHandler struct {
	uc *piggybankuc.UseCase
}

func NewPiggyBankHandler(uc *piggybankuc.UseCase) *PiggyBankHandler {
	return &PiggyBankHandler{uc: uc}
}

func (h *PiggyBankHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 1000 {
		limit = 50
	}
	offset := (page - 1) * limit

	piggyBanks, total, err := h.uc.List(c.Request.Context(), limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(piggyBanks))
	for i, p := range piggyBanks {
		notes, _ := h.uc.GetNotes(c.Request.Context(), p.ID)
		items[i] = transformer.TransformPiggyBank(&p, notes)
	}

	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *PiggyBankHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("piggy_bank"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid piggy bank ID")
		return
	}

	p, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Piggy bank not found")
		return
	}

	notes, _ := h.uc.GetNotes(c.Request.Context(), p.ID)
	resource := transformer.TransformPiggyBank(p, notes)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type storePiggyBankRequest struct {
	Name         string `json:"name" binding:"required"`
	AccountID    uint   `json:"account_id" binding:"required"`
	TargetAmount string `json:"target_amount"`
	StartDate    string `json:"start_date"`
	TargetDate   string `json:"target_date"`
	Order        uint   `json:"order"`
	Active       *bool  `json:"active"`
	Notes        string `json:"notes"`
}

func (h *PiggyBankHandler) Store(c *gin.Context) {
	var req storePiggyBankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	p := &entity.PiggyBank{
		Name:         req.Name,
		AccountID:    req.AccountID,
		TargetAmount: req.TargetAmount,
		Order:        req.Order,
		Active:       active,
	}

	if req.StartDate != "" {
		t, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			response.BadRequest(c, "Invalid start_date format, expected YYYY-MM-DD")
			return
		}
		p.StartDate = &t
	}
	if req.TargetDate != "" {
		t, err := time.Parse("2006-01-02", req.TargetDate)
		if err != nil {
			response.BadRequest(c, "Invalid target_date format, expected YYYY-MM-DD")
			return
		}
		p.TargetDate = &t
	}

	if err := h.uc.Create(c.Request.Context(), p); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	notes, _ := h.uc.GetNotes(c.Request.Context(), p.ID)
	resource := transformer.TransformPiggyBank(p, notes)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type updatePiggyBankRequest struct {
	Name         string `json:"name" binding:"required"`
	AccountID    uint   `json:"account_id"`
	TargetAmount string `json:"target_amount"`
	StartDate    string `json:"start_date"`
	TargetDate   string `json:"target_date"`
	Order        uint   `json:"order"`
	Active       *bool  `json:"active"`
	Notes        string `json:"notes"`
}

func (h *PiggyBankHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("piggy_bank"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid piggy bank ID")
		return
	}

	p, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Piggy bank not found")
		return
	}

	var req updatePiggyBankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	p.Name = req.Name
	p.TargetAmount = req.TargetAmount
	p.Order = req.Order

	if req.AccountID != 0 {
		p.AccountID = req.AccountID
	}
	if req.Active != nil {
		p.Active = *req.Active
	}

	if req.StartDate != "" {
		t, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			response.BadRequest(c, "Invalid start_date format, expected YYYY-MM-DD")
			return
		}
		p.StartDate = &t
	} else {
		p.StartDate = nil
	}
	if req.TargetDate != "" {
		t, err := time.Parse("2006-01-02", req.TargetDate)
		if err != nil {
			response.BadRequest(c, "Invalid target_date format, expected YYYY-MM-DD")
			return
		}
		p.TargetDate = &t
	} else {
		p.TargetDate = nil
	}

	if err := h.uc.Update(c.Request.Context(), p); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	notes, _ := h.uc.GetNotes(c.Request.Context(), p.ID)
	resource := transformer.TransformPiggyBank(p, notes)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *PiggyBankHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("piggy_bank"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid piggy bank ID")
		return
	}
	if err := h.uc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "Piggy bank not found")
		return
	}
	response.NoContent(c)
}

func (h *PiggyBankHandler) ListEvents(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("piggy_bank"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid piggy bank ID")
		return
	}

	events, err := h.uc.ListEvents(c.Request.Context(), uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(events))
	for i, e := range events {
		items[i] = transformer.TransformPiggyBankEvent(&e)
	}

	pg := pagination.NewMeta(len(events), len(events), 1)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

type addPiggyBankEventRequest struct {
	Amount string `json:"amount" binding:"required"`
	Date   string `json:"date"`
}

func (h *PiggyBankHandler) AddEvent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("piggy_bank"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid piggy bank ID")
		return
	}

	var req addPiggyBankEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	eventDate := time.Now()
	if req.Date != "" {
		d, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			response.BadRequest(c, "Invalid date format, expected YYYY-MM-DD")
			return
		}
		eventDate = d
	}

	event := &entity.PiggyBankEvent{
		PiggyBankID: uint(id),
		Amount:      req.Amount,
		Date:        eventDate,
	}

	if err := h.uc.AddEvent(c.Request.Context(), event); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resource := transformer.TransformPiggyBankEvent(event)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

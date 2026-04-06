package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anthropics/firefly-iii-go/internal/adapter/transformer"
	"github.com/anthropics/firefly-iii-go/internal/entity"
	billuc "github.com/anthropics/firefly-iii-go/internal/usecase/bill"
	"github.com/anthropics/firefly-iii-go/pkg/pagination"
	"github.com/anthropics/firefly-iii-go/pkg/response"
	"github.com/gin-gonic/gin"
)

type BillHandler struct {
	uc *billuc.UseCase
}

func NewBillHandler(uc *billuc.UseCase) *BillHandler {
	return &BillHandler{uc: uc}
}

func (h *BillHandler) Index(c *gin.Context) {
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

	bills, total, err := h.uc.List(c.Request.Context(), userGroupID, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(bills))
	for i, b := range bills {
		notes, _ := h.uc.GetNotes(c.Request.Context(), b.ID)
		items[i] = transformer.TransformBill(&b, notes)
	}

	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *BillHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("bill"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid bill ID")
		return
	}

	b, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Bill not found")
		return
	}

	notes, _ := h.uc.GetNotes(c.Request.Context(), b.ID)
	resource := transformer.TransformBill(b, notes)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type storeBillRequest struct {
	Name                  string `json:"name" binding:"required"`
	AmountMin             string `json:"amount_min" binding:"required"`
	AmountMax             string `json:"amount_max" binding:"required"`
	Date                  string `json:"date" binding:"required"`
	RepeatFreq            string `json:"repeat_freq" binding:"required"`
	Skip                  uint   `json:"skip"`
	Automatch             *bool  `json:"automatch"`
	Active                *bool  `json:"active"`
	Order                 uint   `json:"order"`
	TransactionCurrencyID uint   `json:"transaction_currency_id"`
	EndDate               string `json:"end_date"`
	ExtensionDate         string `json:"extension_date"`
	Notes                 string `json:"notes"`
}

func (h *BillHandler) Store(c *gin.Context) {
	var req storeBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		response.BadRequest(c, "Invalid date format, expected YYYY-MM-DD")
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}
	automatch := true
	if req.Automatch != nil {
		automatch = *req.Automatch
	}

	b := &entity.Bill{
		Name:                  req.Name,
		AmountMin:             req.AmountMin,
		AmountMax:             req.AmountMax,
		Date:                  date,
		RepeatFreq:            req.RepeatFreq,
		Skip:                  req.Skip,
		Automatch:             automatch,
		Active:                active,
		Order:                 req.Order,
		TransactionCurrencyID: req.TransactionCurrencyID,
	}

	if req.EndDate != "" {
		t, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			response.BadRequest(c, "Invalid end_date format, expected YYYY-MM-DD")
			return
		}
		b.EndDate = &t
	}
	if req.ExtensionDate != "" {
		t, err := time.Parse("2006-01-02", req.ExtensionDate)
		if err != nil {
			response.BadRequest(c, "Invalid extension_date format, expected YYYY-MM-DD")
			return
		}
		b.ExtensionDate = &t
	}

	if err := h.uc.Create(c.Request.Context(), b); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	notes, _ := h.uc.GetNotes(c.Request.Context(), b.ID)
	resource := transformer.TransformBill(b, notes)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type updateBillRequest struct {
	Name                  string `json:"name" binding:"required"`
	AmountMin             string `json:"amount_min" binding:"required"`
	AmountMax             string `json:"amount_max" binding:"required"`
	Date                  string `json:"date" binding:"required"`
	RepeatFreq            string `json:"repeat_freq" binding:"required"`
	Skip                  uint   `json:"skip"`
	Automatch             *bool  `json:"automatch"`
	Active                *bool  `json:"active"`
	Order                 uint   `json:"order"`
	TransactionCurrencyID uint   `json:"transaction_currency_id"`
	EndDate               string `json:"end_date"`
	ExtensionDate         string `json:"extension_date"`
	Notes                 string `json:"notes"`
}

func (h *BillHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("bill"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid bill ID")
		return
	}

	b, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Bill not found")
		return
	}

	var req updateBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		response.BadRequest(c, "Invalid date format, expected YYYY-MM-DD")
		return
	}

	b.Name = req.Name
	b.AmountMin = req.AmountMin
	b.AmountMax = req.AmountMax
	b.Date = date
	b.RepeatFreq = req.RepeatFreq
	b.Skip = req.Skip
	b.Order = req.Order
	b.TransactionCurrencyID = req.TransactionCurrencyID

	if req.Active != nil {
		b.Active = *req.Active
	}
	if req.Automatch != nil {
		b.Automatch = *req.Automatch
	}

	if req.EndDate != "" {
		t, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			response.BadRequest(c, "Invalid end_date format, expected YYYY-MM-DD")
			return
		}
		b.EndDate = &t
	}
	if req.ExtensionDate != "" {
		t, err := time.Parse("2006-01-02", req.ExtensionDate)
		if err != nil {
			response.BadRequest(c, "Invalid extension_date format, expected YYYY-MM-DD")
			return
		}
		b.ExtensionDate = &t
	}

	if err := h.uc.Update(c.Request.Context(), b); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	notes, _ := h.uc.GetNotes(c.Request.Context(), b.ID)
	resource := transformer.TransformBill(b, notes)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *BillHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("bill"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid bill ID")
		return
	}
	if err := h.uc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "Bill not found")
		return
	}
	response.NoContent(c)
}

package v1

import (
	"net/http"
	"strconv"

	"github.com/anthropics/quillow/internal/adapter/transformer"
	currencyuc "github.com/anthropics/quillow/internal/usecase/currency"
	"github.com/anthropics/quillow/pkg/pagination"
	"github.com/anthropics/quillow/pkg/response"
	"github.com/gin-gonic/gin"
)

type CurrencyHandler struct {
	uc *currencyuc.UseCase
}

func NewCurrencyHandler(uc *currencyuc.UseCase) *CurrencyHandler {
	return &CurrencyHandler{uc: uc}
}

func (h *CurrencyHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	currencies, total, err := h.uc.List(c.Request.Context(), limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	primary, _ := h.uc.GetPrimary(c.Request.Context())
	primaryID := uint(0)
	if primary != nil {
		primaryID = primary.ID
	}

	items := make([]response.Resource, len(currencies))
	for i, curr := range currencies {
		items[i] = transformer.TransformCurrency(&curr, curr.ID == primaryID)
	}

	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *CurrencyHandler) Show(c *gin.Context) {
	code := c.Param("currency_code")
	curr, err := h.uc.GetByCode(c.Request.Context(), code)
	if err != nil {
		response.NotFound(c, "Currency not found")
		return
	}
	primary, _ := h.uc.GetPrimary(c.Request.Context())
	isPrimary := primary != nil && primary.ID == curr.ID
	resource := transformer.TransformCurrency(curr, isPrimary)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *CurrencyHandler) ShowPrimary(c *gin.Context) {
	primary, err := h.uc.GetPrimary(c.Request.Context())
	if err != nil {
		response.NotFound(c, "No primary currency configured")
		return
	}
	resource := transformer.TransformCurrency(primary, true)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type storeCurrencyRequest struct {
	Name          string `json:"name" binding:"required"`
	Code          string `json:"code" binding:"required,min=3,max=32"`
	Symbol        string `json:"symbol" binding:"required"`
	DecimalPlaces int    `json:"decimal_places"`
	Enabled       *bool  `json:"enabled"`
}

func (h *CurrencyHandler) Store(c *gin.Context) {
	var req storeCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	curr, err := h.uc.Create(c.Request.Context(), req.Name, req.Code, req.Symbol, req.DecimalPlaces, enabled)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resource := transformer.TransformCurrency(curr, false)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type updateCurrencyRequest struct {
	Name          string `json:"name" binding:"required"`
	Code          string `json:"code" binding:"required"`
	Symbol        string `json:"symbol" binding:"required"`
	DecimalPlaces int    `json:"decimal_places"`
	Enabled       *bool  `json:"enabled"`
}

func (h *CurrencyHandler) Update(c *gin.Context) {
	code := c.Param("currency_code")
	curr, err := h.uc.GetByCode(c.Request.Context(), code)
	if err != nil {
		response.NotFound(c, "Currency not found")
		return
	}

	var req updateCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	enabled := curr.Enabled
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	updated, err := h.uc.Update(c.Request.Context(), curr.ID, req.Name, req.Code, req.Symbol, req.DecimalPlaces, enabled)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resource := transformer.TransformCurrency(updated, false)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *CurrencyHandler) Enable(c *gin.Context) {
	code := c.Param("currency_code")
	curr, err := h.uc.Enable(c.Request.Context(), code)
	if err != nil {
		response.NotFound(c, "Currency not found")
		return
	}
	resource := transformer.TransformCurrency(curr, false)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *CurrencyHandler) Disable(c *gin.Context) {
	code := c.Param("currency_code")
	curr, err := h.uc.Disable(c.Request.Context(), code)
	if err != nil {
		response.NotFound(c, "Currency not found")
		return
	}
	resource := transformer.TransformCurrency(curr, false)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *CurrencyHandler) Destroy(c *gin.Context) {
	code := c.Param("currency_code")
	curr, err := h.uc.GetByCode(c.Request.Context(), code)
	if err != nil {
		response.NotFound(c, "Currency not found")
		return
	}
	if err := h.uc.Delete(c.Request.Context(), curr.ID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.NoContent(c)
}

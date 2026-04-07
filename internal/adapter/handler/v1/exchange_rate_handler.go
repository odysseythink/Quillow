package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anthropics/quillow/internal/adapter/transformer"
	"github.com/anthropics/quillow/internal/entity"
	currencyuc "github.com/anthropics/quillow/internal/usecase/currency"
	eruc "github.com/anthropics/quillow/internal/usecase/exchangerate"
	"github.com/anthropics/quillow/pkg/pagination"
	"github.com/anthropics/quillow/pkg/response"
	"github.com/gin-gonic/gin"
)

type ExchangeRateHandler struct {
	uc     *eruc.UseCase
	currUC *currencyuc.UseCase
}

func NewExchangeRateHandler(uc *eruc.UseCase, currUC *currencyuc.UseCase) *ExchangeRateHandler {
	return &ExchangeRateHandler{uc: uc, currUC: currUC}
}

func (h *ExchangeRateHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	rates, total, err := h.uc.List(c.Request.Context(), 0, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(rates))
	for i, rate := range rates {
		extra := h.enrichRate(c, &rate)
		items[i] = transformer.TransformExchangeRate(&rate, extra)
	}

	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *ExchangeRateHandler) ShowByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid exchange rate ID")
		return
	}
	rate, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Exchange rate not found")
		return
	}
	extra := h.enrichRate(c, rate)
	resource := transformer.TransformExchangeRate(rate, extra)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *ExchangeRateHandler) ShowByPair(c *gin.Context) {
	fromCode := c.Param("fromCode")
	toCode := c.Param("toCode")

	rates, err := h.uc.ListByPair(c.Request.Context(), fromCode, toCode)
	if err != nil || len(rates) == 0 {
		response.NotFound(c, "Exchange rate not found")
		return
	}

	items := make([]response.Resource, len(rates))
	for i, rate := range rates {
		extra := h.enrichRate(c, &rate)
		items[i] = transformer.TransformExchangeRate(&rate, extra)
	}
	pg := pagination.NewMeta(len(rates), len(rates), 1)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *ExchangeRateHandler) ShowByDate(c *gin.Context) {
	fromCode := c.Param("fromCode")
	toCode := c.Param("toCode")
	dateStr := c.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.BadRequest(c, "Invalid date format, use YYYY-MM-DD")
		return
	}

	rate, err := h.uc.GetByPairAndDate(c.Request.Context(), fromCode, toCode, date)
	if err != nil {
		response.NotFound(c, "Exchange rate not found")
		return
	}
	extra := h.enrichRate(c, rate)
	resource := transformer.TransformExchangeRate(rate, extra)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type storeExchangeRateRequest struct {
	FromCurrencyID uint   `json:"from_currency_id"`
	ToCurrencyID   uint   `json:"to_currency_id"`
	FromCode       string `json:"from_currency_code"`
	ToCode         string `json:"to_currency_code"`
	Rate           string `json:"rate" binding:"required"`
	Date           string `json:"date"`
}

func (h *ExchangeRateHandler) Store(c *gin.Context) {
	var req storeExchangeRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	fromID := req.FromCurrencyID
	toID := req.ToCurrencyID

	if fromID == 0 && req.FromCode != "" {
		if curr, err := h.currUC.GetByCode(c.Request.Context(), req.FromCode); err == nil {
			fromID = curr.ID
		}
	}
	if toID == 0 && req.ToCode != "" {
		if curr, err := h.currUC.GetByCode(c.Request.Context(), req.ToCode); err == nil {
			toID = curr.ID
		}
	}

	date := time.Now()
	if req.Date != "" {
		if d, err := time.Parse("2006-01-02", req.Date); err == nil {
			date = d
		}
	}

	rate, err := h.uc.Create(c.Request.Context(), 0, fromID, toID, req.Rate, date)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	extra := h.enrichRate(c, rate)
	resource := transformer.TransformExchangeRate(rate, extra)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type updateExchangeRateRequest struct {
	Rate string `json:"rate" binding:"required"`
}

func (h *ExchangeRateHandler) UpdateByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid exchange rate ID")
		return
	}
	var req updateExchangeRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rate, err := h.uc.Update(c.Request.Context(), uint(id), req.Rate)
	if err != nil {
		response.NotFound(c, "Exchange rate not found")
		return
	}
	extra := h.enrichRate(c, rate)
	resource := transformer.TransformExchangeRate(rate, extra)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *ExchangeRateHandler) DestroyByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid exchange rate ID")
		return
	}
	if err := h.uc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "Exchange rate not found")
		return
	}
	response.NoContent(c)
}

func (h *ExchangeRateHandler) DestroyByPair(c *gin.Context) {
	fromCode := c.Param("fromCode")
	toCode := c.Param("toCode")
	if err := h.uc.DeleteByPair(c.Request.Context(), fromCode, toCode); err != nil {
		response.NotFound(c, "Exchange rate not found")
		return
	}
	response.NoContent(c)
}

func (h *ExchangeRateHandler) enrichRate(c *gin.Context, rate *entity.CurrencyExchangeRate) transformer.ExchangeRateExtra {
	extra := transformer.ExchangeRateExtra{}
	if from, err := h.currUC.GetByID(c.Request.Context(), rate.FromCurrencyID); err == nil {
		extra.FromCurrency = from
	}
	if to, err := h.currUC.GetByID(c.Request.Context(), rate.ToCurrencyID); err == nil {
		extra.ToCurrency = to
	}
	return extra
}

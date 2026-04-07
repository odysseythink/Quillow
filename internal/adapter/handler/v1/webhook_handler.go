package v1

import (
	"net/http"
	"strconv"

	"github.com/anthropics/quillow/internal/adapter/transformer"
	"github.com/anthropics/quillow/internal/entity"
	webhookuc "github.com/anthropics/quillow/internal/usecase/webhook"
	"github.com/anthropics/quillow/pkg/pagination"
	"github.com/anthropics/quillow/pkg/response"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	uc *webhookuc.UseCase
}

func NewWebhookHandler(uc *webhookuc.UseCase) *WebhookHandler {
	return &WebhookHandler{uc: uc}
}

func (h *WebhookHandler) Index(c *gin.Context) {
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

	webhooks, total, err := h.uc.List(c.Request.Context(), userGroupID, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(webhooks))
	for i, w := range webhooks {
		items[i] = transformer.TransformWebhook(&w)
	}

	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *WebhookHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid webhook ID")
		return
	}

	w, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Webhook not found")
		return
	}

	resource := transformer.TransformWebhook(w)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type storeWebhookRequest struct {
	Active   *bool  `json:"active"`
	Title    string `json:"title" binding:"required"`
	Trigger  int    `json:"trigger" binding:"required"`
	Response int    `json:"response" binding:"required"`
	Delivery int    `json:"delivery" binding:"required"`
	URL      string `json:"url" binding:"required"`
}

func (h *WebhookHandler) Store(c *gin.Context) {
	var req storeWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	w := &entity.Webhook{
		Active:   active,
		Title:    req.Title,
		Trigger:  req.Trigger,
		Response: req.Response,
		Delivery: req.Delivery,
		URL:      req.URL,
	}

	if err := h.uc.Create(c.Request.Context(), w); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resource := transformer.TransformWebhook(w)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type updateWebhookRequest struct {
	Active   *bool  `json:"active"`
	Title    string `json:"title" binding:"required"`
	Trigger  int    `json:"trigger"`
	Response int    `json:"response"`
	Delivery int    `json:"delivery"`
	URL      string `json:"url"`
}

func (h *WebhookHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid webhook ID")
		return
	}

	w, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Webhook not found")
		return
	}

	var req updateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	w.Title = req.Title
	if req.Active != nil {
		w.Active = *req.Active
	}
	if req.Trigger != 0 {
		w.Trigger = req.Trigger
	}
	if req.Response != 0 {
		w.Response = req.Response
	}
	if req.Delivery != 0 {
		w.Delivery = req.Delivery
	}
	if req.URL != "" {
		w.URL = req.URL
	}

	if err := h.uc.Update(c.Request.Context(), w); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resource := transformer.TransformWebhook(w)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *WebhookHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid webhook ID")
		return
	}
	if err := h.uc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "Webhook not found")
		return
	}
	response.NoContent(c)
}

func (h *WebhookHandler) ListMessages(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid webhook ID")
		return
	}

	messages, err := h.uc.ListMessages(c.Request.Context(), uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(messages))
	for i, m := range messages {
		items[i] = transformer.TransformWebhookMessage(&m)
	}

	pg := pagination.NewMeta(len(messages), len(messages), 1)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *WebhookHandler) ListAttempts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid message ID")
		return
	}

	attempts, err := h.uc.ListAttempts(c.Request.Context(), uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(attempts))
	for i, a := range attempts {
		items[i] = transformer.TransformWebhookAttempt(&a)
	}

	pg := pagination.NewMeta(len(attempts), len(attempts), 1)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

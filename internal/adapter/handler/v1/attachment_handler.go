package v1

import (
	"net/http"
	"strconv"

	"github.com/anthropics/quillow/internal/adapter/transformer"
	"github.com/anthropics/quillow/internal/entity"
	"github.com/anthropics/quillow/internal/port"
	"github.com/anthropics/quillow/pkg/pagination"
	"github.com/anthropics/quillow/pkg/response"
	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	repo port.AttachmentRepository
}

func NewAttachmentHandler(repo port.AttachmentRepository) *AttachmentHandler {
	return &AttachmentHandler{repo: repo}
}

func (h *AttachmentHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	atts, total, err := h.repo.List(c.Request.Context(), limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(atts))
	for i, att := range atts {
		items[i] = transformer.TransformAttachment(&att)
	}
	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *AttachmentHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid attachment ID")
		return
	}
	att, err := h.repo.FindByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Attachment not found")
		return
	}
	resource := transformer.TransformAttachment(att)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type storeAttachmentRequest struct {
	Filename       string `json:"filename" binding:"required"`
	AttachableType string `json:"attachable_type" binding:"required"`
	AttachableID   uint   `json:"attachable_id" binding:"required"`
	Title          string `json:"title"`
	Notes          string `json:"notes"`
}

func (h *AttachmentHandler) Store(c *gin.Context) {
	var req storeAttachmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	att := &entity.Attachment{
		UserID:         userID,
		AttachableType: req.AttachableType,
		AttachableID:   req.AttachableID,
		Filename:       req.Filename,
		Title:          req.Title,
		Description:    req.Notes,
		Mime:           "application/octet-stream",
	}
	if err := h.repo.Create(c.Request.Context(), att); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resource := transformer.TransformAttachment(att)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *AttachmentHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid attachment ID")
		return
	}
	if err := h.repo.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "Attachment not found")
		return
	}
	response.NoContent(c)
}

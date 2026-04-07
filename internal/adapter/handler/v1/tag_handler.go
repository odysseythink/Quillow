package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/anthropics/quillow/internal/adapter/transformer"
	"github.com/anthropics/quillow/internal/entity"
	taguc "github.com/anthropics/quillow/internal/usecase/tag"
	"github.com/anthropics/quillow/pkg/pagination"
	"github.com/anthropics/quillow/pkg/response"
	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	uc *taguc.UseCase
}

func NewTagHandler(uc *taguc.UseCase) *TagHandler {
	return &TagHandler{uc: uc}
}

func (h *TagHandler) Index(c *gin.Context) {
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

	tags, total, err := h.uc.List(c.Request.Context(), userGroupID, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(tags))
	for i, t := range tags {
		items[i] = transformer.TransformTag(&t)
	}

	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *TagHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("tag"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid tag ID")
		return
	}

	t, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Tag not found")
		return
	}

	resource := transformer.TransformTag(t)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type storeTagRequest struct {
	Tag         string   `json:"tag" binding:"required"`
	Description string   `json:"description"`
	Date        string   `json:"date"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	ZoomLevel   *int     `json:"zoom_level"`
}

func (h *TagHandler) Store(c *gin.Context) {
	var req storeTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	t := &entity.Tag{
		Tag:         req.Tag,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		ZoomLevel:   req.ZoomLevel,
	}

	if req.Date != "" {
		d, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			response.BadRequest(c, "Invalid date format, expected YYYY-MM-DD")
			return
		}
		t.Date = &d
	}

	if err := h.uc.Create(c.Request.Context(), t); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resource := transformer.TransformTag(t)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type updateTagRequest struct {
	Tag         string   `json:"tag" binding:"required"`
	Description string   `json:"description"`
	Date        string   `json:"date"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	ZoomLevel   *int     `json:"zoom_level"`
}

func (h *TagHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("tag"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid tag ID")
		return
	}

	t, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Tag not found")
		return
	}

	var req updateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	t.Tag = req.Tag
	t.Description = req.Description
	t.Latitude = req.Latitude
	t.Longitude = req.Longitude
	t.ZoomLevel = req.ZoomLevel

	if req.Date != "" {
		d, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			response.BadRequest(c, "Invalid date format, expected YYYY-MM-DD")
			return
		}
		t.Date = &d
	} else {
		t.Date = nil
	}

	if err := h.uc.Update(c.Request.Context(), t); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resource := transformer.TransformTag(t)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *TagHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("tag"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid tag ID")
		return
	}
	if err := h.uc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "Tag not found")
		return
	}
	response.NoContent(c)
}

func (h *TagHandler) Search(c *gin.Context) {
	query := c.Query("query")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	userGroupID := uint(0) // TODO: get from user context in later SP

	tags, err := h.uc.Search(c.Request.Context(), userGroupID, query, limit)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(tags))
	for i, t := range tags {
		items[i] = transformer.TransformTag(&t)
	}

	pg := pagination.NewMeta(len(tags), limit, 1)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *TagHandler) Autocomplete(c *gin.Context) {
	query := c.Query("query")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	userGroupID := uint(0) // TODO: get from user context in later SP

	tags, err := h.uc.Search(c.Request.Context(), userGroupID, query, limit)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	type autocompleteItem struct {
		ID  string `json:"id"`
		Tag string `json:"tag"`
	}
	result := make([]autocompleteItem, len(tags))
	for i, t := range tags {
		result[i] = autocompleteItem{
			ID:  fmt.Sprintf("%d", t.ID),
			Tag: t.Tag,
		}
	}

	c.JSON(http.StatusOK, result)
}

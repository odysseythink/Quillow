package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/anthropics/quillow/internal/adapter/transformer"
	"github.com/anthropics/quillow/internal/entity"
	categoryuc "github.com/anthropics/quillow/internal/usecase/category"
	"github.com/anthropics/quillow/pkg/pagination"
	"github.com/anthropics/quillow/pkg/response"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	uc *categoryuc.UseCase
}

func NewCategoryHandler(uc *categoryuc.UseCase) *CategoryHandler {
	return &CategoryHandler{uc: uc}
}

func (h *CategoryHandler) Index(c *gin.Context) {
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

	categories, total, err := h.uc.List(c.Request.Context(), userGroupID, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(categories))
	for i, cat := range categories {
		notes, _ := h.uc.GetNotes(c.Request.Context(), cat.ID)
		items[i] = transformer.TransformCategory(&cat, notes)
	}

	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *CategoryHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("category"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid category ID")
		return
	}

	cat, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Category not found")
		return
	}

	notes, _ := h.uc.GetNotes(c.Request.Context(), cat.ID)
	resource := transformer.TransformCategory(cat, notes)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type storeCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Notes string `json:"notes"`
}

func (h *CategoryHandler) Store(c *gin.Context) {
	var req storeCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cat := &entity.Category{
		Name: req.Name,
	}

	if err := h.uc.Create(c.Request.Context(), cat); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	notes, _ := h.uc.GetNotes(c.Request.Context(), cat.ID)
	resource := transformer.TransformCategory(cat, notes)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type updateCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Notes string `json:"notes"`
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("category"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid category ID")
		return
	}

	cat, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Category not found")
		return
	}

	var req updateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cat.Name = req.Name

	if err := h.uc.Update(c.Request.Context(), cat); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	notes, _ := h.uc.GetNotes(c.Request.Context(), cat.ID)
	resource := transformer.TransformCategory(cat, notes)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *CategoryHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("category"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid category ID")
		return
	}
	if err := h.uc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "Category not found")
		return
	}
	response.NoContent(c)
}

func (h *CategoryHandler) Search(c *gin.Context) {
	query := c.Query("query")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	userGroupID := uint(0) // TODO: get from user context in later SP

	categories, err := h.uc.Search(c.Request.Context(), userGroupID, query, limit)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(categories))
	for i, cat := range categories {
		notes, _ := h.uc.GetNotes(c.Request.Context(), cat.ID)
		items[i] = transformer.TransformCategory(&cat, notes)
	}

	pg := pagination.NewMeta(len(categories), limit, 1)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *CategoryHandler) Autocomplete(c *gin.Context) {
	query := c.Query("query")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	userGroupID := uint(0) // TODO: get from user context in later SP

	categories, err := h.uc.Search(c.Request.Context(), userGroupID, query, limit)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	type autocompleteItem struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	result := make([]autocompleteItem, len(categories))
	for i, cat := range categories {
		result[i] = autocompleteItem{
			ID:   fmt.Sprintf("%d", cat.ID),
			Name: cat.Name,
		}
	}

	c.JSON(http.StatusOK, result)
}

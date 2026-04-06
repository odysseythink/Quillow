package v1

import (
	"net/http"
	"strconv"

	"github.com/anthropics/firefly-iii-go/internal/adapter/transformer"
	"github.com/anthropics/firefly-iii-go/internal/entity"
	objectgroupuc "github.com/anthropics/firefly-iii-go/internal/usecase/objectgroup"
	"github.com/anthropics/firefly-iii-go/pkg/pagination"
	"github.com/anthropics/firefly-iii-go/pkg/response"
	"github.com/gin-gonic/gin"
)

type ObjectGroupHandler struct {
	uc *objectgroupuc.UseCase
}

func NewObjectGroupHandler(uc *objectgroupuc.UseCase) *ObjectGroupHandler {
	return &ObjectGroupHandler{uc: uc}
}

func (h *ObjectGroupHandler) Index(c *gin.Context) {
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

	objectGroups, total, err := h.uc.List(c.Request.Context(), userGroupID, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(objectGroups))
	for i, og := range objectGroups {
		items[i] = transformer.TransformObjectGroup(&og)
	}

	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *ObjectGroupHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("object_group"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid object group ID")
		return
	}

	og, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Object group not found")
		return
	}

	resource := transformer.TransformObjectGroup(og)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type storeObjectGroupRequest struct {
	Title string `json:"title" binding:"required"`
	Order uint   `json:"order"`
}

func (h *ObjectGroupHandler) Store(c *gin.Context) {
	var req storeObjectGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	og := &entity.ObjectGroup{
		Title: req.Title,
		Order: req.Order,
	}

	if err := h.uc.Create(c.Request.Context(), og); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resource := transformer.TransformObjectGroup(og)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type updateObjectGroupRequest struct {
	Title string `json:"title" binding:"required"`
	Order uint   `json:"order"`
}

func (h *ObjectGroupHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("object_group"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid object group ID")
		return
	}

	og, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Object group not found")
		return
	}

	var req updateObjectGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	og.Title = req.Title
	og.Order = req.Order

	if err := h.uc.Update(c.Request.Context(), og); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resource := transformer.TransformObjectGroup(og)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *ObjectGroupHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("object_group"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid object group ID")
		return
	}
	if err := h.uc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "Object group not found")
		return
	}
	response.NoContent(c)
}

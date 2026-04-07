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

type LinkTypeHandler struct {
	ltRepo   port.LinkTypeRepository
	linkRepo port.TransactionLinkRepository
}

func NewLinkTypeHandler(ltRepo port.LinkTypeRepository, linkRepo port.TransactionLinkRepository) *LinkTypeHandler {
	return &LinkTypeHandler{ltRepo: ltRepo, linkRepo: linkRepo}
}

func (h *LinkTypeHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	items, total, err := h.ltRepo.List(c.Request.Context(), limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resources := make([]response.Resource, len(items))
	for i, lt := range items {
		resources[i] = transformer.TransformLinkType(&lt)
	}
	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(resources, pg))
}

func (h *LinkTypeHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid link type ID")
		return
	}
	lt, err := h.ltRepo.FindByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Link type not found")
		return
	}
	resource := transformer.TransformLinkType(lt)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type storeLinkTypeRequest struct {
	Name    string `json:"name" binding:"required"`
	Inward  string `json:"inward" binding:"required"`
	Outward string `json:"outward" binding:"required"`
}

func (h *LinkTypeHandler) Store(c *gin.Context) {
	var req storeLinkTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	lt := &entity.LinkType{Name: req.Name, Inward: req.Inward, Outward: req.Outward, Editable: true}
	if err := h.ltRepo.Create(c.Request.Context(), lt); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resource := transformer.TransformLinkType(lt)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *LinkTypeHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid link type ID")
		return
	}
	if err := h.ltRepo.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "Link type not found")
		return
	}
	response.NoContent(c)
}

// Transaction Links
func (h *LinkTypeHandler) ListLinks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	links, total, err := h.linkRepo.List(c.Request.Context(), limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(links))
	for i, link := range links {
		items[i] = transformer.TransformTransactionLink(&link)
	}
	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *LinkTypeHandler) ShowLink(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid link ID")
		return
	}
	link, err := h.linkRepo.FindByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Transaction link not found")
		return
	}
	resource := transformer.TransformTransactionLink(link)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type storeTransactionLinkRequest struct {
	LinkTypeID    uint   `json:"link_type_id" binding:"required"`
	InwardID      uint   `json:"inward_id" binding:"required"`
	OutwardID     uint   `json:"outward_id" binding:"required"`
	Notes         string `json:"notes"`
}

func (h *LinkTypeHandler) StoreLink(c *gin.Context) {
	var req storeTransactionLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	link := &entity.TransactionJournalLink{
		LinkTypeID:    req.LinkTypeID,
		SourceID:      req.InwardID,
		DestinationID: req.OutwardID,
		Comment:       req.Notes,
	}
	if err := h.linkRepo.Create(c.Request.Context(), link); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resource := transformer.TransformTransactionLink(link)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *LinkTypeHandler) DestroyLink(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid link ID")
		return
	}
	if err := h.linkRepo.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "Transaction link not found")
		return
	}
	response.NoContent(c)
}

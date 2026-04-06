package v1

import (
	"net/http"
	"strconv"

	"github.com/anthropics/firefly-iii-go/internal/adapter/transformer"
	prefuc "github.com/anthropics/firefly-iii-go/internal/usecase/preference"
	"github.com/anthropics/firefly-iii-go/pkg/pagination"
	"github.com/anthropics/firefly-iii-go/pkg/response"
	"github.com/gin-gonic/gin"
)

type PreferenceHandler struct {
	uc *prefuc.UseCase
}

func NewPreferenceHandler(uc *prefuc.UseCase) *PreferenceHandler {
	return &PreferenceHandler{uc: uc}
}

func (h *PreferenceHandler) Index(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	prefs, total, err := h.uc.List(c.Request.Context(), userID, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(prefs))
	for i, p := range prefs {
		items[i] = transformer.TransformPreference(&p)
	}
	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *PreferenceHandler) Show(c *gin.Context) {
	userID := c.GetUint("user_id")
	name := c.Param("name")
	pref, err := h.uc.GetByName(c.Request.Context(), userID, name)
	if err != nil {
		response.NotFound(c, "Preference not found")
		return
	}
	resource := transformer.TransformPreference(pref)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type storePreferenceRequest struct {
	Name string `json:"name" binding:"required"`
	Data string `json:"data" binding:"required"`
}

func (h *PreferenceHandler) Store(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req storePreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	pref, err := h.uc.Store(c.Request.Context(), userID, req.Name, req.Data)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resource := transformer.TransformPreference(pref)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type updatePreferenceRequest struct {
	Data string `json:"data" binding:"required"`
}

func (h *PreferenceHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	name := c.Param("name")
	var req updatePreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	pref, err := h.uc.Update(c.Request.Context(), userID, name, req.Data)
	if err != nil {
		response.NotFound(c, "Preference not found")
		return
	}
	resource := transformer.TransformPreference(pref)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

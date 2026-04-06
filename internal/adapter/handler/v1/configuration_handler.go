package v1

import (
	"net/http"

	configuc "github.com/anthropics/firefly-iii-go/internal/usecase/configuration"
	"github.com/anthropics/firefly-iii-go/pkg/response"
	"github.com/gin-gonic/gin"
)

type ConfigurationHandler struct {
	uc *configuc.UseCase
}

func NewConfigurationHandler(uc *configuc.UseCase) *ConfigurationHandler {
	return &ConfigurationHandler{uc: uc}
}

func (h *ConfigurationHandler) Index(c *gin.Context) {
	configs, err := h.uc.List(c.Request.Context())
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, configs)
}

func (h *ConfigurationHandler) Show(c *gin.Context) {
	key := c.Param("key")
	config, err := h.uc.GetByName(c.Request.Context(), key)
	if err != nil {
		response.NotFound(c, "Configuration key not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": config})
}

type updateConfigRequest struct {
	Value string `json:"value" binding:"required"`
}

func (h *ConfigurationHandler) Update(c *gin.Context) {
	key := c.Param("key")
	var req updateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	config, err := h.uc.Update(c.Request.Context(), key, req.Value)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": config})
}

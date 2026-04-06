package v1

import (
	"net/http"
	"runtime"

	"github.com/anthropics/firefly-iii-go/internal/adapter/transformer"
	useruc "github.com/anthropics/firefly-iii-go/internal/usecase/user"
	"github.com/anthropics/firefly-iii-go/pkg/config"
	"github.com/anthropics/firefly-iii-go/pkg/response"
	"github.com/gin-gonic/gin"
)

type AboutHandler struct {
	cfg    *config.Config
	userUC *useruc.UseCase
	driver string
}

func NewAboutHandler(cfg *config.Config, userUC *useruc.UseCase, driver string) *AboutHandler {
	return &AboutHandler{cfg: cfg, userUC: userUC, driver: driver}
}

func (h *AboutHandler) About(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":     h.cfg.App.Version,
		"api_version": h.cfg.App.APIVersion,
		"php_version": "N/A (Go backend)",
		"os":          runtime.GOOS + "/" + runtime.GOARCH,
		"driver":      h.driver,
	})
}

func (h *AboutHandler) User(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, role, err := h.userUC.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	resource := transformer.TransformUser(user, role)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

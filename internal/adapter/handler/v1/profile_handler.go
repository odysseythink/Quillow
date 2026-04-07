package v1

import (
	"net/http"

	"github.com/anthropics/quillow/internal/adapter/transformer"
	useruc "github.com/anthropics/quillow/internal/usecase/user"
	"github.com/anthropics/quillow/pkg/response"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	userUC *useruc.UseCase
}

func NewProfileHandler(userUC *useruc.UseCase) *ProfileHandler {
	return &ProfileHandler{userUC: userUC}
}

func (h *ProfileHandler) Show(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, role, err := h.userUC.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}
	resource := transformer.TransformUser(user, role)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.userUC.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

type changeEmailRequest struct {
	Password string `json:"password" binding:"required"`
	NewEmail string `json:"new_email" binding:"required,email"`
}

func (h *ProfileHandler) ChangeEmail(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req changeEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	user, err := h.userUC.ChangeEmail(c.Request.Context(), userID, req.Password, req.NewEmail)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resource := transformer.TransformUser(user, "")
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

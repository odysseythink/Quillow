package v1

import (
	"net/http"
	"strconv"

	"github.com/anthropics/firefly-iii-go/internal/adapter/transformer"
	useruc "github.com/anthropics/firefly-iii-go/internal/usecase/user"
	"github.com/anthropics/firefly-iii-go/pkg/pagination"
	"github.com/anthropics/firefly-iii-go/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	uc *useruc.UseCase
}

func NewUserHandler(uc *useruc.UseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 1000 {
		limit = 50
	}
	offset := (page - 1) * limit

	users, total, err := h.uc.List(c.Request.Context(), limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(users))
	for i, u := range users {
		items[i] = transformer.TransformUser(&u, "")
	}
	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *UserHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}
	user, role, err := h.uc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}
	resource := transformer.TransformUser(user, role)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type createUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	Blocked     bool   `json:"blocked"`
	BlockedCode string `json:"blocked_code"`
	Role        string `json:"role"`
}

func (h *UserHandler) Store(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	user, err := h.uc.Create(c.Request.Context(), req.Email, req.Password, req.Blocked, req.BlockedCode, req.Role)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resource := transformer.TransformUser(user, req.Role)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type updateUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Blocked     bool   `json:"blocked"`
	BlockedCode string `json:"blocked_code"`
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	user, err := h.uc.Update(c.Request.Context(), uint(id), req.Email, req.Blocked, req.BlockedCode)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}
	resource := transformer.TransformUser(user, "")
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *UserHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}
	currentUserID := c.GetUint("user_id")
	if uint(id) == currentUserID {
		response.BadRequest(c, "You cannot delete your own account")
		return
	}
	if err := h.uc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "User not found")
		return
	}
	response.NoContent(c)
}

package middleware

import (
	"github.com/anthropics/quillow/internal/port"
	"github.com/anthropics/quillow/pkg/response"
	"github.com/gin-gonic/gin"
)

func Admin(userRepo port.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		role, err := userRepo.GetRole(c.Request.Context(), userID)
		if err != nil || role != "owner" {
			response.Forbidden(c, "Admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}

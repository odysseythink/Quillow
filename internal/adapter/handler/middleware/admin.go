package middleware

import (
	"github.com/anthropics/firefly-iii-go/internal/port"
	"github.com/anthropics/firefly-iii-go/pkg/response"
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

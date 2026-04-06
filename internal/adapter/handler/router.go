package handler

import (
	"github.com/anthropics/firefly-iii-go/internal/adapter/handler/middleware"
	v1 "github.com/anthropics/firefly-iii-go/internal/adapter/handler/v1"
	"github.com/anthropics/firefly-iii-go/internal/port"
	"github.com/anthropics/firefly-iii-go/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth          *v1.AuthHandler
	About         *v1.AboutHandler
	User          *v1.UserHandler
	Configuration *v1.ConfigurationHandler
	Preference    *v1.PreferenceHandler
}

func SetupRouter(
	r *gin.Engine,
	h Handlers,
	jwtSvc *jwt.Service,
	userRepo port.UserRepository,
) {
	r.Use(middleware.CORS())
	r.Use(middleware.I18n())

	api := r.Group("/api/v1")

	// Public routes
	api.POST("/auth/login", h.Auth.Login)
	api.POST("/auth/refresh", h.Auth.Refresh)
	api.GET("/about", h.About.About)

	// Public config read
	api.GET("/configuration", h.Configuration.Index)
	api.GET("/configuration/:key", h.Configuration.Show)

	// Authenticated routes
	auth := api.Group("")
	auth.Use(middleware.Auth(jwtSvc))
	{
		auth.POST("/auth/logout", h.Auth.Logout)
		auth.GET("/about/user", h.About.User)

		// Preferences
		auth.GET("/preferences", h.Preference.Index)
		auth.POST("/preferences", h.Preference.Store)
		auth.GET("/preferences/:name", h.Preference.Show)
		auth.PUT("/preferences/:name", h.Preference.Update)
	}

	// Admin routes
	admin := auth.Group("")
	admin.Use(middleware.Admin(userRepo))
	{
		admin.GET("/users", h.User.Index)
		admin.POST("/users", h.User.Store)
		admin.GET("/users/:id", h.User.Show)
		admin.PUT("/users/:id", h.User.Update)
		admin.DELETE("/users/:id", h.User.Destroy)

		admin.PUT("/configuration/:key", h.Configuration.Update)
	}
}

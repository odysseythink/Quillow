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
	Account       *v1.AccountHandler
	Currency      *v1.CurrencyHandler
	ExchangeRate  *v1.ExchangeRateHandler
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

		// Accounts
		auth.GET("/accounts", h.Account.Index)
		auth.POST("/accounts", h.Account.Store)
		auth.GET("/accounts/:account", h.Account.Show)
		auth.PUT("/accounts/:account", h.Account.Update)
		auth.DELETE("/accounts/:account", h.Account.Destroy)

		// Autocomplete
		auth.GET("/autocomplete/accounts", h.Account.Autocomplete)

		// Currencies
		auth.GET("/currencies", h.Currency.Index)
		auth.GET("/currencies/primary", h.Currency.ShowPrimary)
		auth.GET("/currencies/default", h.Currency.ShowPrimary)
		auth.GET("/currencies/:currency_code", h.Currency.Show)
		auth.PUT("/currencies/:currency_code", h.Currency.Update)
		auth.POST("/currencies/:currency_code/enable", h.Currency.Enable)
		auth.POST("/currencies/:currency_code/disable", h.Currency.Disable)

		// Exchange Rates
		auth.GET("/exchange-rates", h.ExchangeRate.Index)
		auth.GET("/exchange-rates/:id", h.ExchangeRate.ShowByID)
		auth.POST("/exchange-rates", h.ExchangeRate.Store)
		auth.PUT("/exchange-rates/:id", h.ExchangeRate.UpdateByID)
		auth.DELETE("/exchange-rates/:id", h.ExchangeRate.DestroyByID)
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

		// Admin-only currency operations
		admin.POST("/currencies", h.Currency.Store)
		admin.DELETE("/currencies/:currency_code", h.Currency.Destroy)
	}
}

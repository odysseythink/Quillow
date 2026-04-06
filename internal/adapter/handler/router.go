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
	Transaction   *v1.TransactionHandler
	Attachment    *v1.AttachmentHandler
	LinkType      *v1.LinkTypeHandler
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

	// Public link types
	api.GET("/link-types", h.LinkType.Index)
	api.GET("/link-types/:id", h.LinkType.Show)

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

		// Transactions
		auth.GET("/transactions", h.Transaction.Index)
		auth.POST("/transactions", h.Transaction.Store)
		auth.GET("/transactions/:id", h.Transaction.Show)
		auth.DELETE("/transactions/:id", h.Transaction.Destroy)

		// Search
		auth.GET("/search/transactions", h.Transaction.Search)
		auth.GET("/search/transactions/count", h.Transaction.SearchCount)

		// Attachments
		auth.GET("/attachments", h.Attachment.Index)
		auth.POST("/attachments", h.Attachment.Store)
		auth.GET("/attachments/:id", h.Attachment.Show)
		auth.DELETE("/attachments/:id", h.Attachment.Destroy)

		// Transaction Links
		auth.GET("/transaction-links", h.LinkType.ListLinks)
		auth.POST("/transaction-links", h.LinkType.StoreLink)
		auth.GET("/transaction-links/:id", h.LinkType.ShowLink)
		auth.DELETE("/transaction-links/:id", h.LinkType.DestroyLink)

		// Insight
		auth.GET("/insight/expense/total", h.Transaction.InsightExpense)
		auth.GET("/insight/income/total", h.Transaction.InsightIncome)
		auth.GET("/insight/transfer/total", h.Transaction.InsightTransfer)

		// Summary
		auth.GET("/summary/basic", h.Transaction.Summary)
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

		// Admin-only link type operations
		admin.POST("/link-types", h.LinkType.Store)
		admin.DELETE("/link-types/:id", h.LinkType.Destroy)
	}
}

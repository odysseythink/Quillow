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
	Budget        *v1.BudgetHandler
	Bill          *v1.BillHandler
	Category      *v1.CategoryHandler
	Tag           *v1.TagHandler
	PiggyBank     *v1.PiggyBankHandler
	ObjectGroup   *v1.ObjectGroupHandler
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
	api.GET("/configuration", h.Configuration.Index)
	api.GET("/configuration/:key", h.Configuration.Show)
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
		auth.GET("/autocomplete/categories", h.Category.Autocomplete)
		auth.GET("/autocomplete/tags", h.Tag.Autocomplete)

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

		// Budgets
		auth.GET("/budgets", h.Budget.Index)
		auth.POST("/budgets", h.Budget.Store)
		auth.GET("/budgets/:id", h.Budget.Show)
		auth.PUT("/budgets/:id", h.Budget.Update)
		auth.DELETE("/budgets/:id", h.Budget.Destroy)
		auth.GET("/budgets/:id/limits", h.Budget.ListLimits)
		auth.POST("/budgets/:id/limits", h.Budget.StoreLimits)

		// Bills / Subscriptions
		auth.GET("/bills", h.Bill.Index)
		auth.POST("/bills", h.Bill.Store)
		auth.GET("/bills/:id", h.Bill.Show)
		auth.PUT("/bills/:id", h.Bill.Update)
		auth.DELETE("/bills/:id", h.Bill.Destroy)
		auth.GET("/subscriptions", h.Bill.Index)
		auth.POST("/subscriptions", h.Bill.Store)
		auth.GET("/subscriptions/:id", h.Bill.Show)
		auth.PUT("/subscriptions/:id", h.Bill.Update)
		auth.DELETE("/subscriptions/:id", h.Bill.Destroy)

		// Categories
		auth.GET("/categories", h.Category.Index)
		auth.POST("/categories", h.Category.Store)
		auth.GET("/categories/:id", h.Category.Show)
		auth.PUT("/categories/:id", h.Category.Update)
		auth.DELETE("/categories/:id", h.Category.Destroy)

		// Tags
		auth.GET("/tags", h.Tag.Index)
		auth.POST("/tags", h.Tag.Store)
		auth.GET("/tags/:id", h.Tag.Show)
		auth.PUT("/tags/:id", h.Tag.Update)
		auth.DELETE("/tags/:id", h.Tag.Destroy)

		// Piggy Banks
		auth.GET("/piggy-banks", h.PiggyBank.Index)
		auth.POST("/piggy-banks", h.PiggyBank.Store)
		auth.GET("/piggy-banks/:id", h.PiggyBank.Show)
		auth.PUT("/piggy-banks/:id", h.PiggyBank.Update)
		auth.DELETE("/piggy-banks/:id", h.PiggyBank.Destroy)
		auth.GET("/piggy-banks/:id/events", h.PiggyBank.ListEvents)

		// Object Groups
		auth.GET("/object-groups", h.ObjectGroup.Index)
		auth.POST("/object-groups", h.ObjectGroup.Store)
		auth.GET("/object-groups/:id", h.ObjectGroup.Show)
		auth.PUT("/object-groups/:id", h.ObjectGroup.Update)
		auth.DELETE("/object-groups/:id", h.ObjectGroup.Destroy)
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
		admin.POST("/currencies", h.Currency.Store)
		admin.DELETE("/currencies/:currency_code", h.Currency.Destroy)
		admin.POST("/link-types", h.LinkType.Store)
		admin.DELETE("/link-types/:id", h.LinkType.Destroy)
	}
}

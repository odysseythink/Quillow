package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/anthropics/firefly-iii-go/internal/adapter/handler"
	v1 "github.com/anthropics/firefly-iii-go/internal/adapter/handler/v1"
	"github.com/anthropics/firefly-iii-go/internal/adapter/repository"
	accountuc "github.com/anthropics/firefly-iii-go/internal/usecase/account"
	authuc "github.com/anthropics/firefly-iii-go/internal/usecase/auth"
	billuc "github.com/anthropics/firefly-iii-go/internal/usecase/bill"
	budgetuc "github.com/anthropics/firefly-iii-go/internal/usecase/budget"
	categoryuc "github.com/anthropics/firefly-iii-go/internal/usecase/category"
	configuc "github.com/anthropics/firefly-iii-go/internal/usecase/configuration"
	currencyuc "github.com/anthropics/firefly-iii-go/internal/usecase/currency"
	eruc "github.com/anthropics/firefly-iii-go/internal/usecase/exchangerate"
	objectgroupuc "github.com/anthropics/firefly-iii-go/internal/usecase/objectgroup"
	piggybankuc "github.com/anthropics/firefly-iii-go/internal/usecase/piggybank"
	prefuc "github.com/anthropics/firefly-iii-go/internal/usecase/preference"
	recurrenceuc "github.com/anthropics/firefly-iii-go/internal/usecase/recurrence"
	ruleuc "github.com/anthropics/firefly-iii-go/internal/usecase/rule"
	taguc "github.com/anthropics/firefly-iii-go/internal/usecase/tag"
	txuc "github.com/anthropics/firefly-iii-go/internal/usecase/transaction"
	useruc "github.com/anthropics/firefly-iii-go/internal/usecase/user"
	webhookuc "github.com/anthropics/firefly-iii-go/internal/usecase/webhook"
	"github.com/anthropics/firefly-iii-go/pkg/config"
	"github.com/anthropics/firefly-iii-go/pkg/database"
	"github.com/anthropics/firefly-iii-go/pkg/i18n"
	"github.com/anthropics/firefly-iii-go/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	i18nSvc, err := i18n.NewService("locales")
	if err != nil {
		log.Fatalf("Failed to initialize i18n: %v", err)
	}
	if err := i18nSvc.LoadLocale("en_US"); err != nil {
		log.Printf("Warning: failed to load en_US locale: %v", err)
	}
	_ = i18nSvc

	jwtSvc := jwt.NewService(cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	prefRepo := repository.NewPreferenceRepository(db)
	configRepo := repository.NewConfigurationRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	currencyRepo := repository.NewCurrencyRepository(db)
	exchangeRateRepo := repository.NewExchangeRateRepository(db)
	txRepo := repository.NewTransactionRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)
	linkTypeRepo := repository.NewLinkTypeRepository(db)
	txLinkRepo := repository.NewTransactionLinkRepository(db)
	budgetRepo := repository.NewBudgetRepository(db)
	budgetLimitRepo := repository.NewBudgetLimitRepository(db)
	billRepo := repository.NewBillRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)
	piggyBankRepo := repository.NewPiggyBankRepository(db)
	objectGroupRepo := repository.NewObjectGroupRepository(db)
	ruleGroupRepo := repository.NewRuleGroupRepository(db)
	ruleRepo := repository.NewRuleRepository(db)
	recurrenceRepo := repository.NewRecurrenceRepository(db)
	webhookRepo := repository.NewWebhookRepository(db)

	// Usecases
	authUC := authuc.NewUseCase(userRepo, jwtSvc)
	userUC := useruc.NewUseCase(userRepo)
	prefUC := prefuc.NewUseCase(prefRepo)
	configUC := configuc.NewUseCase(configRepo)
	currUC := currencyuc.NewUseCase(currencyRepo)
	accountUC := accountuc.NewUseCase(accountRepo, currencyRepo)
	erUC := eruc.NewUseCase(exchangeRateRepo, currencyRepo)
	transactionUC := txuc.NewUseCase(txRepo, accountRepo)
	budgetUC := budgetuc.NewUseCase(budgetRepo, budgetLimitRepo)
	billUC := billuc.NewUseCase(billRepo)
	categoryUC := categoryuc.NewUseCase(categoryRepo)
	tagUC := taguc.NewUseCase(tagRepo)
	piggyBankUC := piggybankuc.NewUseCase(piggyBankRepo)
	objectGroupUC := objectgroupuc.NewUseCase(objectGroupRepo)
	ruleUC := ruleuc.NewUseCase(ruleGroupRepo, ruleRepo)
	recurrenceUC := recurrenceuc.NewUseCase(recurrenceRepo)
	webhookUC := webhookuc.NewUseCase(webhookRepo)

	// Handlers
	handlers := handler.Handlers{
		Auth:          v1.NewAuthHandler(authUC),
		About:         v1.NewAboutHandler(cfg, userUC, cfg.Database.Driver),
		User:          v1.NewUserHandler(userUC),
		Configuration: v1.NewConfigurationHandler(configUC),
		Preference:    v1.NewPreferenceHandler(prefUC),
		Account:       v1.NewAccountHandler(accountUC, currUC),
		Currency:      v1.NewCurrencyHandler(currUC),
		ExchangeRate:  v1.NewExchangeRateHandler(erUC, currUC),
		Transaction:   v1.NewTransactionHandler(transactionUC, accountUC),
		Attachment:    v1.NewAttachmentHandler(attachmentRepo),
		LinkType:      v1.NewLinkTypeHandler(linkTypeRepo, txLinkRepo),
		Budget:        v1.NewBudgetHandler(budgetUC),
		Bill:          v1.NewBillHandler(billUC),
		Category:      v1.NewCategoryHandler(categoryUC),
		Tag:           v1.NewTagHandler(tagUC),
		PiggyBank:     v1.NewPiggyBankHandler(piggyBankUC),
		ObjectGroup:   v1.NewObjectGroupHandler(objectGroupUC),
		Rule:          v1.NewRuleHandler(ruleUC),
		Recurrence:    v1.NewRecurrenceHandler(recurrenceUC),
		Webhook:       v1.NewWebhookHandler(webhookUC),
		Cron:          v1.NewCronHandler(""),
	}

	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()
	handler.SetupRouter(r, handlers, jwtSvc, userRepo)

	// Serve embedded frontend (SPA fallback)
	setupFrontend(r)

	addr := ":" + cfg.Server.Port
	log.Printf("Firefly III Go server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

//go:embed all:frontend
var frontendFS embed.FS

func setupFrontend(r *gin.Engine) {
	distFS, err := fs.Sub(frontendFS, "frontend")
	if err != nil {
		log.Printf("Warning: frontend not embedded (build with web/dist): %v", err)
		return
	}
	fileServer := http.FileServer(http.FS(distFS))
	r.NoRoute(gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Try to serve static file first
		path := req.URL.Path
		if f, err := distFS.(fs.ReadFileFS).ReadFile(path[1:]); err == nil {
			_ = f
			fileServer.ServeHTTP(w, req)
			return
		}
		// SPA fallback: serve index.html for all non-API routes
		req.URL.Path = "/"
		fileServer.ServeHTTP(w, req)
	})))
}

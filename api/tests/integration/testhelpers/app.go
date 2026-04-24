package testhelpers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/config"
	"github.com/ajianaz/gofin-full/api/internal/handler"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	"github.com/ajianaz/gofin-full/api/internal/router"
	"github.com/ajianaz/gofin-full/api/internal/service"
	"github.com/ajianaz/gofin-full/api/internal/sse"
)

// TestApp bundles the Fiber app with the database pool and seed data,
// providing a single cleanup entry point.
type TestApp struct {
	App    *fiber.App
	DB     *pgxpool.Pool
	JWTMgr *auth.JWTManager
	Seed   *SeedData
	Cfg    *TestConfig
}

// NewTestApp bootstraps a full Fiber application backed by a real database.
// It connects to the test DB, runs migrations, seeds fixtures, and wires up
// all handlers and routes the same way cmd/server/main.go does.
//
// Returns an error if the database is unreachable (callers should skip tests).
func NewTestApp(cfg *TestConfig) (*TestApp, error) {
	// Connect and migrate
	db, err := SetupTestDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("setup test db: %w", err)
	}

	// Clean any leftover data, then seed fresh fixtures
	TruncateAllTables(db)

	// JWT manager (60 min access, 30 day refresh — matches production defaults)
	jwtMgr := auth.NewJWTManager(cfg.JWTSecret, 60, 30)

	// Seed test data
	seed, err := SeedTestData(db, jwtMgr)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("seed test data: %w", err)
	}

	// Build a minimal config.Config so auth.NewProvider can select the right provider.
	// Falls back to "disabled" for backwards compatibility.
	authProviderName := "disabled"
	if envOr("AUTH_PROVIDER", "") != "" {
		authProviderName = envOr("AUTH_PROVIDER", "disabled")
	}
	prodCfg := &config.Config{
		AuthProvider:          authProviderName,
		AuthAllowRegistration: true, // enable for registration tests
		AuthRefreshExpiry:     30,   // 30 days — matches production default
	}

	// Auth provider — local provider needs DB for bcrypt authentication
	authProvider := auth.NewProvider(prodCfg)
	if ap, ok := authProvider.(interface{ SetDB(*pgxpool.Pool) }); ok {
		ap.SetDB(db)
	}

	// SSE hub (real, but no consumers in tests)
	sseHub := sse.NewHub(zerolog.Nop())

	// Repositories
	userRepo := repository.NewUserRepository(db)
	oauthStateRepo := repository.NewOAuthStateRepository(db)
	groupRepo := repository.NewUserGroupRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)
	txRepo := repository.NewTransactionRepository(db)
	budgetRepo := repository.NewBudgetRepository(db)
	piggyRepo := repository.NewPiggyBankRepository(db)
	ruleGroupRepo := repository.NewRuleGroupRepository(db)
	ruleRepo := repository.NewRuleRepository(db)
	recurrenceRepo := repository.NewRecurrenceRepository(db)
	currencyRepo := repository.NewCurrencyRepository(db)
	billRepo := repository.NewBillRepository(db)
	exchangeRateRepo := repository.NewExchangeRateRepository(db)
	webhookRepo := repository.NewWebhookRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	preferenceRepo := repository.NewPreferenceRepository(db)
	configurationRepo := repository.NewConfigurationRepository(db)
	objectGroupRepo := repository.NewObjectGroupRepository(db)
	noteRepo := repository.NewNoteRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	accountTypeRepo := repository.NewAccountTypeRepository(db)
	walletMemberRepo := repository.NewWalletMemberRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	refreshRepo := repository.NewRefreshTokenRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	// Services
	txService := service.NewTransactionService(txRepo, walletRepo)
	exportService := service.NewExportService(txRepo)
	_ = service.NewNotificationService(notificationRepo, sseHub)
	_ = service.NewExchangeRateService(exchangeRateRepo)

	// Handlers
	healthHandler := handler.NewHealthHandler(db, nil) // nil Redis is fine for tests
	authHandler := handler.NewAuthHandler(jwtMgr, authProvider, prodCfg, userRepo, oauthStateRepo, refreshRepo)
	userHandler := handler.NewUserHandler(userRepo)
	groupHandler := handler.NewUserGroupHandler(groupRepo, userRepo, db, jwtMgr)
	walletHandler := handler.NewWalletHandler(walletRepo)
	categoryHandler := handler.NewCategoryHandler(categoryRepo)
	tagHandler := handler.NewTagHandler(tagRepo)
	txHandler := handler.NewTransactionHandler(txService, txRepo)
	budgetHandler := handler.NewBudgetHandler(budgetRepo)
	piggyHandler := handler.NewPiggyBankHandler(piggyRepo)
	ruleGroupHandler := handler.NewRuleGroupHandler(ruleGroupRepo)
	ruleHandler := handler.NewRuleHandler(ruleRepo)
	recurrenceHandler := handler.NewRecurrenceHandler(recurrenceRepo)
	currencyHandler := handler.NewCurrencyHandler(currencyRepo)
	billHandler := handler.NewBillHandler(billRepo)
	exchangeRateHandler := handler.NewExchangeRateHandler(exchangeRateRepo)
	webhookHandler := handler.NewWebhookHandler(webhookRepo)
	attachmentHandler := handler.NewAttachmentHandler(attachmentRepo)
	notificationHandler := handler.NewNotificationHandler(notificationRepo)
	preferenceHandler := handler.NewPreferenceHandler(preferenceRepo)
	configurationHandler := handler.NewConfigurationHandler(configurationRepo)
	objectGroupHandler := handler.NewObjectGroupHandler(objectGroupRepo)
	noteHandler := handler.NewNoteHandler(noteRepo)
	locationHandler := handler.NewLocationHandler(locationRepo)
	accountTypeHandler := handler.NewAccountTypeHandler(accountTypeRepo)
	walletMemberHandler := handler.NewWalletMemberHandler(walletMemberRepo, userRepo)
	exportHandler := handler.NewExportHandler(exportService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsRepo)
	auditHandler := handler.NewAuditHandler(auditRepo)
	adminHandler := handler.NewAdminHandler(userRepo, configurationRepo)
	apiKeyHandler := handler.NewAPIKeyHandler(apiKeyRepo)
	apiDocHandler := handler.NewAPIDocHandler()
	metricsHandler := handler.NewMetricsHandler()

	// Wire up the Fiber app with all routes
	app := router.New(router.RouterConfig{
		HealthHandler:        healthHandler,
		AuthHandler:          authHandler,
		UserHandler:          userHandler,
		GroupHandler:         groupHandler,
		WalletHandler:        walletHandler,
		CategoryHandler:      categoryHandler,
		TagHandler:           tagHandler,
		TxHandler:            txHandler,
		BudgetHandler:        budgetHandler,
		PiggyHandler:         piggyHandler,
		RuleGroupHandler:     ruleGroupHandler,
		RuleHandler:          ruleHandler,
		RecurrenceHandler:    recurrenceHandler,
		CurrencyHandler:      currencyHandler,
		BillHandler:          billHandler,
		ExchangeRateHandler:  exchangeRateHandler,
		WebhookHandler:       webhookHandler,
		AttachmentHandler:    attachmentHandler,
		NotificationHandler:  notificationHandler,
		PreferenceHandler:    preferenceHandler,
		ConfigurationHandler: configurationHandler,
		ObjectGroupHandler:   objectGroupHandler,
		NoteHandler:          noteHandler,
		LocationHandler:      locationHandler,
		AccountTypeHandler:   accountTypeHandler,
		WalletMemberHandler:  walletMemberHandler,
		ExportHandler:        exportHandler,
		AnalyticsHandler:     analyticsHandler,
		AuditHandler:         auditHandler,
		AdminHandler:         adminHandler,
		APIKeyHandler:        apiKeyHandler,
		KeyLookup:            apiKeyRepo,
		APIDocHandler:        apiDocHandler,
		MetricsHandler:       metricsHandler,
		MemberRepo:           walletMemberRepo,
		RoleLookup:           userRepo,
		TokenVersionLookup:   userRepo,
		JWTManager:           jwtMgr,
		SSEHub:               sseHub,
		CustomMiddleware:     nil, // no logger/recovery overhead in tests
	})

	return &TestApp{
		App:    app,
		DB:     db,
		JWTMgr: jwtMgr,
		Seed:   seed,
		Cfg:    cfg,
	}, nil
}

// Cleanup closes the database pool. Call this from TestMain after m.Run().
func (ta *TestApp) Cleanup() {
	if ta.DB != nil {
		ta.DB.Close()
	}
}

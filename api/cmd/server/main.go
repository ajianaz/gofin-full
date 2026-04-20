package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/azfirazka/gofin-full/api/internal/auth"
	"github.com/azfirazka/gofin-full/api/internal/config"
	"github.com/azfirazka/gofin-full/api/internal/handler"
	"github.com/azfirazka/gofin-full/api/internal/middleware"
	"github.com/azfirazka/gofin-full/api/internal/repository"
	"github.com/azfirazka/gofin-full/api/internal/router"
	"github.com/azfirazka/gofin-full/api/internal/service"
	"github.com/azfirazka/gofin-full/api/internal/sse"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Setup logger
	log := setupLogger(cfg)
	log.Info().Str("env", cfg.AppEnv).Msg("starting gofin server")

	// Production safety checks
	if cfg.IsProduction() {
		if cfg.AuthJWTSecret == "change-me-in-production-32chars!" {
			log.Fatal().Msg("AUTH_JWT_SECRET must be changed from default in production")
		}
		if len(cfg.AuthJWTSecret) < 32 {
			log.Fatal().Msg("AUTH_JWT_SECRET must be at least 32 characters")
		}
		if cfg.StaticCronToken == "PLEASE_REPLACE_WITH_32_CHAR_CODE" {
			log.Fatal().Msg("STATIC_CRON_TOKEN must be changed from default in production")
		}
	}

	// Connect to PostgreSQL
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()
	log.Info().Msg("connected to postgresql")

	// Connect to Redis
	rdb := connectRedis(cfg)
	if rdb != nil {
		log.Info().Msg("connected to redis")
	} else {
		log.Warn().Msg("redis not available, running without cache")
	}

	// Setup auth
	jwtMgr := auth.NewJWTManager(cfg.AuthJWTSecret, cfg.AuthJWTExpiry, cfg.AuthRefreshExpiry)
	authProvider := auth.NewProvider(cfg)
	authProvider.SetDB(db)
	log.Info().Str("provider", authProvider.Name()).Msg("auth provider initialized")

	// Setup SSE hub for real-time notifications
	sseHub := sse.NewHub(log)
	log.Info().Msg("SSE hub initialized")

	// Create repositories
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
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	refreshRepo := repository.NewRefreshTokenRepository(db)
	// Create services
	txService := service.NewTransactionService(txRepo, walletRepo)

	// Create handlers
	healthHandler := handler.NewHealthHandler(db, rdb)
	authHandler := handler.NewAuthHandler(jwtMgr, authProvider, cfg, userRepo, oauthStateRepo, refreshRepo)
	userHandler := handler.NewUserHandler(userRepo)
	groupHandler := handler.NewUserGroupHandler(groupRepo, userRepo, db)
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
	walletMemberHandler := handler.NewWalletMemberHandler(walletMemberRepo)
	notifService := service.NewNotificationService(notificationRepo, sseHub)
	_ = notifService
	_ = service.NewExchangeRateService(exchangeRateRepo)
	exportService := service.NewExportService(txRepo)
	exportHandler := handler.NewExportHandler(exportService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsRepo)
	auditHandler := handler.NewAuditHandler(auditRepo)
	adminHandler := handler.NewAdminHandler(userRepo, configurationRepo)
	apiKeyHandler := handler.NewAPIKeyHandler(apiKeyRepo)
	apiDocHandler := handler.NewAPIDocHandler()
	metricsHandler := handler.NewMetricsHandler()

	// Create router
	app := router.New(router.RouterConfig{
		AppURL:          cfg.AppURL,
		AppEnv:          cfg.AppEnv,
		HealthHandler:   healthHandler,
		AuthHandler:     authHandler,
		UserHandler:     userHandler,
		GroupHandler:    groupHandler,
		WalletHandler:   walletHandler,
		CategoryHandler: categoryHandler,
		TagHandler:      tagHandler,
		TxHandler:       txHandler,
		BudgetHandler:    budgetHandler,
		PiggyHandler:     piggyHandler,
		RuleGroupHandler: ruleGroupHandler,
		RuleHandler:      ruleHandler,
		RecurrenceHandler: recurrenceHandler,
		CurrencyHandler: currencyHandler,
		BillHandler: billHandler,
		ExchangeRateHandler: exchangeRateHandler,
		WebhookHandler: webhookHandler,
		AttachmentHandler: attachmentHandler,
		NotificationHandler: notificationHandler,
		PreferenceHandler: preferenceHandler,
		ConfigurationHandler: configurationHandler,
		ObjectGroupHandler: objectGroupHandler,
		NoteHandler: noteHandler,
		LocationHandler: locationHandler,
		AccountTypeHandler: accountTypeHandler,
		WalletMemberHandler: walletMemberHandler,
		ExportHandler: exportHandler,
		AnalyticsHandler: analyticsHandler,
		AuditHandler: auditHandler,
	AdminHandler: adminHandler,
			APIKeyHandler: apiKeyHandler,
		APIDocHandler: apiDocHandler,
		MetricsHandler: metricsHandler,
		MemberRepo: walletMemberRepo,
			KeyLookup: apiKeyRepo,
		RoleLookup:      userRepo,
		JWTManager:       jwtMgr,
		SSEHub:           sseHub,
		RateLimitMax:       cfg.RateLimitMax,
		RateLimitWindowSec: cfg.RateLimitWindowSeconds,
		DisableMetrics:     cfg.DisablePrometheus,
		RedisClient:        rdb,
		CustomMiddleware: []fiber.Handler{
			middleware.Logger(log),
			middleware.Recovery(log),
		},
	})

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info().Str("addr", cfg.HTTPAddr()).Msg("server listening")
		if err := app.Listen(cfg.HTTPAddr()); err != nil {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	<-quit
	log.Info().Msg("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = app.ShutdownWithContext(ctx)

	log.Info().Msg("server stopped")
}

func setupLogger(cfg *config.Config) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if cfg.AppDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	var output interface{ Write(p []byte) (n int, err error) } = zerolog.ConsoleWriter{Out: os.Stdout}
	if cfg.LogFormat == "json" {
		output = os.Stdout
	}

	return zerolog.New(output).With().Timestamp().Logger()
}

func connectDB(cfg *config.Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	poolCfg.MaxConns = int32(cfg.DBMaxOpenConns)
	poolCfg.MinConns = int32(cfg.DBMaxIdleConns)
	poolCfg.MaxConnLifetime = time.Duration(cfg.DBConnMaxLifetime) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func connectRedis(cfg *config.Config) redis.Cmdable {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil
	}

	return rdb
}

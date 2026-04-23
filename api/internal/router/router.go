package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/handler"
	"github.com/ajianaz/gofin-full/api/internal/middleware"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	"github.com/ajianaz/gofin-full/api/internal/sse"
)

// RouterConfig holds all dependencies needed by the router.
type RouterConfig struct {
	AppURL               string
	AppEnv               string
	HealthHandler        *handler.HealthHandler
	AuthHandler          *handler.AuthHandler
	UserHandler          *handler.UserHandler
	GroupHandler         *handler.UserGroupHandler
	WalletHandler        *handler.WalletHandler
	CategoryHandler      *handler.CategoryHandler
	TagHandler           *handler.TagHandler
	TxHandler            *handler.TransactionHandler
	BudgetHandler        *handler.BudgetHandler
	PiggyHandler         *handler.PiggyBankHandler
	RuleGroupHandler     *handler.RuleGroupHandler
	RuleHandler          *handler.RuleHandler
	RecurrenceHandler    *handler.RecurrenceHandler
	CurrencyHandler      *handler.CurrencyHandler
	BillHandler          *handler.BillHandler
	ExchangeRateHandler  *handler.ExchangeRateHandler
	WebhookHandler       *handler.WebhookHandler
	AttachmentHandler    *handler.AttachmentHandler
	NotificationHandler  *handler.NotificationHandler
	PreferenceHandler    *handler.PreferenceHandler
	ConfigurationHandler *handler.ConfigurationHandler
	ObjectGroupHandler   *handler.ObjectGroupHandler
	NoteHandler          *handler.NoteHandler
	LocationHandler      *handler.LocationHandler
	AccountTypeHandler   *handler.AccountTypeHandler
	WalletMemberHandler  *handler.WalletMemberHandler
	ExportHandler        *handler.ExportHandler
	AnalyticsHandler     *handler.AnalyticsHandler
	AuditHandler         *handler.AuditHandler
	AdminHandler         *handler.AdminHandler
	APIDocHandler        *handler.APIDocHandler
	MetricsHandler       *handler.MetricsHandler
	MemberRepo           *repository.WalletMemberRepository
	APIKeyHandler        *handler.APIKeyHandler
	KeyLookup            auth.KeyLookup
	RoleLookup           auth.RoleLookup
	TokenVersionLookup   auth.TokenVersionLookup
	JWTManager           *auth.JWTManager
	SSEHub               *sse.Hub
	CustomMiddleware     []fiber.Handler
	RateLimitMax         int
	RateLimitWindowSec   int
	DisableMetrics       bool
	RedisClient          redis.Cmdable
}

// New creates a new Fiber app with all routes registered.
func New(cfg RouterConfig) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Global middleware
	app.Use(middleware.RequestID())
	app.Use(middleware.CORS(cfg.AppURL, cfg.AppEnv))
	app.Use(middleware.AcceptHeaders())
	app.Use(middleware.SecurityHeaders())
	if !cfg.DisableMetrics {
		app.Use(middleware.Metrics())
	}

	// Add custom middleware (logger, recovery, auth, etc.)
	for _, m := range cfg.CustomMiddleware {
		app.Use(m)
	}

	// Health check (no auth required)
	app.Get("/health", cfg.HealthHandler.Check)

	// Prometheus metrics (no auth)
	if !cfg.DisableMetrics {
		app.Get("/metrics", cfg.MetricsHandler.Prometheus)
	}

	// API v1 routes
	v1 := app.Group("/api/v1")

	// Auth routes (public)
	authGroup := v1.Group("/auth")
	if cfg.RedisClient != nil && cfg.RateLimitMax > 0 {
		rl := middleware.RateLimit(cfg.RedisClient, cfg.RateLimitMax, time.Duration(cfg.RateLimitWindowSec)*time.Second)
		authGroup.Use(rl)
	}
	authGroup.Get("/provider", cfg.AuthHandler.Provider)
	authGroup.Post("/login", cfg.AuthHandler.Login)
	authGroup.Post("/register", cfg.AuthHandler.Register)
	authGroup.Post("/logout", cfg.AuthHandler.Logout)
	authGroup.Post("/refresh", cfg.AuthHandler.Refresh)

	// OAuth routes (public)
	authGroup.Get("/:provider/url", cfg.AuthHandler.OAuthURL)
	authGroup.Get("/:provider/callback", cfg.AuthHandler.OAuthCallback)

	// Public API info
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Gofin API v1",
			"version": "1.0.0",
		})
	})

	// API documentation
	v1.Get("/docs", cfg.APIDocHandler.APIDocs)
	v1.Get("/openapi.json", cfg.APIDocHandler.OpenAPISpec)

	// Protected routes
	protected := v1.Group("")
	if cfg.KeyLookup != nil {
		protected.Use(auth.APIKeyMiddleware(cfg.KeyLookup))
	}
	// Inject token version lookup into context for AuthMiddleware
	if cfg.TokenVersionLookup != nil {
		protected.Use(func(c *fiber.Ctx) error {
			c.Locals("token_version_lookup", cfg.TokenVersionLookup)
			return c.Next()
		})
	}
	protected.Use(auth.AuthMiddleware(cfg.JWTManager))
	protected.Use(auth.GroupRoleMiddleware(cfg.RoleLookup))

	// Current user
	protected.Get("/users/me", cfg.UserHandler.Show)
	protected.Put("/users/me", cfg.UserHandler.Update)
	protected.Post("/users/me/password", cfg.UserHandler.ChangePassword)

	// User groups
	protected.Get("/groups", cfg.GroupHandler.Index)
	protected.Post("/groups", cfg.GroupHandler.Store)
	protected.Post("/groups/switch", cfg.GroupHandler.Switch)
	protected.Get("/groups/:id", cfg.GroupHandler.Show)
	protected.Put("/groups/:id", auth.RBACMiddleware(auth.RoleOwner), cfg.GroupHandler.Update)
	protected.Delete("/groups/:id", auth.RBACMiddleware(auth.RoleOwner), cfg.GroupHandler.Delete)

	// Wallets — read: no RBAC, write: manage_meta
	protected.Get("/wallets", cfg.WalletHandler.Index)
	protected.Get("/wallets/:id", cfg.WalletHandler.Show)
	protected.Post("/wallets", auth.RBACMiddleware(auth.RoleManageMeta), cfg.WalletHandler.Store)
	protected.Put("/wallets/:id", auth.RBACMiddleware(auth.RoleManageMeta), cfg.WalletHandler.Update)
	protected.Delete("/wallets/:id", auth.RBACMiddleware(auth.RoleManageMeta), cfg.WalletHandler.Delete)

	// Categories — read: no RBAC, write: manage_meta
	protected.Get("/categories", cfg.CategoryHandler.Index)
	protected.Get("/categories/:id", cfg.CategoryHandler.Show)
	protected.Post("/categories", auth.RBACMiddleware(auth.RoleManageMeta), cfg.CategoryHandler.Store)
	protected.Put("/categories/:id", auth.RBACMiddleware(auth.RoleManageMeta), cfg.CategoryHandler.Update)
	protected.Delete("/categories/:id", auth.RBACMiddleware(auth.RoleManageMeta), cfg.CategoryHandler.Delete)

	// Tags — read: no RBAC, write: manage_meta
	protected.Get("/tags", cfg.TagHandler.Index)
	protected.Get("/tags/:id", cfg.TagHandler.Show)
	protected.Post("/tags", auth.RBACMiddleware(auth.RoleManageMeta), cfg.TagHandler.Store)
	protected.Put("/tags/:id", auth.RBACMiddleware(auth.RoleManageMeta), cfg.TagHandler.Update)
	protected.Delete("/tags/:id", auth.RBACMiddleware(auth.RoleManageMeta), cfg.TagHandler.Delete)

	// Transactions — read: no RBAC, write: manage_transactions
	protected.Get("/transactions", cfg.TxHandler.Index)
	protected.Get("/transactions/:id", cfg.TxHandler.Show)
	protected.Post("/transactions", auth.RBACMiddleware(auth.RoleManageTransactions), cfg.TxHandler.Store)
	protected.Post("/transactions/split", auth.RBACMiddleware(auth.RoleManageTransactions), cfg.TxHandler.StoreSplit)
	protected.Put("/transactions/:id", auth.RBACMiddleware(auth.RoleManageTransactions), cfg.TxHandler.Update)
	protected.Delete("/transactions/:id", auth.RBACMiddleware(auth.RoleManageTransactions), cfg.TxHandler.Delete)

	// Budgets — read: no RBAC, write: manage_budgets
	protected.Get("/budgets", cfg.BudgetHandler.Index)
	protected.Get("/budgets/:id", cfg.BudgetHandler.Show)
	protected.Post("/budgets", auth.RBACMiddleware(auth.RoleManageBudgets), cfg.BudgetHandler.Store)
	protected.Put("/budgets/:id", auth.RBACMiddleware(auth.RoleManageBudgets), cfg.BudgetHandler.Update)
	protected.Delete("/budgets/:id", auth.RBACMiddleware(auth.RoleManageBudgets), cfg.BudgetHandler.Delete)

	// Piggy banks (nested under wallets) — read: no RBAC, write: manage_piggy_banks
	protected.Get("/wallets/:wallet_id/piggy_banks", cfg.PiggyHandler.Index)
	protected.Get("/wallets/:wallet_id/piggy_banks/:id", cfg.PiggyHandler.Show)
	protected.Post("/wallets/:wallet_id/piggy_banks", auth.RBACMiddleware(auth.RoleManagePiggyBanks), cfg.PiggyHandler.Store)
	protected.Put("/wallets/:wallet_id/piggy_banks/:id", auth.RBACMiddleware(auth.RoleManagePiggyBanks), cfg.PiggyHandler.Update)
	protected.Delete("/wallets/:wallet_id/piggy_banks/:id", auth.RBACMiddleware(auth.RoleManagePiggyBanks), cfg.PiggyHandler.Delete)
	protected.Post("/wallets/:wallet_id/piggy_banks/:id/add-money", auth.RBACMiddleware(auth.RoleManagePiggyBanks), cfg.PiggyHandler.AddMoney)
	protected.Post("/wallets/:wallet_id/piggy_banks/:id/remove-money", auth.RBACMiddleware(auth.RoleManagePiggyBanks), cfg.PiggyHandler.RemoveMoney)
	// Alias routes for frontend compatibility
	protected.Post("/wallets/:wallet_id/piggy_banks/:id/add", auth.RBACMiddleware(auth.RoleManagePiggyBanks), cfg.PiggyHandler.AddMoney)
	protected.Post("/wallets/:wallet_id/piggy_banks/:id/remove", auth.RBACMiddleware(auth.RoleManagePiggyBanks), cfg.PiggyHandler.RemoveMoney)

	// Rule groups — read: no RBAC, write: manage_rules
	protected.Get("/rule-groups", cfg.RuleGroupHandler.Index)
	protected.Get("/rule-groups/:id", cfg.RuleGroupHandler.Show)
	protected.Post("/rule-groups", auth.RBACMiddleware(auth.RoleManageRules), cfg.RuleGroupHandler.Store)
	protected.Put("/rule-groups/:id", auth.RBACMiddleware(auth.RoleManageRules), cfg.RuleGroupHandler.Update)
	protected.Delete("/rule-groups/:id", auth.RBACMiddleware(auth.RoleManageRules), cfg.RuleGroupHandler.Delete)

	// Rules — read: no RBAC, write: manage_rules
	protected.Get("/rules", cfg.RuleHandler.Index)
	protected.Get("/rules/:id", cfg.RuleHandler.Show)
	protected.Post("/rules", auth.RBACMiddleware(auth.RoleManageRules), cfg.RuleHandler.Store)
	protected.Put("/rules/:id", auth.RBACMiddleware(auth.RoleManageRules), cfg.RuleHandler.Update)
	protected.Delete("/rules/:id", auth.RBACMiddleware(auth.RoleManageRules), cfg.RuleHandler.Delete)

	// Recurrences — read: no RBAC, write: manage_recurring
	protected.Get("/recurrences", cfg.RecurrenceHandler.Index)
	protected.Get("/recurrences/:id", cfg.RecurrenceHandler.Show)
	protected.Post("/recurrences", auth.RBACMiddleware(auth.RoleManageRecurring), cfg.RecurrenceHandler.Store)
	protected.Put("/recurrences/:id", auth.RBACMiddleware(auth.RoleManageRecurring), cfg.RecurrenceHandler.Update)
	protected.Delete("/recurrences/:id", auth.RBACMiddleware(auth.RoleManageRecurring), cfg.RecurrenceHandler.Delete)

	// Currencies (reference data) — cacheable, read-only
	protected.Get("/currencies", middleware.CacheControl(300), cfg.CurrencyHandler.Index)
	protected.Get("/currencies/:code", cfg.CurrencyHandler.Show)

	// Wallet types (reference data) — cacheable, read-only
	protected.Get("/wallet-types", middleware.CacheControl(300), cfg.AccountTypeHandler.Index)

	// Bills — read: no RBAC, write: manage_transactions
	protected.Get("/bills", cfg.BillHandler.Index)
	protected.Get("/bills/:id", cfg.BillHandler.Show)
	protected.Post("/bills", auth.RBACMiddleware(auth.RoleManageTransactions), cfg.BillHandler.Store)
	protected.Put("/bills/:id", auth.RBACMiddleware(auth.RoleManageTransactions), cfg.BillHandler.Update)
	protected.Delete("/bills/:id", auth.RBACMiddleware(auth.RoleManageTransactions), cfg.BillHandler.Delete)

	// Exchange rates — read: no RBAC, write: manage_currencies
	protected.Get("/exchange-rates", cfg.ExchangeRateHandler.Index)
	protected.Get("/exchange-rates/rate", cfg.ExchangeRateHandler.Show)
	protected.Post("/exchange-rates", auth.RBACMiddleware(auth.RoleManageCurrencies), cfg.ExchangeRateHandler.Store)
	protected.Delete("/exchange-rates/:id", auth.RBACMiddleware(auth.RoleManageCurrencies), cfg.ExchangeRateHandler.Delete)

	// Webhooks — read: no RBAC, write: manage_webhooks
	protected.Get("/webhooks", cfg.WebhookHandler.Index)
	protected.Get("/webhooks/:id", cfg.WebhookHandler.Show)
	protected.Get("/webhooks/:id/messages", cfg.WebhookHandler.Messages)
	protected.Post("/webhooks", auth.RBACMiddleware(auth.RoleManageWebhooks), cfg.WebhookHandler.Store)
	protected.Put("/webhooks/:id", auth.RBACMiddleware(auth.RoleManageWebhooks), cfg.WebhookHandler.Update)
	protected.Delete("/webhooks/:id", auth.RBACMiddleware(auth.RoleManageWebhooks), cfg.WebhookHandler.Delete)

	// Attachments (polymorphic) — read: no RBAC, write: manage_meta
	protected.Get("/attachments", cfg.AttachmentHandler.Index)
	protected.Get("/attachments/:id", cfg.AttachmentHandler.Show)
	protected.Post("/attachments", auth.RBACMiddleware(auth.RoleManageMeta), cfg.AttachmentHandler.Store)
	protected.Delete("/attachments/:id", auth.RBACMiddleware(auth.RoleManageMeta), cfg.AttachmentHandler.Delete)

	// Notifications — user-scoped, no RBAC needed
	protected.Get("/notifications", cfg.NotificationHandler.Index)
	protected.Get("/notifications/unread", cfg.NotificationHandler.Unread)
	protected.Put("/notifications/:id/read", cfg.NotificationHandler.MarkRead)
	protected.Put("/notifications/read-all", cfg.NotificationHandler.MarkAllRead)

	// Real-time notifications (SSE)
	sseMW := sse.NewMiddleware(cfg.SSEHub, zerolog.Nop())
	protected.Get("/notifications/stream", sseMW.RateLimiter(5, time.Minute), sseMW.Stream())

	// Preferences (user-scoped) — no RBAC needed
	protected.Get("/preferences", cfg.PreferenceHandler.Index)
	protected.Get("/preferences/:name", cfg.PreferenceHandler.Show)
	protected.Post("/preferences", cfg.PreferenceHandler.Set)
	protected.Delete("/preferences/:name", cfg.PreferenceHandler.Delete)

	// Configurations (system-level) — read: view_memberships, write: owner
	protected.Get("/configurations", cfg.ConfigurationHandler.Index)
	protected.Get("/configurations/:name", cfg.ConfigurationHandler.Show)
	protected.Post("/configurations", auth.RBACMiddleware(auth.RoleOwner), cfg.ConfigurationHandler.Set)

	// Object groups — read: no RBAC, write: manage_meta
	protected.Get("/object-groups", cfg.ObjectGroupHandler.Index)
	protected.Get("/object-groups/:id", cfg.ObjectGroupHandler.Show)
	protected.Post("/object-groups", auth.RBACMiddleware(auth.RoleManageMeta), cfg.ObjectGroupHandler.Store)
	protected.Put("/object-groups/:id", auth.RBACMiddleware(auth.RoleManageMeta), cfg.ObjectGroupHandler.Update)
	protected.Delete("/object-groups/:id", auth.RBACMiddleware(auth.RoleManageMeta), cfg.ObjectGroupHandler.Delete)

	// Notes (polymorphic) — read: no RBAC, write: manage_meta
	protected.Get("/notes", cfg.NoteHandler.Index)
	protected.Post("/notes", auth.RBACMiddleware(auth.RoleManageMeta), cfg.NoteHandler.Store)
	protected.Put("/notes/:id", auth.RBACMiddleware(auth.RoleManageMeta), cfg.NoteHandler.Update)
	protected.Delete("/notes/:id", auth.RBACMiddleware(auth.RoleManageMeta), cfg.NoteHandler.Delete)

	// Locations (polymorphic) — read: no RBAC, write: manage_meta
	protected.Get("/locations", cfg.LocationHandler.Show)
	protected.Post("/locations", auth.RBACMiddleware(auth.RoleManageMeta), cfg.LocationHandler.Store)
	protected.Delete("/locations/:id", auth.RBACMiddleware(auth.RoleManageMeta), cfg.LocationHandler.Delete)

	// Wallet members (sharing)
	protected.Get("/wallets/:wallet_id/members", cfg.WalletMemberHandler.Index)
	protected.Post("/wallets/:wallet_id/members", auth.RBACMiddleware(auth.RoleOwner), cfg.WalletMemberHandler.Store)
	protected.Put("/wallets/:wallet_id/members/:id", auth.RBACMiddleware(auth.RoleOwner), cfg.WalletMemberHandler.Update)
	protected.Delete("/wallets/:wallet_id/members/:user_id", auth.RBACMiddleware(auth.RoleOwner), cfg.WalletMemberHandler.Delete)

	// Export — read: no RBAC
	protected.Get("/export/csv", cfg.ExportHandler.CSV)
	protected.Get("/export/ofx", cfg.ExportHandler.OFX)
	protected.Post("/export/reconcile", cfg.ExportHandler.Reconcile)

	// Analytics — read: view_reports
	protected.Get("/analytics/spending-by-category", auth.RBACMiddleware(auth.RoleViewReports), cfg.AnalyticsHandler.SpendingByCategory)
	protected.Get("/analytics/spending-by-period", auth.RBACMiddleware(auth.RoleViewReports), cfg.AnalyticsHandler.SpendingByPeriod)
	protected.Get("/analytics/net-worth", auth.RBACMiddleware(auth.RoleViewReports), cfg.AnalyticsHandler.NetWorth)

	// Audit trail — view_memberships
	protected.Get("/audit-logs", auth.RBACMiddleware(auth.RoleViewMemberships), cfg.AuditHandler.Index)

	// API keys (user-scoped, no RBAC needed)
	if cfg.APIKeyHandler != nil {
		protected.Get("/api-keys", cfg.APIKeyHandler.List)
		protected.Post("/api-keys", cfg.APIKeyHandler.Create)
		protected.Delete("/api-keys/:id", cfg.APIKeyHandler.Delete)
	}

	// Admin endpoints
	protected.Get("/admin/users", auth.AdminMiddleware(), cfg.AdminHandler.ListUsers)
	protected.Post("/admin/users", auth.AdminMiddleware(), cfg.AdminHandler.CreateUser)
	protected.Get("/admin/feature-flags", auth.AdminMiddleware(), cfg.AdminHandler.FeatureFlags)
	protected.Post("/admin/feature-flags", auth.AdminMiddleware(), cfg.AdminHandler.SetFeatureFlag)

	return app
}

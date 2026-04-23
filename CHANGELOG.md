# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [Unreleased]

### Added
- **Real API integration for all 4 reports pages** — reports overview, net-worth, spending-by-category, spending-by-period now use `reportService` API calls instead of mock data
- `reportService` in `web/src/lib/services/reports.ts` (spendingByCategory, spendingByPeriod, netWorth methods)
- **Real API integration for groups page** — replaced mock data with `groupService.list()`, added switch group functionality via `groupService.switch()`
- `groupService` in `web/src/lib/services/groups.ts` (list, create, switch methods)
- i18n key `groups.switch` for English and Indonesian
- **Real API integration for currencies page** — replaced mock data with `currencyService.list()` call to `GET /api/v1/currencies`
- **Real API integration for exchange rates page** — replaced mock data with `currencyService.exchangeRates()` call to `GET /api/v1/exchange-rates`
- `currencyService` in `web/src/lib/services/currencies.ts`
- **Full CRUD via Web UI** — Create, Read, Update, Delete for all 9 resources (wallets, transactions, categories, budgets, bills, tags, piggy banks, recurring transactions, rule groups)
- `update()` and `delete()` methods in all service files
- Delete buttons with confirm dialogs on all list pages
- Error states and empty states on all list pages and dashboard
- `aria-label` on all icon-only buttons (delete, copy)
- i18n keys: `common.loading`, `common.saving`, `common.error`, `common.errorSave`
- 21 E2E tests covering full CRUD lifecycle (register → 9 creates → list pages → 9 deletes)

### Changed
- **Reports pages use real API** — replaced `mockTransactions`, `mockWallets`, `mockCategories`, `mockBudgets` with `reportService` and `walletService` calls, added loading/error/empty states
- **Currencies page uses real API** — replaced `mockCurrencies` with `currencyService.list()` including loading/error states
- **Exchange rates page uses real API** — replaced `mockExchangeRates` with `currencyService.exchangeRates()` including loading/error/empty states
- **Button hover feedback** — default, secondary, destructive buttons now use lightness shift instead of barely-perceptible opacity change
- **Sidebar hover contrast** — dark mode sidebar-accent lightness increased for better visibility
- **All hardcoded text replaced with i18n** — loading, saving, error, confirm dialogs, empty states now use `t()` keys
- Wallet create form simplified to only API-accepted fields (removed unused balance/virtualBalance/iban/currency fields)

### Fixed
- Date fields serialized as ISO 8601 (RFC3339) for Go `time.Time` compatibility (tags, transactions, recurring)
- Numeric fields (`amount_min`, `target_amount`) stringified for Go string-typed JSON fields (bills, piggy banks)
- **Dark mode dropdown menus** — `cn-menu-translucent` now uses `var(--popover)` instead of hardcoded white (fixes white-on-white invisible text)
- Playwright dialog handler collision across serial tests (switched to `page.once`)
- `api-keys` page ghost button had `hover:text-destructive` duplicating base class (no hover feedback)

## [0.1.0] - 2026-04-22

### Added
- Go API backend with Fiber v2 — 125+ endpoints covering personal finance domain
- SvelteKit 5 frontend with Tailwind CSS v4, shadcn-svelte, dark mode, i18n (id/en)
- Double-entry bookkeeping with split transactions
- Hierarchical RBAC with 21 group-level roles and 3 wallet-level roles
- UUID v7 primary keys for all entities
- Real-time notifications via Server-Sent Events (SSE)
- JWT authentication with refresh tokens
- OAuth2 login (Google, GitHub)
- Optional Keycloak (OIDC) integration
- Budget tracking, bill management, piggy banks, recurring transactions
- Rules engine, currency management, exchange rates
- CSV/OFX export and reconciliation
- Analytics: spending by category, period, net worth
- Audit trail and webhook support
- Prometheus metrics endpoint
- Self-hosted Docker Compose deployment with Caddy (HTTPS), PostgreSQL 17, Redis 7
- Automated daily database backup with retention
- Auto-seed admin user on first startup
- E2E tests (Playwright) and unit/integration tests

### Documentation
- Comprehensive README (monorepo overview, self-host quick start)
- API and Web READMEs
- Architecture overview, production runbook
- OpenAPI 3.0 specification
- GitHub issue and PR templates

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [Unreleased]

### Added
- **E2E tests for API integration validation** — 31 Playwright tests covering settings (api-keys, preferences, notifications, profile), wallet members, rules group detail, and full API endpoint validation
- **E2E tests for all remaining pages** — 62 Playwright tests covering dashboard, wallets (list+create), transactions (list+create), categories (list+create), budgets (list+create), bills (list+create), recurring (list+create), piggy-banks (list+create), rules (list+create), tags (list+create), groups, currencies, exchange-rates, export, reports (4 pages), admin (users+audit-log), and full API endpoint validation for all 15 service endpoints

### Fixed
- **Null-safety in service helpers** — `JsonApiMany.data` now accepts `null` (Go backend returns `null` for empty slices), `unwrapMany` handles null gracefully
- **Unit tests for all 19 service layer modules** (127 tests) — currencies, groups, reports, export, admin, wallets, transactions, categories, budgets, bills, tags, piggy-banks, recurring, rules, auth, api-keys, preferences, notifications, wallet-members
- **Real API integration for settings/api-keys page** — `apiKeyService` (list, create, delete), loading/error states, copy-to-clipboard, create dialog showing raw key
- **Real API integration for settings/preferences page** — `preferenceService` (list, get, set, delete), local preference config map for type/options metadata, optimistic updates with rollback
- **Real API integration for settings/notifications page** — `notificationService` (list, markRead, markAllRead), read/unread badge display
- **Real API integration for wallets/[id]/members page** — `walletMemberService` (list, add, updateRole, remove), role badges (owner/editor/viewer)
- **Real API integration for rules/[groupId] page** — `ruleService.listRules()` and `ruleService.getGroup()`, simplified rule display (title, active status, priority)
- **Real API integration for settings/profile page** — `authService.getMe()` for user data, `authService.updateProfile()` and `authService.changePassword()`, initials derived from name, password validation
- **Complete ruleService** — added `listGroups`, `getGroup`, `listRules`, `getRule`, `createRule`, `updateRule`, `deleteRule` methods
- **Additional service methods** — `piggyBankService.addMoney/removeMoney`, `transactionService.split`, `groupService.get/update/delete`, `authService.updateProfile/changePassword`
- New domain types: `Notification`, `WalletMember`, `ApiKeyListItem`, `ApiKeyCreateResponse`, `PreferenceItem`
- i18n keys: `settings.profile.saveSuccess`, `settings.profile.passwordMismatch`, `settings.profile.passwordRequired`, `settings.profile.passwordChanged`
- Tests verify JSON:API response unwrapping, field mapping (backend attribute names to frontend fields), default values for missing fields, error handling, and query string construction
- **Vitest unit test framework** — vitest + @vitest/coverage-v8 + jsdom for frontend unit testing
- `vitest.config.ts` with jsdom environment, path aliases ($lib, $components, $app), and globals
- Test scripts: `test` (single run), `test:watch` (watch mode), `test:coverage` (with v8 coverage)
- **Real API integration for admin pages** — admin users and audit log pages now use `adminService` API calls instead of mock data
- `adminService` in `web/src/lib/services/admin.ts` (listUsers, listAuditLogs methods)
- **Real API integration for export page** — CSV/OFX download via `exportService`, wallet dropdown populated from API
- `exportService` in `web/src/lib/services/export.ts` (downloadCSV, downloadOFX methods)
- i18n key `export.exporting` for English and Indonesian
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

### Removed
- 17 unused mock data files after all pages converted to real API (mock-audit-log, mock-bills, mock-budgets, mock-categories, mock-currencies, mock-exchange-rates, mock-groups, mock-notifications, mock-piggy-banks, mock-recurring, mock-tags, mock-transactions, mock-users, mock-wallets, mock-api-keys, mock-preferences, mock-rules)

### Changed
- **Settings/api-keys page uses real API** — replaced `mockApiKeys` with `apiKeyService.list()`, create shows raw key in alert, copy-to-clipboard for key prefix
- **Settings/preferences page uses real API** — replaced `mockPreferences` with `preferenceService.list()`, changes persist immediately via `preferenceService.set()`
- **Settings/notifications page uses real API** — replaced inline mock notifications with `notificationService.list()`, mark-all-read calls real API
- **Settings/profile page uses real API** — replaced hardcoded user data with `authService.getMe()`, save and change-password wired to API
- **Wallets/[id]/members page uses real API** — replaced inline mock members with `walletMemberService.list()`
- **Rules/[groupId] page uses real API** — replaced `mockRuleGroups` and `mockRules` with `ruleService.getGroup()` and `ruleService.listRules()`
- **Rules parent page** — updated `ruleService.list()` call to `ruleService.listGroups()` after service method rename
- **Admin users page uses real API** — replaced `mockUsers` with `adminService.listUsers()`, added loading/error states
- **Audit log page uses real API** — replaced `mockAuditLog` with `adminService.listAuditLogs()`, entity filter re-fetches from API, action filter is client-side
- **Export page uses real API** — replaced `mockWallets` with `walletService.list()`, form submit triggers actual file download via `exportService`
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

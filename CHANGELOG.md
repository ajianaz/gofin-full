# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [Unreleased]

### Changed
- **License updated from MIT to Apache 2.0** — project license changed to Apache License, Version 2.0

### Added
- **JWT token version/invalidation mechanism (C6)** — `token_version` column on users table tracks token invalidation; JWT access and refresh tokens include `token_version` claim; `AuthMiddleware` compares claim version against DB and rejects mismatched tokens (401 "Token has been invalidated"); `IncrementTokenVersion` repository method for callers to invalidate all tokens for a user (called on logout, password change, group membership removal)
- **Database migration** — `000012_token_version.up.sql` adds `token_version INTEGER NOT NULL DEFAULT 0` to users table
- **Change password endpoint** — `POST /users/me/password` validates current password, hashes new password (min 8 chars), and invalidates all existing JWT tokens via `IncrementTokenVersion`
- **Piggy bank alias routes** — `/wallets/:wallet_id/piggy_banks/:id/add` and `/remove` as aliases for `/add-money` and `/remove-money` for frontend compatibility
- **Database migration** — `000013_notes_locations_ownership.up.sql` adds `user_id` and `user_group_id` columns to notes and locations tables for ownership tracking

### Fixed
- **JWT re-issue on group switch** — `POST /api/v1/groups/switch` now returns a new JWT token pair with the updated `group_id` claim, so subsequent requests operate on the correct group without requiring re-authentication
- **Withdrawal balance check (H8)** — `CreateTransaction` now verifies source wallet has sufficient `virtual_balance` before processing withdrawal; returns error with current and required amounts
- **Split transaction balance check** — `CreateSplitTransaction` aggregates debits per source wallet and validates each has sufficient balance before processing
- **Piggy bank withdraw balance check (H9)** — `RemoveMoney` now verifies piggy bank has sufficient current amount (from events) before allowing withdrawal
- **Piggy bank add wallet balance check (H9b)** — `AddMoney` now verifies source wallet has sufficient `virtual_balance` before moving money to piggy bank
- **Transaction split URL mismatch** — frontend `transactionService.split()` now sends correct payload (`type`, `date`, `journals`) to `POST /transactions/split` matching the backend API

### Fixed
- **Group update restricted to owner role (C1)** — `PUT /groups/:id` now requires `RoleOwner` RBAC middleware, matching the delete route
- **Wallet member Index group verification (C3)** — `GET /wallets/:wallet_id/members` now verifies the requesting user has a role on the wallet via `GetWalletRole`, returns 403 if unauthorized
- **Attachment Show/Delete user ownership (C4+H1)** — `Show` and `Delete` attachment handlers now verify the attachment belongs to the requesting user via `UserID` comparison, preventing cross-user access
- **DB transactions for money operations (C5)** — `CreateTransaction`, `CreateSplitTransaction`, and `DeleteTransaction` now wrap all multi-step operations (create group, journal, transactions, balance updates) in a single database transaction, preventing inconsistent state on partial failure
- **Webhook URL validation (H3)** — webhook create/update now validates URLs: blocks internal IPs (loopback, private ranges, cloud metadata endpoints), requires http/https scheme
- **JWT secret startup validation (H4)** — `config.Load()` rejects default JWT secret and secrets under 32 characters in production; development/testing/local environments are exempted
- **Rate limiter fail-secure on Redis failure (H5)** — when Redis is unavailable, the rate limiter now falls back to an in-memory sliding window limiter instead of allowing all requests through
- **Account lockout after failed logins (H6)** — after 5 consecutive failed login attempts, the account is temporarily locked for 15 minutes with a clear error message; counter resets on successful login; uses Redis INCR with TTL

### Fixed (Security Audit)

#### Critical
- **Notes group/user ownership scoping (C1)** — notes handler and repository now require user authentication, pass `user_id` and `group_id` on create, and verify ownership before update/delete
- **Locations group/user ownership scoping (C2)** — locations handler and repository now require user authentication, pass `user_id` and `group_id` on create, and verify ownership before delete
- **Configurations endpoint privilege escalation (C3)** — `POST /configurations` changed from `RoleOwner` to `AdminMiddleware()` to prevent non-admin users from setting system configurations
- **Attachments index group isolation (C4)** — attachment list now filters by `user_id`, preventing cross-user attachment visibility

#### High
- **Webhook messages IDOR (H1)** — `GET /webhooks/:id/messages` now verifies the webhook belongs to the requesting user's active group before returning messages
- **Audit logs using actual groupID (H3)** — audit log handler now passes the real `groupID` instead of hardcoded `0`
- **JWT invalidation on member removal (H4)** — removing a wallet member or deleting a group now calls `IncrementTokenVersion` to invalidate affected users' JWT tokens immediately
- **SSRF DNS rebinding prevention (H5)** — webhook URL validation now resolves DNS at validation time (not just string matching), preventing attackers from using DNS rebinding to reach internal services
- **Piggy bank TOCTOU race condition (H6)** — `AddMoney` now uses a DB transaction with `SELECT ... FOR UPDATE` row locking to prevent concurrent balance modifications
- **Request size limit middleware (H8)** — fixed the no-op `RequestSizeLimit` middleware to actually check `Content-Length` header and return 413; wired it into the router with `MaxRequestBodyBytes` config

#### Medium
- **CSV injection sanitization (M2)** — `sanitizeCSV` now trims whitespace before checking for formula prefixes, catching `\r\n=cmd` patterns
- **JWT secret validation in all environments (M3)** — removed the dev/local/testing bypass; JWT secret length and default value are always validated
- **In-memory rate limiter memory leak (M4)** — added periodic eviction (every 5 minutes) to remove stale rate limit entries, preventing unbounded memory growth
- **Export service token handling (M5)** — export service now uses consistent token retrieval pattern matching the central API client
- **IPv6 SSRF bypass (M6)** — webhook URL validation now explicitly checks IPv6 private ranges (`::1/128`, `fe80::/10`, `fc00::/7`)
- **CSV injection prevention (M3)** — user-controlled fields in CSV export (description, category, wallet names, notes, tags) are prefixed with a single quote when they start with formula characters (`=`, `+`, `-`, `@`, tab, carriage return)
- **Exchange rate delete group filter (L3)** — `DELETE /exchange-rates/:id` now filters by `user_group_id`, preventing cross-group deletion

### Fixed
- **DB transaction wrapping for money operations** — `CreateTransaction`, `CreateSplitTransaction`, and `DeleteTransaction` now execute all DB writes (group + journal + transactions + balance updates + soft deletes) in a single database transaction, preventing inconsistent state on partial failure

### Added
- **E2E tests for API integration validation** — 31 Playwright tests covering settings (api-keys, preferences, notifications, profile), wallet members, rules group detail, and full API endpoint validation
- **E2E tests for all remaining pages** — 62 Playwright tests covering dashboard, wallets (list+create), transactions (list+create), categories (list+create), budgets (list+create), bills (list+create), recurring (list+create), piggy-banks (list+create), rules (list+create), tags (list+create), groups, currencies, exchange-rates, export, reports (4 pages), admin (users+audit-log), and full API endpoint validation for all 15 service endpoints

### Fixed
- **JWT secret startup validation** — `config.Load()` now rejects the default `AUTH_JWT_SECRET` value and secrets shorter than 32 characters in all non-development environments (previously only a warning in production)
- **Webhook URL SSRF prevention** — webhook create and update endpoints now validate URLs, blocking internal IPs (localhost, 127.0.0.1, private ranges), cloud metadata endpoints (AWS, GCP, GKE), and requiring http/https scheme
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

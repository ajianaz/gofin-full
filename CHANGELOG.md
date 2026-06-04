# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [Unreleased]

### Backend
- **Pagination deterministic ordering** — added `id` tiebreaker to all ListPaginated ORDER BY clauses to prevent row shuffling across pages
- **Pagination added** to all list endpoints — GET /wallets, /bills, /budgets, /categories, /tags, /webhooks, /notifications, /rules, /rule-groups, /piggy-banks now accept `page` (default 1) and `per_page` (default 20, max 100) query params and return `meta.pagination` with total, count, per_page, current_page, total_pages (closes #196)

### Frontend
- **Dashboard error state** — shows error message with retry button instead of empty page when data loading fails; fires error toast automatically (closes #202)
- **Export toast feedback** — shows success toast (CSV/OFX downloaded) and error toast (via global toast system) instead of inline error text (closes #203)

### i18n
- Added `common.errorLoading`, `common.errorLoadingDescription`, `common.retry` keys (en + id)
- Added `export.csvDownloaded`, `export.ofxDownloaded` keys (en + id)

## [0.1.8] - 2026-06-01

### Security
- **CRITICAL: SSE auth fix** — `/notifications/stream` now uses JWT auth context instead of URL param for user_id (prevents IDOR) (closes #192)
- **refresh_token removed from response body** — Login/register/refresh/group-switch no longer expose refresh_token in JSON; only in httpOnly cookie (closes #195)
- **Cookie Secure flag configurable** — Uses `!cfg.IsLocal()` so local HTTP dev works without HTTPS (closes #194)
- **Content-Security-Policy header** added to frontend responses — strict in production, relaxed for Vite HMR in dev (closes #188)
- **Register rate limiting** — 5 requests/minute per IP for POST /auth/register, layered on existing auth group limit (closes #201)

### Backend
- Admin ListUsers returns `name` field correctly instead of email (closes #193)
- Logout handler no longer requires Content-Type header (closes #191)
- Removed stale `ALLOW_2FA_BYPASS` from `.env.example` (closes #208)
- Admin role check optimized — single `HasAnyGlobalRole` query instead of double `HasGlobalRole` (closes #206)
- UserGroup handler accepts `secureCookies` param for consistent cookie Secure flag
- Added `PublicTokenResponse` type for safe token responses

### Frontend
- Register page auto-login after success — redirects to dashboard instead of login page (closes #189)
- Error feedback when auto-login after registration fails

### Tooling
- Cora CLI upgraded to v0.1.7 — deterministic review output (fixes cora-cli#97)
- `.cora.yaml` updated with `temperature: 0`, `cache_ttl: 1440` for reproducible reviews

## [0.1.7] - 2026-06-01

### Security
- **httpOnly cookies**: Migrated token storage from localStorage to httpOnly, Secure, SameSite=Lax cookies — tokens no longer accessible via JavaScript (XSS protection) (closes #128)
- Tokens removed from OAuth redirect URL fragment (closes #143)
- OAuth callback uses cookie-based session verification via `restoreSession()`
- Logout clears httpOnly cookies + increments token_version

### Changed
- Backend: Replaced all 63 `log.Printf` calls with zerolog structured logging (closes #138)
- Backend: Wired validation framework — Required, MinLength, Email, PasswordStrength validators (closes #137)
- Backend: Cached token_version in Redis (eliminates per-request DB query) (closes #139)
- Backend: Fixed N+1 queries in transaction detail & admin user list (closes #140)
- Backend: Added pagination to all unpaginated list endpoints (closes #150)
- Backend: RBAC hierarchy gaps — reviewed and fixed semantic role relationships (closes #154)
- Backend: Added FK from wallets to user_groups table (closes #155)
- Backend: Removed legacy unused columns on users (reset_token, remember_token) (closes #156)
- Backend: Disabled auth provider hardened — explicit config-based enforcement (closes #132)
- Backend: Feature flags sourced from config instead of env vars (closes #158)
- Backend: SecurityHeaders receives config from RouterConfig instead of env vars (closes #153)
- Frontend: Unified token state management — authStore as single source of truth (closes #145)
- Frontend: Export service uses central API client with token refresh (closes #134)
- Frontend: Added global toast/notification system for API errors (closes #146)
- Frontend: SvelteKit load functions for auth guards instead of onMount (closes #147)
- Frontend: Vendor chunk splitting — ui-vendor + data-vendor for better caching (closes #148)
- Frontend: Generated TypeScript interfaces from OpenAPI spec — 78 typed interfaces (closes #133)
- Frontend: CSP headers tightened in hooks.server.ts (closes #129)
- Frontend: Admin layout guard — unauthenticated users redirected from admin routes (closes #130)
- Frontend: Sidebar DOM manipulation replaced with Svelte `$effect` (closes #131)
- Frontend: StatusBadge i18n — hardcoded English labels replaced with translations (closes #161)
- Frontend: Export download deduplicated — extracted shared helper (closes #162)
- Frontend: NavSection component extracted for consistent nav rendering (closes #159)
- Frontend: Dark mode CSS tokens added to component variables (closes #160)
- Housekeeping: Renamed CLAUDE.md → AGENT.md, updated issue templates, overhauled labels

### Added
- Backend: Unit tests for validation framework — 9 test cases (closes #157)
- Backend: API key last_used goroutine uses request context for tracing (closes #163)
- Backend: UserRepository.Update simplified — removed buggy double query construction (closes #151)

### Removed
- 2FA dead code: unused TwoFAService, config field, preference schema entry (closes #142)

## [0.1.5]

### Added
- Email verification on registration — optional, auto-verify without SMTP (closes #109 P3)
- Migration 000015: `verified` column on users table (existing users = true)
- `POST /auth/verify-email` — validate token, set user verified
- `POST /auth/resend-verification` — resend verification email (protected)
- `AUTH_REQUIRE_VERIFICATION` config (default: false) — block unverified login
- FE `/verify-email` page + login not-verified notice

## [0.1.4] - 2026-05-27

### Added
- Google OAuth login — FE buttons on login/register pages, OAuth callback handler (closes #109 P1)
- Forgot/reset password — SMTP mail service, reset token via Redis, BE endpoints, FE forms (closes #109 P2)

## [0.1.3] - 2026-05-27

### Fixed
- User header hydration — waits for client mount before showing auth data (closes #101)
- `/savings` → `/piggy-banks` redirect (308) for back-compatibility (closes #103)
- Report sub-routes — `/reports/category`, `/reports/income-expense`, `/reports/period` (closes #104)
- Period filter date range fix on Transactions page (closes #88)
- Rules form missing field bindings (closes #89)
- Dashboard multi-currency per wallet card (closes #90)
- Settings tabs sidebar navigation (closes #93)
- Logout cleanup — clear auth store, redirect to login (closes #94)
- `PUT /users/me` no longer 500 when only name provided — dynamic SET clause (closes #108)
- Token invalidation after logout — `TokenVersionLookup` injected in production (closes #106)
- Transaction date accepts YYYY-MM-DD, RFC3339, RFC3339Z formats (closes #107)
- Export: replace native `<option>` with shadcn `<Select>` (closes #87)

## [0.1.2] - 2026-05-27

### Fixed
- Security: GroupRoleMiddleware fails closed (401) on DB error (closes #80)
- Security: Config read endpoints require admin role (closes #81)
- Security: Preference API rejects unknown keys — only 10 defined keys accepted (closes #76)
- Security: HTML-escape user text across 11 handlers — prevents stored XSS (closes #96)
- Security: `POST /auth/logout` moved to protected group (closes #97)
- Security: API key restricted to read-only (GET/HEAD) — write ops 403 (closes #98)
- Security: Login rate limiting enabled by default (closes #100)
- Service layer: preserve API response values instead of hardcoding defaults (closes #86)
- formatCurrency: default decimal places 0→2 (closes #91)
- i18n: add missing keys, replace hardcoded English strings (closes #92)

## [0.1.1] - 2026-05-27

### Added
- Multi-currency support — Phase 1 (wallet currency selector) + Phase 2 (exchange rates UI)
- Dashboard: empty state with CTA buttons for new users (closes #30)
- Dashboard: individual wallet cards with name, currency, and balance (closes #40)
- Wallet creation: optional opening balance field (closes #34)
- Wallet list: edit wallet name and currency inline (closes #29)
- Exchange rates: add and delete from UI (closes #31)
- Docker: `docker-compose.traefik.yml` for external Traefik deployment
- CI: GHCR publish + Trivy security scanning (CRITICAL/HIGH)
- CI: SSH auto-deploy workflow on develop push
- Configurable rate limiting — global API middleware + login lockout with max attempts and lockout duration
- Page titles (dynamic per route), favicon (SVG logo), meta description, theme-color (closes #65)
- i18n: missing locale keys for categories, currencies, tags pages (closes #67)
- SelectValue in all SelectTrigger instances — fixes blank dropdown text on 10+ pages

### Changed
- Currency consistency: all pages use wallet/transaction currency symbol instead of locale default
- Dashboard: savings card calculated (income - expense), removed fake trend badges
- Removed all remaining hardcoded currency symbols

### Fixed
- CurrencyResolver handles both UUID and currency code inputs
- Transaction default currency resolved from source wallet instead of hardcoded EUR (closes #33)
- Auth: `setupGroup()` after login/register/restore — fixes "no active group" 400 errors (closes #36-#40)
- API Key auth: fallback to first group_membership when `user_group_id` is NULL (closes #44)
- Preferences: show default values for new users (closes #43)
- Security: API Key auth can no longer create/delete API keys — JWT required (closes #47)
- Input validation: wallet name max 100 chars, opening_balance non-negative, piggy bank target_amount non-negative, exchange rate positive (closes #48-#50)
- CRUD IDOR: Update & Delete return 404 for non-existent resources (closes #52-#55)
- IDOR existence checks on all remaining CRUD handlers (closes #59)
- Auth: login 422 for empty fields, group switch 422/404 for invalid/non-existent group (closes #61)
- Category: name max length validation, trim whitespace (closes #64)
- Svelte hydration error on Transactions page — Select components guarded behind client-only render (closes #21)
- Transactions/Budgets/Exchange Rates "Failed to load data" — resolved by setupGroup fix (closes #37-#39)
- Blank dropdown text on Transactions, PiggyBanks pages — added `SelectValue`, removed duplicate ChevronDown (closes #69, #72)
- Wallet creation: `opening_balance` number→string conversion — fixes 422 from UI (closes #77)
- User profile: `name` field — migration + API support (closes #75)
- E2E test flakiness and assertion mismatches
- CI: Trivy action versioning and deploy workflow fixes

## [0.1.0] - 2026-04-22

### Added
- First open source release

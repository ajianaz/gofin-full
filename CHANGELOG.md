# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [Unreleased]

## [0.1.6] - 2026-05-27

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

## [0.1.5] - 2026-05-26

### Fixed
- GroupRoleMiddleware fails closed (401) on DB error (closes #80)
- Config read endpoints require admin role (closes #81)
- Preference API rejects unknown keys — only 10 defined keys accepted (closes #76)
- HTML-escape user text across 11 handlers — prevents stored XSS (closes #96)
- `POST /auth/logout` moved to protected group (closes #97)
- API key restricted to read-only (GET/HEAD) — write ops 403 (closes #98)
- Login rate limiting enabled by default (closes #100)
- Token invalidation after logout — `TokenVersionLookup` injected in production (closes #106)
- Transaction date accepts YYYY-MM-DD, RFC3339, RFC3339Z formats (closes #107)
- Service layer: preserve API response values instead of hardcoding defaults (closes #86)
- formatCurrency: default decimal places 0→2 (closes #91)
- i18n: add missing keys, replace hardcoded English strings (closes #92)
- Export: replace native `<option>` with shadcn `<Select>` (closes #87)

## [0.1.4] - 2026-05-24

### Added
- Page titles (dynamic per route), favicon (SVG logo), meta description, theme-color (closes #65)
- i18n: missing locale keys for categories, currencies, tags pages (closes #67)
- SelectValue in all SelectTrigger instances — fixes blank dropdown text on 10+ pages

### Changed
- Removed all remaining hardcoded currency symbols

### Fixed
- Blank dropdown text on Transactions list (4 filters), PiggyBanks list, PiggyBanks create — added `SelectValue`, removed duplicate ChevronDown (closes #69, #72)
- Wallet creation: `opening_balance` number→string conversion — fixes 422 from UI (closes #77)
- User profile: `name` field — migration + API support (closes #75)

## [0.1.3] - 2026-05-22

### Added
- Configurable rate limiting — global API middleware + login lockout with max attempts and lockout duration

### Fixed
- Auth: `setupGroup()` after login/register/restore — fixes "no active group" 400 errors (closes #36-#40)
- API Key auth: fallback to first group_membership when `user_group_id` is NULL (closes #44)
- Preferences: show default values for new users (closes #43)
- Security: API Key auth can no longer create/delete API keys — JWT required (closes #47)
- Input validation: wallet name max 100 chars, opening_balance non-negative, piggy bank target_amount non-negative, exchange rate positive (closes #48-#50)
- CRUD IDOR: Update & Delete return 404 for non-existent resources (closes #52-#55)
- IDOR existence checks on all remaining CRUD handlers (closes #59)
- Auth: login 422 for empty fields, group switch 422/404 for invalid/non-existent group (closes #61-#62)
- Category: name max length validation, trim whitespace (closes #64)

## [0.1.2] - 2026-05-20

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

### Changed
- Currency consistency: all pages use wallet/transaction currency symbol instead of locale default
- Dashboard: savings card calculated (income - expense), removed fake trend badges

### Fixed
- CurrencyResolver handles both UUID and currency code inputs
- Transaction default currency resolved from source wallet instead of hardcoded EUR (closes #33)
- E2E test flakiness and assertion mismatches
- CI: Trivy action versioning and deploy workflow fixes

## [0.1.1] - 2026-05-19

### Fixed
- Auth: setupGroup after login/register/restore
- Dashboard: wallet cards
- API Key auth fallback for NULL user_group_id
- Preferences default values for new users
- API Key auth restriction for key management
- Input validation for wallets, piggy banks, exchange rates
- CRUD IDOR protection
- Exchange rate validation
- Configurable rate limiting

## [0.1.0] - 2026-04-22

### Added
- First open source release

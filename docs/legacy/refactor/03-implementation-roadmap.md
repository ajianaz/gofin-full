# Implementation Roadmap — Go API Rewrite

---

> Urutan implementasi per phase untuk Go API rewrite Firefly III.
> Target: API-only, Fiber + sqlx + PostgreSQL + Redis.

## Design Decisions (Locked)

| Decision | Value |
|----------|-------|
| Framework | Fiber v2 |
| DB access | sqlx (raw SQL, no ORM) |
| Database | PostgreSQL |
| Cache/Rate limit | Redis |
| Auth | Keycloak (OAuth2) |
| Monetary | `shopspring/decimal` (string, bukan float64) |
| Naming | Account → **Wallet**, AccountType → **WalletType** |
| API format | JSON (flat, bukan JSON:API) |
| Double-entry | Triple-layer: TransactionGroup → TransactionJournal → Transaction |

---

## Phase 0 — Foundation (Week 1-2)

Target: Project skeleton, bisa run, health check.

### 0.1 Project Setup

- [ ] Init Go module, Fiber app, directory structure
- [ ] Docker Compose: PostgreSQL + Redis + app
- [ ] Environment config (viper/envconfig)
- [ ] Structured logging (zerolog/slog)
- [ ] Health check endpoint: `GET /health`

### 0.2 Database

- [ ] Connection pool (sqlx + pgxpool)
- [ ] Migration tool (goose/golang-migrate)
- [ ] Base schema: `users`, `user_groups`, `group_memberships`, `user_roles`
- [ ] Seed: user_roles (22 roles from doc 17)

### 0.3 Middleware Stack

- [ ] Request ID / correlation ID
- [ ] Panic recovery
- [ ] CORS
- [ ] Request logger
- [ ] Error handler (doc 12)

### Deliverables

```
GET /health → {"status": "ok"}
POST /auth/register → create user + group
POST /auth/login → JWT token
```

---

## Phase 1 — Auth & Users (Week 3-4)

Target: Keycloak integration, user management, group membership.

### 1.1 Keycloak Integration

- [ ] Keycloak client setup (realm, client-id, client-secret)
- [ ] OAuth2 token exchange middleware
- [ ] JWT validation middleware (RS256)
- [ ] Token refresh flow
- [ ] Personal access tokens (PAT) table

### 1.2 User API

- [ ] `POST /auth/register` — create user + group + owner membership
- [ ] `POST /auth/login` — Keycloak auth
- [ ] `POST /auth/token/refresh` — refresh token
- [ ] `POST /auth/password/email` — reset request
- [ ] `POST /auth/password/reset` — reset execute
- [ ] `GET /api/v1/about` — version info
- [ ] `GET /api/v1/about/user` — current user info
- [ ] `GET /api/v1/configuration` — user config
- [ ] `PUT /api/v1/configuration/{key}` — update config

### 1.3 User Group & RBAC

- [ ] Group membership management
- [ ] Role check middleware (22 roles from doc 17)
- [ ] Admin middleware
- [ ] Demo user protection middleware

### 1.4 Preferences

- [ ] `GET /api/v1/preferences` — list
- [ ] `POST /api/v1/preferences` — create
- [ ] `GET /api/v1/preferences/{name}` — get
- [ ] `PUT /api/v1/preferences/{name}` — update

### Referensi

- Doc 04: Auth & RBAC
- Doc 12: Error responses
- Doc 17: Configuration

---

## Phase 2 — Wallets (Accounts) (Week 5-6)

Target: Full wallet CRUD dengan semua type constraints.

### 2.1 Schema

- [ ] `wallets` table (rename from `accounts`)
- [ ] `account_types` table (seed 14 types)
- [ ] `account_meta` table (EAV)

### 2.2 Wallet Type System

- [ ] WalletType enum (14 types, doc 17)
- [ ] Account type constraint matrices:
  - `account_to_transaction` (source+dest → type)
  - `source_dests` (type → valid pairs)
  - `allowed_opposing_types`
  - `allowed_transaction_types`
  - `dynamic_creation_allowed`
- [ ] Auto-create system wallets (expense, revenue, initial_balance, reconciliation)

### 2.3 Wallet API

- [ ] `GET /api/v1/wallets` — list (filter: type, active, name)
- [ ] `POST /api/v1/wallets` — create
- [ ] `GET /api/v1/wallets/{id}` — show
- [ ] `PUT /api/v1/wallets/{id}` — update
- [ ] `DELETE /api/v1/wallets/{id}` — soft delete

### 2.4 Wallet Features

- [ ] Virtual balance (foreign currency)
- [ ] Opening balance (Asset, Loan, Debt, Mortgage)
- [ ] IBAN, account number, BIC
- [ ] Account roles (defaultAsset, savingAsset, ccAsset, cashWalletAsset)
- [ ] Liability fields (interest, interest_period, liability_direction)
- [ ] Credit card fields (cc_type, cc_monthly_payment_date)

### 2.5 Native Amount Observer

- [ ] `wallet.virtual_balance` → `wallet.native_virtual_balance` conversion

### 2.6 Cascade Delete

- [ ] Delete wallet → delete all transactions, piggy banks, attachments, notes

### Referensi

- Doc 06: Database schema
- Doc 17: Configuration (section 7)
- Doc 18: Observer side effects
- Doc 21: Account→Wallet redesign

---

## Phase 3 — Currencies & Exchange Rates (Week 7)

Target: Multi-currency support, exchange rate management.

### 3.1 Schema

- [ ] `currencies` table
- [ ] `currency_exchange_rates` table

### 3.2 Currency API

- [ ] `GET /api/v1/currencies` — list
- [ ] `POST /api/v1/currencies` — create (admin)
- [ ] `GET /api/v1/currencies/{code}` — show
- [ ] `PUT /api/v1/currencies/{code}` — update
- [ ] `DELETE /api/v1/currencies/{code}` — disable (admin)
- [ ] `GET /api/v1/currencies/{code}/exchange-rates` — list rates

### 3.3 Exchange Rate API

- [ ] `POST /api/v1/currencies/{code}/exchange-rates` — create rate
- [ ] `DELETE /api/v1/currencies/{code}/exchange-rates/{id}` — delete rate
- [ ] `POST /api/v1/currencies/exchange-rate` — create by codes
- [ ] `POST /api/v1/currencies/exchange-rate/date` — create by date

### 3.4 Exchange Rate Logic

- [ ] Rate lookup order (doc 07 section 7)
- [ ] `convertToPrimary` preference check
- [ ] External rate download cron (Azure Blob Storage, doc 17 section 6.4)

### Referensi

- Doc 07: Business logic (multi-currency)
- Doc 17: Configuration (exchange rates)
- Doc 18: ProcessesExchangeRates listener

---

## Phase 4 — Categories, Tags, Bills (Week 8-9)

Target: Supporting models sebelum transactions.

### 4.1 Categories

- [ ] Schema: `categories`, `category_transaction_journal` pivot
- [ ] CRUD: `GET/POST/PUT/DELETE /api/v1/categories`
- [ ] Cascade delete (attachments, notes)

### 4.2 Tags

- [ ] Schema: `tags`, `tag_transaction_journal` pivot
- [ ] CRUD: `GET/POST/PUT/DELETE /api/v1/tags`
- [ ] Cascade delete (attachments, notes, locations)

### 4.3 Bills (Subscriptions/Recurring)

- [ ] Schema: `bills`, `bill_meta`
- [ ] CRUD: `GET/POST/PUT/DELETE /api/v1/bills`
- [ ] Bill periods: daily, weekly, monthly, quarterly, half-year, yearly
- [ ] Native amount conversion (amount_min → native_amount_min, amount_max → native_amount_max)
- [ ] Cascade delete (attachments, notes)

### 4.4 Object Groups

- [ ] Schema: `object_groups`, `object_group_members`
- [ ] CRUD: `GET/POST/PUT/DELETE /api/v1/object-groups`

### Referensi

- Doc 06: Database schema
- Doc 10: Response fields
- Doc 11: Validation rules
- Doc 18: Cascade delete observers

---

## Phase 5 — Transactions (Core) (Week 10-12)

Target: **Paling penting.** Full double-entry bookkeeping.

### 5.1 Schema

- [ ] `transaction_groups`
- [ ] `transaction_journals`
- [ ] `transactions` (2 per journal: source + destination)
- [ ] `transaction_types` (seed 7 types)
- [ ] `journal_meta` (EAV)
- [ ] `transaction_journal_links`

### 5.2 Transaction Type Validation

- [ ] Account type constraint check (Phase 2 matrices)
- [ ] Auto-derive transaction type from source+dest
- [ ] Auto-create system wallets (expense, revenue) if not exist

### 5.3 Transaction API

- [ ] `GET /api/v1/transactions` — list (pagination, filter, sort)
- [ ] `POST /api/v1/transactions` — create (single + split)
- [ ] `GET /api/v1/transactions/{id}` — show (group + journals + transactions)
- [ ] `PUT /api/v1/transactions/{id}` — update
- [ ] `DELETE /api/v1/transactions/{id}` — soft delete

### 5.4 Transaction Features

- [ ] Split transactions (multiple journals per group)
- [ ] Budget, category, bill, tag linking
- [ ] Journal meta fields (SEPA, dates, external IDs, import hash)
- [ ] Notes (polymorphic)
- [ ] Attachments (polymorphic)
- [ ] Piggy bank linking

### 5.5 Native Amount Conversion

- [ ] TransactionObserver: `amount` → `native_amount`
- [ ] TransactionObserver: `foreign_amount` → `native_foreign_amount`
- [ ] `convertToPrimary` preference check

### 5.6 Side Effects (after create/update/delete)

- [ ] Apply rules (Phase 8)
- [ ] Recalculate credit/liability balances
- [ ] Fire webhooks (Phase 7)
- [ ] Remove period statistics cache
- [ ] Recalculate running balance

### 5.7 Import Deduplication

- [ ] SHA-256 hash computation (doc 15)
- [ ] Duplicate check via journal_meta
- [ ] Atomic batch rollback on duplicate
- [ ] 422 error response on duplicate

### 5.8 Bulk Operations

- [ ] `POST /api/v1/data/bulk/transactions` — bulk update
- [ ] `DELETE /api/v1/data/destroy` — soft delete
- [ ] `DELETE /api/v1/data/purge` — hard delete

### 5.9 Transaction Links

- [ ] `GET /api/v1/transaction-links` — list
- [ ] `POST /api/v1/transaction-links` — create
- [ ] `GET /api/v1/link-types` — list types

### Referensi

- Doc 06: Database schema
- Doc 07: Business logic (transaction flow)
- Doc 08: API format
- Doc 09: Business flows (transaction lifecycle)
- Doc 10: Response fields
- Doc 11: Validation rules (transaction store/update)
- Doc 12: Error responses (410 gone, 422 duplicate)
- Doc 15: Import deduplication
- Doc 18: Observer side effects

---

## Phase 6 — Search & Autocomplete (Week 13-14)

Target: Full search system dengan 90+ operators.

### 6.1 Search Parser

- [ ] Query parser: `field:value` syntax
- [ ] Bare word search (description)
- [ ] Negation: `-field:value`
- [ ] Value quoting: `field:"value with spaces"`
- [ ] Subquery grouping: `(field1:value1 field2:value2)`

### 6.2 Search Operators (90+)

- [ ] Description operators (is, starts, ends, contains)
- [ ] Account operators (both, source, destination, account number)
- [ ] Amount operators (is, less, more, foreign)
- [ ] Date operators (on, before, after + components)
- [ ] Meta date operators (post-filter)
- [ ] Category, budget, bill operators
- [ ] Tag operators (SQL + post-filter)
- [ ] Notes, currency, type operators
- [ ] External ID/URL/reference operators
- [ ] Attachment operators (post-filter)
- [ ] Reconciliation, ID, existence operators

### 6.3 Search API

- [ ] `POST /api/v1/search/transactions` — search transactions
- [ ] `POST /api/v1/search/transactions/count` — count only

### 6.4 Autocomplete

- [ ] 16 autocomplete endpoints (`/api/v1/autocomplete/*`)

### Referensi

- Doc 13: Search operators (complete catalog)

---

## Phase 7 — Piggy Banks & Recurring (Week 15-16)

Target: Savings goals dan recurring transactions.

### 7.1 Piggy Banks

- [ ] Schema: `piggy_banks`, `account_piggy_bank`, `piggy_bank_events`, `piggy_bank_repetitions`
- [ ] CRUD: `GET/POST/PUT/DELETE /api/v1/piggy-banks`
- [ ] Add/remove money events
- [ ] Target amount → native_target_amount conversion
- [ ] Only Asset, Loan, Debt, Mortgage wallets

### 7.2 Recurring Transactions

- [ ] Schema: `recurrences`, `recurrence_transactions`, `recurrence_meta`, `recurrence_repetitions`
- [ ] CRUD: `GET/POST/PUT/DELETE /api/v1/recurrences`
- [ ] Repetition types: daily, weekly, monthly, quarterly, half-year, yearly
- [ ] Cron job: generate transactions from recurrences
- [ ] Manual trigger: `POST /api/v1/recurrences/{id}/trigger`

### 7.3 Available Budgets

- [ ] Schema: `available_budgets`
- [ ] `GET /api/v1/available-budgets`
- [ ] Recalculation on budget limit change

### Referensi

- Doc 07: Business logic
- Doc 09: Business flows (piggy bank, recurring)
- Doc 18: Observer side effects

---

## Phase 8 — Webhooks (Week 17)

Target: Outbound webhook system.

### 8.1 Schema

- [ ] `webhooks`, `webhook_messages`, `webhook_attempts`

### 8.2 Webhook API

- [ ] CRUD: `GET/POST/PUT/DELETE /api/v1/webhooks`
- [ ] `POST /api/v1/webhooks/{id}/messages` — trigger test
- [ ] `GET /api/v1/webhooks/{id}/messages` — list messages

### 8.3 Webhook Delivery

- [ ] Message queue (async via goroutine/worker)
- [ ] HMAC-SHA3-256 signing
- [ ] Retry logic (exponential backoff: immediate, 30s, 5min — improvement from FF3)
- [ ] Dead letter queue after max attempts
- [ ] Cleanup: delete sent messages after 14 days

### 8.4 Trigger Types

- [ ] STORE_TRANSACTION, UPDATE_TRANSACTION, DESTROY_TRANSACTION
- [ ] STORE_BUDGET, UPDATE_BUDGET, DESTROY_BUDGET
- [ ] STORE_UPDATE_BUDGET_LIMIT

### Referensi

- Doc 14: Webhook payload format

---

## Phase 9 — Rule Engine (Week 18-20)

Target: 22 active rule actions, 90+ triggers.

### 9.1 Schema

- [ ] `rule_groups`, `rules`, `rule_triggers`, `rule_actions`

### 9.2 Rule API

- [ ] CRUD: `GET/POST/PUT/DELETE /api/v1/rules`
- [ ] CRUD: `GET/POST/PUT/DELETE /api/v1/rule-groups`
- [ ] `POST /api/v1/rules/{id}/test` — test trigger
- [ ] `POST /api/v1/rules/trigger` — execute rules on transactions
- [ ] `GET /api/v1/rules/validate-expression` — expression validation

### 9.3 Rule Triggers (90+)

- [ ] Implement all search operators as triggers (reuse Phase 6)
- [ ] Strict vs non-strict mode
- [ ] Stop processing chain

### 9.4 Rule Actions (22 active)

- [ ] set_category, clear_category
- [ ] set_budget, clear_budget
- [ ] add_tag, remove_tag, remove_all_tags
- [ ] set_description, set_notes, clear_notes
- [ ] set_source_account, set_destination_account
- [ ] link_to_bill
- [ ] convert_withdrawal, convert_deposit, convert_transfer
- [ ] switch_accounts
- [ ] update_piggy (add/remove from piggy bank)
- [ ] delete_transaction
- [ ] set_source_to_cash, set_destination_to_cash
- [ ] set_amount

### 9.5 Rule Execution

- [ ] Fire on transaction create/update
- [ ] `fireWebhooks` flag control
- [ ] Rule action failure notifications

### Referensi

- Doc 07: Business logic (rule engine)
- Doc 09: Business flows (rule engine user flow)
- Doc 17: Configuration (rule actions)

---

## Phase 10 — Attachments (Week 21)

Target: File upload/download system.

### 10.1 Schema

- [ ] `attachments` (polymorphic: 9 model types)

### 10.2 Attachment API

- [ ] `GET/POST/PUT/DELETE /api/v1/attachments`
- [ ] `POST /api/v1/attachments/{id}/upload` — upload file
- [ ] `GET /api/v1/attachments/{id}/download` — download file
- [ ] Sub-resource endpoints (per parent model)

### 10.3 Upload System

- [ ] File storage: `storage/upload/at-{id}.data`
- [ ] MIME validation (50+ types)
- [ ] Size validation (1GB max)
- [ ] MD5 hash computation
- [ ] 2-step upload (create metadata → upload content)

### Referensi

- Doc 19: Attachment handling

---

## Phase 11 — Reports, Charts, Insights (Week 22-23)

Target: Dashboard analytics.

### 11.1 Chart API

- [ ] `GET /api/v1/chart/balance/balance`
- [ ] `GET /api/v1/chart/account/overview`
- [ ] `GET /api/v1/chart/budget/overview`
- [ ] `GET /api/v1/chart/category/overview`

### 11.2 Insight API

- [ ] Expense insights (11 endpoints)
- [ ] Income insights (7 endpoints)
- [ ] Transfer insights (6 endpoints)

### 11.3 Summary API

- [ ] `GET /api/v1/summary/basic` — dashboard boxes

### 11.4 Export API

- [ ] `POST /api/v1/export` — generate export
- [ ] `GET /api/v1/export/{id}/download` — download CSV
- [ ] 9 export types (accounts, bills, budgets, categories, piggy banks, recurring, rules, tags, transactions)

### Referensi

- Doc 07: Business logic (reports)
- Doc 09: Business flows (reporting, export)

---

## Phase 12 — Rate Limiting & Notifications (Week 24)

Target: Production-ready security & notifications.

### 12.1 Rate Limiting

- [ ] Global IP rate limit (1000/min, 10000/hour)
- [ ] Per-endpoint rate limits (auth, CRUD, search, export, webhook, upload)
- [ ] Progressive login escalation
- [ ] Redis-backed (sliding window, token bucket)
- [ ] X-RateLimit headers
- [ ] 429 error response

### 12.2 Notifications

- [ ] Email channel (SMTP)
- [ ] Slack channel (webhook)
- [ ] Pushover channel
- [ ] 6 configurable notification types
- [ ] 7 non-configurable security notification types
- [ ] User notification preferences

### 12.3 Cron Jobs

- [ ] Exchange rate download
- [ ] Recurring transaction generation
- [ ] Webhook message delivery
- [ ] Bill reminders
- [ ] Version check

### Referensi

- Doc 16: Rate limiting strategy
- Doc 20: Notification system

---

## Phase 13 — Wallet Sharing (Week 25-26)

Target: Multi-user wallet collaboration.

### 13.1 Schema

- [ ] `wallet_members` table (wallet_id, user_id, role, invited_by)

### 13.2 Sharing API

- [ ] `GET /api/v1/wallets/{id}/members`
- [ ] `POST /api/v1/wallets/{id}/members`
- [ ] `PUT /api/v1/wallets/{id}/members/{userId}`
- [ ] `DELETE /api/v1/wallets/{id}/members/{userId}`
- [ ] `DELETE /api/v1/wallets/{id}/members/me` (leave)
- [ ] `GET /api/v1/wallets/shared-with-me`

### 13.3 Access Control

- [ ] Role hierarchy: owner > editor > viewer
- [ ] Modify all transaction queries (scope to accessible wallets)
- [ ] Modify search, reports, charts (scope to accessible wallets)
- [ ] System wallets cannot be shared

### 13.4 Notifications

- [ ] Wallet shared invitation
- [ ] Member role changed
- [ ] Member removed

### Referensi

- Doc 21: Account sharing analysis

---

## Phase 14 — Import System (Week 27-28)

Target: Native file import (CSV, OFX).

### 14.1 Import API

- [ ] `POST /api/v1/import` — start import job
- [ ] `GET /api/v1/import/{id}/status` — check progress
- [ ] Supported formats: CSV, OFX, CAMT.053

### 14.2 Import Logic

- [ ] File parsing
- [ ] Column mapping
- [ ] Transaction creation via existing store endpoint
- [ ] Deduplication (doc 15)
- [ ] Import account lifecycle

---

## Phase 15 — 2FA, Admin, Polish (Week 29-30)

Target: Production readiness.

### 15.1 Two-Factor Authentication

- [ ] TOTP setup (Google Authenticator compatible)
- [ ] Backup codes
- [ ] 2FA verification middleware
- [ ] MFA events + notifications

### 15.2 Admin API

- [ ] User management CRUD
- [ ] System configuration
- [ ] Data destruction/purge

### 15.3 Polish

- [ ] Audit log
- [ ] API versioning
- [ ] OpenAPI/Swagger documentation
- [ ] Performance benchmarks
- [ ] Load testing

---

## Summary

| Phase | Topic | Weeks | Depends On |
|-------|-------|-------|-----------|
| 0 | Foundation | 1-2 | — |
| 1 | Auth & Users | 3-4 | Phase 0 |
| 2 | Wallets | 5-6 | Phase 1 |
| 3 | Currencies & Exchange Rates | 7 | Phase 1 |
| 4 | Categories, Tags, Bills | 8-9 | Phase 1 |
| **5** | **Transactions (Core)** | **10-12** | **Phase 2, 3, 4** |
| 6 | Search & Autocomplete | 13-14 | Phase 5 |
| 7 | Piggy Banks & Recurring | 15-16 | Phase 5 |
| 8 | Webhooks | 17 | Phase 5 |
| 9 | Rule Engine | 18-20 | Phase 5, 6 |
| 10 | Attachments | 21 | Phase 5 |
| 11 | Reports, Charts, Export | 22-23 | Phase 5 |
| 12 | Rate Limiting & Notifications | 24 | Phase 1 |
| 13 | Wallet Sharing | 25-26 | Phase 5 |
| 14 | Import System | 27-28 | Phase 5 |
| 15 | 2FA, Admin, Polish | 29-30 | All |

### Critical Path

```
Phase 0 → 1 → 2 + 3 + 4 (parallel) → 5 → 6 + 7 + 8 + 9 + 10 (parallel) → 11 → 12 → 13 → 14 → 15
```

### MVP Definition (Phase 0-5)

**After Phase 5** (~12 weeks), system sudah bisa:
- Register/login dengan Keycloak
- CRUD wallets dengan semua type constraints
- CRUD categories, tags, bills
- Create/read/update/delete transactions (full double-entry)
- Multi-currency support
- Search transactions

# Firefly III - Research Documentation

Dokumentasi riset lengkap arsitektur project Firefly III. Digunakan sebagai referensi untuk membangun project baru di Golang (Fiber + sqlx).

## docs/research/ — Arsitektur & Spesifikasi

| # | File | Topik |
|---|------|-------|
| 01 | [01-overview.md](01-overview.md) | Project overview, tech stack, directory structure, design patterns, triple-layer transaction |
| 02 | [02-database-models-relations.md](02-database-models-relations.md) | 51 model, relationships, enums, pivot tables, ER map |
| 03 | [03-api-architecture.md](03-api-architecture.md) | API routes (900 lines), 141 controllers, middleware, auth |
| 04 | [04-auth-rbac.md](04-auth-rbac.md) | Auth guards, two-tier role system, user groups, enforcement mechanisms |
| 05 | [05-refactor-notes-rbac-wallet.md](05-refactor-notes-rbac-wallet.md) | Gap analysis, RBAC per account proposal, phased strategy |
| 06 | [06-database-schema.md](06-database-schema.md) | **Complete DB schema** — 47 tabel, column types, FK, indexes |
| 07 | [07-business-logic.md](07-business-logic.md) | **Core logic** — balance calc, transaction flow, recurring, budgets, piggy banks, rule engine (30 actions, 90+ operators), multi-currency, encryption |
| 08 | [08-api-format.md](08-api-format.md) | **API request/response examples** — JSON:API format, account & transaction payloads, pagination, filtering, errors, auth flow |
| 09 | [09-business-flows.md](09-business-flows.md) | **14 business flows** — onboarding, transaction lifecycle, budget lifecycle, bill/subscription, piggy bank, reconciliation, rule engine user flow, cron/automation, reporting, multi-currency, export, group management, audit trail, import |
| 10 | [10-api-response-fields.md](10-api-response-fields.md) | **Transformer field definitions** — 18 transformers, exact fields/Go types, cross-cutting patterns (CurrencyBlock, PrimaryCurrency, ObjectGroupRef, pc_ prefix, Link), ID type convention |
| 11 | [11-api-validation-rules.md](11-api-validation-rules.md) | **API request validation** — 22 endpoint validation rule sets, conditional rules, custom after-validation chains, 20+ custom rules reference for Go implementation |
| 12 | [12-api-error-responses.md](12-api-error-responses.md) | **Error response catalog** — HTTP status codes (400-500), 5 JSON body shapes, Go struct definitions, demo user protection, accept header enforcement |
| 13 | [13-search-operators.md](13-search-operators.md) | **Search operator catalog** — 90+ operators, SQL WHERE equivalents, 23 categories (description, account, amount, date, category, budget, bill, tag, notes, etc.), alias map |
| 14 | [14-webhook-payload-format.md](14-webhook-payload-format.md) | **Webhook specification** — trigger types, JSON payload, HMAC-SHA3-256 signing, retry logic, message lifecycle, Go implementation examples |
| 15 | [15-import-deduplication.md](15-import-deduplication.md) | **Import deduplication** — SHA-256 hash algorithm, field list, detection flow, duplicate handling (422), atomic rollback, Go implementation |
| 16 | [16-rate-limiting-strategy.md](16-rate-limiting-strategy.md) | **Rate limiting strategy** — vulnerability analysis, limits per category, progressive login escalation, Redis key design, middleware stack, env config, Go implementation |
| 17 | [17-configuration-reference.md](17-configuration-reference.md) | **Complete config reference** — 80+ env vars, feature flags, account type constraints (source→dest matrix), rule actions (22 active), date ranges, API filters/sort |
| 18 | [18-observer-side-effects.md](18-observer-side-effects.md) | **Observer & event side effects** — 9 native amount conversion observers, 15 cascade delete observers, 7 critical event listeners, Go implementation pattern |
| 19 | [19-attachment-handling.md](19-attachment-handling.md) | **Attachment specification** — file storage (at-{id}.data), 50+ MIME whitelist, 1GB max, 2-step upload API, polymorphic to 9 models, cascade delete |
| 20 | [20-notification-system.md](20-notification-system.md) | **Notification & event system** — 4 channels (email/slack/pushover/ntfy), 18 notification types, 51 events, 37 listeners, critical event→listener wiring |
| 21 | [21-account-sharing-analysis.md](21-account-sharing-analysis.md) | **Account→Wallet redesign** — rename terminology, wallet sharing model (owner/editor/viewer), wallet_members table, sharing API, access control, phased implementation |
| 22 | [22-database-diagrams.md](22-database-diagrams.md) | **Mermaid ERD diagrams** — full 47-table ERD, triple-layer transaction detail, UserGroup RBAC, wallet sharing model, polymorphic relations, webhook lifecycle, rule engine flow, side effects, auth sequence, domain grouping |

## docs/refactor/ — Analisis Rewrite ke Golang

| # | File | Topik |
|---|------|-------|
| 01 | [01-golang-fiber-sqlx-analysis.md](01-golang-fiber-sqlx-analysis.md) | Go vs PHP comparison, arsitektur yang diusulkan, library mapping, tantangan (Passport compat, rule engine, bcmath), estimasi effort, strategi migration |
| 02 | [02-oauth-token-compatibility.md](02-oauth-token-compatibility.md) | **OAuth token structure** — JWT (RS256, RSA 4096-bit), DB schema (6 tables), expiry settings, guard config, grant types, scopes, revocation, Go compatibility strategy (3 options) |
| 03 | [03-implementation-roadmap.md](03-implementation-roadmap.md) | **Implementation roadmap** — 15 phases, 30 weeks, critical path, MVP definition (Phase 0-5), deliverables per phase, dependency graph |

## Temuan Kunci

- **UserGroup** adalah unit scoping utama — semua data milik group, bukan user langsung
- **RBAC sudah ada** di level group (22 permissions), tapi **belum ada di level account**
- **Account saat ini single-owner** (`user_id`), tidak ada sharing mechanism
- **Double-entry bookkeeping** — triple-layer: TransactionGroup → TransactionJournal → Transaction (2 per journal)
- **bcmath everywhere** — semua amount sebagai string, jangan pernah pakai float64
- **Multi-currency pervasive** — setiap monetary model punya `amount` + `native_amount` (auto-calc via observer)
- **Rule engine** — 30 action types, 90+ search operators, strict/non-strict mode, stop processing chain
- **Encryption** — shadow column pattern (plaintext + encrypted backup via AES-256-CBC)
- **JSON:API** — League Fractal serialization, nested transaction groups
- **Business flows** — 14 alur bisnis terdokumentasi lengkap (bukan hanya teknikal)
- **OAuth token** — JWT (RS256, RSA 4096-bit) + DB storage (database-backed JWT), 14-day expiry, no scopes defined
- **Transformer fields** — 18 transformers dengan exact field definitions siap untuk Go response structs
- **Validation rules** — 22 endpoint validation sets dengan 20+ custom rules yang perlu di-port ke Go
- **ID type inconsistency** — beberapa transformer pakai `string`, lainnya `int` untuk ID field
- **Error responses** — 5 JSON body shapes, 422 validation paling umum, demo user 403, rule-deleted transaction 410
- **Search operators** — 90+ operators, 23 kategori, sebagian SQL-based dan sebagian post-filter (meta dates, tags, attachments, balances)
- **Webhook signing** — HMAC-SHA3-256, format `t={timestamp},v1={hmac}`, retry max 3x tanpa exponential backoff
- **Import dedup** — SHA-256 exact-match hash, disabled by default (`error_if_duplicate_hash=false`), atomic batch rollback
- **Rate limiting** — Firefly III hampir tidak punya rate limiting (hanya 3 endpoint), komprehensif strategy dirancang untuk Go API
- **Configuration** — 80+ env vars, account type constraint matrices (account_to_transaction, source_dests, allowed_opposing_types), feature flags
- **Observer side effects** — 9 native amount conversion, 15 cascade delete, 7 critical event listeners — HARUS replicate exact behavior di Go
- **Attachment** — stored as `at-{id}.data`, 50+ MIME types, 1GB max, 2-step API (create metadata → upload file)
- **Notifications** — 4 channels, 18 notification types, 51 events, 37 listeners
- **Account→Wallet** — rename terminology, sharing model (owner/editor/viewer), wallet_members table proposal

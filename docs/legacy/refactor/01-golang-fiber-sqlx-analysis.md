# Refactor ke Golang: Analisis & Diskusi

## Pertanyaan Fundamental

> **Apakah ini rewrite dari nol, atau incremental migration?**

Ini adalah pertanyaan paling penting. Jawabannya menentukan semua keputusan arsitektural selanjutnya.

---

## 1. Kenapa Golang? (Why)

### Keuntungan

| Aspek | PHP/Laravel Saat Ini | Golang (Fiber + sqlx) |
|-------|---------------------|----------------------|
| **Performance** | ~500-1000 req/s (shared nothing) | ~10,000-50,000 req/s |
| **Memory** | ~30-80MB per worker | ~10-30MB total (binary) |
| **Deployment** | PHP-FPM + Nginx + Composer | Single binary |
| **Concurrency** | Synchronous (worker-based) | Goroutines (native async) |
| **Startup time** | ~500ms-2s | ~10ms |
| **Binary size** | ~100MB vendor + PHP runtime | ~15-30MB single binary |
| **Docker image** | ~200-400MB | ~20-50MB (scratch/alpine) |
| **Type safety** | Dynamic + PHPStan | Static + compiler |

### Kerugian

| Aspek | PHP/Laravel | Golang |
|-------|-------------|--------|
| **Development speed** | Cepat (convention over config) | Lebih lambat (boilerplate) |
| **ORM/Eloquent** | Mature, feature-rich | sqlx = raw SQL (no ORM) |
| **Ecosystem** | Massive (Composer) | Growing tapi lebih kecil |
| **Developer pool** | Besar | Lebih kecil |
| **Rapid prototyping** | Sangat cepat | Lebih rigid |
| **Template engine** | Twig (mature) | Perlu pilih (html/template, templ, etc.) |

---

## 2. Arsitektur yang Diusulkan

### Layer Diagram

```
┌─────────────────────────────────────────────────────┐
│                    HTTP Layer                        │
│              Fiber (framework)                        │
│         Routes / Middleware / Handlers                │
├─────────────────────────────────────────────────────┤
│                  Transport Layer                      │
│           Request/Response DTO                       │
│         JSON:API Serialization                       │
│         Validation (go-playground/validator)          │
├─────────────────────────────────────────────────────┤
│                   Service Layer                      │
│           Business Logic (pure Go)                   │
│         Interfaces (ports)                           │
├─────────────────────────────────────────────────────┤
│                 Repository Layer                     │
│           sqlx (database access)                     │
│         Raw SQL queries + scanning                   │
├─────────────────────────────────────────────────────┤
│                  Infrastructure                       │
│         PostgreSQL / MySQL                           │
│         Redis (cache/queue)                          │
└─────────────────────────────────────────────────────┘
```

### Project Structure

```
firefly-iii-go/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── config/                  # Configuration
│   │   └── config.go
│   ├── middleware/               # Fiber middleware
│   │   ├── auth.go              # JWT/Passport-compatible auth
│   │   ├── rbac.go              # Role-based access control
│   │   ├── ratelimit.go         # Rate limiting
│   │   └── binder.go            # Route model binding
│   ├── domain/                  # Domain models (pure structs, no deps)
│   │   ├── user.go
│   │   ├── account.go
│   │   ├── transaction.go
│   │   ├── budget.go
│   │   ├── bill.go
│   │   └── usergroup.go
│   ├── dto/                     # Data Transfer Objects
│   │   ├── request/
│   │   │   ├── account.go
│   │   │   ├── transaction.go
│   │   │   └── auth.go
│   │   └── response/
│   │       ├── jsonapi.go       # JSON:API serializer
│   │       ├── account.go
│   │       └── transaction.go
│   ├── repository/              # Data access (sqlx)
│   │   ├── postgres/
│   │   │   ├── user_repo.go
│   │   │   ├── account_repo.go
│   │   │   ├── transaction_repo.go
│   │   │   └── migrations/
│   │   └── repository.go        # Interfaces
│   ├── service/                 # Business logic
│   │   ├── auth_service.go
│   │   ├── account_service.go
│   │   ├── transaction_service.go
│   │   ├── budget_service.go
│   │   └── rbac_service.go
│   ├── handler/                 # HTTP handlers (Fiber)
│   │   ├── account_handler.go
│   │   ├── transaction_handler.go
│   │   ├── auth_handler.go
│   │   └── usergroup_handler.go
│   └── router/
│       └── router.go            # Route registration
├── pkg/                         # Shared utilities
│   ├── pagination/
│   ├── currency/
│   ├── bcrypt/
│   └── jwt/
├── migrations/                  # SQL migration files
├── go.mod
├── go.sum
├── Makefile
└── Dockerfile
```

---

## 3. Mapping: Laravel → Golang

### Model Mapping

| Laravel (Eloquent) | Golang (struct) | Catatan |
|--------------------|-----------------|---------|
| `$fillable` | Struct fields + json tags | Manual definition |
| `$casts` | Custom types | Misal: `Decimal string`, `Time time.Time` |
| `$hidden` | `json:"-"` tag | Hide sensitive fields |
| `SoftDeletes` | `deleted_at *time.Time` + WHERE filter | Manual di query |
| `HasMany` | `[]RelatedStruct` | Manual join/query |
| `BelongsTo` | `*RelatedStruct` + `ID int` | Manual join/query |
| `BelongsToMany` | `[]RelatedStruct` via pivot query | Manual join/query |
| `MorphMany` | Interface + type column | Manual polymorphic query |
| `Scopes` | Query builder methods | Manual WHERE clauses |
| `Observers` | Service layer hooks | Pre/post save logic |

### Contoh: User Model

```go
// internal/domain/user.go
type User struct {
    ID          int64          `json:"id" db:"id"`
    Email       string         `json:"email" db:"email"`
    Password    string         `json:"-" db:"password"`
    Blocked     bool           `json:"blocked" db:"blocked"`
    BlockedCode string         `json:"blocked_code,omitempty" db:"blocked_code"`
    UserGroupID int64          `json:"user_group_id" db:"user_group_id"`
    CreatedAt   time.Time      `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
    DeletedAt   *time.Time     `json:"-" db:"deleted_at"`
}
```

### Contoh: Transaction Triple-Layer

```go
// internal/domain/transaction.go
type TransactionGroup struct {
    ID         int64    `json:"id" db:"id"`
    UserID     int64    `json:"user_id" db:"user_id"`
    UserGroupID int64   `json:"user_group_id" db:"user_group_id"`
    Title      string   `json:"title" db:"title"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
    // relations
    Journals   []TransactionJournal `json:"journals,omitempty"`
}

type TransactionJournal struct {
    ID                    int64     `json:"id" db:"id"`
    TransactionGroupID    int64     `json:"transaction_group_id" db:"transaction_group_id"`
    TransactionTypeID     int64     `json:"transaction_type_id" db:"transaction_type_id"`
    TransactionCurrencyID int64     `json:"transaction_currency_id" db:"transaction_currency_id"`
    BillID                *int64    `json:"bill_id,omitempty" db:"bill_id"`
    Description           string    `json:"description" db:"description"`
    Date                  time.Time `json:"date" db:"date"`
    Order                 int       `json:"order" db:"order"`
    Completed             bool      `json:"completed" db:"completed"`
    // relations
    Transactions          []Transaction `json:"transactions,omitempty"`
    Tags                  []Tag         `json:"tags,omitempty"`
    Budgets               []Budget      `json:"budgets,omitempty"`
    Categories            []Category    `json:"categories,omitempty"`
}

type Transaction struct {
    ID                    int64   `json:"id" db:"id"`
    AccountID             int64   `json:"account_id" db:"account_id"`
    TransactionJournalID  int64   `json:"transaction_journal_id" db:"transaction_journal_id"`
    TransactionCurrencyID int64   `json:"transaction_currency_id" db:"transaction_currency_id"`
    ForeignCurrencyID     *int64  `json:"foreign_currency_id,omitempty" db:"foreign_currency_id"`
    Amount                string  `json:"amount" db:"amount"`               // bcmath string
    NativeAmount          string  `json:"native_amount" db:"native_amount"` // auto-calc
    ForeignAmount         *string `json:"foreign_amount,omitempty" db:"foreign_amount"`
    Reconciled            bool    `json:"reconciled" db:"reconciled"`
    Description           string  `json:"description" db:"description"`
}
```

### Repository Mapping

| Laravel | Golang (sqlx) |
|---------|--------------|
| `Account::where('user_id', $id)->get()` | `sqlx.Select(&accounts, "SELECT * FROM accounts WHERE user_id = $1", id)` |
| `$account->transactions()->paginate(50)` | Manual pagination + `LIMIT/OFFSET` |
| `$user->accounts()->create([...])` | `sqlx.NamedExec("INSERT INTO accounts (...) VALUES (...)")` |
| `$account->load('accountType')` | Manual JOIN atau separate query |
| `$account->delete()` | `sqlx.Exec("UPDATE accounts SET deleted_at = NOW() WHERE id = $1", id)` |

### Contoh Repository

```go
// internal/repository/postgres/account_repo.go
type AccountRepository struct {
    db *sqlx.DB
}

func (r *AccountRepository) FindByID(ctx context.Context, id int64) (*domain.Account, error) {
    var account domain.Account
    err := r.db.GetContext(ctx, &account,
        `SELECT * FROM accounts WHERE id = $1 AND deleted_at IS NULL`, id)
    if err != nil {
        return nil, fmt.Errorf("account not found: %w", err)
    }
    return &account, nil
}

func (r *AccountRepository) FindByUserGroup(ctx context.Context, userGroupID int64, page, limit int) ([]domain.Account, int, error) {
    var total int
    err := r.db.GetContext(ctx, &total,
        `SELECT COUNT(*) FROM accounts WHERE user_group_id = $1 AND deleted_at IS NULL`, userGroupID)
    if err != nil {
        return nil, 0, err
    }

    var accounts []domain.Account
    offset := (page - 1) * limit
    err = r.db.SelectContext(ctx, &accounts,
        `SELECT * FROM accounts WHERE user_group_id = $1 AND deleted_at IS NULL
         ORDER BY id ASC LIMIT $2 OFFSET $3`, userGroupID, limit, offset)
    return accounts, total, err
}
```

---

## 4. Komponen Library Golang

| Kebutuhan | Library | Status |
|-----------|---------|--------|
| **HTTP Framework** | `github.com/gofiber/fiber/v2` | Mature, Express-like, very fast |
| **SQL** | `github.com/jmoiron/sqlx` | Mature, lightweight SQL wrapper |
| **DB Driver** | `github.com/lib/pq` (Postgres) / `github.com/go-sql-driver/mysql` | Standard |
| **Validation** | `github.com/go-playground/validator/v10` | Most popular Go validator |
| **JWT** | `github.com/golang-jwt/jwt/v5` | Standard |
| **BCrypt** | `golang.org/x/crypto/bcrypt` | Standard library |
| **Config** | `github.com/spf13/viper` | Popular config management |
| **Migrations** | `github.com/golang-migrate/migrate/v4` | DB-agnostic migrations |
| **Logging** | `log/slog` (stdlib Go 1.21+) | Built-in structured logging |
| **Testing** | `testing` + `testify` + `testcontainers-go` | Standard + mocks |
| **UUID** | `github.com/google/uuid` | Standard |
| **Decimal** | `github.com/shopspring/decimal` | Critical untuk financial (bcmath equivalent) |
| **Time/Period** | `github.com/spf13/cast` atau custom | Perlu replikasi `spatie/period` |
| **Cron** | `github.com/robfig/cron/v3` | Cron scheduling |

---

## 5. Tantangan Besar

### 5.1 Bcrypt Password Hash Compatibility

**Problem**: Password di-hash dengan Laravel's bcrypt. Go harus bisa verify hash yang sama.

**Solution**: `golang.org/x/crypto/bcrypt` menggunakan format hash yang sama ($2y$). Langsung compatible.

```go
err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainPassword))
```

### 5.2 Passport OAuth2 Token Compatibility

**Problem**: Existing clients menggunakan Passport tokens. Apakah harus compatible?

**Options**:
- **A) Full compatibility**: Implement Passport's token storage & validation di Go (baca tabel `oauth_access_tokens` dari DB)
- **B) Migration**: Minta semua client re-authenticate dengan JWT baru (breaking change)
- **C) Hybrid**: Support both Passport (read existing) + JWT (new tokens) selama transition

**Rekomendasi**: Option A untuk transition, lalu gradual ke JWT native.

### 5.3 Triple-Layer Transaction

**Problem**: Firefly III punya 3-layer transaction (Group → Journal → Transaction). Kompleks tapi fundamental.

**Solution**: Implement sebagai service layer dengan transaction (DB transaction):

```go
func (s *TransactionService) Create(ctx context.Context, input CreateTransactionInput) (*TransactionGroup, error) {
    tx, err := s.db.BeginTxx(ctx, nil)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    // 1. Create TransactionGroup
    group := &TransactionGroup{Title: input.Title, ...}
    // INSERT INTO transaction_groups ...

    // 2. Create TransactionJournal
    journal := &TransactionJournal{TransactionGroupID: group.ID, ...}
    // INSERT INTO transaction_journals ...

    // 3. Create Transactions (source + destination)
    source := &Transaction{JournalID: journal.ID, AccountID: input.SourceID, Amount: "-50.00", ...}
    dest := &Transaction{JournalID: journal.ID, AccountID: input.DestID, Amount: "50.00", ...}
    // INSERT INTO transactions ...

    // 4. Calculate native amounts (currency conversion)
    // ...

    return tx.Commit()
}
```

### 5.4 Rule Engine (ExpressionLanguage)

**Problem**: Firefly III menggunakan Symfony ExpressionLanguage untuk rule engine. ~31 action implementations.

**Options**:
- **A) Port ke Go**: Implement expression evaluator di Go (complex)
- **B) Embed PHP**: Gunakan embedded PHP interpreter (go-php) untuk rule evaluation
- **C) Replace**: Implement rule engine baru di Go dengan syntax berbeda
- **D) WASM**: Compile PHP rules ke WebAssembly dan jalankan di Go

**Rekomendasi**: Option C — tulis ulang rule engine di Go. Rules bisa di-migrate dari storage format yang sama.

### 5.5 Bcmath (Decimal Precision)

**Problem**: PHP `bcmath` menyimpan amount sebagai string. Go perlu `shopspring/decimal` untuk presisi financial.

```go
import "github.com/shopspring/decimal"

amount, _ := decimal.NewFromString("-50.00")
nativeAmount := amount.Mul(exchangeRate) // currency conversion
```

### 5.6 JSON:API Serialization

**Problem**: API saat ini menggunakan League Fractal + JsonApiSerializer. Client mengharapkan format ini.

**Solution**: Implement JSON:API serializer di Go:

```go
type JSONAPIResponse struct {
    Data    interface{} `json:"data"`
    Meta    *Meta       `json:"meta,omitempty"`
    Links   *Links      `json:"links,omitempty"`
    Errors  []Error     `json:"errors,omitempty"`
}

type JSONAPIResource struct {
    Type       string      `json:"type"`
    ID         string      `json:"id"`
    Attributes interface{} `json:"attributes"`
    Links      *SelfLinks  `json:"links,omitempty"`
}
```

### 5.7 Soft Deletes Everywhere

**Problem**: Hampir semua model menggunakan soft deletes. Setiap query perlu `WHERE deleted_at IS NULL`.

**Solution**: Query builder helper:

```go
func withSoftDelete(baseQuery string) string {
    return baseQuery + " AND deleted_at IS NULL"
}
```

---

## 6. Strategi Migrasi

### Option A: Big Bang Rewrite (Tidak Direkomendasikan)

- Tulis ulang semua dari nol di Go
- Switch sekaligus
- **Risk**: Sangat tinggi, banyak bug, downtime lama

### Option B: Strangler Fig Pattern (Direkomendasikan)

- Go service jalan **berdampingan** dengan Laravel
- Migrate endpoint per endpoint
- Gunakan **shared database** (sama)
- Reverse proxy route traffic berdasarkan maturity

```
Client
  │
  ├── /api/v1/accounts ──→ Go service (migrated)
  ├── /api/v1/budgets ───→ Go service (migrated)
  ├── /api/v1/rules ─────→ Laravel (belum migrated)
  └── /api/v1/webhooks ──→ Laravel (belum migrated)
```

**Urutan Migrasi** (dari paling simple ke paling complex):

```
Phase 1: Auth & User Management
  ├── POST /oauth/token (Passport compatibility)
  ├── GET /api/v1/about
  ├── GET /api/v1/user-groups
  └── GET /api/v1/preferences

Phase 2: Read Endpoints (CRUD tanpa complex logic)
  ├── GET/POST/PUT/DELETE /api/v1/accounts
  ├── GET/POST/PUT/DELETE /api/v1/categories
  ├── GET/POST/PUT/DELETE /api/v1/tags
  └── GET/POST/PUT/DELETE /api/v1/currencies

Phase 3: Transaction Core
  ├── GET/POST/PUT/DELETE /api/v1/transactions
  ├── GET /api/v1/transaction-journals
  └── GET /api/v1/autocomplete/*

Phase 4: Complex Domains
  ├── GET/POST/PUT/DELETE /api/v1/budgets (dengan limits)
  ├── GET/POST/PUT/DELETE /api/v1/bills
  ├── GET/POST/PUT/DELETE /api/v1/piggy-banks
  └── GET/POST/PUT/DELETE /api/v1/recurrences

Phase 5: Advanced Features
  ├── Rule Engine
  ├── Webhooks
  ├── Charts & Insights
  ├── Data Export/Import
  └── Search

Phase 6: Cleanup
  └── Decommission Laravel
```

### Option C: API Gateway + Microservices

- Go sebagai API gateway
- Laravel sebagai legacy service
- Gateway route ke service yang sesuai
- **Overkill** untuk project ini, tapi scalable ke depan

---

## 7. Estimasi Effort

| Area | Kompleksitas | Estimasi |
|------|-------------|----------|
| Project setup + config | Low | 1-2 hari |
| Auth (Passport compat + JWT) | Medium | 3-5 hari |
| Domain models (51 structs) | Low | 2-3 hari |
| Account CRUD + RBAC | Medium | 3-5 hari |
| Transaction (triple-layer) | **High** | 7-10 hari |
| Budget + Bill + PiggyBank | Medium | 5-7 hari |
| Rule Engine | **Very High** | 10-15 hari |
| JSON:API serialization | Medium | 3-5 hari |
| Charts & Insights | Medium | 5-7 hari |
| Webhooks + Cronjobs | Medium | 3-5 hari |
| Testing (unit + integration) | Medium | Ongoing |
| **Total estimasi** | | **~50-70 hari kerja** |

---

## 8. Database: Tetap Sama atau Migrate?

### Opsi: Tetap Pakai Database yang Sama

**Direkomendasikan**. Schema tidak perlu berubah. Go service baca/tulis ke database yang sama dengan Laravel.

Keuntungan:
- Strangler Fig pattern langsung bisa diterapkan
- Data tidak perlu migration
- Bisa A/B testing per endpoint

Yang perlu dihandle:
- `encrypted` columns (AES encryption) — implement di Go
- `deleted_at` soft deletes — query filter
- `created_at` / `updated_at` — auto-manage di Go
- Laravel-specific timestamp format

---

## 9. Pertanyaan untuk Diskusi

### Kritis

1. **Apakah ini full rewrite atau strangler fig migration?** — Ini menentukan timeline dan risk.

2. **Apakah API contract harus 100% compatible?** — Jika iya, JSON:API format harus di-maintain. Jika tidak, bisa simpler.

3. **Apakah web UI (Twig) ikut di-rewrite?** — Atau Go hanya handle API, dan frontend tetap pakai existing V1/V2?

4. **Rule Engine: port atau replace?** — Ini adalah bagian paling complex. ~31 action types + expression language.

5. **Apakah multi-DB support tetap diperlukan?** — Saat ini support MySQL/PostgreSQL/SQLite. di Go dengan sqlx bisa, tapi lebih banyak maintenance.

### Teknis

6. **Go version?** — Rekomendasi Go 1.22+ (generics, slog, improved tooling)

7. **Deployment target?** — Docker, Kubernetes, bare metal? Ini mempengaruhi binary build strategy.

8. **Testing strategy?** — Integration tests dengan testcontainers-go? Atau mock-based?

9. **Apakah perlu API versioning baru (v2)?** — Ini kesempatan untuk clean break dari JSON:API ke format simpler.

10. **Apakah RBAC per account (dari docs/research) diimplementasikan bareng refactor?** — Atau terpisah?

---

## 10. Rekomendasi

1. **Gunakan Strangler Fig Pattern** — Jangan big bang rewrite
2. **Tetap pakai database yang sama** — Shared DB antara Go dan Laravel selama transition
3. **Go hanya untuk API** — Frontend tetap consume API, tidak perlu port Twig views
4. **Implement RBAC per account bareng refactor** — Lebih mudah di Go karena typesafe
5. **Mulai dari auth + simple CRUD** — Bangun confidence sebelum tackle complex domains
6. **Consider API v2** — Kesempatan untuk meninggalkan JSON:API complexity
7. **`shopspring/decimal` wajib** — Jangan pernah pakai float64 untuk uang
8. **sqlx over ORM** — Sesuai request, dan memang lebih cocok untuk financial app yang perlu kontrol query precision

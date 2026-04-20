# Firefly III - Project Overview & Architecture

## Apa itu Firefly III?

Firefly III adalah **personal finance manager** yang open-source dan self-hosted. Dibangun dengan PHP/Laravel, aplikasi ini membantu pengguna melacak pemasukan, pengeluaran, anggaran (budget), dan laporan keuangan.

- **Repository**: https://github.com/firefly-iii/firefly-iii
- **Lisensi**: AGPL-3.0-or-later
- **Versi**: 6.5.9
- **Dokumentasi**: https://docs.firefly-iii.org/

---

## Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| **Backend** | PHP 8.5+, Laravel 12 |
| **Database** | MySQL / PostgreSQL / SQLite |
| **Cache/Queue** | Redis (opsional), sync (default) |
| **API Auth** | Laravel Passport (OAuth2) |
| **API Serialization** | League Fractal (JSON:API spec) |
| **Template Engine** | Twig (server-side rendering) |
| **Frontend V1** | Vue.js 2.7 + Bootstrap 3 + Chart.js (Laravel Mix) |
| **Frontend V2** | Alpine.js 3 + Bootstrap 5 + AdminLTE 4 + Chart.js 4 (Vite) |
| **Rule Engine** | Symfony ExpressionLanguage |
| **Testing** | PHPUnit 12, PHPStan 2 |

---

## Struktur Direktori

```
app/
├── Api/V1/                  # API controllers, requests, middleware
│   ├── Controllers/         # 141 controller files (CRUD per model)
│   ├── Requests/            # 72 FormRequest validation classes
│   └── Middleware/          # ApiDemoUser middleware
├── Console/Commands/        # Artisan commands (correction, export, integrity)
├── Enums/                   # 11 enum files (domain bounded contexts)
├── Events/                  # Domain events (Model, Security, Admin, System)
├── Handlers/
│   ├── ExchangeRate/        # Currency conversion logic
│   └── Observer/            # 19 Eloquent model observers
├── Helpers/                 # Attachments, Collector, Fiscal, Report
├── Http/
│   ├── Controllers/         # Web controllers (33 subdirectories)
│   ├── Middleware/          # Auth, Admin, Binder, Demo, etc.
│   └── Requests/            # 37 web FormRequest classes
├── Jobs/                    # 7 queued jobs
├── Listeners/               # Event listeners (mirrors event structure)
├── Models/                  # 50 Eloquent model files
├── Repositories/            # 22 repository subdirectories (Interface + Impl)
├── Rules/                   # 28 custom Laravel validation rules
├── Services/                # Update checking, Password breach, Webhook
├── Support/                 # 35 utility classes (Steam, Amount, Navigation, etc.)
│   ├── Binder/              # Route model binding classes
│   ├── Cronjobs/            # 6 cronjob implementations
│   ├── Facades/             # 11 facades
│   ├── Http/Api/            # ValidatesUserGroupTrait
│   └── Request/             # ChecksLogin trait
└── Transformers/            # 26 Fractal transformers (JSON:API)

config/                      # 37 configuration files
database/
├── migrations/              # 60 migration files
├── seeders/                 # Database seeders
└── factories/               # Model factories
resources/
├── assets/v1/               # Vue.js frontend
├── assets/v2/               # Alpine.js frontend
├── views/                   # Twig templates (43+ subdirectories)
└── lang/                    # 35+ locale translations
routes/
├── api.php                  # REST API routes (~900 lines)
├── web.php                  # Web UI routes (~1900 lines)
├── console.php              # Artisan console routes
└── channels.php             # Broadcasting channels
```

---

## Design Patterns

### 1. Repository Pattern
Setiap domain punya `Interface` + `Implementation`, di-bind via Service Provider dengan user context injection.

```php
// app/Providers/AccountServiceProvider.php
$this->app->bind(static function (Application $app): AccountRepositoryInterface {
    $repository = app(AccountRepository::class);
    if ($app->auth->check()) {
        $repository->setUser(auth()->user());
    }
    return $repository;
});
```

### 2. Observer Pattern
19 Eloquent observers untuk cross-cutting concerns:
- `TransactionObserver` - auto-convert `amount` ke `native_amount` (currency conversion)
- `DeletedTransactionJournalObserver` - cascade cleanup
- `DeletedAccountObserver`, `DeletedTransactionGroupObserver` - cascade cleanup

### 3. Strategy Pattern (Rule Engine)
`ActionFactory` memetakan string action types ke concrete action classes. 31 action implementations.

### 4. Template Method Pattern (Cronjobs)
`AbstractCronjob` mendefinisikan template; subclass mengimplementasikan `fire()`.

### 5. Facade Pattern
11 facades wrapping support classes: `Steam`, `Amount`, `Preferences`, `Navigation`, dll.

### 6. Route Model Binding with Authorization
Model mendefinisikan `routeBinder()` static methods yang sekaligus melakukan authorization check.

---

## Service Provider Architecture

22 service providers, masing-masing per-domain:

| Provider | Bind |
|----------|------|
| `AccountServiceProvider` | `AccountRepositoryInterface`, `OperationsRepositoryInterface`, `AccountTaskerInterface` |
| `BudgetServiceProvider` | Budget repositories |
| `CategoryServiceProvider` | Category repositories |
| `JournalServiceProvider` | Journal repositories (CRUD, API, CLI variants) |
| `BillServiceProvider` | Bill repositories |
| `PiggyBankServiceProvider` | Piggy bank repositories |
| `RecurringServiceProvider` | Recurrence repositories |
| `RuleServiceProvider` | Rule repositories |
| `FireflyServiceProvider` | 30+ services, singletons (ExpressionLanguage, etc.) |

---

## Core Domain Model: Triple-Layer Transaction

Ini adalah konsep arsitektural paling penting di Firefly III:

```
TransactionGroup (1) ──has-many──> (N) TransactionJournal (1) ──has-many──> (2+) Transaction
     │                                     │                                      │
     │  title                              │  date, description, type             │  account_id, amount
     │  user_id                            │  transaction_currency_id            │  transaction_currency_id
     │  user_group_id                      │  tags, budgets, categories           │  foreign_currency_id
     │                                     │                                      │  reconciled
```

**Contoh: Withdrawal $50 dari "Checking" ke "Groceries"**

1. **TransactionGroup**: title = "Grocery shopping"
2. **TransactionJournal**: type = Withdrawal, date = today, description = "Grocery shopping"
3. **Transaction 1**: account = Checking (Asset), amount = "-50.00"
4. **Transaction 2**: account = Groceries (Expense), amount = "50.00"

**Contoh: Transfer $100 dari "Checking" ke "Savings"**

1. **TransactionGroup**: title = "Transfer to savings"
2. **TransactionJournal 1**: type = Withdrawal, source = Checking, amount = -100
3. **TransactionJournal 2**: type = Deposit, destination = Savings, amount = +100

---

## Event System

### Events (`app/Events/`)

| Kategori | Lokasi | Contoh |
|----------|--------|--------|
| **Model** | `Events/Model/` (10 subdirs) | `TransactionGroup/CreatedSingleTransactionGroup` |
| **Security** | `Events/Security/System/` + `User/` | `UserLoggedIn`, `UserFailedLoginAttempt` |
| **Admin** | `Events/Admin/` | `InvitationCreated` |
| **System** | `Events/System/` | System-level events |

### Listeners (`app/Listeners/`)
Mirror event structure. Contoh: `ProcessesNewTransactionGroup`, `SendsWebhookMessages`, `HandlesNewUserRegistration`.

### Jobs (`app/Jobs/`)
7 queued jobs: `CreateAutoBudgetLimits`, `CreateRecurringTransactions`, `DownloadExchangeRates`, `MailError`, `SendWebhookMessage`, `WarnAboutBills`.

### Cronjobs (`app/Support/Cronjobs/`)
6 implementations: AutoBudget, BillWarning, ExchangeRates, Recurring, UpdateCheck, Webhook. Default interval: 12 hours.

# API Architecture

## API Versioning

URL-based versioning dengan satu versi: **v1**.

- Semua route di `routes/api.php` di-prefix `v1/`
- Base URL: `/api/v1/{resource}`
- Fractal serializer base URL: `{scheme}://{host}/api/v1`

---

## Authentication: Laravel Passport (OAuth2)

```php
// config/auth.php
'api' => [
    'driver'   => 'passport',
    'provider' => 'users',
],
```

**OAuth Routes** (`routes/web.php`):
- `POST /oauth/token` — token issuance (throttled)
- `GET /oauth/authorize` — authorization (user-full-auth)
- Token refresh, client management, personal access tokens

> **Catatan**: Tidak ada API scopes yang didefinisikan atau di-enforce. Token OAuth diterbitkan tanpa scope restrictions.

---

## Middleware Stack

### Global `api` middleware group
```
AcceptHeaders → auth:api (Passport) → Binder
```

### `api-admin` middleware group
```
IsAdminApi (checks 'owner' role)
```

### Individual Auth Middleware

| Middleware | File | Purpose |
|-----------|------|---------|
| `Authenticate` | `app/Http/Middleware/Authenticate.php` | Validasi user login, cek blocked |
| `IsAdminApi` | `app/Http/Middleware/IsAdminApi.php` | Requires global `owner` role (throws AuthorizationException) |
| `Binder` | `app/Http/Middleware/Binder.php` | Route model binding dengan authorization check |
| `AcceptHeaders` | `app/Http/Middleware/AcceptHeaders.php` | Validasi Accept/Content-Type headers |

---

## API Endpoints

Semua route di `routes/api.php` (~900 lines).

### Cron & System
| Method | Endpoint | Middleware | Purpose |
|--------|----------|-----------|---------|
| GET | `/v1/cron/{cliToken}` | Binder, AcceptHeaders (no auth) | Trigger cron jobs |
| GET | `/v1/about` | api | System info + current user |
| GET | `/v1/configuration` | api | Read config |
| PUT | `/v1/configuration` | **api-admin** | Update config |
| GET/POST/PUT/DELETE | `/v1/users/*` | **api-admin** | User CRUD |

### Autocomplete
| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/v1/autocomplete/accounts` | Account dropdown |
| GET | `/v1/autocomplete/bills` | Bill dropdown |
| GET | `/v1/autocomplete/budgets` | Budget dropdown |
| GET | `/v1/autocomplete/categories` | Category dropdown |
| GET | `/v1/autocomplete/currencies` | Currency dropdown |
| GET | `/v1/autocomplete/piggy-banks` | Piggy bank dropdown |
| GET | `/v1/autocomplete/tags` | Tag dropdown |
| GET | `/v1/autocomplete/transaction-journals` | Transaction dropdown |
| GET | `/v1/autocomplete/transaction-types` | Transaction type dropdown |

### Charts
| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/v1/chart/balance` | Balance chart |
| GET | `/v1/chart/account` | Account overview chart |
| GET | `/v1/chart/budget` | Budget chart |
| GET | `/v1/chart/category` | Category chart |

### Insights & Summary
| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/v1/insight/expense/*` | Expense analytics |
| GET | `/v1/insight/income/*` | Income analytics |
| GET | `/v1/insight/transfer/*` | Transfer analytics |
| GET | `/v1/summary/*` | Basic summary |

### CRUD Resources
| Method | Endpoint | Resource |
|--------|----------|----------|
| GET/POST | `/v1/accounts` | Accounts |
| GET/PUT/DELETE | `/v1/accounts/{account}` | Single account |
| GET/POST | `/v1/attachments` | Attachments |
| GET/POST/PUT/DELETE | `/v1/bills/*` | Bills |
| GET/POST/PUT/DELETE | `/v1/subscriptions/*` | Subscriptions |
| GET | `/v1/available-budgets` | Available budgets |
| GET/POST | `/v1/budgets` | Budgets |
| GET/PUT/DELETE | `/v1/budgets/{budget}` | Single budget |
| GET | `/v1/budget-limits` | Budget limits |
| GET/POST | `/v1/categories` | Categories |
| GET/PUT/DELETE | `/v1/categories/{category}` | Single category |
| GET/PUT | `/v1/object-groups/*` | Object groups |
| GET/POST | `/v1/piggy-banks` | Piggy banks |
| GET/PUT/DELETE | `/v1/piggy-banks/{piggyBank}` | Single piggy bank |
| GET/POST | `/v1/recurrences` | Recurring transactions |
| GET/PUT/DELETE | `/v1/recurrences/{recurrence}` | Single recurrence |
| GET/POST | `/v1/rules` | Rules |
| GET/PUT/DELETE | `/v1/rules/{rule}` | Single rule |
| GET/POST | `/v1/rule-groups` | Rule groups |
| GET/PUT/DELETE | `/v1/rule-groups/{ruleGroup}` | Single rule group |
| GET/POST | `/v1/tags` | Tags |
| GET/PUT/DELETE | `/v1/tags/{tag}` | Single tag |
| GET/POST | `/v1/transactions` | Transactions |
| GET/PUT/DELETE | `/v1/transactions/{transaction}` | Single transaction |
| GET/DELETE | `/v1/transaction-journals/{journal}` | Transaction journals |
| GET | `/v1/currencies` | Currencies |
| POST/DELETE | `/v1/currencies` | **Admin only** |
| GET/POST/PUT | `/v1/transaction-links` | Transaction links |
| GET | `/v1/link-types` | Link types |
| POST/PUT/DELETE | `/v1/link-types` | **Admin only** |
| GET | `/v1/exchange-rates` | Exchange rates |
| POST/PUT/DELETE | `/v1/exchange-rates` | Exchange rate CRUD |
| GET | `/v1/search` | Search |
| GET/POST/PUT/DELETE | `/v1/webhooks/*` | Webhooks |
| GET/POST/PUT | `/v1/preferences/*` | User preferences |
| POST | `/v1/batch/*` | Batch operations |
| GET | `/v1/user-groups` | User groups |
| GET/PUT | `/v1/user-groups/{userGroup}` | Single user group |

### Sub-Resource Endpoints
Beberapa resource punya sub-endpoint:
- `/v1/accounts/{account}/piggy-banks` — piggy banks di account
- `/v1/accounts/{account}/transactions` — transactions di account
- `/v1/bills/{bill}/attachments`, `/rules`, `/transactions`
- `/v1/budgets/{budget}/transactions`, `/attachments`, `/limits`
- `/v1/categories/{category}/transactions`, `/attachments`
- `/v1/piggy-banks/{piggyBank}/events`, `/attachments`, `/accounts`
- `/v1/tags/{tag}/transactions`, `/attachments`
- `/v1/transactions/{transaction}/attachments`, `/piggy-bank-events`
- `/v1/webhooks/{webhook}/messages`, `/trigger-transaction`, `/messages/{message}/attempts`

---

## Authorization Per Endpoint

Setiap API controller mendefinisikan `$acceptedRoles`:

| Controller Type | Required Role |
|----------------|---------------|
| Chart controllers | `READ_ONLY` |
| Autocomplete | `READ_ONLY` atau domain-specific read |
| Data export | `READ_ONLY` |
| Data destroy/purge | `FULL` |
| Bulk transactions | `MANAGE_TRANSACTIONS` |
| Transaction store/update | `MANAGE_TRANSACTIONS` |
| Currency exchange rate CRUD | `OWNER` |
| Most CRUD (accounts, budgets, bills, etc.) | `MANAGE_TRANSACTIONS` |
| User group show/update | `VIEW_MEMBERSHIPS` |
| Preferences | `READ_ONLY` |

**Cara kerja**: `ValidatesUserGroupTrait` di setiap controller mengecek:
1. User ter-autentikasi
2. User adalah member dari `user_group_id` yang diminta
3. User punya minimal satu `$acceptedRoles` di group tersebut

---

## API Controllers (141 files)

```
app/Api/V1/Controllers/
├── Autocomplete/     (13 files) - Dropdown data
├── Chart/            (5 files)  - Chart data
├── Data/             (4 files)  - Export, destroy, purge, bulk
├── Insight/
│   ├── Expense/      (6 files)
│   ├── Income/       (4 files)
│   └── Transfer/     (4 files)
├── Models/
│   ├── Account/      (5 files)
│   ├── Attachment/   (4 files)
│   ├── AvailableBudget/ (1 file)
│   ├── Bill/         (5 files)
│   ├── Budget/       (5 files)
│   ├── BudgetLimit/  (5 files)
│   ├── Category/     (5 files)
│   ├── CurrencyExchangeRate/ (5 files)
│   ├── ObjectGroup/  (4 files)
│   ├── PiggyBank/    (5 files)
│   ├── Recurrence/   (6 files)
│   ├── Rule/         (6 files)
│   ├── RuleGroup/    (6 files)
│   ├── Tag/          (4 files)
│   ├── Transaction/  (5 files)
│   ├── TransactionCurrency/ (4 files)
│   ├── TransactionLink/ (3 files)
│   ├── TransactionLinkType/ (3 files)
│   ├── UserGroup/    (3 files)
│   └── Webhook/      (5 files)
├── Search/           (2 files)
├── Summary/          (1 file)
├── System/           (6 files)
└── User/             (1 file)
```

---

## Request Validation (72 files)

Semua extends `ApiRequest` yang extends `FormRequest` dengan trait `ChecksLogin`.

**Pattern**:
```php
class StoreRequest extends ApiRequest
{
    protected array $acceptedRoles = [UserRoleEnum::MANAGE_TRANSACTIONS];

    public function rules(): array { /* ... */ }
}
```

`ChecksLogin` trait melakukan authorization check di `authorize()` sebelum controller di-reach.

---

## API Response Format: JSON:API

Menggunakan **League Fractal** + **JsonApiSerializer** (bukan Laravel API Resources).

**Content-Type**: `application/vnd.api+json` (atau `application/json`)

**Base Controller methods**:
- `jsonApiList()` — paginated collection
- `jsonApiObject()` — single item

### Transformers (26 files)

```
app/Transformers/
├── AccountTransformer, AttachmentTransformer, AvailableBudgetTransformer
├── BillTransformer, BudgetLimitTransformer, BudgetTransformer
├── CategoryTransformer, CurrencyTransformer, ExchangeRateTransformer
├── LinkTypeTransformer, ObjectGroupTransformer
├── PiggyBankEventTransformer, PiggyBankTransformer
├── PreferenceTransformer, RecurrenceTransformer
├── RuleGroupTransformer, RuleTransformer
├── TagTransformer, TransactionGroupTransformer, TransactionLinkTransformer
├── UserGroupTransformer, UserTransformer
├── WebhookAttemptTransformer, WebhookMessageTransformer, WebhookTransformer
```

---

## Rate Limiting

> **Tidak ada API rate limiting** pada endpoint API. Hanya OAuth token issuance yang di-throttle.

---

## API Scopes

> **Tidak ada API scopes** yang didefinisikan. Token diterbitkan tanpa scope restrictions. Ini berarti setiap token OAuth punya akses penuh ke semua data user sesuai role group-nya.

---

## Multi-Tenancy via User Groups

API endpoints menerima optional `user_group_id` parameter:
```php
// app/Support/Http/Api/ValidatesUserGroupTrait.php
if ($request->has('user_group_id')) {
    $groupId = (int) $request->get('user_group_id');
} else {
    $groupId = (int) $user->user_group_id;
}
```

Data di-scope ke group ini. Binder classes memastikan model yang di-resolve dari route parameters juga milik group yang aktif.

---

## Data Scoping Pattern

Setiap controller di constructor-nya memanggil:
```php
$this->repository->setUser(auth()->user());
```

Ini memastikan semua query repository di-scope ke authenticated user + active user group.

# API Request & Response Format Reference

---

## 1. JSON:API Response Format

Menggunakan League Fractal + JsonApiSerializer. Content-Type: `application/vnd.api+json` (atau `application/json`).

### List Response

```json
{
  "data": [
    {
      "type": "accounts",
      "id": "1",
      "attributes": { /* ... */ },
      "links": { "self": "/accounts/1" }
    }
  ],
  "meta": {
    "pagination": {
      "total": 42,
      "count": 50,
      "per_page": 50,
      "current_page": 1,
      "total_pages": 1,
      "links": { "next": null, "previous": null }
    }
  },
  "links": {
    "self": "https://example.com/api/v1/accounts?page=1"
  }
}
```

### Single Object Response

```json
{
  "data": {
    "type": "accounts",
    "id": "1",
    "attributes": { /* ... */ },
    "links": { "self": "/accounts/1" }
  },
  "meta": {},
  "links": {
    "self": "https://example.com/api/v1/accounts/1"
  }
}
```

### Account Attributes (dari AccountTransformer)

```json
{
  "type": "accounts",
  "id": "1",
  "attributes": {
    "created_at": "2024-01-15T10:30:00+00:00",
    "updated_at": "2024-01-15T10:30:00+00:00",
    "active": true,
    "order": 1,
    "name": "Checking Account",
    "type": "asset",
    "account_role": "defaultAsset",
    "object_group_id": null,
    "object_group_order": null,
    "object_group_title": null,
    "object_has_currency_setting": true,
    "currency_id": "1",
    "currency_name": "US Dollar",
    "currency_code": "USD",
    "currency_symbol": "$",
    "currency_decimal_places": 2,
    "primary_currency_id": "1",
    "primary_currency_name": "US Dollar",
    "primary_currency_code": "USD",
    "primary_currency_symbol": "$",
    "primary_currency_decimal_places": 2,
    "current_balance": "1234.56",
    "pc_current_balance": "1234.56",
    "opening_balance": "0",
    "pc_opening_balance": "0",
    "virtual_balance": "0",
    "pc_virtual_balance": "0",
    "debt_amount": null,
    "pc_debt_amount": null,
    "balance_difference": "1234.56",
    "pc_balance_difference": "1234.56",
    "current_balance_date": "2024-06-15T12:00:00+00:00",
    "notes": null,
    "monthly_payment_date": null,
    "credit_card_type": null,
    "account_number": null,
    "iban": "NL91ABNA0417164300",
    "bic": null,
    "opening_balance_date": null,
    "liability_type": null,
    "liability_direction": null,
    "interest": null,
    "interest_period": null,
    "include_net_worth": true,
    "longitude": null,
    "latitude": null,
    "zoom_level": null,
    "last_activity": "2024-06-01T08:00:00+00:00"
  },
  "links": { "self": "/accounts/1" }
}
```

> `pc_` prefix = primary currency conversion

### Transaction Group Response (nested)

```json
{
  "data": {
    "type": "transactions",
    "id": "42",
    "attributes": {
      "created_at": "2024-01-15T10:30:00+00:00",
      "updated_at": "2024-01-15T10:30:00+00:00",
      "user": "1",
      "user_group": "1",
      "group_title": "Groceries at Walmart",
      "transactions": [
        {
          "user": "1",
          "transaction_journal_id": "100",
          "type": "withdrawal",
          "date": "2024-01-15T10:30:00+00:00",
          "order": 0,
          "object_has_currency_setting": true,
          "currency_id": "1",
          "currency_code": "USD",
          "currency_name": "US Dollar",
          "currency_symbol": "$",
          "currency_decimal_places": 2,
          "foreign_currency_id": null,
          "foreign_currency_code": null,
          "foreign_currency_name": null,
          "foreign_currency_symbol": null,
          "foreign_currency_decimal_places": null,
          "primary_currency_id": "1",
          "primary_currency_code": "USD",
          "primary_currency_name": "US Dollar",
          "primary_currency_symbol": "$",
          "primary_currency_decimal_places": 2,
          "amount": "42.50",
          "pc_amount": "42.50",
          "foreign_amount": null,
          "pc_foreign_amount": null,
          "source_balance_after": "1192.06",
          "pc_source_balance_after": null,
          "destination_balance_after": "0",
          "pc_destination_balance_after": null,
          "source_balance_dirty": false,
          "destination_balance_dirty": false,
          "description": "Groceries",
          "source_id": "1",
          "source_name": "Checking Account",
          "source_iban": "NL91ABNA0417164300",
          "source_type": "Asset account",
          "destination_id": "5",
          "destination_name": "Groceries",
          "destination_iban": null,
          "destination_type": "Expense account",
          "budget_id": null,
          "budget_name": null,
          "category_id": "3",
          "category_name": "Groceries",
          "bill_id": null,
          "bill_name": null,
          "reconciled": false,
          "notes": null,
          "tags": ["shopping", "groceries"],
          "internal_reference": null,
          "external_id": null,
          "original_source": "ff3-v6.1.0",
          "recurrence_id": null,
          "recurrence_total": null,
          "recurrence_count": null,
          "external_url": null,
          "import_hash_v2": null,
          "sepa_cc": null,
          "sepa_ct_op": null,
          "sepa_ct_id": null,
          "sepa_db": null,
          "sepa_country": null,
          "sepa_ep": null,
          "sepa_ci": null,
          "sepa_batch_id": null,
          "interest_date": null,
          "book_date": null,
          "process_date": null,
          "due_date": null,
          "payment_date": null,
          "invoice_date": null,
          "longitude": null,
          "latitude": null,
          "zoom_level": null,
          "has_attachments": false
        }
      ],
      "links": [
        { "rel": "self", "uri": "/transactions/42" }
      ]
    },
    "links": {
      "self": "https://example.com/api/v1/transactions/42"
    }
  }
}
```

---

## 2. Request Format

JSON body (bukan form data). Accepts: `application/json` atau `application/vnd.api+json`.

### Create Account

```
POST /api/v1/accounts
```

```json
{
  "name": "Checking Account",
  "type": "asset",
  "iban": "NL91ABNA0417164300",
  "bic": "ABNANL2A",
  "account_number": "123456789",
  "opening_balance": "1000",
  "opening_balance_date": "2024-01-01",
  "virtual_balance": "500",
  "currency_id": 1,
  "currency_code": "USD",
  "active": true,
  "include_net_worth": true,
  "account_role": "defaultAsset",
  "credit_card_type": "monthlyFull",
  "monthly_payment_date": "2024-01-28",
  "liability_type": "loan",
  "liability_direction": "credit",
  "liability_amount": "5000",
  "liability_start_date": "2024-01-01",
  "interest": "3.5",
  "interest_period": "monthly",
  "notes": "My main checking account"
}
```

**Validation**:
- `name`: required, max:1024, unique per user
- `type`: required, valid account type
- `iban`: valid IBAN, unique, nullable
- `bic`: valid BIC, nullable
- `opening_balance`: numeric, required_with opening_balance_date
- `virtual_balance`: numeric, nullable
- `account_role`: required if type=asset
- `liability_type`: required if type=liability (loan/debt/mortgage)
- `liability_direction`: required if type=liability (credit/debit)
- `interest`: numeric, 0-100

### Create Transaction

```
POST /api/v1/transactions
```

```json
{
  "group_title": "Groceries at Walmart",
  "error_if_duplicate_hash": false,
  "apply_rules": true,
  "fire_webhooks": true,
  "transactions": [
    {
      "type": "withdrawal",
      "date": "2024-01-15",
      "description": "Groceries",
      "amount": "42.50",
      "currency_id": 1,
      "currency_code": "USD",
      "foreign_amount": null,
      "foreign_currency_id": null,
      "source_id": 1,
      "source_name": "Checking Account",
      "source_iban": "NL91ABNA0417164300",
      "destination_id": 5,
      "destination_name": "Groceries",
      "budget_id": 1,
      "budget_name": "Monthly Budget",
      "category_id": 3,
      "category_name": "Groceries",
      "bill_id": null,
      "piggy_bank_id": null,
      "tags": ["shopping", "groceries"],
      "notes": "Weekly groceries",
      "reconciled": false,
      "order": 0
    }
  ]
}
```

**Validation per transaction**:
- `type`: required, one of: withdrawal, deposit, transfer, opening-balance, reconciliation
- `date`: required, valid date/datetime
- `amount`: required, positive, cannot be 0
- `description`: max:1000
- `source_id`/`destination_id`: must belong to user
- `budget_id`/`category_id`/`bill_id`/`piggy_bank_id`: must exist, must belong to user
- `tags`: array of strings, max 255 chars each
- `notes`: max:32768

**Additional validation**:
- Must submit at least 1 transaction
- All journals must have description
- All transaction types must be equal
- If >1 journal, `group_title` is required (split transaction)

---

## 3. Pagination

### Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number (min:1, max:131337) |
| `limit` | int | 50 (user preference) | Per page (min:1, max:131337) |
| `sort` | string | null | Sort: `field` (asc), `-field` (desc) |

### Date Range

| Parameter | Type | Description |
|-----------|------|-------------|
| `start` | date | Start date (required with end) |
| `end` | date | End date (must be after start) |

### Account Type Filter

| Value | Maps To |
|-------|---------|
| `all` | All types |
| `normal` | asset, expense, revenue, loan, debt, mortgage |
| `asset` | asset + default |
| `cash` | cash only |
| `expense` | expense + beneficiary |
| `revenue` | revenue only |
| `liability` | debt, loan, mortgage, credit-card |
| `hidden` | initial-balance, import, reconciliation |
| Individual: `asset`, `cash`, `credit-card`, `loan`, `mortgage`, `debt`, etc. |

---

## 4. Error Responses

### Validation Error (422)
```json
{
  "message": "The given data was invalid.",
  "errors": {
    "name": ["The name field is required."],
    "type": ["The selected type is invalid."]
  }
}
```

### Not Found (404)
```json
{
  "message": "Resource not found",
  "exception": "NotFoundHttpException"
}
```

### Unauthorized (401)
```json
{
  "message": "Unauthenticated.",
  "exception": "AuthenticationException"
}
```

### Internal Error (500) — Production
```json
{
  "message": "Internal Firefly III Exception: Error message here",
  "exception": "UndisclosedException"
}
```

### HTTP Status Codes

| Code | When |
|------|------|
| 200 | Successful GET, PUT, DELETE |
| 201 | Successful POST (create) |
| 400 | Bad request (invalid headers, etc.) |
| 401 | Unauthenticated / Unauthorized |
| 404 | Resource not found |
| 405 | Method not allowed |
| 422 | Validation failure |
| 500 | Internal error |

---

## 5. Auth Flow (OAuth2)

### Token Request
```
POST /oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=password&client_id=1&client_secret=xxx
&username=user@example.com&password=xxx
```

### Token Response
```json
{
  "token_type": "Bearer",
  "expires_in": 31536000,
  "access_token": "eyJ0eXAiOiJKV1...",
  "refresh_token": "def502..."
}
```

### Using Token
```
GET /api/v1/accounts
Authorization: Bearer eyJ0eXAiOiJKV1...
Accept: application/vnd.api+json
```

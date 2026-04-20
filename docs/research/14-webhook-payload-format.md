# Webhook Payload Format

---

> Lengkap spesifikasi webhook: trigger types, JSON payload, signing, retry logic, lifecycle.

## 1. Trigger Types

| Enum | Value | Description |
|------|-------|-------------|
| `ANY` | 50 | Fires on ANY event |
| `STORE_TRANSACTION` | 100 | New transaction group created |
| `UPDATE_TRANSACTION` | 110 | Transaction group updated |
| `DESTROY_TRANSACTION` | 120 | Transaction group destroyed |
| `STORE_BUDGET` | 200 | Budget created |
| `UPDATE_BUDGET` | 210 | Budget updated |
| `DESTROY_BUDGET` | 220 | Budget destroyed |
| `STORE_UPDATE_BUDGET_LIMIT` | 230 | Budget limit created/updated/destroyed |

> `ANY` = catch-all, match dengan trigger apapun

## 2. Response Types (Content Selection)

| Enum | Value | Description |
|------|-------|-------------|
| `TRANSACTIONS` | 200 | Full transaction group data |
| `ACCOUNTS` | 210 | All accounts from transaction |
| `BUDGET` | 230 | Budget or budget limit data |
| `RELEVANT` | 240 | Auto-pick based on object type |
| `NONE` | 220 | No content |

### RELEVANT Resolution

| Model Class | Resolves To |
|-------------|------------|
| `TransactionGroup` | `TRANSACTIONS` |
| `Budget` / `BudgetLimit` | `BUDGET` |

## 3. Delivery Type

| Enum | Value | Notes |
|------|-------|-------|
| `JSON` | 300 | Hanya JSON yang supported |

## 4. JSON Payload Structure

### Envelope (semua webhook)

```json
{
  "uuid": "<string: UUIDv4>",
  "user_id": "<integer>",
  "user_group_id": "<integer>",
  "trigger": "<string: STORE_TRANSACTION>",
  "response": "<string: TRANSACTIONS>",
  "url": "<string: webhook URL>",
  "version": "<string: v0>",
  "content": { ... }
}
```

> `version` selalu `"v0"`

### Content: TRANSACTIONS Response

Single object (bukan array):

```json
{
  "id": "<integer>",
  "created_at": "<ISO 8601>",
  "updated_at": "<ISO 8601>",
  "user": "<integer: user_id>",
  "group_title": "<string|null>",
  "transactions": [
    {
      "user": "<integer>",
      "transaction_journal_id": "<string>",
      "type": "<string: withdrawal|deposit|transfer>",
      "date": "<ISO 8601>",
      "order": "<integer>",
      "currency_id": "<string>",
      "currency_code": "<string>",
      "currency_symbol": "<string>",
      "currency_decimal_places": "<integer>",
      "foreign_currency_id": "<string|null>",
      "foreign_currency_code": "<string|null>",
      "foreign_currency_symbol": "<string|null>",
      "foreign_currency_decimal_places": "<integer|null>",
      "amount": "<string: positive decimal>",
      "foreign_amount": "<string|null>",
      "description": "<string>",
      "source_id": "<string>",
      "source_name": "<string>",
      "source_iban": "<string|null>",
      "source_type": "<string>",
      "destination_id": "<string>",
      "destination_name": "<string>",
      "destination_iban": "<string|null>",
      "destination_type": "<string>",
      "budget_id": "<string|null>",
      "budget_name": "<string|null>",
      "category_id": "<string|null>",
      "category_name": "<string|null>",
      "bill_id": "<string|null>",
      "bill_name": "<string|null>",
      "reconciled": "<boolean>",
      "notes": "<string|null>",
      "tags": ["<string>"],
      "internal_reference": "<string|null>",
      "external_id": "<string|null>",
      "original_source": "<string|null>",
      "recurrence_id": "<string|null>",
      "bunq_payment_id": "<string|null>",
      "import_hash_v2": "<string|null>",
      "sepa_cc": "<string|null>",
      "sepa_ct_op": "<string|null>",
      "sepa_ct_id": "<string|null>",
      "sepa_db": "<string|null>",
      "sepa_country": "<string|null>",
      "sepa_ep": "<string|null>",
      "sepa_ci": "<string|null>",
      "sepa_batch_id": "<string|null>",
      "interest_date": "<ISO 8601|null>",
      "book_date": "<ISO 8601|null>",
      "process_date": "<ISO 8601|null>",
      "due_date": "<ISO 8601|null>",
      "payment_date": "<ISO 8601|null>",
      "invoice_date": "<ISO 8601|null>",
      "longitude": "<float|null>",
      "latitude": "<float|null>",
      "zoom_level": "<integer|null>"
    }
  ],
  "links": [{ "rel": "self", "uri": "/transactions/<id>" }]
}
```

### Content: ACCOUNTS Response

Array of account objects — field lengkapnya sama dengan `AccountTransformer` (lihat `10-api-response-fields.md`).

### Content: BUDGET Response

Single object — field lengkapnya sama dengan `BudgetTransformer`.

### Content: NONE Response

```json
{}
```

## 5. HTTP Request Details

```
POST {webhook_url}
Content-Type: application/json
Accept: application/json
Signature: t=<timestamp>,v1=<hmac_hex>
User-Agent: FireflyIII/<version>

{json_body}
```

| Header | Value |
|--------|-------|
| `Content-Type` | `application/json` |
| `Accept` | `application/json` |
| `Signature` | `t={timestamp},v1={hmac_hex}` |
| `User-Agent` | `FireflyIII/6.x.x` |
| Timeout | 10 detik (connect: 3.14s) |

## 6. Webhook Signing

### Algorithm

**HMAC-SHA3-256**

### Signing Process

```
1. timestamp = Carbon::now()->getTimestamp()
2. payload = "{timestamp}.{json_body}"
3. signature = hash_hmac('sha3-256', payload, webhook_secret)
4. header = "t={timestamp},v1={signature}"
```

### Verification (di receiving end)

```
1. Parse Signature header → extract timestamp dan signature
2. Reconstruct: "{timestamp}.{raw_request_body}"
3. Compute: hash_hmac('sha3-256', reconstructed, your_webhook_secret)
4. Timing-safe comparison
```

### Secret

Random 24-character string, generated via `Str::random(24)` saat webhook creation.

### URL Validation

Sebelum kirim, URL di-validasi:
- Hostname resolved ke IP
- IPv4 loopback (`127.0.0.0/8`) di-**izinkan**
- Reserved IP ranges (non-public) di-**tolak**

## 7. Webhook Message Lifecycle

### State Machine

```
Created (pending) → Queued → Sending → Success (sent=true)
                                    → Failure (sent=false, errored=true)
                                    → Retry (max 3 attempts)
                                    → Cleanup (auto-delete after 14 days)
```

### Database Tables

| Table | Key Fields |
|-------|------------|
| `webhooks` | id, user_id, title, secret (32 chars), active, trigger, response, delivery, url (1024 chars) |
| `webhook_messages` | id, webhook_id, sent (bool), errored (bool), uuid (UUIDv4), message (JSON longText) |
| `webhook_attempts` | id, webhook_message_id, status_code (smallint), logs (longText), response (longText) |

### Retry Logic

| Aspect | Value |
|--------|-------|
| Max attempts | **3 total** (2 retries + 1 original) |
| Backoff | **Tidak ada exponential backoff** — retry setiap cron run |
| Cron interval | **10 menit** minimum |
| Messages per cron run | **5 max** |
| Messages per webhook | **3 max** |

### Cleanup

Messages dengan `sent=true` dan `created_at < 14 hari ago` auto-deleted.

## 8. Trigger-to-Event Mapping

| Trigger | Source Event | Listener |
|---------|-------------|----------|
| `STORE_TRANSACTION` | `CreatedSingleTransactionGroup` | `ProcessesNewTransactionGroup` |
| `UPDATE_TRANSACTION` | `UpdatedSingleTransactionGroup` | `ProcessesUpdatedTransactionGroup` |
| `DESTROY_TRANSACTION` | `DestroyedSingleTransactionGroup` | `ProcessesDestroyedTransactionGroup` |
| `STORE_BUDGET` | `CreatedBudget` | `ProcessesBudgets` |
| `UPDATE_BUDGET` | `UpdatedBudget` | `ProcessesBudgets` |
| `DESTROY_BUDGET` | `DestroyingBudget` | `ProcessesBudgets` |
| `STORE_UPDATE_BUDGET_LIMIT` | `CreatedBudgetLimit` / `UpdatedBudgetLimit` / `DestroyedBudgetLimit` | `ProcessesBudgetLimits` |

> `fireWebhooks` flag pada event mengontrol apakah webhooks di-fire. Jika `false`, `createWebhookMessages()` di-skip.

## 9. Feature Flags

Dua config harus `true` untuk webhooks aktif:
1. `config('firefly.feature_flags.webhooks')`
2. `FireflyConfig::get('allow_webhooks', config('firefly.allow_webhooks'))->data`

## 10. Go Implementation Notes

### Signing (Go)

```go
import (
    "crypto/hmac"
    "crypto/sha3"
    "encoding/hex"
    "strconv"
    "time"
)

func SignWebhook(body []byte, secret string) string {
    timestamp := strconv.FormatInt(time.Now().Unix(), 10)
    payload := timestamp + "." + string(body)
    mac := hmac.New(sha3.New256, []byte(secret))
    mac.Write([]byte(payload))
    sig := hex.EncodeToString(mac.Sum(nil))
    return fmt.Sprintf("t=%s,v1=%s", timestamp, sig)
}
```

### Verification (Go)

```go
func VerifyWebhook(body []byte, signature string, secret string) bool {
    // Parse "t=...,v1=..."
    parts := strings.SplitN(signature, ",", 2)
    tPart := strings.TrimPrefix(parts[0], "t=")
    v1Part := strings.TrimPrefix(parts[1], "v1=")

    payload := tPart + "." + string(body)
    mac := hmac.New(sha3.New256, []byte(secret))
    mac.Write([]byte(payload))
    expected := hex.EncodeToString(mac.Sum(nil))

    return hmac.Equal([]byte(v1Part), []byte(expected))
}
```

### Retry Strategy Recommendation untuk Go

Implementasi **exponential backoff** (improvement dari Firefly III):
- Attempt 1: immediately
- Attempt 2: 30 detik delay
- Attempt 3: 5 menit delay
- Max 5 attempts (bukan 3)
- Dead letter queue setelah max attempts

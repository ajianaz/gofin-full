# API Response Fields — Transformer Definitions

---

> Referensi lengkap field yang di-return oleh setiap transformer. Digunakan untuk membuat Go response structs.
> Account & Transaction sudah didokumentasi di `08-api-format.md`. Dokumen ini fokus ke transformer lainnya.

## Cross-Cutting Patterns

### Date Format

Semua datetime field menggunakan ISO 8601 / RFC 3339 (AtomString):
```
"2024-01-15T12:00:00+00:00"
```

### HATEOAS Links

Hampir semua transformer mengembalikan array `links`. Go type:
```go
type Link struct {
    Rel string `json:"rel"`
    URI string `json:"uri"`
}
```

### Currency Block Pattern

Banyak transformer mengulang blok 5-field ini:
```go
type CurrencyBlock struct {
    ID             string `json:"currency_id"`
    Name           string `json:"currency_name"`
    Code           string `json:"currency_code"`
    Symbol         string `json:"currency_symbol"`
    DecimalPlaces  int    `json:"currency_decimal_places"`
}
```

### Primary Currency Block Pattern

Identik di Budget, Bill, Category, PiggyBank, AvailableBudget, BudgetLimit, UserGroup:
```go
type PrimaryCurrencyBlock struct {
    ID             string `json:"primary_currency_id"`
    Name           string `json:"primary_currency_name"`
    Code           string `json:"primary_currency_code"`
    Symbol         string `json:"primary_currency_symbol"`
    DecimalPlaces  int    `json:"primary_currency_decimal_places"`
}
```

### Object Group Ref Pattern

Di Budget, Bill, PiggyBank:
```go
type ObjectGroupRef struct {
    ID    *int    `json:"object_group_id"`
    Order *int    `json:"object_group_order"`
    Title *string `json:"object_group_title"`
}
```

### `pc_` Prefix Convention

Field dengan prefix `pc_` = nilai yang sama dikonversi ke primary currency user. Bisa `null` jika `convertToPrimary` setting disabled.

---

## 1. BudgetTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `string` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `active` | `bool` | |
| `name` | `string` | |
| `order` | `int` | |
| `notes` | `*string` | Dari meta, nullable |
| `auto_budget_type` | `*string` | `"reset"`, `"rollover"`, `"adjusted"`, atau `null` |
| `auto_budget_period` | `*string` | Period string, atau `null` |
| `object_group_id` | `*int` | Dari meta |
| `object_group_order` | `*int` | Dari meta |
| `object_group_title` | `*string` | Dari meta |
| `object_has_currency_setting` | `bool` | `true` jika budget punya currency setting |
| `currency_id` | `*string` | `null` jika tidak ada currency setting |
| `currency_code` | `*string` | |
| `currency_name` | `*string` | |
| `currency_symbol` | `*string` | |
| `currency_decimal_places` | `*int` | |
| `primary_currency_id` | `string` | Selalu ada |
| `primary_currency_name` | `string` | |
| `primary_currency_code` | `string` | |
| `primary_currency_symbol` | `string` | |
| `primary_currency_decimal_places` | `int` | |
| `auto_budget_amount` | `*string` | `null` jika tidak ada auto budget |
| `pc_auto_budget_amount` | `*string` | `null` jika convertToPrimary=false |
| `spent` | `*[]CurrencySum` | `null` jika tidak ada data |
| `pc_spent` | `*[]CurrencySum` | `null` jika convertToPrimary=false |
| `links` | `[]Link` | |

---

## 2. BillTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `int` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `name` | `string` | |
| `object_has_currency_setting` | `bool` | Selalu `true` |
| `currency_id` | `string` | |
| `currency_name` | `string` | |
| `currency_code` | `string` | |
| `currency_symbol` | `string` | |
| `currency_decimal_places` | `int` | |
| `primary_currency_id` | `string` | |
| `primary_currency_name` | `string` | |
| `primary_currency_code` | `string` | |
| `primary_currency_symbol` | `string` | |
| `primary_currency_decimal_places` | `int` | |
| `amount_min` | `string` | Decimal string |
| `pc_amount_min` | `string` | Primary-converted |
| `amount_max` | `string` | Decimal string |
| `pc_amount_max` | `string` | Primary-converted |
| `amount_avg` | `string` | Decimal string |
| `pc_amount_avg` | `string` | Primary-converted |
| `date` | `string` | AtomString, next expected date |
| `end_date` | `*string` | Nullable AtomString |
| `extension_date` | `*string` | Nullable AtomString |
| `repeat_freq` | `string` | `"monthly"`, `"weekly"`, dll. |
| `skip` | `int` | Jumlah skip (0 = none) |
| `active` | `bool` | |
| `order` | `int` | |
| `notes` | `*string` | Dari meta |
| `object_group_id` | `*int` | Dari meta |
| `object_group_order` | `*int` | Dari meta |
| `object_group_title` | `*string` | Dari meta |
| `paid_dates` | `[]string` | Array of AtomString |
| `pay_dates` | `[]string` | Array of AtomString |
| `next_expected_match` | `*string` | Nullable AtomString |
| `next_expected_match_diff` | `*int` | Hari sampai next match |
| `links` | `[]Link` | |

---

## 3. CategoryTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `int` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `name` | `string` | |
| `notes` | `*string` | Dari meta |
| `object_has_currency_setting` | `bool` | Selalu `false` |
| `primary_currency_id` | `string` | |
| `primary_currency_name` | `string` | |
| `primary_currency_code` | `string` | |
| `primary_currency_symbol` | `string` | |
| `primary_currency_decimal_places` | `int` | |
| `spent` | `*[]CurrencySum` | |
| `pc_spent` | `*[]CurrencySum` | |
| `earned` | `*[]CurrencySum` | |
| `pc_earned` | `*[]CurrencySum` | |
| `transferred` | `*[]CurrencySum` | |
| `pc_transferred` | `*[]CurrencySum` | |
| `links` | `[]Link` | |

### CurrencySum Sub-Structure

Digunakan di Category, AvailableBudget:

```go
type CurrencySum struct {
    ID             int    `json:"currency_id"`
    Code           string `json:"currency_code"`
    Name           string `json:"currency_name"`
    Symbol         string `json:"currency_symbol"`
    DecimalPlaces  int    `json:"currency_decimal_places"`
    Sum            string `json:"sum"` // bc-rounded decimal string
}
```

---

## 4. TagTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `int` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `tag` | `string` | Label tag |
| `date` | `*string` | Format `"Y-m-d"`, nullable |
| `description` | `*string` | `null` jika empty string |
| `longitude` | `*float64` | Dari Location relation |
| `latitude` | `*float64` | Dari Location relation |
| `zoom_level` | `*int` | Dari Location relation |
| `links` | `[]Link` | |

---

## 5. PiggyBankTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `string` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `name` | `string` | |
| `percentage` | `*int` | 0-100, `null` jika target=null atau current=0 |
| `start_date` | `*string` | Nullable AtomString |
| `target_date` | `*string` | Nullable AtomString |
| `order` | `int` | |
| `active` | `bool` | Selalu `true` (hardcoded) |
| `notes` | `*string` | Dari meta |
| `object_group_id` | `*int` | Dari meta |
| `object_group_order` | `*int` | Dari meta |
| `object_group_title` | `*string` | Dari meta |
| `accounts` | `[]AccountRef` | Dari meta — linked accounts |
| `object_has_currency_setting` | `bool` | Selalu `true` |
| `currency_id` | `string` | |
| `currency_name` | `string` | |
| `currency_code` | `string` | |
| `currency_symbol` | `string` | |
| `currency_decimal_places` | `int` | |
| `primary_currency_id` | `string` | |
| `primary_currency_name` | `string` | |
| `primary_currency_code` | `string` | |
| `primary_currency_symbol` | `string` | |
| `primary_currency_decimal_places` | `int` | |
| `target_amount` | `string` | Decimal |
| `pc_target_amount` | `string` | Primary-converted |
| `current_amount` | `string` | Decimal |
| `pc_current_amount` | `string` | Primary-converted |
| `left_to_save` | `string` | Decimal |
| `pc_left_to_save` | `string` | Primary-converted |
| `save_per_month` | `string` | Decimal |
| `pc_save_per_month` | `string` | Primary-converted |
| `links` | `[]Link` | |

---

## 6. RuleTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `string` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `rule_group_id` | `string` | |
| `rule_group_title` | `string` | |
| `title` | `string` | |
| `description` | `string` | |
| `order` | `int` | |
| `active` | `bool` | |
| `strict` | `bool` | ALL triggers must match |
| `stop_processing` | `bool` | Stop evaluating further rules |
| `trigger` | `string` | User action trigger value (`user_action` type) |
| `triggers` | `[]RuleTrigger` | Semua triggers KECUALI `user_action` type |
| `actions` | `[]RuleAction` | |
| `links` | `[]Link` | |

### RuleTrigger Sub-Structure

```go
type RuleTrigger struct {
    ID              string `json:"id"`
    CreatedAt       string `json:"created_at"`
    UpdatedAt       string `json:"updated_at"`
    Type            string `json:"type"`       // Leading `-` stripped
    Value           string `json:"value"`      // "true" jika operator needs_context=false
    Prohibited      bool   `json:"prohibited"` // true jika original type dimulai `-`
    Order           int    `json:"order"`
    Active          bool   `json:"active"`
    StopProcessing  bool   `json:"stop_processing"`
}
```

### RuleAction Sub-Structure

```go
type RuleAction struct {
    ID              string `json:"id"`
    CreatedAt       string `json:"created_at"`
    UpdatedAt       string `json:"updated_at"`
    Type            string `json:"type"`
    Value           string `json:"value"`
    Order           int    `json:"order"`
    Active          bool   `json:"active"`
    StopProcessing  bool   `json:"stop_processing"`
}
```

---

## 7. RuleGroupTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `int` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `title` | `string` | |
| `description` | `string` | |
| `order` | `int` | |
| `active` | `bool` | |
| `links` | `[]Link` | |

---

## 8. RecurrenceTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `string` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `type` | `string` | Transaction type (withdrawal/deposit/transfer) |
| `title` | `string` | |
| `description` | `string` | |
| `first_date` | `string` | AtomString, selalu ada |
| `latest_date` | `*string` | Nullable AtomString |
| `repeat_until` | `*string` | Nullable AtomString |
| `apply_rules` | `bool` | |
| `active` | `bool` | |
| `nr_of_repetitions` | `*int` | `null` jika repetitions=0 |
| `notes` | `*string` | Dari meta |
| `repetitions` | `[]RecurrenceRepetition` | Dari meta |
| `transactions` | `[]RecurrenceTransaction` | Dari meta |
| `links` | `[]Link` | |

---

## 9. UserGroupTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `int` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `in_use` | `bool` | `true` jika current user belongs to group |
| `can_see_members` | `bool` | User punya VIEW_MEMBERSHIPS role atau owner |
| `title` | `string` | |
| `primary_currency_id` | `string` | |
| `primary_currency_name` | `string` | |
| `primary_currency_code` | `string` | |
| `primary_currency_symbol` | `string` | |
| `primary_currency_decimal_places` | `int` | |
| `members` | `[]Member` | Kosong jika user tidak bisa lihat members |

### Member Sub-Structure

```go
type Member struct {
    UserID    string   `json:"user_id"`
    UserEmail string   `json:"user_email"`
    You       bool     `json:"you"`
    Roles     []string `json:"roles"` // e.g. ["owner", "admin"]
}
```

> Members merged by email — user dengan multiple roles dapat single entry dengan semua roles digabung.

---

## 10. UserTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `int` | (NOT cast ke string) |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `email` | `string` | |
| `blocked` | `bool` | `true` jika blocked==1 |
| `blocked_code` | `*string` | `null` jika empty string |
| `role` | `string` | Dari repository |
| `links` | `[]Link` | |

---

## 11. AvailableBudgetTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `string` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `object_has_currency_setting` | `bool` | Selalu `true` |
| `currency_id` | `string` | |
| `currency_name` | `string` | |
| `currency_code` | `string` | |
| `currency_symbol` | `string` | |
| `currency_decimal_places` | `int` | |
| `primary_currency_id` | `string` | |
| `primary_currency_name` | `string` | |
| `primary_currency_code` | `string` | |
| `primary_currency_symbol` | `string` | |
| `primary_currency_decimal_places` | `int` | |
| `amount` | `string` | bc-rounded |
| `pc_amount` | `*string` | `null` jika convertToPrimary=false |
| `start` | `string` | AtomString |
| `end` | `string` | AtomString (end of day) |
| `spent_in_budgets` | `map[string]string` | Budget ID → amount |
| `pc_spent_in_budgets` | `map[string]string` | |
| `spent_outside_budgets` | `map[string]string` | |
| `pc_spent_outside_budgets` | `map[string]string` | |
| `links` | `[]Link` | |

---

## 12. ObjectGroupTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `string` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `title` | `string` | |
| `order` | `int` | |
| `links` | `[]Link` | |

---

## 13. PreferenceTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `int` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `user_group_id` | `*int` | `null` jika user_group_id=0 |
| `name` | `string` | |
| `data` | `interface{}` | Arbitrary JSON (string, int, bool, array) |

> **Tidak punya `links` field.**

---

## 14. AttachmentTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `string` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `attachable_id` | `string` | |
| `attachable_type` | `string` | `FireflyIII\Models\` prefix stripped |
| `hash` | `string` | MD5 file hash |
| `filename` | `string` | |
| `download_url` | `string` | Full URL |
| `upload_url` | `string` | Full URL |
| `title` | `string` | |
| `notes` | `*string` | |
| `mime` | `string` | MIME type |
| `size` | `int` | Bytes |
| `links` | `[]Link` | |

---

## 15. WebhookTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `int` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `active` | `bool` | |
| `title` | `string` | |
| `secret` | `string` | Signing secret |
| `triggers` | `[]WebhookTrigger` | Dari meta |
| `deliveries` | `[]WebhookDelivery` | Dari meta |
| `responses` | `[]WebhookResponse` | Dari meta |
| `url` | `string` | Target URL |
| `links` | `[]Link` | |

---

## 16. CurrencyTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `int` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `native` | `bool` | Group's native currency |
| `default` | `bool` | Alias dari `native` |
| `primary` | `bool` | Alias dari `native` |
| `enabled` | `bool` | |
| `name` | `string` | Full name (e.g. "US Dollar") |
| `code` | `string` | ISO 4217 (e.g. "USD") |
| `symbol` | `string` | (e.g. "$") |
| `decimal_places` | `int` | |
| `links` | `[]Link` | |

---

## 17. LinkTypeTransformer

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `int` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `name` | `string` | |
| `inward` | `string` | (e.g. "Relates to") |
| `outward` | `string` | (e.g. "Is related to") |
| `editable` | `bool` | User bisa buat link type ini |
| `links` | `[]Link` | |

---

## 18. BudgetLimitTransformer (has available include)

| Field | Go Type | Notes |
|-------|---------|-------|
| `id` | `string` | |
| `created_at` | `string` | AtomString |
| `updated_at` | `string` | AtomString |
| `start` | `string` | AtomString |
| `end` | `string` | AtomString (end of day) |
| `budget_id` | `string` | |
| `object_has_currency_setting` | `bool` | Selalu `true` |
| `currency_id` | `string` | |
| `currency_name` | `string` | |
| `currency_code` | `string` | |
| `currency_symbol` | `string` | |
| `currency_decimal_places` | `int` | |
| `primary_currency_id` | `string` | |
| `primary_currency_name` | `string` | |
| `primary_currency_code` | `string` | |
| `primary_currency_symbol` | `string` | |
| `primary_currency_decimal_places` | `int` | |
| `amount` | `string` | bc-rounded |
| `pc_amount` | `*string` | `null` jika convertToPrimary=false |
| `period` | `string` | |
| `spent` | `map[string]string` | Dari meta |
| `pc_spent` | `map[string]string` | Dari meta |
| `notes` | `*string` | Dari meta |
| `links` | `[]Link` | |

> **Available include:** `budget` — returns full BudgetTransformer result.

---

## Ringkasan ID Type Convention

Penting untuk Go struct design — ID type TIDAK konsisten:

| Transformer | ID Type | Catatan |
|-------------|---------|---------|
| Account | `string` | |
| TransactionGroup | `string` | |
| Budget | `string` | |
| PiggyBank | `string` | |
| Rule | `string` | |
| RuleGroup | `int` | |
| Bill | `int` | |
| Category | `int` | |
| Tag | `int` | |
| User | `int` | |
| UserGroup | `int` | |
| AvailableBudget | `string` | |
| ObjectGroup | `string` | |
| Attachment | `string` | |
| Webhook | `int` | |
| Currency | `int` | |
| LinkType | `int` | |
| BudgetLimit | `string` | |

> **Go recommendation:** Gunakan `int64` untuk semua ID, dan cast di serialization layer. Atau tetap `string` untuk JSON:API compliance.

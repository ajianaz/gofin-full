# API Request Validation Rules

---

> Lengkap validation rules per endpoint untuk membuat Go request validation structs.
> Custom rules yang perlu diimplementasi di Go juga tercantum di bagian bawah.

## 1. Budget Store

**Class:** `BudgetFormStoreRequest`

| Field | Rules |
|-------|-------|
| `name` | `required`, `min:1`, `max:255`, `uniqueObjectForUser:budgets,name` |
| `active` | `numeric`, `min:0`, `max:1` |
| `auto_budget_type` | `numeric`, `integer`, `gte:0`, `lte:3` |
| `auto_budget_currency_id` | `exists:transaction_currencies,id` |
| `auto_budget_amount` | `required_if:auto_budget_type,1,2`, `IsValidPositiveAmount` |
| `auto_budget_period` | `in:daily,weekly,monthly,quarterly,half_year,yearly` |
| `notes` | `min:1`, `max:32768`, `nullable` |

**Custom after-validation:** `validateAutoBudgetAmount` — cek jika auto_budget_type ≠ 0, amount dan currency_id harus ada.

## 2. Budget Update

**Class:** `BudgetFormUpdateRequest`

| Field | Rules |
|-------|-------|
| `name` | `required`, `min:1`, `max:255`, `uniqueObjectForUser:budgets,name,{budget_id}` |
| `active` | `numeric`, `min:0`, `max:1` |
| `auto_budget_type` | `numeric`, `integer`, `gte:0`, `lte:31` |
| `auto_budget_currency_id` | `exists:transaction_currencies,id` |
| `auto_budget_amount` | `required_if:auto_budget_type,1,2`, `numeric`, `IsValidPositiveAmount` |
| `auto_budget_period` | `in:daily,weekly,monthly,quarterly,half_year,yearly` |
| `notes` | `min:1`, `max:32768`, `nullable` |

> `auto_budget_type` max `lte:31` pada update vs `lte:3` pada store.

## 3. Bill Store

**Class:** `BillStoreRequest`

| Field | Rules |
|-------|-------|
| `name` | `required`, `min:1`, `max:255`, `uniqueObjectForUser:bills,name` |
| `amount_min` | `required`, `IsValidPositiveAmount` |
| `amount_max` | `required`, `IsValidPositiveAmount` |
| `transaction_currency_id` | `required`, `exists:transaction_currencies,id` |
| `date` | `required`, `date` |
| `bill_end_date` | `nullable`, `date` |
| `extension_date` | `nullable`, `date` |
| `repeat_freq` | `required`, `in:daily,weekly,monthly,quarterly,half-year,yearly` |
| `skip` | `required`, `integer`, `gte:0`, `lte:31` |
| `active` | `boolean` |
| `notes` | `min:1`, `max:32768`, `nullable` |

## 4. Bill Update

**Class:** `BillUpdateRequest` — sama dengan Store, kecuali `name` excludes current `{bill_id}`.

## 5. Piggy Bank Store

**Class:** `PiggyBankStoreRequest`

| Field | Rules |
|-------|-------|
| `name` | `required`, `min:1`, `max:255`, `uniquePiggyBankForUser` |
| `accounts` | `required`, `array` |
| `accounts.*` | `required`, `belongsToUser:accounts` |
| `target_amount` | `nullable`, `IsValidPositiveAmount` |
| `start_date` | `date` |
| `target_date` | `date`, `nullable` |
| `order` | `integer`, `min:1` |
| `object_group` | `min:0`, `max:255` |
| `notes` | `min:1`, `max:32768`, `nullable` |

**Custom after-validation:** Cek linked accounts punya matching currency, dan account types dalam: `asset`, `loan`, `debt`, `mortgage`.

## 6. Piggy Bank Update

**Class:** `PiggyBankUpdateRequest` — sama dengan Store, kecuali:
- `name` excludes `{piggyBank_id}`
- Tambah `transaction_currency_id`: `exists:transaction_currencies,id`
- `order` punya `max:32768`

## 7. Category Store / Update

**Class:** `CategoryFormRequest`

| Field | Rules |
|-------|-------|
| `name` | `required`, `min:1`, `max:255`, `uniqueObjectForUser:categories,name` |
| `notes` | `min:1`, `max:32768`, `nullable` |

> Update: `name` excludes `{category_id}`.

## 8. Tag Store / Update

**Class:** `TagFormRequest`

| Field | Rules |
|-------|-------|
| `tag` | `required`, `max:1024`, `min:1`, `uniqueObjectForUser:tags,tag` |
| `id` | (empty / `belongsToUser:tags` on update) |
| `description` | `max:32768`, `min:1`, `nullable` |
| `date` | `date`, `nullable`, `after:1970-01-02`, `before:2038-01-17` |
| `latitude` | `numeric`, `min:-90`, `max:90`, `nullable`, `required_with:longitude` |
| `longitude` | `numeric`, `min:-180`, `max:180`, `nullable`, `required_with:latitude` |
| `zoom_level` | `numeric`, `min:0`, `max:80`, `nullable`, `required_with:latitude` |

> Update: `tag` excludes `{tag_id}`.

## 9. Recurrence Store / Update

**Class:** `RecurrenceFormRequest`

### Base Rules

| Field | Rules |
|-------|-------|
| `title` | `required`, `min:1`, `max:255`, `uniqueObjectForUser:recurrences,title` |
| `first_date` | `required`, `date` (store: `before:today+25y`, `after:today`) |
| `repetition_type` | `required`, `ValidRecurrenceRepetitionValue`, `ValidRecurrenceRepetitionType`, `min:1`, `max:32` |
| `skip` | `required`, `numeric`, `integer`, `gte:0`, `lte:31` |
| `notes` | `min:1`, `max:32768`, `nullable` |
| `recurring_description` | `min:0`, `max:32768` |
| `active` | `numeric`, `min:0`, `max:1` |
| `apply_rules` | `numeric`, `min:0`, `max:1` |
| `transaction_description` | `required`, `min:1`, `max:255` |
| `transaction_type` | `required`, `in:withdrawal,deposit,transfer` |
| `transaction_currency_id` | `required`, `exists:transaction_currencies,id` |
| `amount` | `required`, `IsValidPositiveAmount` |
| `source_id` | `numeric`, `belongsToUser:accounts,id`, `nullable` |
| `source_name` | `min:1`, `max:255`, `nullable` |
| `destination_id` | `numeric`, `belongsToUser:accounts,id`, `nullable` |
| `destination_name` | `min:1`, `max:255`, `nullable` |
| `foreign_amount` | `nullable`, `IsValidPositiveAmount` |
| `budget_id` | `mustExist:budgets,id`, `belongsToUser:budgets,id`, `nullable` |
| `bill_id` | `mustExist:bills,id`, `belongsToUser:bills,id`, `nullable` |
| `category` | `min:1`, `max:255`, `nullable` |
| `tags` | `min:1`, `max:255`, `nullable` |

### Conditional Rules (by transaction_type)

| Type | Source | Destination |
|------|--------|-------------|
| `withdrawal` | `source_id`: required, exists, belongsToUser | `destination_name`: nullable |
| `deposit` | `source_name`: nullable | `destination_id`: required, exists, belongsToUser |
| `transfer` | `source_id`: required, exists, belongsToUser, different:destination_id | `destination_id`: required, exists, belongsToUser, different:source_id |

### Conditional Rules (by repetition_end)

| repetition_end | Additional Rules |
|----------------|-----------------|
| `times` | `repetitions`: required, numeric, min:0, max:255 |
| `until_date` | `repeat_until`: required, date, after:tomorrow |

### Conditional Rules (foreign currency)

| Condition | Rules |
|-----------|-------|
| `foreign_currency_id > 0` | `foreign_currency_id`: exists:transaction_currencies,id |
| `foreign_amount` not null | `foreign_currency_id`: exists, different:transaction_currency_id |

## 10. Rule Store / Update

**Class:** `RuleFormRequest`

| Field | Rules |
|-------|-------|
| `title` | `required`, `min:1`, `max:255`, `uniqueObjectForUser:rules,title` |
| `description` | `min:1`, `max:32768`, `nullable` |
| `stop_processing` | `boolean` |
| `rule_group_id` | `required`, `belongsToUser:rule_groups` |
| `trigger` | `required`, `in:store-journal,update-journal,manual-activation` |
| `triggers.*.type` | `required`, `in:{validTriggers}` (dynamic) |
| `triggers.*.value` | `required_if:triggers.*.type,{contextTriggers}`, `max:1024`, `min:1` |
| `actions.*.type` | `required`, `in:{validActions}` (from config) |
| `actions.*.value` | `required_if:actions.*.type,{contextActions}`, `min:0`, `max:1024` |
| `strict` | `in:0,1` |
| `run_after_form` | `in:0,1` |

## 11. RuleGroup Store / Update

**Class:** `RuleGroupFormRequest`

| Field | Rules |
|-------|-------|
| `title` | `required`, `min:1`, `max:255`, `uniqueObjectForUser:rule_groups,title` |
| `description` | `min:1`, `max:32768`, `nullable` |
| `active` | `IsBoolean` |

## 12. ObjectGroup Store / Update

**Class:** `ObjectGroupFormRequest`

| Field | Rules |
|-------|-------|
| `title` | `required`, `min:1`, `max:255`, `uniqueObjectGroup` |

## 13. Webhook Create

**Class:** `Webhook\CreateRequest`

| Field | Rules |
|-------|-------|
| `title` | `required`, `min:1`, `max:255`, `uniqueObjectForUser:webhooks,title` |
| `active` | `IsBoolean` |
| `triggers` | `required`, `array`, `min:1`, `max:10` |
| `triggers.*` | `required`, `in:50,100,110,120,200,210,220,230` |
| `responses` | `required`, `array`, `min:1`, `max:1` |
| `responses.*` | `required`, `in:200,210,230,240,220` |
| `deliveries` | `required`, `array`, `min:1`, `max:1` |
| `deliveries.*` | `required`, `in:300` |
| `url` | `required`, `url`, `IsValidWebhookUrl` |

### Webhook Enum Values

| Enum | Values | Meaning |
|------|--------|---------|
| **Trigger** | 50=ANY, 100=STORE_TX, 110=UPDATE_TX, 120=DESTROY_TX, 200=STORE_BUDGET, 210=UPDATE_BUDGET, 220=DESTROY_BUDGET, 230=STORE_UPDATE_BUDGET_LIMIT | |
| **Response** | 200=TRANSACTIONS, 210=ACCOUNTS, 220=NONE, 230=BUDGET, 240=RELEVANT | |
| **Delivery** | 300=JSON | Hanya JSON |

## 14. Webhook Update

**Class:** `Webhook\UpdateRequest` — sama dengan Create, kecuali:
- `title` tidak required, excludes `{webhook_id}`
- `url` tidak required, tambah `uniqueExistingWebhook:{webhook_id}`

## 15. UserGroup Update

**Class:** `UserGroup\UpdateRequest`

| Field | Rules |
|-------|-------|
| `title` | `required`, `min:1`, `max:255` |
| `primary_currency_id` | `exists:transaction_currencies,id` |
| `primary_currency_code` | `exists:transaction_currencies,code` |

## 16. Preference Store

**Class:** `PreferenceStoreRequest`

| Field | Rules |
|-------|-------|
| `name` | `required` |
| `data` | `required` (any type — string, bool, numeric) |

> `data` coercion: `"true"` → `true`, `"false"` → `false`, numeric → `float`.

## 17. Preference Update

**Class:** `PreferenceUpdateRequest`

| Field | Rules |
|-------|-------|
| `data` | `required` (any type) |

## 18. CurrencyExchangeRate Store

**Class:** `CurrencyExchangeRate\StoreRequest`

| Field | Rules |
|-------|-------|
| `date` | `required`, `date`, `after:1970-01-02`, `before:2038-01-17` |
| `rate` | `required`, `numeric`, `gt:0` |
| `from` | `required`, `exists:transaction_currencies,code` |
| `to` | `required`, `exists:transaction_currencies,code` |

## 19. CurrencyExchangeRate Update

**Class:** `CurrencyExchangeRate\UpdateRequest`

| Field | Rules |
|-------|-------|
| `date` | `date`, `after:1970-01-02`, `before:2038-01-17` |
| `rate` | `required`, `numeric`, `gt:0` |
| `from` | `nullable`, `exists:transaction_currencies,code` |
| `to` | `nullable`, `exists:transaction_currencies,code` |

## 20. Account Store / Update

**Class:** `AccountFormRequest`

| Field | Rules |
|-------|-------|
| `name` | `required`, `max:1024`, `min:1`, `uniqueAccountForUser` |
| `opening_balance` | `nullable`, `IsValidAmount` |
| `opening_balance_date` | `date`, `required_with:opening_balance`, `nullable` |
| `iban` | `iban`, `nullable`, `UniqueIban` |
| `BIC` | `bic`, `nullable` |
| `virtual_balance` | `nullable`, `IsValidAmount` |
| `currency_id` | `exists:transaction_currencies,id` |
| `account_number` | `min:1`, `max:255`, `uniqueAccountNumberForUser`, `nullable` |
| `account_role` | `in:defaultAsset,sharedAsset,savingAsset,ccAsset,cashWalletAsset` |
| `active` | `boolean` |
| `cc_type` | `in:monthlyFull` |
| `interest_period` | `in:daily,monthly,yearly` |
| `notes` | `min:1`, `max:32768`, `nullable` |
| `latitude` | `numeric`, `min:-90`, `max:90`, `nullable`, `required_with:longitude` |
| `longitude` | `numeric`, `min:-180`, `max:180`, `nullable`, `required_with:latitude` |
| `zoom_level` | `numeric`, `min:0`, `max:80`, `nullable`, `required_with:latitude` |

## 21. Transaction Store (API v1)

**Class:** `Transaction\StoreRequest`

### Group-level

| Field | Rules |
|-------|-------|
| `group_title` | `min:1`, `max:1000`, `nullable` |
| `error_if_duplicate_hash` | `IsBoolean` |
| `fire_webhooks` | `IsBoolean` |
| `apply_rules` | `IsBoolean` |

### Per-transaction (`transactions[]`)

| Field | Rules |
|-------|-------|
| `type` | `required`, `in:withdrawal,deposit,transfer,opening-balance,reconciliation` |
| `date` | `required`, `IsDateOrTime` |
| `order` | `numeric`, `min:0` |
| `currency_id` | `numeric`, `exists:transaction_currencies,id`, `nullable` |
| `currency_code` | `min:3`, `max:51`, `exists:transaction_currencies,code`, `nullable` |
| `foreign_currency_id` | `numeric`, `exists:transaction_currencies,id`, `nullable` |
| `foreign_currency_code` | `min:3`, `max:51`, `exists:transaction_currencies,code`, `nullable` |
| `amount` | `required`, `IsValidPositiveAmount` |
| `foreign_amount` | `nullable`, `IsValidZeroOrMoreAmount` |
| `description` | `nullable`, `min:1`, `max:1000` |
| `source_id` | `numeric`, `nullable`, `BelongsUser` |
| `source_name` | `min:1`, `max:255`, `nullable` |
| `source_iban` | `min:1`, `max:255`, `nullable`, `iban` |
| `source_number` | `min:1`, `max:255`, `nullable` |
| `source_bic` | `min:1`, `max:255`, `nullable`, `bic` |
| `destination_id` | `numeric`, `nullable`, `BelongsUser` |
| `destination_name` | `min:1`, `max:255`, `nullable` |
| `destination_iban` | `min:1`, `max:255`, `nullable`, `iban` |
| `destination_number` | `min:1`, `max:255`, `nullable` |
| `destination_bic` | `min:1`, `max:255`, `nullable`, `bic` |
| `budget_id` | `mustExist:budgets,id`, `BelongsUser` |
| `budget_name` | `min:1`, `max:255`, `nullable`, `BelongsUser` |
| `category_id` | `mustExist:categories,id`, `BelongsUser`, `nullable` |
| `category_name` | `min:1`, `max:255`, `nullable` |
| `bill_id` | `numeric`, `nullable`, `mustExist:bills,id`, `BelongsUser` |
| `bill_name` | `min:1`, `max:255`, `nullable`, `BelongsUser` |
| `piggy_bank_id` | `numeric`, `nullable`, `mustExist:piggy_banks,id`, `BelongsUser` |
| `piggy_bank_name` | `min:1`, `max:255`, `nullable`, `BelongsUser` |
| `reconciled` | `IsBoolean` |
| `notes` | `min:1`, `max:32768`, `nullable` |
| `tags` | `min:0`, `max:255` |
| `tags.*` | `min:0`, `max:255` |
| `internal_reference` | `min:1`, `max:255`, `nullable` |
| `external_id` | `min:1`, `max:255`, `nullable` |
| `external_url` | `min:1`, `max:255`, `nullable`, `url` |
| `sepa_cc` | `min:1`, `max:255`, `nullable` |
| `sepa_ct_op` | `min:1`, `max:255`, `nullable` |
| `sepa_ct_id` | `min:1`, `max:255`, `nullable` |
| `sepa_db` | `min:1`, `max:255`, `nullable` |
| `sepa_country` | `min:1`, `max:255`, `nullable` |
| `sepa_ep` | `min:1`, `max:255`, `nullable` |
| `sepa_ci` | `min:1`, `max:255`, `nullable` |
| `sepa_batch_id` | `min:1`, `max:255`, `nullable` |
| `interest_date` | `date`, `nullable` |
| `book_date` | `date`, `nullable` |
| `process_date` | `date`, `nullable` |
| `due_date` | `date`, `nullable` |
| `payment_date` | `date`, `nullable` |
| `invoice_date` | `date`, `nullable` |
| `latitude` | `numeric`, `min:-90`, `max:90`, `nullable`, `required_with:longitude` |
| `longitude` | `numeric`, `min:-180`, `max:180`, `nullable`, `required_with:latitude` |
| `zoom_level` | `numeric`, `min:0`, `max:80`, `nullable`, `required_with:latitude` |

### Custom After-Validation (sequential)

1. **validateTransactionArray** — must be valid array
2. **validateOneTransaction** — minimal 1 transaction
3. **validateDescriptions** — semua journal harus punya description
4. **validateTransactionTypes** — semua type dalam group harus sama
5. **validateForeignCurrencyInformation** — validasi foreign currency
6. **validateAccountInformation** — validasi source/destination accounts
7. **validateEqualAccounts** — validasi kesamaan source/destination by type
8. **validateGroupDescription** — group_title wajib jika >1 journal

## 22. Transaction Update (API v1)

**Class:** `Transaction\UpdateRequest`

### Group-level

| Field | Rules |
|-------|-------|
| `group_title` | `min:1`, `max:1000`, `nullable` |
| `apply_rules` | `IsBoolean` |

### Per-transaction — sama dengan Store, kecuali:
- `type`, `date`, `amount` **tidak required**
- Tambah `transaction_journal_id`: `nullable`, `numeric`, `BelongsUser`
- **Tidak ada** IBAN/BIC/number/piggy_bank fields
- **Tidak ada** location fields

### Custom After-Validation (sequential)

1. **validateJournalIds** — verify journal IDs jika >1 transaction
2. **validateGroupDescription** — group_title wajib jika >1 journal
3. **validateTransactionTypesForUpdate** — semua type harus sama
4. **preventUpdateReconciled** — reconciled transaction tidak bisa ubah source/destination/amount
5. **validateEqualAccountsForUpdate** — validasi source/destination by type
6. **validateAccountInformationUpdate** — validasi accounts dan currency

---

## Custom Rules Reference

Rules khusus yang perlu diimplementasi di Go:

| Custom Rule | Used By | Go Implementation |
|-------------|---------|-----------------|
| `IsValidPositiveAmount` | Budget, Bill, PiggyBank, Recurrence, Transaction | Regex `^\d+(\.\d+)?$` atau parse ke `decimal.Decimal` |
| `IsValidZeroOrMoreAmount` | Transaction (foreign_amount) | Parse ke `decimal.Decimal`, cek ≥ 0 |
| `IsValidAmount` | Account (opening_balance, virtual_balance) | Parse ke `decimal.Decimal`, bisa negatif |
| `IsBoolean` | RuleGroup, Webhook, Transaction | Accept `true/false/1/0` |
| `IsDateOrTime` | Transaction (date) | Parse ISO 8601 date or datetime |
| `BelongsUser` / `belongsToUser` | Transaction, Rule | Query DB: entity.user_group_id = current user group |
| `UniqueIban` | Account | Cek IBAN uniqueness dalam user group |
| `uniqueObjectForUser` | Budget, Bill, Category, Tag, Recurrence, Rule, RuleGroup, Webhook | Cek uniqueness dalam user group scope |
| `uniqueAccountForUser` | Account | Cek account name uniqueness dalam user group |
| `uniquePiggyBankForUser` | PiggyBank | Cek piggy bank name uniqueness dalam user group |
| `uniqueObjectGroup` | ObjectGroup | Cek object group title uniqueness |
| `IsValidWebhookUrl` | Webhook | Validasi URL format + tidak boleh localhost/private IP |
| `IsValidActionExpression` | Rule | Validasi syntax expression rule action |
| `ruleTriggerValue` | Rule | Validasi trigger value terhadap operator type |
| `ruleActionValue` | Rule | Validasi action value terhadap action type |
| `ValidRecurrenceRepetitionType` | Recurrence | Validasi repetition type string format |
| `ValidRecurrenceRepetitionValue` | Recurrence | Validasi repetition value |
| `IsValidSlackOrDiscordUrl` | Preferences | Validasi Slack/Discord webhook URL |

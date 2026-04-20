# Search & Rule Engine — Operator Catalog

---

> Lengkap catalog semua search operators dengan SQL WHERE clause equivalents.
> Digunakan untuk implementasi search system dan rule engine di Go.

## Architecture Overview

```
Query string → QueryParser → FieldNode (operator:value) + StringNode (bare words)
  → OperatorQuerySearch → GroupCollector → SQL WHERE clauses
```

### Key Tables

| Table | Alias | Digunakan Untuk |
|-------|-------|-----------------|
| `transaction_journals` | — | description, date, bill_id, created_at, updated_at |
| `transaction_groups` | — | group title |
| `transactions` | `source`, `destination` | account_id, amount, foreign_amount, reconciled |
| `journal_meta` | — | external_id, external_url, internal_reference, recurrence_id, sepa_ct_id, date fields |
| `accounts` | — | name, IBAN, account_number |
| `categories` | — | via `category_transaction_journal` pivot |
| `budgets` | — | via `budget_transaction_journal` pivot |
| `tags` | — | via `tag_transaction_journal` pivot |
| `notes` | — | polymorphic on `transaction_journals` |
| `attachments` | — | polymorphic on `transaction_journals` |

### Negation

Semua operator support negation via `-` prefix: `-description_contains:coffee`

### Value Quoting

Values dengan spaces: `description_contains:"coffee shop"`

---

## 1. Description Operators

| Operator | Aliases | SQL WHERE |
|----------|---------|-----------|
| `description_is` | `description` | `description = 'value' OR title = 'value'` |
| `description_starts` | — | `description LIKE 'value%' OR title LIKE 'value%'` |
| `description_ends` | — | `description LIKE '%value' OR title LIKE '%value'` |
| `description_contains` | — | `description LIKE '%value%' OR title LIKE '%value%'` |

> Searches both `transaction_journals.description` and `transaction_groups.title`

## 2. Account Name Operators

### Both Source & Destination

| Operator | Aliases | SQL WHERE |
|----------|---------|-----------|
| `account_is` | — | `source.account_id IN (...) OR destination.account_id IN (...)` |
| `account_contains` | — | Same pattern |
| `account_starts` | — | Same pattern |
| `account_ends` | — | Same pattern |
| `account_id` | — | Same pattern (comma-separated IDs) |
| `account_is_cash` | — | `source.account_id = {cash_id} OR destination.account_id = {cash_id}` |

### Source Only

| Operator | Aliases | SQL WHERE |
|----------|---------|-----------|
| `source_account_is` | `from_account_is` | `source.account_id IN (...)` |
| `source_account_contains` | `from_account_contains`, `from`, `source` | `source.account_id IN (...)` |
| `source_account_starts` | `from_account_starts` | `source.account_id IN (...)` |
| `source_account_ends` | `from_account_ends` | `source.account_id IN (...)` |
| `source_account_id` | — | `source.account_id = {id}` |
| `source_is_cash` | — | `source.account_id = {cash_id}` |

### Destination Only

| Operator | Aliases | SQL WHERE |
|----------|---------|-----------|
| `destination_account_is` | `to_account_is` | `destination.account_id IN (...)` |
| `destination_account_contains` | `to_account_contains`, `to`, `destination` | `destination.account_id IN (...)` |
| `destination_account_starts` | `to_account_starts` | `destination.account_id IN (...)` |
| `destination_account_ends` | `to_account_ends` | `destination.account_id IN (...)` |
| `destination_account_id` | — | `destination.account_id = {id}` |
| `destination_is_cash` | — | `destination.account_id = {cash_id}` |

> Account operators pertama resolve account by name/IBAN/number via repository, lalu filter by ID

## 3. Account Number / IBAN Operators

| Operator | SQL WHERE |
|----------|-----------|
| `account_nr_is` | `source.account_id IN (...) OR destination.account_id IN (...)` |
| `account_nr_contains` | Same |
| `account_nr_starts` | Same |
| `account_nr_ends` | Same |
| `source_account_nr_is` | `source.account_id IN (...)` |
| `source_account_nr_contains` | Same |
| `source_account_nr_starts` | Same |
| `source_account_nr_ends` | Same |
| `destination_account_nr_is` | `destination.account_id IN (...)` |
| `destination_account_nr_contains` | Same |
| `destination_account_nr_starts` | Same |
| `destination_account_nr_ends` | Same |

## 4. Amount Operators

| Operator | Aliases | SQL WHERE | Notes |
|----------|---------|-----------|-------|
| `amount_is` | `amount`, `amount_exactly` | `source.amount = -value` | Negated karena source negatif |
| `amount_less` | `amount_max`, `less` | `destination.amount <= value` | |
| `amount_more` | `amount_min`, `more` | `destination.amount >= value` | |
| `foreign_amount_is` | `foreign_amount` | `source.foreign_amount = -value` (WHERE NOT NULL) | |
| `foreign_amount_less` | `foreign_amount_max` | `destination.foreign_amount <= value` (WHERE NOT NULL) | |
| `foreign_amount_more` | `foreign_amount_min` | `destination.foreign_amount >= value` (WHERE NOT NULL) | |

> Amount values: `str_replace(',', '.', $value)` lalu `Steam::positive()`

## 5. Date Operators (transaction_journals.date)

### Supported Formats

- Exact: `2024-01-15`, `today`, `yesterday`
- Range: `2024-01-01..2024-01-31`
- Components: `2024` (year), `2024-01` (month), `01` (day)

| Operator | Aliases | SQL WHERE |
|----------|---------|-----------|
| `date_on` | `date`, `date_is`, `on` | `date >= 'value 00:00:00' AND date <= 'value 23:59:59'` |
| `date_before` | `before` | `date <= 'value 23:59:59'` |
| `date_after` | `after` | `date >= 'value 00:00:00'` |

### Component Sub-Operators

| Component | Is | Is Not | Before | After |
|-----------|-----|--------|--------|-------|
| Year | `YEAR(date) = v` | `YEAR(date) != v` | `YEAR(date) <= v` | `YEAR(date) >= v` |
| Month | `MONTH(date) = v` | `MONTH(date) != v` | `MONTH(date) <= v` | `MONTH(date) >= v` |
| Day | `DAY(date) = v` | `DAY(date) != v` | `DAY(date) <= v` | `DAY(date) >= v` |

## 6. Meta Date Operators (journal_meta — post-filter)

Fields: `interest_date`, `book_date`, `process_date`, `due_date`, `payment_date`, `invoice_date`

| Operator Pattern | Filter |
|-----------------|--------|
| `{field}_on` | `metaDate >= start AND metaDate <= end` |
| `{field}_before` | `metaDate <= date` |
| `{field}_after` | `metaDate >= date` |

> SQL hanya filter: `WHERE journal_meta.name = '{field}' AND data IS NOT NULL`
> Selisihnya post-filter di PHP/Go

## 7. Object Date Operators (columns on transaction_journals)

Fields: `created_at`, `updated_at`

| Operator Pattern | SQL WHERE |
|-----------------|-----------|
| `created_at_on` | `created_at >= 'value 00:00:00' AND created_at <= 'value 23:59:59'` |
| `created_at_before` | `created_at <= 'value 00:00:00'` |
| `created_at_after` | `created_at >= 'value 00:00:00'` |
| `updated_at_on` | Same pattern |
| `updated_at_before` | Same pattern |
| `updated_at_after` | Same pattern |

## 8. Category Operators

| Operator | Aliases | SQL WHERE |
|----------|---------|-----------|
| `category_is` | — | `categories.id = {id}` |
| `category_contains` | `category` | `categories.id IN (...)` |
| `category_starts` | — | `categories.id IN (...)` |
| `category_ends` | — | `categories.id IN (...)` |
| `has_any_category` | — | `category_id IS NOT NULL` |
| `has_no_category` | — | `category_id IS NULL AND type != 'Opening balance'` |

## 9. Budget Operators

| Operator | Aliases | SQL WHERE |
|----------|---------|-----------|
| `budget_is` | — | `budgets.id = {id}` |
| `budget_contains` | `budget` | `budgets.id IN (...)` |
| `budget_starts` | — | `budgets.id IN (...)` |
| `budget_ends` | — | `budgets.id IN (...)` |
| `has_any_budget` | — | `budget_id IS NOT NULL` |
| `has_no_budget` | — | `budget_id IS NULL` |

## 10. Bill Operators

| Operator | Aliases | SQL WHERE |
|----------|---------|-----------|
| `bill_is` | `subscription_is` | `bill_id = {id}` |
| `bill_contains` | `bill`, `subscription` | `bill_id IN (...)` |
| `bill_starts` | `subscription_starts` | `bill_id IN (...)` |
| `bill_ends` | `subscription_ends` | `bill_id IN (...)` |
| `has_any_bill` | `has_any_subscription` | `bill_id IS NOT NULL` |
| `has_no_bill` | `has_no_subscription` | `bill_id IS NULL` |

## 11. Tag Operators (SQL + post-filter)

| Operator | Negation | Logic |
|----------|----------|-------|
| `tag_is` | `tag_is_not` | Post-filter: ALL specified tags must be present (AND) |
| `tag_contains` | `-tag_contains` | 1 result: AND; >1: OR |
| `tag_starts` | `-tag_starts` | Same pattern |
| `tag_ends` | `-tag_ends` | Same pattern |
| `has_any_tag` | `has_no_tag` | `tag_id IS NOT NULL` |
| `has_no_tag` | `has_any_tag` | `tag_id IS NULL` |

## 12. Notes Operators

| Operator | Aliases | SQL WHERE |
|----------|---------|-----------|
| `notes_contains` | `notes_contain`, `notes` | `notes.text LIKE '%value%'` |
| `notes_starts` | `notes_start` | `notes.text LIKE 'value%'` |
| `notes_ends` | `notes_end` | `notes.text LIKE '%value'` |
| `notes_is` | `notes_are` | `notes.text = 'value'` |
| `any_notes` | `has_any_notes`, `has_notes` | `notes.text IS NOT NULL` |
| `no_notes` | — | `notes.text IS NULL OR notes.text = ''` |

## 13. Currency Operators

| Operator | SQL WHERE |
|----------|-----------|
| `currency_is` | `source.transaction_currency_id = {id} OR source.foreign_currency_id = {id}` |
| `foreign_currency_is` | `source.foreign_currency_id = {id}` |

## 14. Transaction Type Operator

| Operator | Aliases | SQL WHERE |
|----------|---------|-----------|
| `transaction_type` | `type` | `transaction_types.type IN ('Value')` |

Values: `withdrawal`, `deposit`, `transfer`, `opening-balance`, `reconciliation`, `liability-credit`

## 15. External ID / URL / Internal Reference

| Operator | SQL WHERE |
|----------|-----------|
| `external_id_is` | `name='external_id' AND data = '"value"'` |
| `external_id_contains` | `external_id` → `data LIKE '%value%'` |
| `external_id_starts` | `data LIKE '"value%'` |
| `external_id_ends` | `data LIKE '%value"'` |
| `external_url_is` | `name='external_url' AND data = '"value"'` |
| `external_url_contains` | `external_url` → `data LIKE '%value%'` |
| `internal_reference_is` | `name='internal_reference' AND data = '"value"'` |
| `internal_reference_contains` | `internal_reference` → `data LIKE '%value%'` |
| `any_external_url` | `name='external_url' AND data IS NOT NULL` |
| `no_external_url` | Complex: name != 'external_url' OR data IS NULL |
| `any_external_id` | `name='external_id' AND data IS NOT NULL` |
| `no_external_id` | Complex: same pattern |
| `recurrence_id` | `name='recurrence_id' AND data = '"value"'` |
| `sepa_ct_is` | `name='sepa_ct_id' AND data = '"value"'` |

> Values di `journal_meta.data` disimpan JSON-encoded: `"value"`

## 16. Attachment Operators (post-filter)

| Operator | SQL WHERE | Post-filter |
|----------|-----------|-------------|
| `has_attachments` | `attachable_id IS NOT NULL` | — |
| `attachment_name_is` | `attachment`, `attachment_is` | `filename === value OR title === value` |
| `attachment_name_contains` | — | `str_contains(lower(filename), lower(value))` |
| `attachment_name_starts` | — | `str_starts_with(lower(filename), lower(value))` |
| `attachment_name_ends` | — | `str_ends_with(lower(filename), lower(value))` |
| `attachment_notes_are` | `attachment_notes` | `notes.text === value` |
| `attachment_notes_contains` | `attachment_notes_contain` | `str_contains(lower(notes.text), lower(value))` |

## 17. Reconciliation Operator

| Operator | SQL WHERE |
|----------|-----------|
| `reconciled` | `source.reconciled = 1 AND destination.reconciled = 1` |
| `-reconciled` | `source.reconciled = 0 AND destination.reconciled = 0` |

## 18. ID Operators

| Operator | SQL WHERE |
|----------|-----------|
| `id` | `transaction_groups.id IN (...)` (comma-separated) |
| `journal_id` | `transaction_journals.id IN (...)` |
| `-id` | `transaction_groups.id NOT IN (...)` |
| `-journal_id` | `transaction_journals.id NOT IN (...)` |

## 19. Existence Operator

| Operator | SQL WHERE |
|----------|-----------|
| `exists` | `deleted_at IS NULL AND type NOT IN ('liability-credit', 'opening-balance', 'reconciliation')` |
| `-exists` | `id = -1` (always empty) |

## 20. Account Balance Operators (post-filter, VERY SLOW)

| Operator | Post-filter |
|----------|-------------|
| `source_balance_is` | `bccomp(balance, value) == 0` |
| `source_balance_gt` | `bccomp(balance, value) > 0` |
| `source_balance_gte` | `bccomp(balance, value) >= 0` |
| `source_balance_lt` | `bccomp(balance, value) < 0` |
| `source_balance_lte` | `bccomp(balance, value) <= 0` |
| `destination_balance_*` | Same pattern |

> Menggunakan `Steam::accountsBalancesOptimized()` — **sangat lambat**, hindari di production

## 21. Bare Word Search

Kata tanpa `field:` prefix = description search:

```
coffee shop → WHERE (description LIKE '%coffee%' AND description LIKE '%shop%')
           OR (title LIKE '%coffee%' AND title LIKE '%shop%')
```

Negasi: `-coffee` → `WHERE NOT LIKE '%coffee%'`

## 22. Subquery Grouping

Parentheses membuat subquery:
```
(source_account_is:Checking destination_account_is:Savings)
```

## 23. Boolean Auto-Flip

Operator dengan `needs_context=false`:
- `has_attachments:false` → sama dengan `-has_attachments`
- `has_attachments:true` dengan `-` prefix → flip ke positive

---

## Complete Alias Map

| Alias | Resolves To |
|-------|------------|
| `type` | `transaction_type` |
| `tag` | `tag_is` |
| `description` | `description_is` |
| `notes`, `notes_are`, `notes_contain` | `notes_contains` |
| `from_account_is` | `source_account_is` |
| `from`, `source` | `source_account_contains` |
| `to`, `destination` | `destination_account_contains` |
| `category` | `category_contains` |
| `budget` | `budget_contains` |
| `bill` | `bill_contains` |
| `subscription_is` | `bill_is` |
| `subscription` | `bill_contains` |
| `external_id` | `external_id_contains` |
| `internal_reference` | `internal_reference_contains` |
| `external_url` | `external_url_contains` |
| `attachment`, `attachment_is` | `attachment_name_is` |
| `attachment_notes` | `attachment_notes_are` |
| `date`, `date_is`, `on` | `date_on` |
| `before` | `date_before` |
| `after` | `date_after` |
| `amount`, `amount_exactly` | `amount_is` |
| `less` | `amount_less` |
| `more` | `amount_more` |
| `created_on*` | `created_at*` |
| `updated_on*` | `updated_at*` |

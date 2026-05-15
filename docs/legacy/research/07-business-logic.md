# Business Logic Reference

---

## 1. Balance Calculation (`Steam.php`)

### Core Method: `accountsBalancesOptimized()`

Query `SUM(transactions.amount)` grouped by `account_id` dan `currency_code`, di-join dengan `transaction_journals` dan `transaction_currencies`.

Filter: `transaction_journals.date <= $date` (atau `<` saat `inclusive=false`).

Hasil per account:
- `"balance"` — saldo di currency account sendiri
- `"pc_balance"` — saldo dikonversi ke primary currency user
- Currency code keys (misal `"EUR"`, `"USD"`) untuk setiap currency yang punya aktivitas

### Virtual Balance vs Actual Balance

`virtual_balance` selalu **ditambahkan** di atas calculated actual balance:
```
total_balance = actual_balance + account.virtual_balance
```

### Running Balance (`AccountBalanceCalculator.php`)

Menyimpan `balance_before` dan `balance_after` langsung di setiap row `Transaction`:
- Iterasi semua transactions secara kronologis
- Order: `date ASC, order DESC, id ASC, description ASC, amount ASC`
- `balance_before` = latest known balance
- `balance_after` = balance_before + amount
- Flag `balance_dirty` untuk skip transactions yang sudah dihitung

---

## 2. Transaction Creation Flow

### Chain: Store → Factory → Event → Listener

```
Controller
  → TransactionGroupRepository::store()
    → TransactionGroupFactory::create()
      → TransactionGroup (save)
      → TransactionJournalFactory::create() (per journal)
        → Resolve transaction type
        → Resolve currencies
        → AccountValidator (validate source/destination)
        → Account auto-creation (counterpart accounts)
        → TransactionJournal (save, completed=false)
        → TransactionFactory::create() × 2 (double-entry)
          → Transaction source (negative amount)
          → Transaction destination (positive amount)
    → Fire CreatedSingleTransactionGroup event
      → ProcessesNewTransactionGroup listener:
        1. processRules() — fire rule engine
        2. recalculateCredit() — credit/liability accounts
        3. createWebhookMessages() — send webhooks
        4. removePeriodStatistics() — clear cache
        5. recalculateRunningBalance() — update balance fields
        6. Mark journals completed
```

### Account Auto-Creation

Untuk **withdrawal**: source dibuat dulu, destination auto-created sebagai Expense account jika belum ada.
Untuk **deposit**: destination dibuat dulu, source auto-created sebagai Revenue account.
Untuk **transfer**: source dibuat dulu, destination auto-created sebagai Beneficiary account.

Jika destination masih NULL setelah creation, diganti dengan **cash account**.

### Double-Entry Enforcement

Setiap journal **harus punya tepat 2 transactions**:
- Source account: `amount = -X` (negatif)
- Destination account: `amount = +X` (positif)
- Total per journal = 0

### Foreign Currency Handling

Untuk transfer dengan foreign currency:
- Destination transaction **menukar** currency/foreign:
  - `transaction_currency_id` = foreign currency
  - `foreign_currency_id` = original currency
  - `amount` = foreign amount
  - `foreign_amount` = original amount

---

## 3. Recurring Transactions

### Data Model
- `Recurrence` — title, first_date, repeat_until, repetitions, active
- `RecurrenceRepetition` — type, moment, skip, weekend
- `RecurrenceTransaction` — template data untuk generate transaction

### Repetition Types

| Type | `moment` | Contoh |
|------|----------|--------|
| `daily` | — | Setiap hari, skip=1 → setiap 2 hari |
| `weekly` | "1,7" | Hari Senin dan Minggu |
| `monthly` | "15" | Tanggal 15 setiap bulan |
| `ndom` | "last,7" | Hari Minggu terakhir bulan |
| `yearly` | "2025-06-15" | Setiap tanggal 15 Juni |

### Weekend Handling
| Value | Behavior |
|-------|----------|
| 1 | Skip ke Jumat |
| 2 | Skip ke Senin |
| 3 | Ignore (execute anyway) |

### Execution Flow
1. Fetch semua active recurrences
2. Filter: belum exceed repetitions, repeat_until belum lewat, first_date sudah lewat
3. Cek belum fire hari ini (via latest_date)
4. Hitung occurrences via `getOccurrencesInRange()`
5. Untuk occurrences yang match `today`:
   - Cek duplicate
   - Bangun transaction data dari RecurrenceTransaction templates
   - Call `TransactionGroupRepository::store()`
   - Update `recurrence.latest_date`

---

## 4. Budget System

### Budget Limit Per Period

`available = budget_limit.amount - SUM(transaction.amount)` untuk withdrawal/expenses yang match budget dalam date range.

### Auto-Budget Types

| Type | Behavior |
|------|----------|
| `AUTO_BUDGET_RESET` (1) | Setiap periode dapat fixed amount |
| `AUTO_BUDGET_ROLLOVER` (2) | Unused budget rolls over: `new = leftover + auto_amount` |
| `AUTO_BUDGET_ADJUSTED` (3) | Overspent carries forward: `new = max(prev_limit + spent + auto_amount, 0)` |

### Auto-Budget Periods

| Period | Magic Day |
|--------|-----------|
| daily | Setiap hari |
| weekly | Senin |
| monthly | Tanggal 1 |
| quarterly | 1 Jan, 1 Apr, 1 Jul, 1 Okt |
| half_year | 1 Jan, 1 Jul |
| yearly | 1 Jan |

Job `CreateAutoBudgetLimits` berjalan via cron, cek apakah hari ini magic day untuk setiap auto-budget.

---

## 5. Piggy Banks

### Tracking Money

`current_amount` disimpan di **pivot table** `account_piggy_bank` (bukan di piggy_banks sendiri).

```
total_saved = SUM(account_piggy_bank.current_amount) across all linked accounts
```

### addAmount() / removeAmount()

1. Get current pivot `current_amount`
2. `new_amount = bcadd(current, amount)` atau `bcsub(current, amount)`
3. Convert ke `native_current_amount` jika currency berbeda
4. Save pivot
5. Fire `PiggyBankAmountIsChanged` event

### Piggy Bank Events

`PiggyBankEvent` mencatat setiap addition/removal:
- `amount` selalu positif
- `transaction_journal_id` link ke journal yang trigger
- `native_amount` auto-calc

---

## 6. Rule Engine

### Architecture

```
RuleGroup (stop_processing flag)
  └── Rule (strict/non-strict, stop_processing)
        ├── RuleTrigger (stop_processing) × N
        └── RuleAction (stop_processing) × N
```

### Trigger vs Action

**Trigger**: Kondisi yang dicocokkan terhadap transactions (90+ operators dari search system).
**Action**: Perubahan yang diterapkan jika trigger match (30 action types).

### Strict vs Non-Strict

- **Strict**: ALL triggers must match (AND logic). All values combined into single query.
- **Non-Strict**: EACH trigger searched independently, results UNION-ed (OR logic).

### Stop Processing Chain

```
Group.stop_processing → stop all rules in group
Rule.stop_processing → stop all actions in rule
Trigger.stop_processing → stop all triggers in rule
Action.stop_processing → stop all remaining actions
```

### 30 Action Types

| Action | Purpose |
|--------|---------|
| `set_category` | Set transaction category |
| `clear_category` | Remove category |
| `set_budget` | Set transaction budget |
| `clear_budget` | Remove budget |
| `add_tag` | Add tag |
| `remove_tag` | Remove specific tag |
| `remove_all_tags` | Remove all tags |
| `set_description` | Overwrite description |
| `append_description` | Append to description |
| `prepend_description` | Prepend to description |
| `set_notes` | Overwrite notes |
| `append_notes` | Append to notes |
| `prepend_notes` | Prepend to notes |
| `clear_notes` | Remove all notes |
| `set_source_account` | Change source account |
| `set_destination_account` | Change destination account |
| `switch_accounts` | Swap source and destination |
| `convert_withdrawal` | Change type to withdrawal |
| `convert_deposit` | Change type to deposit |
| `convert_transfer` | Change type to transfer |
| `set_amount` | Change transaction amount |
| `link_to_bill` | Associate with bill/subscription |
| `update_piggy` | Add/remove from piggy bank |
| `delete_transaction` | Delete the transaction |
| `append_descr_to_notes` | Append description to notes |
| `append_notes_to_descr` | Append notes to description |
| `move_descr_to_notes` | Replace notes with description |
| `move_notes_to_descr` | Replace description with notes |
| `set_source_to_cash` | Set source to cash account |
| `set_destination_to_cash` | Set destination to cash account |

### Trigger Operators (90+)

| Category | Operators |
|----------|-----------|
| Account | `account_is`, `account_contains`, `source_account_*`, `destination_account_*`, `account_id` |
| Description | `description_is`, `description_contains`, `description_starts`, `description_ends` |
| Amount | `amount_is`, `amount_less`, `amount_more`, `foreign_amount_*` |
| Category/Budget/Bill | `*_is`, `*_contains`, `*_starts`, `*_ends` |
| Tags | `tag_is`, `tag_is_not`, `tag_contains`, `tag_starts`, `tag_ends` |
| Date | `date_on`, `date_before`, `date_after`, `interest_date_*`, `book_date_*`, etc. |
| Boolean | `reconciled`, `has_attachments`, `has_any_category`, `has_no_budget`, `source_is_cash` |
| Special | `journal_id`, `id`, `recurrence_id`, `external_id_*`, `internal_reference_*` |

---

## 7. Multi-Currency

### Exchange Rate Storage

Per UserGroup, per date. Query: `WHERE date <= $date ORDER BY date DESC` untuk rate terbaru.

### Rate Lookup Order

1. In-memory cache
2. Laravel Cache (forever)
3. Direct DB: `from_currency → to_currency`
4. Reverse DB: `to_currency → from_currency`, lalu `1/rate`
5. Cross-rate via EUR: `from → EUR → to`

### native_amount Auto-Calculation

**Setiap kali Transaction di create/update**, observer menghitung:
```
native_amount = amount × exchange_rate (ke primary currency)
```

Pattern ini berlaku untuk: `Transaction`, `PiggyBankEvent`, `BudgetLimit`, `AutoBudget`, `Bill`, `AvailableBudget`.

### Primary Currency

Disimpan per **UserGroup** (bukan per User). Diakses via `Amount::getPrimaryCurrencyByUserGroup($userGroup)`.

---

## 8. Encryption

### Shadow Column Pattern

Kolom sensitive punya dua versi:
- `column_name` — plaintext (untuk query)
- `column_name_encrypted` — AES-256-CBC encrypted (backup)

### Models dengan Encrypted Columns

| Model | Encrypted Columns |
|-------|-------------------|
| Bill | `name_encrypted`, `amount_min_encrypted`, `amount_max_encrypted`, `match_encrypted` |
| PiggyBankEvent | `amount_encrypted` |

### Encryption Key

Menggunakan `APP_KEY` dari Laravel config via `Crypt` facade (AES-256-CBC). Tidak ada key management terpisah.

---

## 9. Search System

### Query Format

```
operator:value operator:"value with spaces" free_text_word
```

Multiple operators di-AND. Negation dengan `-operator:value`.

### Contoh

```
description_contains:coffee amount_more:100 date_after:2024-01-01
```

### Account Search

`GET /api/v1/search/accounts?query=...&field=all|id|name|iban|number&type=...`

### Transaction Search

`POST /api/v1/search/transactions`
```json
{"query": "description_starts:Groceries date_after:2024-01-01", "page": 1, "limit": 50}
```

### Transaction Count

`POST /api/v1/search/transactions/count`
```json
{"external_identifier": "ext-123", "description": "search string"}
```

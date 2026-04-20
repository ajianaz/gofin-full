# Business Flows — Alur Bisnis Lengkap

---

## 1. New User Onboarding

### Phase 1: Registration

1. User kunjungi halaman registrasi
2. Jika `single_user_mode=true` dan sudah ada user → registrasi diblokir (kecuali punya invite code)
3. User submit: email, password, invite code (opsional)
4. `User` record dibuat (email + bcrypt password)
5. Fire event `NewUserRegistered`

### Phase 2: System Setup (async, background job)

1. Kirim email welcome ke user (jika `notification_user_new_reg` = true)
2. Kirim notifikasi ke admin (jika `notification_admin_new_reg` = true)
3. Jika user pertama → attach global role `owner`
4. Buat `UserGroup` (title = email user)
5. Buat `GroupMembership` (user + group + role OWNER)
6. Set `user.user_group_id` = group baru
7. Seed exchange rates

### Phase 3: Setup Wizard

Jika user punya **0 asset accounts** → redirect ke wizard:

1. Pilih bahasa
2. Pilih currency → di-enable & set sebagai primary currency
3. Masukkan nama bank & saldo
4. (Opsional) Masukkan nama rekening tabungan & saldo

**Auto-created accounts:**

| Account | Role | Saldo |
|---------|------|-------|
| "[Nama Bank]" | `defaultAsset` | Saldo user |
| "[Nama Bank] Savings" | `savingAsset` | Saldo tabungan (opsional) |
| "Cash wallet" | `cashWalletAsset` | 0 |

Semua 3 account ID disimpan di preference `frontpageAccounts` (muncul di dashboard).

### Yang TIDAK Auto-created

- Budget, Category, Tag, Bill, Piggy Bank, Rule — semua harus manual

---

## 2. Transaction Lifecycle

### Creating a Transaction

```
User input (form/API)
  │
  ├── group_title (wajib jika split)
  ├── apply_rules (default: true)
  ├── fire_webhooks (default: true)
  └── transactions[] ← minimal 1
        ├── type: withdrawal|deposit|transfer|opening-balance|reconciliation
        ├── date, description, amount
        ├── currency_id / currency_code
        ├── source_id / source_name / source_iban
        ├── destination_id / destination_name / destination_iban
        ├── budget_id, category_id, bill_id, piggy_bank_id
        └── tags[]

Factory Chain:
  1. TransactionGroupFactory → buat TransactionGroup
  2. TransactionJournalFactory → per transaction:
     a. Resolve transaction type
     b. Resolve currencies
     c. Validate source/destination accounts
     d. Auto-create counterpart account jika belum ada
        - Withdrawal: destination = Expense account
        - Deposit: source = Revenue account
        - Transfer: destination = Beneficiary account
     e. Jika destination masih NULL → pakai Cash account
     f. Buat TransactionJournal (completed=false)
     g. Buat 2 Transaction (double-entry):
        - Source: amount negatif
        - Destination: amount positif
  3. Fire event CreatedSingleTransactionGroup

Post-Processing (event listener, sequential):
  1. processRules() → fire rule engine
  2. recalculateCredit() → credit/liability accounts
  3. createWebhookMessages() → kirim webhook
  4. removePeriodStatistics() → clear cache
  5. recalculateRunningBalance() → update balance fields
  6. Mark semua journals completed
```

### Transaction Types & Counterpart Accounts

| Type | Source Account Type | Destination Account Type | Contoh |
|------|--------------------|-------------------------|--------|
| Withdrawal | Asset/Cash | Expense (auto-created) | Belanja |
| Deposit | Revenue (auto-created) | Asset/Cash | Gaji |
| Transfer | Asset/Cash | Asset/Cash/Beneficiary | Transfer antar bank |
| Opening Balance | Initial Balance | Asset/Cash (auto-created) | Saldo awal |
| Reconciliation | Asset | Reconciliation (auto-created) | Koreksi saldo |

### Split Transaction

- >1 transaction dalam array `transactions[]`
- `group_title` wajib
- Semua journal harus punya description
- Semua journal harus type yang sama
- Share satu TransactionGroup

### Update Transaction

- Update field yang diizinkan (amount, description, date, category, budget, dll.)
- Fire event `UpdatedSingleTransactionGroup`
- Re-run rule engine, recalculate balances, fire webhooks
- Audit log entry dibuat untuk setiap perubahan

### Delete Transaction

- Soft delete (set `deleted_at`)
- Fire event `DestroyedSingleTransactionGroup`
- Clear period statistics cache
- Recalculate running balances

---

## 3. Budget Lifecycle

### A. Membuat Budget

```
User → Budgets → Create
  ├── Name (wajib, unik per user)
  ├── Active toggle
  ├── Auto-budget type:
  │   ├── 0 = Tidak ada auto-budget (manual saja)
  │   ├── 1 = Reset → setiap periode fixed amount
  │   ├── 2 = Rollover → unused budget carry forward
  │   └── 3 = Adjusted → overspent reduce next period
  ├── Auto-budget amount (wajib jika type ≠ 0)
  ├── Auto-budget period: daily/weekly/monthly/quarterly/half_year/yearly
  └── Notes, Attachments
```

### B. Budget Limit

```
User → Budget Show → Add Limit
  ├── Currency
  ├── Start date
  ├── End date
  └── Amount (limit)
```

- Jika limit sudah ada untuk budget/currency/period yang sama → update
- Jika amount = 0 → delete limit

### C. Tracking Spending

```
Available = budget_limit.amount - SUM(withdrawal.amount)
```

- Hanya transaction withdrawal yang linked ke budget yang dihitung
- Filtered per date range (start_date → end_date)
- Tidak ada blocking — overspend hanya informatif

### D. Auto-Budget (via Cron, setiap 12 jam)

```
Cek: hari ini = "magic day" untuk period?

| Period | Magic Day |
|--------|-----------|
| Daily | Setiap hari |
| Weekly | Senin |
| Monthly | Tanggal 1 |
| Quarterly | 1 Jan, 1 Apr, 1 Jul, 1 Okt |
| Half-year | 1 Jan, 1 Jul |
| Yearly | 1 Jan |

Jika belum ada budget limit untuk periode ini → buat:

| Type | Formula |
|------|---------|
| Reset | limit = auto_amount |
| Rollover | leftover = prev_limit + spent; limit = leftover + auto_amount |
| Adjusted | available = prev_limit + auto_amount + spent; limit = max(available, 1) |
```

---

## 4. Bill/Subscription Lifecycle

### A. Membuat Bill

```
User → Bills → Create
  ├── Name (wajib, unik)
  ├── Amount min & max (wajib, range — bill bisa variabel)
  ├── Currency
  ├── Start date (wajib)
  ├── Repeat freq: daily/weekly/monthly/quarterly/half-year/yearly
  ├── Skip (berapa periode dilewati, default 0)
  ├── End date (opsional)
  ├── Extension date (opsional)
  └── Active toggle
```

### B. Auto-Matching (via Rule Engine)

1. Setelah buat bill → user diminta buat rule otomatis
2. Rule triggers: `description_contains:[bill name]`, `amount_is_between:[min]:[max]`
3. Rule action: `link_to_bill` → set `bill_id` pada TransactionJournal
4. Matching hanya untuk withdrawal transactions

### C. Bill Warnings (via Cron, setiap 12 jam)

| Kondisi | Aksi |
|---------|------|
| Bill overdue >6 hari | Kirim notifikasi `SubscriptionsAreOverdueForPayment` |
| End date dalam [90, 30, 14, 7, 0] hari | Kirim notifikasi `SubscriptionNeedsExtensionOrRenewal` |
| Extension date dalam periode yang sama | Kirim notifikasi |

### D. Bill TIDAK Auto-Buat Transaksi

- User harus manual create payment transaction
- Atau pakai Recurring Transaction system

### E. Pay Date Calculation

```
Mulai dari start_date → advance per repeat_freq → akumulasi skip
→ filter hanya tanggal dalam range → hanya setelah last paid date
→ handle end-of-month edge case (e.g., tanggal 30 di bulan 28 hari)
```

---

## 5. Piggy Bank Lifecycle

### A. Membuat Savings Goal

```
User → Piggy Banks → Create
  ├── Name
  ├── Target amount
  ├── Currency (default: primary currency user)
  ├── Start date (default: today)
  ├── Target date (opsional)
  └── Account (linked asset account)
```

### B. Add/Remove Money

```
User → Piggy Bank → Add/Remove
  ├── Pilih amount
  ├── Cek: cukup saldo di source account? (untuk add)
  ├── Cek: cukup saved di piggy bank? (untuk remove)
  ├── Cek: amount ≤ left to save? (cap ke remaining)
  └── Buat PiggyBankEvent + update pivot current_amount
```

### C. Tracking Progress

- `current_amount` di tabel pivot `account_piggy_bank`
- `PiggyBankEvent` mencatat setiap add/remove
- `PiggyBankRepetition` tracks progress per date range
- History chart: line chart of savings over time

### D. Saat Target Tercapai

- Tidak ada automation — purely display-level
- UI hitung: `current_amount >= target_amount`
- User bisa manually mark `active = false`

---

## 6. Reconciliation Flow

### Konsep

Reconciliation = mencocokkan saldo di Firefly III dengan saldo di bank.

### Alur

```
1. User → Account → Reconcile (hanya untuk asset accounts)

2. Pilih date range

3. System hitung:
   - Start balance = saldo account sebelum start date
   - End balance = saldo account di end date
   - Load transactions (start-3 hari → end+3 hari)

4. User pilih transactions yang sudah "clear" di bank

5. Real-time calculation:
   difference = (startBalance - endBalance) + clearedAmount + selectedAmount

6. Submit:
   - Mark selected journals → reconciled = true
   - Jika ada difference → buat Reconciliation type transaction
     (antara asset account ↔ reconciliation account)
```

### Reconciliation Account

- Tipe: `AccountTypeEnum::RECONCILIATION` (hidden system account)
- Auto-created on-demand saat reconcilasi
- Nama: "Reconciliation account [AccountName] ([Currency])"
- Tidak muncul di account list normal

---

## 7. Rule Engine User Flow

### Struktur

```
RuleGroup (container, order, stop_processing)
  └── Rule (title, strict mode, stop_processing)
        ├── RuleTrigger × N (condition: operator:value)
        └── RuleAction × N (action: type:value)
```

### Membuat Rule

```
1. Buat Rule Group (atau pakai default yang auto-created)
2. Dalam group, buat Rule
3. Tambahkan Triggers:
   - description_contains:"coffee"
   - amount_more:"100"
   - category_is:"Food"
4. Tambahkan Actions:
   - set_category:"Coffee"
   - add_tag:"morning"
   - set_budget:"Morning Budget"
```

### Mode: Strict vs Non-Strict

| Mode | Behavior |
|------|----------|
| **Strict** (default) | ALL triggers must match (AND) |
| **Non-Strict** | EACH trigger independently (OR), results UNION-ed |

### Stop Processing Chain

```
Group.stop_processing → berhenti semua rule di group ini
Rule.stop_processing → berhenti semua action di rule ini
Trigger.stop_processing → berhenti semua trigger di rule ini
Action.stop_processing → berhenti action selanjutnya
```

### Test Trigger

- User bisa test triggers → system tampilkan 20 transaction yang match
- Bisa test per rule atau per rule group

### Execute Rule

- **Otomatis**: Setiap create/update transaction → fire rules
- **Manual**: User pilih rule/group + account + date range → execute
- Result: count of modified transactions

---

## 8. Cron/Automation System

### Semua Cron Jobs

| Job | Interval | Fungsi |
|-----|----------|--------|
| **Recurring** | 12 jam | Buat transaksi dari recurring yang due hari ini |
| **Auto Budget** | 12 jam | Buat budget limits untuk periode baru |
| **Bill Warning** | 12 jam | Kirim notifikasi bill overdue/expiring |
| **Exchange Rates** | 12 jam | Download kurs dari external API |
| **Webhooks** | **10 menit** | Kirim pending webhook messages |
| **Update Check** | **7 hari** | Cek versi baru Firefly III |

### Trigger Mechanism

```
1. Cek last-run timestamp dari config table
2. Jika elapsed > min interval → fire job
3. Jika force=true → skip timing check
4. Update timestamp setelah sukses
5. Mark preferences cache dirty
```

---

## 9. Reporting System

### Jenis Report

| Report | Konten |
|--------|--------|
| **Default** | Net Worth chart + Operations (income vs expense) |
| **Budget** | Expense per budget, pie chart per category |
| **Category** | Income/expense per category, pie chart breakdown |
| **Tag** | Income/expense per tag |
| **Expense Account** | Income/expense per expense/revenue account |
| **Double Account** | Flow through counterpart account |

### Net Worth Calculation

```
1. Ambil semua accounts (exclude: include_net_worth=false)
2. Hitung balance per account via Steam::accountsBalancesOptimized()
3. Convert ke primary currency (opsional)
4. SUBTRACT virtual_balance dari setiap account
5. Group by currency, sum per currency
6. Build time series (week by week) untuk chart
```

### Report Data Includes

- Income & expense per period (daily/weekly/monthly/quarterly/yearly)
- Pie charts: category/budget/tag/source/destination distribution
- Bar charts: period-over-period comparison
- Line charts: cumulative totals
- Drill-down: klik chart segment → tampilkan transactions

---

## 10. Multi-Currency Flow

### Setting Primary Currency

- Per UserGroup (bukan per user)
- Via preferences atau group settings

### Exchange Rate Management

```
Cron (12 jam) → DownloadExchangeRates
  → Fetch dari external API (Guzzle HTTP)
  → Store ke currency_exchange_rates table
  → Per UserGroup, per date
```

### Rate Lookup (saat konversi)

```
1. In-memory cache
2. Laravel Cache (forever)
3. Direct DB: from_currency → to_currency (date <= requested date)
4. Reverse DB: to_currency → from_currency, rate = 1/rate
5. Cross-rate via EUR: from → EUR → to
```

### Impact di Transaction

```
Setiap Transaction (create/update):
  → TransactionObserver
    → Converts amount ke native_amount (primary currency)
    → Stores both amount (original) dan native_amount (converted)
```

### Impact di Semua Monetary Model

Semua model punya pasangan amount:
- Transaction: `amount` + `native_amount`
- PiggyBankEvent: `amount` + `native_amount`
- BudgetLimit: `amount` + `native_amount`
- Bill: `amount_min` + `native_amount_min`
- Account: `virtual_balance` + `native_virtual_balance`

---

## 11. Data Export Flow

### Format: CSV Only

### Endpoint & Data

| Endpoint | Kolom Utama |
|----------|-----------|
| `/export/accounts` | id, name, type, iban, balance, currency, role, active |
| `/export/bills` | id, name, amount_min/max, currency, repeat_freq, skip, active |
| `/export/budgets` | id, name, active, start_date, end_date, currency, amount |
| `/export/categories` | id, name |
| `/export/piggy-banks` | id, name, account, currency, target_amount, current_amount |
| `/export/recurring` | id, title, type, repeat, skip, source/dest, amount, tags |
| `/export/rules` | id, title, triggers, actions |
| `/export/tags` | id, tag, date, description |
| `/export/transactions` | id, type, date, amount, currency, source, destination, category, budget, bill, tags, notes |

### Transaction Export Detail

- Default range: 1 tahun terakhir
- Bisa filter per account
- Amounts: negatif untuk withdrawal, positif untuk deposit/transfer
- Tags: comma-merged ke single field
- Include: SEPA fields, date fields, import hashes, external ID, recurrence info

---

## 12. Administration (Group) Management

### Membuat Group

```
User → /administrations → Create
  └── Hanya butuh title
```

### Invite User ke Group

```
Admin → Settings → Users → Invite
  1. Masukkan email
  2. Buat InvitedUser (code, expiry 2 hari)
  3. Kirim email invitation
  4. User register via invite link → auto-join group
```

### Manage Roles

```
Admin → /administrations → Edit
  └── Ubah role member (via group_memberships table)
```

### Switch Group

```
User → /administrations → Select group
  → user.user_group_id = selected_group_id (save ke DB)
  → Semua query otomatis scope ke group baru
```

### API Switch

```
GET /api/v1/accounts?user_group_id=2
  → scope ke group 2 untuk request ini saja
```

---

## 13. Audit Trail

### Apa yang Di-log

- **Transaction changes**: update amount, description, category, budget, source/destination, tags, date, notes, dll.
- **TransactionGroup changes**: update title

### Format

| Field | Value |
|-------|-------|
| auditable_type | `TransactionJournal`, `TransactionGroup` |
| auditable_id | ID model |
| changer_type | `User` |
| changer_id | User ID |
| action | `update_amount`, `update_description`, dll. |
| before | Old value (JSON) |
| after | New value (JSON) |

### Apa yang TIDAK Di-log

- Account changes
- Budget changes
- Category changes
- Piggy bank changes
- Rule changes
- User admin changes
- Configuration changes
- Login/logout events

---

## 14. Data Import

> **Catatan penting**: Import functionality saat ini **BUKAN bagian dari codebase ini**. Import di-handle oleh package terpisah (`firefly-iii-import-routine`).

### Yang Bisa Diinfer dari Codebase

- Tipe account `Import` ada (temporary holding account)
- Transaction punya `import_hash` dan `import_hash_v2` untuk deduplication
- Account resolution selama create transaction: by ID, IBAN, account number, atau name
- Historical: support CSV, OFX, dan banking file formats

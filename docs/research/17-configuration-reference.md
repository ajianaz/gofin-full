# Configuration & Environment Variables — Complete Reference

---

> Semua environment variables dan config values yang diperlukan Go API rewrite.
> Hanya mencakup yang critical/medium untuk API behavior.

## 1. Application Core

| Variable | Default | Purpose | Go Priority |
|----------|---------|---------|-------------|
| `APP_ENV` | `production` | Environment (local/production/testing) | HIGH |
| `APP_DEBUG` | `false` | Debug mode — controls error verbosity | HIGH |
| `APP_URL` | `http://localhost` | External URL for URL generation, emails | HIGH |
| `APP_KEY` | (none) | 32-char AES-256-CBC encryption key | HIGH |
| `TZ` | `UTC` | Application timezone — all date/time calculations | HIGH |
| `DEFAULT_LANGUAGE` | `en_US` | Default language for new users | MEDIUM |
| `DEFAULT_LOCALE` | `equal` | Number formatting locale | MEDIUM |

## 2. Authentication

### 2.1 Auth Guards

| Variable | Default | Purpose |
|----------|---------|---------|
| `AUTHENTICATION_GUARD` | `web` | Auth driver: `web` (session) or `remote_user_guard` (Authelia) |
| `AUTHENTICATION_GUARD_HEADER` | `REMOTE_USER` | HTTP header for remote user auth |
| `AUTHENTICATION_GUARD_EMAIL` | (none) | Email header for remote user auth |
| `CUSTOM_LOGOUT_URL` | (empty) | Custom logout redirect URL |

### 2.2 OAuth / Passport

| Variable | Default | Purpose |
|----------|---------|---------|
| `PASSPORT_PRIVATE_KEY` | (none) | RSA private key for JWT RS256 signing |
| `PASSPORT_PUBLIC_KEY` | (none) | RSA public key for JWT RS256 verification |

### 2.3 2FA (Google2FA)

| Variable | Default | Purpose |
|----------|---------|---------|
| `google2fa.enabled` | `true` | Whether 2FA is active |
| `google2fa.lifetime` | `0` (eternal) | Minutes before re-prompting 2FA |
| `google2fa.keep_alive` | `true` | Renew 2FA lifetime on each request |
| `google2fa.window` | `1` | OTP verification window (tolerance) |
| `google2fa.otp_secret_column` | `mfa_secret` | DB column for TOTP secret |

### 2.4 Password

| Variable | Default | Purpose |
|----------|---------|---------|
| Password reset expire | 60 min | Token validity window |
| Password reset throttle | 300 sec | Rate limit per token request |
| Password confirmation timeout | 10800 sec (3hr) | Re-auth prompt interval |

## 3. Database

| Variable | Default | Purpose |
|----------|---------|---------|
| `DB_CONNECTION` | `mysql` | Driver (mysql/pgsql/sqlite) |
| `DB_HOST` | `db` | Database host |
| `DB_PORT` | `3306` | Database port |
| `DB_DATABASE` | `firefly` | Database name |
| `DB_USERNAME` | `firefly` | Database user |
| `DB_PASSWORD` | (none) | Database password |

### MySQL SSL

| Variable | Default | Purpose |
|----------|---------|---------|
| `MYSQL_USE_SSL` | `false` | Enable MySQL SSL |
| `MYSQL_SSL_VERIFY_SERVER_CERT` | `true` | Verify server cert |
| `MYSQL_SSL_CAPATH` | `/etc/ssl/certs/` | CA certificate directory |

### PostgreSQL SSL

| Variable | Default | Purpose |
|----------|---------|---------|
| `PGSQL_SSL_MODE` | `prefer` | SSL mode |
| `PGSQL_SCHEMA` | `public` | Database schema (PG 15+) |

## 4. Redis / Cache / Queue

### 4.1 Redis

| Variable | Default | Purpose |
|----------|---------|---------|
| `REDIS_SCHEME` | `tcp` | Connection scheme (tcp/unix) |
| `REDIS_HOST` | `127.0.0.1` | Redis host |
| `REDIS_PORT` | `6379` | Redis port |
| `REDIS_USERNAME` | (empty) | Redis ACL username (Redis 6+) |
| `REDIS_PASSWORD` | (empty) | Redis password |
| `REDIS_DB` | `0` | Default Redis database index |
| `REDIS_CACHE_DB` | `1` | Cache-specific Redis database index |

### 4.2 Cache

| Variable | Default | Purpose |
|----------|---------|---------|
| `CACHE_DRIVER` | `file` | Cache backend (file/redis/database) |
| `CACHE_PREFIX` | `firefly` | Key prefix for shared caches |

### 4.3 Queue

| Variable | Default | Purpose |
|----------|---------|---------|
| `QUEUE_CONNECTION` | `sync` | Queue driver (sync/database/redis) |

## 5. Feature Flags (`config/firefly.php`)

| Flag | Default | Purpose |
|------|---------|---------|
| `feature_flags.export` | `true` | Enable data export |
| `feature_flags.webhooks` | `true` | Enable webhook system |
| `feature_flags.handle_debts` | `true` | Enable debt/loan handling |
| `feature_flags.expression_engine` | `true` | Enable expression engine for rules |
| `feature_flags.running_balance_column` | `true` | Show running balance |

## 6. Business Logic Config

### 6.1 Upload & Attachment

| Setting | Default | Purpose |
|----------|---------|---------|
| `maxUploadSize` | `1073741824` (1 GB) | Maximum upload file size |
| `valid_attachment_models` | 9 models | Account, Bill, Budget, Category, PiggyBank, Tag, Transaction, TransactionJournal, Recurrence |

### 6.2 Cron

| Variable | Default | Purpose |
|----------|---------|---------|
| `STATIC_CRON_TOKEN` | `PLEASE_REPLACE_WITH_32_CHAR_CODE` | Auth token for `GET /api/v1/cron/{token}` |

### 6.3 External Services

| Variable | Default | Purpose |
|----------|---------|---------|
| `ENABLE_EXTERNAL_RATES` | `false` | Download exchange rates from external source |
| `ENABLE_EXCHANGE_RATES` | `false` | Enable exchange rate UI |
| `ALLOW_WEBHOOKS` | `false` | Allow webhook delivery (default only) |
| `WEBHOOK_MAX_ATTEMPTS` | `3` | Max webhook retry attempts |

### 6.4 Exchange Rate Source

| Setting | Value |
|----------|-------|
| External URL | `https://ff3exchangerates.z6.web.core.windows.net/{year}/{isoWeek}/{currencyCode}.json` |
| Rate file format | `{"date": "YYYY-MM-DD", "rates": {"EUR": 1.0, "USD": 1.13, ...}}` |
| Per-user storage | Rates stored per user, per group, per date |
| Enabled currencies only | Only `currency.enabled = true` downloaded |

### 6.5 Default User Preferences

| Preference | Default | Purpose |
|------------|---------|---------|
| `anonymous` | `false` | User is anonymous |
| `frontpageAccounts` | `[]` | Accounts shown on front page |
| `listPageSize` | `50` | Default list page size |
| `currencyPreference` | `EUR` | Default currency for user |
| `language` | `en_US` | User language |
| `locale` | `equal` | Number formatting locale |
| `convertToPrimary` | `false` | Convert foreign amounts to primary currency |

## 7. Account Type Constraints

### 7.1 AccountTypeEnum (14 values)

| Enum | Value | Can Be Source? | Can Be Destination? |
|------|-------|---------------|-------------------|
| `ASSET` | Asset account | Yes | Yes |
| `DEFAULT` | Default account | Yes | Yes |
| `BENEFICIARY` | Beneficiary account | — | — |
| `CASH` | Cash account | Yes | Yes |
| `CREDITCARD` | Credit card | — | — |
| `DEBT` | Debt | Yes | Yes |
| `EXPENSE` | Expense account | **No** | Yes |
| `IMPORT` | Import account | — | — |
| `INITIAL_BALANCE` | Initial balance account | Yes | Yes |
| `LIABILITY_CREDIT` | Liability credit account | Yes | **No** |
| `LOAN` | Loan | Yes | Yes |
| `MORTGAGE` | Mortgage | Yes | Yes |
| `RECONCILIATION` | Reconciliation account | Yes | Yes |
| `REVENUE` | Revenue account | Yes | **No** |

### 7.2 `account_to_transaction` Matrix (Source → Dest → Type)

| Source | Destination | Transaction Type |
|--------|-------------|-----------------|
| Asset | Asset | Transfer |
| Asset | Expense | Withdrawal |
| Asset | Loan/Debt/Mortgage | Withdrawal |
| Asset | Cash | Withdrawal |
| Asset | InitialBalance | Opening balance |
| Asset | Reconciliation | Reconciliation |
| Revenue | Asset/Loan/Debt/Mortgage | Deposit |
| Cash | Asset/Loan/Debt/Mortgage | Deposit |
| Loan/Debt/Mortgage | Asset | Deposit |
| Loan/Debt/Mortgage | Loan/Debt/Mortgage | Transfer |
| Loan/Debt/Mortgage | Expense | Withdrawal |
| Loan/Debt/Mortgage | InitialBalance | Opening balance |
| Reconciliation | Asset | Reconciliation |
| InitialBalance | Asset/Loan/Debt/Mortgage | Opening balance |
| LiabilityCredit | Loan/Debt/Mortgage | Liability credit |

### 7.3 `source_dests` — Valid Pairs per Transaction Type

| Type | Valid Sources | Valid Destinations |
|------|--------------|-------------------|
| Withdrawal | Asset, Loan, Debt, Mortgage | Expense, Loan, Debt, Mortgage, Cash |
| Deposit | Revenue, Cash, Loan, Debt, Mortgage | Asset, Loan, Debt, Mortgage |
| Transfer | Asset, Loan, Debt, Mortgage | Same family (Asset↔Asset, Liability↔Liability) |
| Opening balance | Asset, Loan, Debt, Mortgage | InitialBalance |
| Reconciliation | Asset | Reconciliation |
| Liability credit | Loan, Debt, Mortgage | LiabilityCredit |

### 7.4 `dynamic_creation_allowed` — Auto-Created Account Types

Expense, Revenue, Initial balance, Reconciliation, Liability credit — auto-created if not exist.

### 7.5 Special Capabilities

| Capability | Account Types |
|------------|---------------|
| Can have virtual/foreign amounts | Asset only |
| Can have opening balance | Asset, Loan, Debt, Mortgage |
| Can have currency set | Asset, Loan, Debt, Mortgage, Cash, InitialBalance, LiabilityCredit, Reconciliation |
| Can hold piggy banks | Asset, Loan, Debt, Mortgage |

### 7.6 Account Roles (Asset only)

| Role | Purpose |
|------|---------|
| `defaultAsset` | Standard checking/current account |
| `sharedAsset` | Shared bank account (UI label only, NO sharing logic) |
| `savingAsset` | Savings account |
| `ccAsset` | Credit card asset |
| `cashWalletAsset` | Cash wallet |

### 7.7 Valid Meta Fields per Account Type

| Account Type | Valid Fields |
|-------------|-------------|
| Asset (basic) | `account_role`, `account_number`, `currency_id`, `BIC`, `include_net_worth` |
| Credit card | `account_role`, `cc_monthly_payment_date`, `cc_type`, `account_number`, `currency_id`, `BIC`, `include_net_worth` |
| Loan/Debt/Mortgage | `account_number`, `currency_id`, `BIC`, `interest`, `interest_period`, `include_net_worth`, `liability_direction` |

## 8. TransactionTypeEnum (7 values)

| Enum | Value |
|------|-------|
| `DEPOSIT` | Deposit |
| `WITHDRAWAL` | Withdrawal |
| `TRANSFER` | Transfer |
| `OPENING_BALANCE` | Opening balance |
| `RECONCILIATION` | Reconciliation |
| `LIABILITY_CREDIT` | Liability credit |
| `INVALID` | Invalid (error state) |

## 9. Rule Actions (22 active)

| Action | Class | Purpose |
|--------|-------|---------|
| `set_category` | SetCategory | Set transaction category |
| `clear_category` | ClearCategory | Remove category |
| `set_budget` | SetBudget | Set transaction budget |
| `clear_budget` | ClearBudget | Remove budget |
| `add_tag` | AddTag | Add a tag |
| `remove_tag` | RemoveTag | Remove a tag |
| `remove_all_tags` | RemoveAllTags | Remove all tags |
| `set_description` | SetDescription | Set description |
| `set_source_account` | SetSourceAccount | Change source account |
| `set_destination_account` | SetDestinationAccount | Change destination account |
| `set_notes` | SetNotes | Set notes |
| `clear_notes` | ClearNotes | Remove notes |
| `link_to_bill` | LinkToBill | Link to a bill |
| `convert_withdrawal` | ConvertToWithdrawal | Convert to withdrawal |
| `convert_deposit` | ConvertToDeposit | Convert to deposit |
| `convert_transfer` | ConvertToTransfer | Convert to transfer |
| `switch_accounts` | SwitchAccounts | Swap source/dest |
| `update_piggy` | UpdatePiggyBank | Add/remove from piggy bank |
| `delete_transaction` | DeleteTransaction | Delete the transaction |
| `set_source_to_cash` | SetSourceToCashAccount | Set source to cash |
| `set_destination_to_cash` | SetDestinationToCashAccount | Set dest to cash |
| `set_amount` | SetAmount | Set transaction amount |

### Commented Out (class exists but disabled in config)

`append_description`, `prepend_description`, `append_notes`, `prepend_notes`, `append_descr_to_notes`, `append_notes_to_descr`, `move_descr_to_notes`, `move_notes_to_descr`, `set_foreign_amount`, `set_foreign_currency`

### Context-Only Actions (limited set for context-triggered rules)

`set_category`, `set_budget`, `add_tag`, `remove_tag`, `set_description`, `append_description`, `prepend_description`, `set_source_account`, `set_destination_account`, `set_notes`, `append_notes`, `prepend_notes`, `link_to_bill`, `convert_transfer`

## 10. Date & View Ranges

| Category | Values |
|----------|--------|
| Valid view ranges | `1D`, `1W`, `1M`, `3M`, `6M`, `1Y` |
| Dynamic date ranges | `last7`, `last30`, `last90`, `last365`, `MTD`, `QTD`, `YTD` |
| Preselected account lists | `all`, `assets`, `liabilities` |
| Bill periods | `daily`, `weekly`, `monthly`, `quarterly`, `half-year`, `yearly` |
| Interest periods | Same as bill periods |
| Bill reminder periods | `[90, 30, 14, 7, 0]` days before due |
| Credit card types | `monthlyFull` |
| Range → Repeat frequency | `1D→weekly`, `1W→weekly`, `1M→monthly`, `3M→quarterly`, `6M→half-year`, `1Y→yearly`, `custom→custom` |

## 11. Journal Meta Fields (valid names)

**SEPA:** `sepa_cc`, `sepa_ct_op`, `sepa_ct_id`, `sepa_db`, `sepa_country`, `sepa_ep`, `sepa_ci`, `sepa_batch_id`

**Dates:** `interest_date`, `book_date`, `process_date`, `due_date`, `payment_date`, `invoice_date`

**References:** `external_id`, `external_url`, `internal_reference`, `bunq_payment_id`, `recurrence_id`

**Import:** `import_hash`, `import_hash_v2`, `original_source`

## 12. API Filters & Sort

### Allowed API Filters

| Entity | Filters |
|--------|---------|
| accounts | `name` (string), `active` (boolean), `iban` (string), `balance` (numeric), `last_activity` (date), `balance_difference` (numeric) |

### Allowed Sort Columns

| Entity | Sort Columns |
|--------|-------------|
| transactions | `description`, `amount` |
| accounts | `name`, `active`, `iban`, `order`, `account_number`, `balance`, `last_activity`, `balance_difference`, `current_debt` |

## 13. Notification Channels

| Channel | Configurable | Settings Storage |
|---------|-------------|-----------------|
| `email` | No (always on) | Mail config |
| `slack` | Yes | User pref: `slack_webhook_url` |
| `pushover` | Yes | User pref: `pushover_app_token`, `pushover_user_token` |
| `ntfy` | Yes (commented out) | User pref: `ntfy_server`, `ntfy_topic` |

## 14. Error Handling

| Variable | Default | Purpose |
|----------|---------|---------|
| `SEND_ERROR_MESSAGE` | `true` | Email errors to site owner |
| `SITE_OWNER` | (none) | Owner email for error reports |

## 15. CORS

| Setting | Default |
|----------|---------|
| paths | `['api/*']` |
| allowed_origins | `['*']` |
| supports_credentials | `false` |

## 16. Supported Languages (35)

`ar_SA`, `bg_BG`, `cs_CZ`, `da_DK`, `de_DE`, `el_GR`, `en_GB`, `en_US`, `es_ES`, `ca_ES`, `fa_IR`, `fi_FI`, `fr_FR`, `hu_HU`, `id_ID`, `it_IT`, `ja_JP`, `ko_KR`, `nb_NO`, `nn_NO`, `nl_NL`, `pl_PL`, `pt_BR`, `pt_PT`, `ro_RO`, `ru_RU`, `sk_SK`, `sl_SI`, `sv_SE`, `tr_TR`, `uk_UA`, `vi_VN`, `zh_TW`, `zh_CN`

## 17. Filesystem Disks

| Disk | Root Path | Purpose |
|------|-----------|---------|
| `upload` | `storage/upload` | User file uploads |
| `export` | `storage/export` | Data exports |

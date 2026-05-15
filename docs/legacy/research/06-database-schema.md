# Database Schema (Complete)

Source: 52 migration files di `database/migrations/`.

**Convention**: `id()` = `BIGINT UNSIGNED AUTO_INCREMENT PK`. `timestamps()` = `created_at, updated_at TIMESTAMP`. `softDeletes()` = `deleted_at TIMESTAMP NULL`. Semua amount disimpan sebagai `TEXT` (bcmath string).

---

## 1. users

| Column | Type | Nullable | Default | Notes |
|--------|------|----------|---------|-------|
| id | INT UNSIGNED PK | NO | auto | |
| objectguid | UUID | YES | NULL | LDAP objectGuid |
| email | VARCHAR(255) | NO | | |
| password | VARCHAR(60) | NO | | bcrypt |
| remember_token | VARCHAR(100) | YES | NULL | |
| reset | VARCHAR(32) | YES | NULL | |
| blocked | TINYINT UNSIGNED | NO | 0 | |
| blocked_code | VARCHAR(25) | YES | NULL | |
| mfa_secret | VARCHAR(50) | YES | NULL | 2FA secret |
| domain | VARCHAR | YES | NULL | LDAP domain |
| user_group_id | BIGINT UNSIGNED | YES | NULL | FK → user_groups.id |

**FK**: `user_group_id` → `user_groups.id` ON DELETE SET NULL

---

## 2. user_groups

| Column | Type | Nullable | Default |
|--------|------|----------|---------|
| id | BIGINT UNSIGNED PK | NO | auto |
| title | VARCHAR(255) | NO | |
| created_at | TIMESTAMP | NO | |
| updated_at | TIMESTAMP | NO | |
| deleted_at | TIMESTAMP | YES | NULL |

---

## 3. user_roles

| Column | Type | Nullable | Default |
|--------|------|----------|---------|
| id | BIGINT UNSIGNED PK | NO | auto |
| title | VARCHAR(255) | NO | unique |
| created_at | TIMESTAMP | NO | |
| updated_at | TIMESTAMP | NO | |
| deleted_at | TIMESTAMP | YES | NULL |

**Enum values** (dari UserRoleEnum): `ro`, `mng_trx`, `mng_meta`, `read_budgets`, `read_piggies`, `read_subscriptions`, `read_rules`, `read_recurring`, `read_webhooks`, `read_currencies`, `mng_budgets`, `mng_piggies`, `mng_subscriptions`, `mng_rules`, `mng_recurring`, `mng_webhooks`, `mng_currencies`, `view_reports`, `view_memberships`, `full`, `owner`

---

## 4. group_memberships

| Column | Type | Nullable | Default |
|--------|------|----------|---------|
| id | BIGINT UNSIGNED PK | NO | auto |
| user_id | BIGINT UNSIGNED | NO | |
| user_group_id | BIGINT UNSIGNED | NO | |
| user_role_id | BIGINT UNSIGNED | NO | |

**Unique**: `(user_id, user_group_id, user_role_id)`
**FK**: `user_id` → `users.id`, `user_group_id` → `user_groups.id`, `user_role_id` → `user_roles.id`

---

## 5. roles

| Column | Type | Nullable |
|--------|------|----------|
| id | INT UNSIGNED PK | NO |
| name | VARCHAR(255) | NO, unique |
| display_name | VARCHAR(255) | YES |
| description | VARCHAR(255) | YES |
| created_at | TIMESTAMP | NO |
| updated_at | TIMESTAMP | NO |

**Values**: `owner`, `demo`

## 6. role_user (pivot)

| Column | Type |
|--------|------|
| user_id | INT UNSIGNED FK → users.id |
| role_id | INT UNSIGNED FK → roles.id |

**PK**: `(user_id, role_id)`

---

## 7. accounts

| Column | Type | Nullable | Default | Notes |
|--------|------|----------|---------|-------|
| id | BIGINT UNSIGNED PK | NO | auto | |
| user_id | BIGINT UNSIGNED | NO | | FK → users.id |
| user_group_id | BIGINT UNSIGNED | NO | | FK → user_groups.id |
| account_type_id | INT UNSIGNED | NO | | FK → account_types.id |
| name | VARCHAR(255) | NO | | |
| active | TINYINT UNSIGNED | NO | 1 | |
| virtual_balance | TEXT | YES | NULL | bcmath string |
| iban | VARCHAR(255) | YES | NULL | |
| encrypted | TINYINT UNSIGNED | NO | 0 | flag for encrypted shadow cols |
| native_virtual_balance | TEXT | YES | NULL | converted to primary currency |
| created_at | TIMESTAMP | NO | |
| updated_at | TIMESTAMP | NO | |
| deleted_at | TIMESTAMP | YES | NULL | |

**FK**: `user_id` → `users.id`, `user_group_id` → `user_groups.id`, `account_type_id` → `account_types.id`

## 8. account_types

| Column | Type |
|--------|------|
| id | INT UNSIGNED PK |
| type | VARCHAR(255), unique |

**Values**: `Asset account`, `Beneficiary account`, `Cash account`, `Credit card`, `Debt`, `Default account`, `Expense account`, `Import account`, `Initial balance account`, `Liability credit account`, `Loan`, `Mortgage`, `Reconciliation account`, `Revenue account`

## 9. account_meta

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| account_id | BIGINT UNSIGNED FK → accounts.id |
| name | VARCHAR(255) |
| data | TEXT (JSON) |

**Unique**: `(account_id, name)`

---

## 10. transaction_groups

| Column | Type | Nullable | Default |
|--------|------|----------|---------|
| id | BIGINT UNSIGNED PK | NO | auto |
| user_id | BIGINT UNSIGNED | NO | |
| user_group_id | BIGINT UNSIGNED | NO | |
| title | VARCHAR(255) | NO | |
| created_at | TIMESTAMP | NO | |
| updated_at | TIMESTAMP | NO | |
| deleted_at | TIMESTAMP | YES | NULL |

## 11. transaction_journals

| Column | Type | Nullable | Default |
|--------|------|----------|---------|
| id | BIGINT UNSIGNED PK | NO | auto |
| transaction_group_id | BIGINT UNSIGNED | NO | |
| transaction_type_id | INT UNSIGNED | NO | |
| bill_id | BIGINT UNSIGNED | YES | NULL |
| transaction_currency_id | INT UNSIGNED | NO | |
| description | VARCHAR(65535) | NO | |
| completed | TINYINT UNSIGNED | NO | 0 |
| order | INT | NO | 0 |
| date | DATE | NO | |
| date_tz | VARCHAR(255) | YES | NULL |
| user_id | BIGINT UNSIGNED | NO | |
| user_group_id | BIGINT UNSIGNED | NO | |
| tag_count | INT | NO | 0 |
| encrypted | TINYINT UNSIGNED | NO | 0 |
| created_at | TIMESTAMP | NO | |
| updated_at | TIMESTAMP | NO | |
| deleted_at | TIMESTAMP | YES | NULL |

## 12. transactions

| Column | Type | Nullable | Default | Notes |
|--------|------|----------|---------|-------|
| id | BIGINT UNSIGNED PK | NO | auto | |
| transaction_journal_id | BIGINT UNSIGNED | NO | | FK → transaction_journals.id |
| account_id | BIGINT UNSIGNED | NO | | FK → accounts.id |
| transaction_currency_id | INT UNSIGNED | NO | | FK → transaction_currencies.id |
| foreign_currency_id | INT UNSIGNED | YES | NULL | FK → transaction_currencies.id |
| amount | TEXT | NO | | bcmath string |
| native_amount | TEXT | YES | NULL | auto-calc to primary currency |
| foreign_amount | TEXT | YES | NULL | bcmath string |
| reconciled | TINYINT UNSIGNED | NO | 0 | |
| identifier | INT | YES | NULL | |
| description | VARCHAR(65535) | YES | NULL | |
| encrypted | TINYINT UNSIGNED | NO | 0 | |
| created_at | TIMESTAMP | NO | |
| updated_at | TIMESTAMP | NO | |
| deleted_at | TIMESTAMP | YES | NULL |

## 13. transaction_types

| Column | Type |
|--------|------|
| id | INT UNSIGNED PK |
| type | VARCHAR(255), unique |

**Values**: `Deposit`, `Invalid`, `Liability credit`, `Opening balance`, `Reconciliation`, `Transfer`, `Withdrawal`

## 14. transaction_currencies

| Column | Type | Nullable |
|--------|------|----------|
| id | INT UNSIGNED PK | NO |
| code | VARCHAR(3), unique | NO |
| name | VARCHAR(255) | YES |
| symbol | VARCHAR(255) | YES |
| decimal_places | INT UNSIGNED | NO, default 2 |
| enabled | TINYINT UNSIGNED | NO, default 1 |
| created_at | TIMESTAMP | NO |
| updated_at | TIMESTAMP | NO |
| deleted_at | TIMESTAMP | YES | NULL |

## 15. transaction_journal_meta

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| transaction_journal_id | BIGINT UNSIGNED FK |
| name | VARCHAR(255) |
| data | TEXT (JSON) |

## 16. journal_links

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| source_id | BIGINT UNSIGNED FK → transaction_journals.id |
| destination_id | BIGINT UNSIGNED FK → transaction_journals.id |
| link_type_id | INT UNSIGNED FK → link_types.id |
| comment | VARCHAR(255) |
| deleted_at | TIMESTAMP, nullable |

**Unique**: `(source_id, destination_id, link_type_id)`

## 17. link_types

| Column | Type |
|--------|------|
| id | INT UNSIGNED PK |
| name | VARCHAR(255), unique |
| outward | VARCHAR(255) |
| inward | VARCHAR(255) |

---

## 18. budgets

| Column | Type | Nullable |
|--------|------|----------|
| id | BIGINT UNSIGNED PK | NO |
| user_id | BIGINT UNSIGNED | NO |
| user_group_id | BIGINT UNSIGNED | NO |
| name | VARCHAR(255) | NO |
| active | TINYINT UNSIGNED | NO, default 1 |
| order | INT | NO, default 0 |
| encrypted | TINYINT UNSIGNED | NO, default 0 |
| created_at | TIMESTAMP | NO |
| updated_at | TIMESTAMP | NO |
| deleted_at | TIMESTAMP | YES | NULL |

## 19. budget_limits

| Column | Type | Nullable | Notes |
|--------|------|----------|-------|
| id | BIGINT UNSIGNED PK | NO | |
| budget_id | BIGINT UNSIGNED | NO | FK → budgets.id |
| start_date | DATE | NO | |
| end_date | DATE | NO | |
| amount | TEXT | NO | bcmath string |
| native_amount | TEXT | YES | NULL | auto-calc |
| transaction_currency_id | INT UNSIGNED | NO | |
| created_at | TIMESTAMP | NO |
| updated_at | TIMESTAMP | NO |
| deleted_at | TIMESTAMP | YES | NULL |

## 20. auto_budgets

| Column | Type | Nullable | Notes |
|--------|------|----------|-------|
| id | BIGINT UNSIGNED PK | NO | |
| budget_id | BIGINT UNSIGNED | NO | |
| amount | TEXT | NO | bcmath string |
| period | VARCHAR(255) | NO | daily/weekly/monthly/quarterly/half_year/yearly |
| auto_budget_type | TINYINT UNSIGNED | NO | 1=reset, 2=rollover, 3=adjusted |

## 21. available_budgets

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| user_id | BIGINT UNSIGNED |
| user_group_id | BIGINT UNSIGNED |
| transaction_currency_id | INT UNSIGNED |
| start_date | DATE |
| end_date | DATE |
| amount | TEXT |
| native_amount | TEXT, nullable |

---

## 22. bills

| Column | Type | Nullable | Notes |
|--------|------|----------|-------|
| id | BIGINT UNSIGNED PK | NO | |
| name_encrypted | TEXT | YES | NULL | encrypted shadow |
| amount_min_encrypted | TEXT | YES | NULL | encrypted shadow |
| amount_max_encrypted | TEXT | YES | NULL | encrypted shadow |
| match_encrypted | TEXT | YES | NULL | encrypted shadow |
| name | VARCHAR(255) | NO | plaintext |
| amount_min | TEXT | NO | bcmath string |
| amount_max | TEXT | NO | bcmath string |
| match | VARCHAR(255) | NO | |
| user_id | BIGINT UNSIGNED | NO | |
| user_group_id | BIGINT UNSIGNED | NO | |
| transaction_currency_id | INT UNSIGNED | NO | |
| date | DATE | NO | |
| end_date | DATE | YES | NULL |
| repeat_freq | VARCHAR(255) | NO | weekly/monthly/quarterly/yearly |
| skip | INT | NO, default 0 | |
| active | TINYINT UNSIGNED | NO, default 1 | |
| automatch | BOOLEAN | NO, default 0 | |
| created_at | TIMESTAMP | NO | |
| updated_at | TIMESTAMP | NO | |
| deleted_at | TIMESTAMP | YES | NULL |

---

## 23. categories

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| user_id | BIGINT UNSIGNED |
| user_group_id | BIGINT UNSIGNED |
| name | VARCHAR(255) |
| encrypted | TINYINT UNSIGNED (default 0) |
| created_at, updated_at | TIMESTAMP |
| deleted_at | TIMESTAMP, nullable |

## 24. tags

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| user_id | BIGINT UNSIGNED |
| user_group_id | BIGINT UNSIGNED |
| tag | VARCHAR(255) |
| date | DATE |
| date_tz | VARCHAR(255), nullable |
| description | TEXT, nullable |
| tag_mode | VARCHAR(255), nullable |
| latitude/longitude/zoom_level | TEXT, nullable |
| encrypted | TINYINT UNSIGNED (default 0) |
| created_at, updated_at | TIMESTAMP |
| deleted_at | TIMESTAMP, nullable |

## 25. piggy_banks

| Column | Type | Nullable | Notes |
|--------|------|----------|-------|
| id | BIGINT UNSIGNED PK | NO | |
| account_id | BIGINT UNSIGNED | NO | FK → accounts.id |
| name | VARCHAR(255) | NO | |
| order | INT | NO, default 0 | |
| target_amount | TEXT | NO | bcmath string |
| start_date | DATE | YES | NULL | |
| target_date | DATE | YES | NULL | |
| start_date_tz | VARCHAR(255), nullable | | |
| target_date_tz | VARCHAR(255), nullable | | |
| active | TINYINT UNSIGNED | NO, default 1 | |
| transaction_currency_id | INT UNSIGNED | NO | |
| native_target_amount | TEXT | YES | NULL | auto-calc |
| created_at, updated_at | TIMESTAMP | |
| deleted_at | TIMESTAMP, nullable |

## 26. piggy_bank_events

| Column | Type | Nullable |
|--------|------|----------|
| id | BIGINT UNSIGNED PK | NO |
| piggy_bank_id | BIGINT UNSIGNED | NO |
| transaction_journal_id | BIGINT UNSIGNED | NO |
| amount_encrypted | TEXT | YES | encrypted shadow |
| date | DATE | NO |
| amount | TEXT | NO | bcmath string |
| native_amount | TEXT | YES | NULL |
| created_at, updated_at | TIMESTAMP |

## 27. account_piggy_bank (pivot)

| Column | Type |
|--------|------|
| account_id | BIGINT UNSIGNED FK → accounts.id |
| piggy_bank_id | BIGINT UNSIGNED FK → piggy_banks.id |
| current_amount | TEXT |
| native_current_amount | TEXT, nullable |

**PK**: `(account_id, piggy_bank_id)`

---

## 28. rule_groups

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| user_id | BIGINT UNSIGNED |
| user_group_id | BIGINT UNSIGNED |
| title | VARCHAR(255) |
| description | TEXT, nullable |
| stop_processing | TINYINT UNSIGNED (default 0) |
| active | TINYINT UNSIGNED (default 1) |
| order | INT (default 0) |
| created_at, updated_at | TIMESTAMP |
| deleted_at | TIMESTAMP, nullable |

## 29. rules

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| rule_group_id | BIGINT UNSIGNED |
| title | VARCHAR(255) |
| description | TEXT, nullable |
| priority | INT (default 0) |
| active | TINYINT UNSIGNED (default 1) |
| strict | BOOLEAN (default 1) |
| stop_processing | BOOLEAN (default 0) |
| user_id | BIGINT UNSIGNED |
| user_group_id | BIGINT UNSIGNED |
| deleted_at | TIMESTAMP, nullable |
| created_at, updated_at | TIMESTAMP |

## 30. rule_triggers

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| rule_id | BIGINT UNSIGNED |
| trigger_type | VARCHAR(255) |
| trigger_value | TEXT |
| stop_processing | BOOLEAN (default 0) |
| order | INT (default 0) |
| active | BOOLEAN (default 1) |
| created_at, updated_at | TIMESTAMP |

## 31. rule_actions

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| rule_id | BIGINT UNSIGNED |
| action_type | VARCHAR(255) |
| action_value | TEXT |
| stop_processing | BOOLEAN (default 0) |
| order | INT (default 0) |
| active | BOOLEAN (default 1) |
| created_at, updated_at | TIMESTAMP |

---

## 32. recurrences

| Column | Type | Nullable |
|--------|------|----------|
| id | BIGINT UNSIGNED PK | NO |
| user_id | BIGINT UNSIGNED | NO |
| user_group_id | BIGINT UNSIGNED | NO |
| transaction_type_id | INT UNSIGNED | NO |
| title | VARCHAR(255) | NO |
| description | TEXT | YES | NULL |
| first_date | DATE | NO |
| repeat_until | DATE | YES | NULL |
| latest_date | DATE | YES | NULL |
| repetitions | INT | YES | NULL |
| active | TINYINT UNSIGNED (default 1) |
| apply_rules | BOOLEAN (default 1) |
| automatic | BOOLEAN (default 0) |
| created_at, updated_at | TIMESTAMP |
| deleted_at | TIMESTAMP, nullable |

## 33. recurrence_repetitions

| Column | Type | Notes |
|--------|------|-------|
| id | BIGINT UNSIGNED PK | |
| recurrence_id | BIGINT UNSIGNED FK | |
| type | VARCHAR(255) | daily/weekly/monthly/ndom/yearly |
| moment | VARCHAR(255) | e.g. "1,7" (Mon,Sun) or "15" (day of month) |
| skip | INT (default 0) | skip every N+1 |
| weekend | INT (default 1) | 1=skip to Friday, 2=skip to Monday, 3=ignore |

## 34. recurrence_transactions

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| recurrence_id | BIGINT UNSIGNED FK |
| transaction_currency_id | INT UNSIGNED |
| foreign_currency_id | INT UNSIGNED, nullable |
| source_account_id | BIGINT UNSIGNED |
| source_account_name | VARCHAR(255) |
| destination_account_id | BIGINT UNSIGNED |
| destination_account_name | VARCHAR(255) |
| amount | TEXT |
| foreign_amount | TEXT, nullable |
| description | TEXT |
| budget_id | BIGINT UNSIGNED, nullable |
| category_id | BIGINT UNSIGNED, nullable |
| piggy_bank_id | BIGINT UNSIGNED, nullable |
| bill_id | BIGINT UNSIGNED, nullable |
| tags | TEXT, nullable (JSON array) |

---

## 35. currency_exchange_rates

| Column | Type | Notes |
|--------|------|-------|
| id | BIGINT UNSIGNED PK | |
| user_id | BIGINT UNSIGNED | |
| user_group_id | BIGINT UNSIGNED | |
| from_currency_id | INT UNSIGNED | |
| to_currency_id | INT UNSIGNED | |
| rate | TEXT | bcmath string |
| date | DATE | |

## 36. user_currency (pivot)

| Column | Type |
|--------|------|
| user_id | BIGINT UNSIGNED FK |
| transaction_currency_id | INT UNSIGNED FK |
| user_default | TINYINT UNSIGNED |

## 37. currency_user_group (pivot)

| Column | Type |
|--------|------|
| user_group_id | BIGINT UNSIGNED FK |
| transaction_currency_id | INT UNSIGNED FK |
| group_default | TINYINT UNSIGNED |

---

## 38. notes

| Column | Type | Notes |
|--------|------|-------|
| id | BIGINT UNSIGNED PK | |
| noteable_id | BIGINT UNSIGNED | morph ID |
| noteable_type | VARCHAR(255) | morph type |
| title | VARCHAR(255), nullable | |
| text | TEXT | |
| created_at, updated_at | TIMESTAMP |

## 39. attachments

| Column | Type | Notes |
|--------|------|-------|
| id | BIGINT UNSIGNED PK | |
| attachable_id | BIGINT UNSIGNED | morph ID |
| attachable_type | VARCHAR(255) | morph type |
| user_id | BIGINT UNSIGNED | |
| user_group_id | BIGINT UNSIGNED | |
| md5 | VARCHAR(255) | |
| filename | VARCHAR(255) | |
| title | VARCHAR(255), nullable | |
| description | TEXT, nullable | |
| mime | VARCHAR(255) | |
| size | BIGINT UNSIGNED | |
| uploaded | TINYINT UNSIGNED (default 0) | |
| created_at, updated_at | TIMESTAMP |
| deleted_at | TIMESTAMP, nullable |

## 40. locations

| Column | Type | Notes |
|--------|------|-------|
| id | BIGINT UNSIGNED PK | |
| locatable_id | BIGINT UNSIGNED | morph ID |
| locatable_type | VARCHAR(255) | morph type |
| latitude | VARCHAR(255) | |
| longitude | VARCHAR(255) | |
| zoom_level | INT, nullable | |
| created_at, updated_at | TIMESTAMP |

## 41. object_groups

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| user_id | BIGINT UNSIGNED |
| user_group_id | BIGINT UNSIGNED |
| title | VARCHAR(255) |
| order | INT (default 0) |
| created_at, updated_at | TIMESTAMP |
| deleted_at | TIMESTAMP, nullable |

## 42. object_groupables (polymorphic pivot)

| Column | Type |
|--------|------|
| object_group_id | BIGINT UNSIGNED FK |
| object_groupable_id | BIGINT UNSIGNED | morph ID |
| object_groupable_type | VARCHAR(255) | morph type (Account/Bill/PiggyBank) |

## 43. preferences

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| user_id | BIGINT UNSIGNED |
| name | VARCHAR(255) |
| data | TEXT (JSON) |
| created_at, updated_at | TIMESTAMP |

## 44. configuration

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| name | VARCHAR(255), unique |
| value | TEXT |

---

## 45. webhooks

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| user_id | BIGINT UNSIGNED |
| user_group_id | BIGINT UNSIGNED |
| title | VARCHAR(255) |
| url | VARCHAR(1024) |
| active | TINYINT UNSIGNED (default 1) |
| trigger | VARCHAR(255) | store-journal, update-journal, etc. |
| response | VARCHAR(255) | none, log, 201, error |
| delivery | VARCHAR(255) | none, individual, batch |
| created_at, updated_at | TIMESTAMP |
| deleted_at | TIMESTAMP, nullable |

## 46. webhook_messages

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| webhook_id | BIGINT UNSIGNED FK |
| transaction_journal_id | BIGINT UNSIGNED |
| user_id | BIGINT UNSIGNED |
| user_group_id | BIGINT UNSIGNED |
| message | TEXT |
| uuid | VARCHAR(255) |
| created_at, updated_at | TIMESTAMP |
| deleted_at | TIMESTAMP, nullable |

## 47. webhook_attempts

| Column | Type |
|--------|------|
| id | BIGINT UNSIGNED PK |
| webhook_message_id | BIGINT UNSIGNED FK |
| status | VARCHAR(255) | success/error |
| response_code | VARCHAR(255) |
| created_at | TIMESTAMP |

---

## Pivot Tables (Many-to-Many)

| Table | From | To | Extra Columns |
|-------|------|----|---------------|
| `role_user` | users | roles | — |
| `tag_transaction_journal` | tags | transaction_journals | — |
| `budget_transaction_journal` | budgets | transaction_journals | — |
| `category_transaction_journal` | categories | transaction_journals | — |
| `budget_transaction` | budgets | transactions | — |
| `category_transaction` | categories | transactions | — |
| `account_piggy_bank` | accounts | piggy_banks | `current_amount`, `native_current_amount` |
| `user_currency` | users | transaction_currencies | `user_default` |
| `currency_user_group` | user_groups | transaction_currencies | `group_default` |
| `group_memberships` | users+user_groups | user_roles | — |
| `object_groupables` | object_groups | Account/Bill/PiggyBank | — |

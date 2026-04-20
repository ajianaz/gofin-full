# Database Models, Schema & Relasi

## Daftar Semua Model (51 files)

### Core Financial
| Model | File | Table |
|-------|------|-------|
| **User** | `app/User.php` | `users` |
| **UserGroup** | `app/Models/UserGroup.php` | `user_groups` |
| **GroupMembership** | `app/Models/GroupMembership.php` | `group_memberships` |
| **UserRole** | `app/Models/UserRole.php` | `user_roles` |
| **Role** | `app/Models/Role.php` | `roles` |
| **Account** | `app/Models/Account.php` | `accounts` |
| **AccountType** | `app/Models/AccountType.php` | `account_types` |
| **AccountMeta** | `app/Models/AccountMeta.php` | `account_meta` |
| **TransactionGroup** | `app/Models/TransactionGroup.php` | `transaction_groups` |
| **TransactionJournal** | `app/Models/TransactionJournal.php` | `transaction_journals` |
| **Transaction** | `app/Models/Transaction.php` | `transactions` |
| **TransactionType** | `app/Models/TransactionType.php` | `transaction_types` |
| **TransactionCurrency** | `app/Models/TransactionCurrency.php` | `transaction_currencies` |
| **TransactionJournalLink** | `app/Models/TransactionJournalLink.php` | `journal_links` |
| **TransactionJournalMeta** | `app/Models/TransactionJournalMeta.php` | `transaction_journal_meta` |
| **LinkType** | `app/Models/LinkType.php` | `link_types` |

### Budget
| Model | File | Table |
|-------|------|-------|
| **Budget** | `app/Models/Budget.php` | `budgets` |
| **BudgetLimit** | `app/Models/BudgetLimit.php` | `budget_limits` |
| **AutoBudget** | `app/Models/AutoBudget.php` | `auto_budgets` |
| **AvailableBudget** | `app/Models/AvailableBudget.php` | `available_budgets` |

### Financial Lainnya
| Model | File | Table |
|-------|------|-------|
| **Bill** | `app/Models/Bill.php` | `bills` |
| **Category** | `app/Models/Category.php` | `categories` |
| **PiggyBank** | `app/Models/PiggyBank.php` | `piggy_banks` |
| **PiggyBankEvent** | `app/Models/PiggyBankEvent.php` | `piggy_bank_events` |
| **PiggyBankRepetition** | `app/Models/PiggyBankRepetition.php` | `piggy_bank_repetitions` |
| **CurrencyExchangeRate** | `app/Models/CurrencyExchangeRate.php` | `currency_exchange_rates` |

### Rule Engine
| Model | File | Table |
|-------|------|-------|
| **RuleGroup** | `app/Models/RuleGroup.php` | `rule_groups` |
| **Rule** | `app/Models/Rule.php` | `rules` |
| **RuleTrigger** | `app/Models/RuleTrigger.php` | `rule_triggers` |
| **RuleAction** | `app/Models/RuleAction.php` | `rule_actions` |

### Recurring
| Model | File | Table |
|-------|------|-------|
| **Recurrence** | `app/Models/Recurrence.php` | `recurrences` |
| **RecurrenceMeta** | `app/Models/RecurrenceMeta.php` | `recurrence_meta` |
| **RecurrenceRepetition** | `app/Models/RecurrenceRepetition.php` | `recurrence_repetitions` |
| **RecurrenceTransaction** | `app/Models/RecurrenceTransaction.php` | `recurrence_transactions` |
| **RecurrenceTransactionMeta** | `app/Models/RecurrenceTransactionMeta.php` | `recurrence_transaction_meta` |

### Supporting
| Model | File | Table |
|-------|------|-------|
| **Tag** | `app/Models/Tag.php` | `tags` |
| **Note** | `app/Models/Note.php` | `notes` |
| **Attachment** | `app/Models/Attachment.php` | `attachments` |
| **Location** | `app/Models/Location.php` | `locations` |
| **ObjectGroup** | `app/Models/ObjectGroup.php` | `object_groups` |
| **Preference** | `app/Models/Preference.php` | `preferences` |
| **Configuration** | `app/Models/Configuration.php` | `configuration` |
| **PeriodStatistic** | `app/Models/PeriodStatistic.php` | `period_statistics` |

### Webhook
| Model | File | Table |
|-------|------|-------|
| **Webhook** | `app/Models/Webhook.php` | `webhooks` |
| **WebhookAttempt** | `app/Models/WebhookAttempt.php` | `webhook_attempts` |
| **WebhookDelivery** | `app/Models/WebhookDelivery.php` | `webhook_deliveries` |
| **WebhookMessage** | `app/Models/WebhookMessage.php` | `webhook_messages` |
| **WebhookResponse** | `app/Models/WebhookResponse.php` | `webhook_responses` |
| **WebhookTrigger** | `app/Models/WebhookTrigger.php` | `webhook_triggers` |

### Lainnya
| Model | File | Table |
|-------|------|-------|
| **InvitedUser** | `app/Models/InvitedUser.php` | `invited_users` |
| **AuditLogEntry** | `app/Models/AuditLogEntry.php` | `audit_log_entries` |

---

## Enums Penting

### AccountTypeEnum (`app/Enums/AccountTypeEnum.php`)
```
Asset account, Beneficiary account, Cash account, Credit card, Debt,
Default account, Expense account, Import account, Initial balance account,
Liability credit account, Loan, Mortgage, Reconciliation account, Revenue account
```

**Kategori:**
- **User-managed**: Asset, Cash, Credit card, Loan, Mortgage, Debt
- **System/Counterpart**: Expense, Revenue, Beneficiary, Import, Initial balance, Reconciliation

### TransactionTypeEnum (`app/Enums/TransactionTypeEnum.php`)
```
Deposit, Invalid, Liability credit, Opening balance, Reconciliation, Transfer, Withdrawal
```

### UserRoleEnum (`app/Enums/UserRoleEnum.php`)
```
READ_ONLY, MANAGE_TRANSACTIONS, MANAGE_META,
READ_BUDGETS, READ_PIGGY_BANKS, READ_SUBSCRIPTIONS, READ_RULES, READ_RECURRING, READ_WEBHOOKS, READ_CURRENCIES,
MANAGE_BUDGETS, MANAGE_PIGGY_BANKS, MANAGE_SUBSCRIPTIONS, MANAGE_RULES, MANAGE_RECURRING, MANAGE_WEBHOOKS, MANAGE_CURRENCIES,
VIEW_REPORTS, VIEW_MEMBERSHIPS, FULL, OWNER
```

---

## Entity Relationship Map

### Triple-Layer Transaction

```
TransactionGroup ──(1:N)──> TransactionJournal ──(1:N)──> Transaction
     │                           │                            │
     │ user_id                   │ user_id                    │ account_id
     │ user_group_id             │ user_group_id              │ amount (string, bcmath)
     │ title                     │ transaction_type_id        │ native_amount (auto-calc)
     │                           │ transaction_currency_id    │ foreign_currency_id
     │                           │ bill_id                    │ foreign_amount
     │                           │ description                │ reconciled
     │                           │ date                       │
     │                           │ tags (BelongsToMany)       │
     │                           │ budgets (BelongsToMany)    │
     │                           │ categories (BelongsToMany) │
```

### Account Hub

```
User ──(1:N)──> Account ──(1:N)──> Transaction
                   │
                   ├── account_type_id ──> AccountType
                   ├── user_id ──> User
                   ├── user_group_id ──> UserGroup
                   ├── piggyBanks (BelongsToMany via account_piggy_bank)
                   ├── accountMeta (HasMany)
                   ├── attachments (MorphMany)
                   ├── notes (MorphMany)
                   ├── locations (MorphMany)
                   └── objectGroups (MorphToMany)
```

### User & RBAC

```
UserGroup ──(1:N)──> GroupMembership <──(N:1)── UserRole
     │                    │
     │                    ├── user_id ──> User
     │                    ├── user_group_id ──> UserGroup
     │                    └── user_role_id ──> UserRole
     │
     ├── HasMany: accounts, budgets, bills, categories, tags,
     │   transactionGroups, transactionJournals, ruleGroups, rules,
     │   recurrences, webhooks, availableBudgets, objectGroups, ...
     │
     └── BelongsToMany: TransactionCurrency (with group_default)

User ──(1:N)──> GroupMembership  (user bisa di multiple groups)
User ──(M:N)──> Role (via role_user pivot — GLOBAL roles)
```

### Polymorphic Relationships

| Polymorph | Note | Attachment | Location | ObjectGroup | AuditLog |
|-----------|------|------------|----------|-------------|----------|
| Account | Y | Y | Y | Y | |
| Budget | Y | Y | | | |
| Bill | Y | Y | | Y | |
| PiggyBank | Y | Y | | Y | |
| TransactionJournal | Y | Y | Y | | Y |
| TransactionJournalLink | Y | | | | |
| Recurrence | Y | Y | | | |
| Attachment | Y | | | | |
| BudgetLimit | Y | | | | |
| Tag | | | Y | | |

---

## Pivot Tables (Many-to-Many)

| Pivot Table | From | To | Extra Columns |
|-------------|------|----|---------------|
| `role_user` | User | Role | (composite PK) |
| `tag_transaction_journal` | Tag | TransactionJournal | |
| `budget_transaction_journal` | Budget | TransactionJournal | |
| `category_transaction_journal` | Category | TransactionJournal | |
| `budget_transaction` | Budget | Transaction | |
| `category_transaction` | Category | Transaction | |
| `account_piggy_bank` | Account | PiggyBank | `current_amount`, `native_current_amount` |
| `user_currency` | User | TransactionCurrency | `user_default` |
| `currency_user_group` | UserGroup | TransactionCurrency | `group_default` |
| `group_memberships` | User+UserGroup | UserRole | (composite unique) |
| `object_groupables` | ObjectGroup | Account/Bill/PiggyBank (morph) | |

---

## Detail Model Penting

### User (`app/User.php`)

**Fillable**: `email`, `password`, `blocked`, `blocked_code`, `user_group_id`
**Traits**: `HasApiTokens`, `Notifiable`, `ReturnsIntegerIdTrait`

**Relationships**:
| Relation | Type | Model |
|----------|------|-------|
| `userGroup()` | BelongsTo | UserGroup |
| `roles()` | BelongsToMany | Role (global) |
| `groupMemberships()` | HasMany | GroupMembership |
| `accounts()` | HasMany | Account |
| `transactionGroups()` | HasMany | TransactionGroup |
| `transactionJournals()` | HasMany | TransactionJournal |
| `transactions()` | HasManyThrough | Transaction (via TransactionJournal) |
| `budgets()` | HasMany | Budget |
| `bills()` | HasMany | Bill |
| `categories()` | HasMany | Category |
| `tags()` | HasMany | Tag |
| `piggyBanks()` | HasManyThrough | PiggyBank (via Account) |
| `rules()` | HasMany | Rule |
| `ruleGroups()` | HasMany | RuleGroup |
| `recurrences()` | HasMany | Recurrence |
| `webhooks()` | HasMany | Webhook |
| `preferences()` | HasMany | Preference |
| `currencies()` | BelongsToMany | TransactionCurrency |

**Key Methods**:
- `hasRole(string $role): bool` — cek global role
- `hasRoleInGroupOrOwner(UserGroup, UserRoleEnum): bool` — cek group-level role (cascade ke FULL/OWNER)
- `hasSpecificRoleInGroup(UserGroup, UserRoleEnum): bool` — exact role check only

### Account (`app/Models/Account.php`)

**Fillable**: `user_id`, `user_group_id`, `account_type_id`, `name`, `active`, `virtual_balance`, `iban`, `native_virtual_balance`
**Traits**: `SoftDeletes`, `HasFactory`, `ReturnsIntegerIdTrait`

**Relationships**:
| Relation | Type | Model |
|----------|------|-------|
| `user()` | BelongsTo | User |
| `userGroup()` | BelongsTo | UserGroup |
| `accountType()` | BelongsTo | AccountType |
| `transactions()` | HasMany | Transaction |
| `accountMeta()` | HasMany | AccountMeta |
| `piggyBanks()` | BelongsToMany | PiggyBank |
| `attachments()` | MorphMany | Attachment |
| `notes()` | MorphMany | Note |
| `locations()` | MorphMany | Location |
| `objectGroups()` | MorphToMany | ObjectGroup |

### UserGroup (`app/Models/UserGroup.php`)

**Fillable**: `title`

**HasMany ke semua domain model**: accounts, budgets, bills, categories, tags, transactionGroups, transactionJournals, ruleGroups, rules, recurrences, webhooks, availableBudgets, objectGroups, preferences, currencyExchangeRates, attachments, groupMemberships, periodStatistics, piggyBanks (via Account).

**BelongsToMany**: TransactionCurrency (with group_default pivot).

---

## Semua Database Tables

```
users, roles, role_user, permissions, permission_role,
user_groups, user_roles, group_memberships,
account_types, accounts, account_meta, account_balances,
transaction_types, transaction_currencies, currency_exchange_rates,
transaction_groups, transaction_journals, transaction_journal_meta, journal_links, transactions,
budgets, budget_limits, auto_budgets, available_budgets,
bills, categories, tags,
piggy_banks, piggy_bank_events, piggy_bank_repetitions, account_piggy_bank,
rule_groups, rules, rule_triggers, rule_actions,
recurrences, recurrence_meta, recurrence_repetitions, recurrence_transactions, recurrence_transaction_meta,
notes, attachments, locations, object_groups, object_groupables,
preferences, configuration,
invited_users, audit_log_entries, period_statistics,
webhooks, webhook_attempts, webhook_deliveries, webhook_messages, webhook_responses, webhook_triggers,
notifications,
jobs, sessions, password_resets,
(OAuth tables dari Passport)
```

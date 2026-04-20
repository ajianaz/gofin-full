# Observer & Event Side Effects — Complete Catalog

---

> Semua side effects yang terjadi saat create/update/delete model.
> Go API HARUS replicate exact same behavior. Tanpa ini, data integrity akan rusak.

## 1. Observer Registration

Semua observer menggunakan PHP 8 attribute `#[ObservedBy([...])]` di model class. Tidak ada registration di EventServiceProvider.

## 2. Native Amount Conversion Observers

Convert `amount` → `native_amount` pada `created` dan `updated` events. Menggunakan `ConvertsAmountToPrimaryAmount::convert()`.

### Conversion Logic

```
1. IF user has disabled conversion (convertToPrimary=false) → set native to NULL
2. IF original currency is NULL → skip
3. IF original currency == primary currency → skip (no conversion needed)
4. IF amount is empty or zero → set both to NULL
5. ELSE → lookup exchange rate → convert → save via saveQuietly()
```

### Complete Conversion Map

| Observer | Model | Source → Target | Context (User) |
|----------|-------|----------------|----------------|
| `TransactionObserver` | Transaction | `amount` → `native_amount` | `transaction→journal→date` |
| `TransactionObserver` | Transaction | `foreign_amount` → `native_foreign_amount` | same |
| `BillObserver` | Bill | `amount_min` → `native_amount_min` | `bill→user` |
| `BillObserver` | Bill | `amount_max` → `native_amount_max` | same |
| `BudgetLimitObserver` | BudgetLimit | `amount` → `native_amount` | `budget→user` |
| `AutoBudgetObserver` | AutoBudget | `amount` → `native_amount` | `autoBudget→budget→user` |
| `AvailableBudgetObserver` | AvailableBudget | `amount` → `native_amount` | `availableBudget→user` |
| `PiggyBankObserver` | PiggyBank | `target_amount` → `native_target_amount` | `piggyBank→account→user` |
| `PiggyBankEventObserver` | PiggyBankEvent | `amount` → `native_amount` | `event→piggyBank→account→user` |

### Go Implementation

```go
func ConvertToPrimary(db *sqlx.DB, model interface{}, userID int64) error {
    pref, _ := GetUserPreference(db, userID, "convertToPrimary")
    if pref == "false" {
        setNativeNull(model)
        return nil
    }
    primaryCurrencyID := getUserPrimaryCurrency(db, userID)
    modelCurrencyID := getModelCurrency(model)
    if modelCurrencyID == nil || modelCurrencyID == primaryCurrencyID {
        return nil
    }
    rate := lookupExchangeRate(db, userID, *modelCurrencyID, getDate(model))
    if rate == nil {
        return nil
    }
    setNativeAmount(model, convert(getAmount(model), rate))
    return nil
}
```

## 3. Cascade Delete Observers

Semua `Deleted*` observer hook ke event `deleting` (sebelum record dihapus dari DB).

### `DeletedAccountObserver` — Paling Kompleks

```
1. DELETE FROM account_piggy_bank WHERE account_id = ?
2. Destroy semua attachments (file + DB)
3. Collect: Transaction IDs → Journal IDs → Group IDs
4. DELETE FROM transactions WHERE journal_id IN (...)
5. DELETE FROM transaction_journals WHERE id IN (...)
6. DELETE FROM transaction_groups WHERE id IN (...)
7. DELETE FROM notes WHERE noteable_type='Account' AND noteable_id=?
8. DELETE FROM locations WHERE locatable_type='Account' AND locatable_id=?
```

### `DeletedTransactionJournalObserver` — Kedua Paling Kompleks

```
1. DELETE FROM transactions WHERE transaction_journal_id = ? (tanpa events)
2. DELETE FROM transaction_journal_links WHERE source_id = ? OR destination_id = ?
3. UPDATE piggy_bank_events SET transaction_journal_id = NULL WHERE transaction_journal_id = ?
4. DELETE FROM budget_transaction_journal WHERE transaction_journal_id = ?
5. DELETE FROM category_transaction_journal WHERE transaction_journal_id = ?
6. DELETE FROM tag_transaction_journal WHERE transaction_journal_id = ?
7. Destroy semua attachments (file + DB)
8. DELETE FROM journal_meta WHERE transaction_journal_id = ?
9. DELETE FROM locations WHERE locatable_type='TransactionJournal'
10. DELETE FROM notes WHERE noteable_type='TransactionJournal'
11. DELETE FROM audit_logs WHERE auditable_type='TransactionJournal'
```

### `DeletedTransactionGroupObserver`

```
1. DELETE FROM transaction_journals WHERE transaction_group_id = ?
   (masing-masing trigger DeletedTransactionJournalObserver)
```

### `DeletedRecurrenceObserver`

```
1. Destroy attachments (file + DB)
2. DELETE FROM recurrence_repetitions WHERE recurrence_id = ?
3. DELETE FROM recurrence_meta WHERE recurrence_id = ?
4. DELETE FROM recurrence_transactions WHERE recurrence_id = ?
   (masing-masing trigger DeletedRecurrenceTransactionObserver)
5. DELETE FROM notes
```

### `DeletedRecurrenceTransactionObserver`

```
1. DELETE FROM recurrence_transaction_meta WHERE recurrence_transaction_id = ?
```

### `DeletedRuleGroupObserver`

```
1. DELETE FROM rules WHERE rule_group_id = ?
   (masing-masing trigger DeletedRuleObserver)
```

### `DeletedRuleObserver`

```
1. DELETE FROM rule_actions WHERE rule_id = ?
2. DELETE FROM rule_triggers WHERE rule_id = ?
```

### `DeletedCategoryObserver`

```
1. Destroy attachments (file + DB)
2. DELETE FROM notes
```

### `DeletedTagObserver`

```
1. Destroy attachments (file + DB)
2. DELETE FROM tag_locations WHERE tag_id = ?
```

### `DeletedAttachmentObserver`

```
1. DELETE FROM notes WHERE noteable_type='Attachment'
```

### `DeletedWebhookObserver`

```
1. DELETE FROM webhook_messages WHERE webhook_id = ?
   (masing-masing trigger DeletedWebhookMessageObserver)
```

### `DeletedWebhookMessageObserver`

```
1. DELETE FROM webhook_attempts WHERE webhook_message_id = ?
```

### `PiggyBankObserver` (deleting)

```
1. Destroy attachments (file + DB)
2. DELETE FROM piggy_bank_events WHERE piggy_bank_id = ?
3. DELETE FROM piggy_bank_repetitions WHERE piggy_bank_id = ?
4. DELETE FROM notes
```

## 4. Event-Listener Side Effects (Critical for API)

### Transaction Lifecycle (CRITICAL — core business logic)

| Listener | Event | Side Effects |
|----------|-------|-------------|
| `ProcessesNewTransactionGroup` | `CreatedSingleTransactionGroup` | Apply rules → recalc credit → fire webhooks → remove period stats → recalc running balance |
| `ProcessesUpdatedTransactionGroup` | `UpdatedSingleTransactionGroup` | Unify source/dest across splits → apply rules → recalc credit → fire webhooks |
| `ProcessesDestroyedTransactionGroup` | `DestroyedSingleTransactionGroup` | Recalc credit → fire webhooks → remove period stats → recalc running balance |

### Account Lifecycle

| Listener | Event | Side Effects |
|----------|-------|-------------|
| `UpdatesAccountInformation` | `CreatedNewAccount` | Recalc credit → convert virtual balance to primary |
| `UpdatesAccountInformation` | `UpdatedExistingAccount` | Recalc credit → convert virtual balance → rename rule triggers/actions if name changed |

### Budget Lifecycle

| Listener | Event | Side Effects |
|----------|-------|-------------|
| `ProcessesBudgets` | `CreatedBudget` / `DestroyingBudget` / `UpdatedBudget` | Generate webhook messages |
| `ProcessesBudgetLimits` | `CreatedBudgetLimit` / `DestroyedBudgetLimit` / `UpdatedBudgetLimit` | **Recalc available budgets** + generate webhook messages |

### Exchange Rate Lifecycle

| Listener | Event | Side Effects |
|----------|-------|-------------|
| `ProcessesExchangeRates` | `Created/Destroyed/UpdatedCurrencyExchangeRate` | **Recalc ALL primary currency amounts** for affected group+currency |

### Currency Change

| Listener | Event | Side Effects |
|----------|-------|-------------|
| `RecalculatesPrimaryCurrencyAmounts` | `UserGroupChangedPrimaryCurrency` | **Recalc ALL primary currency amounts** in entire system |

### Registration

| Listener | Event | Side Effects |
|----------|-------|-------------|
| `HandlesNewUserRegistration` | `NewUserRegistered` | Create UserGroup + membership + seed exchange rates + send welcome/admin emails |

### Piggy Bank

| Listener | Event | Side Effects |
|----------|-------|-------------|
| `CreatesPiggyBankEventForChangedAmount` | `PiggyBankAmountIsChanged` | Create PiggyBankEvent record |

### Non-Critical Listeners (logging/notifications only)

| Listener | Event | Side Effect |
|----------|-------|-------------|
| `StoresNewIpAddress` | `UserSuccessfullyLoggedIn` | Record IP in preference |
| `NotifiesUserAboutNewIpAddress` | `UserLoggedInFromNewIpAddress` | Send notification |
| `NotifiesUserAboutFailedLogin` | `UserFailedLoginAttempt` | Send notification |
| `NotifiesUserAboutNewAccessToken` | `AccessTokenCreated` | Send notification |
| `NotifiesUserAboutFailedRuleAction` | `RuleActionFailed*` | Send notification |
| `NotifiesAboutExtensionOrRenewal` | `SubscriptionNeedsExtensionOrRenewal` | Send bill reminder |
| `NotifiesAboutOverdueSubscriptions` | `SubscriptionsAreOverdueForPayment` | Send overdue reminder |
| `ChecksForNewVersion` | `SystemRequestedVersionCheck` | Check GitHub for update |
| `UpdatesRulesForChangedBill` | `UpdatedExistingBill` | Rename rule references |
| `UpdatesRulesForChangedPiggyBankName` | `PiggyBankNameIsChanged` | Rename rule references |

## 5. Go Implementation Pattern

```go
// Transaction create side effects (in order)
func AfterTransactionCreate(ctx context.Context, tx *sqlx.Tx, group *TransactionGroup) error {
    // 1. Apply rules
    if err := ApplyRules(ctx, tx, group); err != nil { return err }

    // 2. Recalculate credit/liability balances
    if err := RecalculateCredit(ctx, tx, group); err != nil { return err }

    // 3. Fire webhooks
    if err := FireWebhooks(ctx, tx, "STORE_TRANSACTION", group); err != nil { return err }

    // 4. Invalidate period statistics cache
    if err := RemovePeriodStats(ctx, tx, group); err != nil { return err }

    // 5. Recalculate running balance
    if err := RecalculateRunningBalance(ctx, tx, group); err != nil { return err }

    return nil
}
```

## 6. Summary Table

| Category | Count | Critical for Go? |
|----------|-------|-----------------|
| Native amount conversion observers | 9 | YES |
| Cascade delete observers | 15 | YES |
| Critical event listeners | 7 | YES |
| Non-critical listeners (notifications) | 15+ | No (Phase 2) |
| **Total side effects** | **46+** | |

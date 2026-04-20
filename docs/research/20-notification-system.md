# Notification & Event System — Complete Catalog

---

> Notification channels, types, user preferences, event→listener wiring.

## 1. Notification Channels

| Channel | Always On | UI Configurable | Settings (user preference) |
|---------|-----------|----------------|---------------------------|
| `email` | Yes | No | — |
| `slack` | No | Yes | `slack_webhook_url` |
| `pushover` | No | Yes | `pushover_app_token`, `pushover_user_token` |
| `ntfy` | No | Yes (commented out) | `ntfy_server`, `ntfy_topic` |

### Channel Activation Logic

```
mail: selalu aktif
slack: aktif jika user punya slack_webhook_url preference
pushover: aktif jika user punya pushover_app_token AND pushover_user_token
ntfy: aktif jika user punya ntfy_server AND ntfy_topic
```

## 2. User Notification Types

### Configurable (user can enable/disable)

| Key | Trigger | Default Channels |
|-----|---------|-----------------|
| `bill_reminder` | Bill due/expiring | mail, slack, pushover |
| `transaction_creation` | New transaction stored | mail only |
| `rule_action_failures` | Rule action fails | slack, pushover (NOT mail) |
| `new_access_token` | OAuth access token created | mail, slack, pushover |
| `user_login` | Login from new IP address | mail, slack, pushover |
| `login_failure` | Failed login attempt | mail, slack, pushover |

### Non-Configurable (always sent)

| Key | Trigger | Channel |
|-----|---------|---------|
| `new_password` | Password reset requested | mail |
| `enabled_mfa` | 2FA enabled | mail |
| `disabled_mfa` | 2FA disabled | mail |
| `few_left_mfa` | Few MFA backup codes left | mail |
| `no_left_mfa` | No MFA backup codes left | mail |
| `many_failed_mfa` | Repeated MFA failures | mail |
| `new_backup_codes` | Backup codes regenerated | mail |

### Preference Storage

User preferences stored in `preferences` table:

```
notification_bill_reminder = true/false
notification_transaction_creation = true/false
notification_rule_action_failures = true/false
notification_new_access_token = true/false
notification_user_login = true/false
notification_login_failure = true/false
```

Default: all `true`.

## 3. Owner/Admin Notification Types

| Key | Trigger | Channel |
|-----|---------|---------|
| `admin_new_reg` | New admin registered | mail |
| `user_new_reg` | New user registered (welcome) | mail |
| `new_version` | New FF3 version available | mail |
| `invite_created` | User invitation created | mail |
| `invite_redeemed` | Invitation redeemed | mail |
| `unknown_user_attempt` | Unknown email tried login | mail |

> Owner = site owner (`SITE_OWNER` env var), bukan group owner.

## 4. Event Catalog (51 events)

### Security Events — User (12)

| Event | Trigger |
|-------|---------|
| `UserChangedEmailAddress` | User changes email |
| `UserFailedLoginAttempt` | Login fails |
| `UserHasDisabledMFA` | 2FA disabled |
| `UserHasEnabledMFA` | 2FA enabled |
| `UserHasFewMFABackupCodesLeft` | ≤3 backup codes left |
| `UserHasGeneratedNewBackupCodes` | New codes generated |
| `UserHasNoMFABackupCodesLeft` | 0 backup codes |
| `UserHasUsedBackupCode` | Backup code used |
| `UserKeepsFailingMFA` | Repeated 2FA failures |
| `UserLoggedInFromNewIpAddress` | New IP detected |
| `UserRequestedNewPassword` | Password reset |
| `UserSuccessfullyLoggedIn` | Successful login |

### Security Events — System (5)

| Event | Trigger |
|-------|---------|
| `NewInvitationCreated` | User invited |
| `NewUserRegistered` | New user signs up |
| `SystemFoundNewVersionOnline` | Update available |
| `SystemRequestedVersionCheck` | Cron checks for update |
| `UnknownUserTriedLogin` | Unknown email login attempt |

### Model Events — Transaction (8)

| Event | Trigger |
|-------|---------|
| `CreatedSingleTransactionGroup` | Transaction group created |
| `UpdatedSingleTransactionGroup` | Transaction group updated |
| `DestroyedSingleTransactionGroup` | Transaction group deleted |
| `TransactionGroupEventFlags` | Event flags for webhooks |
| `TransactionGroupEventObjects` | Event objects for webhooks |
| `TransactionGroupRequestsAuditLogEntry` | Audit log requested |
| `TransactionGroupsRequestedReporting` | Transaction report requested |
| `UserRequestedBatchProcessing` | Bulk update |

### Model Events — Budget (4)

| Event | Trigger |
|-------|---------|
| `CreatedBudget` | Budget created |
| `DestroyingBudget` | Budget being destroyed |
| `DestroyedBudget` | Budget destroyed |
| `UpdatedBudget` | Budget updated |

### Model Events — BudgetLimit (3)

| Event | Trigger |
|-------|---------|
| `CreatedBudgetLimit` | Budget limit created |
| `DestroyedBudgetLimit` | Budget limit destroyed |
| `UpdatedBudgetLimit` | Budget limit updated |

### Model Events — Currency (3)

| Event | Trigger |
|-------|---------|
| `CreatedCurrencyExchangeRate` | Exchange rate created |
| `DestroyedCurrencyExchangeRate` | Exchange rate deleted |
| `UpdatedCurrencyExchangeRate` | Exchange rate updated |

### Model Events — Other (10)

| Event | Trigger |
|-------|---------|
| `CreatedNewAccount` | Account created |
| `UpdatedExistingAccount` | Account updated |
| `UpdatedExistingBill` | Bill updated |
| `PiggyBankAmountIsChanged` | Piggy bank amount changed |
| `PiggyBankNameIsChanged` | Piggy bank renamed |
| `RuleActionFailedOnArray` | Rule action failed (array) |
| `RuleActionFailedOnObject` | Rule action failed (object) |
| `SubscriptionNeedsExtensionOrRenewal` | Bill needs renewal |
| `SubscriptionsAreOverdueForPayment` | Bills overdue |
| `UserGroupChangedPrimaryCurrency` | Group currency changed |

### Model Events — Webhook (1)

| Event | Trigger |
|-------|---------|
| `WebhookMessagesRequestSending` | Webhook messages ready to send |

## 5. Critical Event→Listener Wiring (for Go API)

### Must Implement (data integrity)

```
CreatedSingleTransactionGroup  → ProcessesNewTransactionGroup  [rules + credit + webhooks + stats]
UpdatedSingleTransactionGroup  → ProcessesUpdatedTransactionGroup  [unify + rules + credit + webhooks]
DestroyedSingleTransactionGroup → ProcessesDestroyedTransactionGroup [credit + webhooks + stats]
CreatedBudgetLimit             → ProcessesBudgetLimits  [recalc available budgets + webhooks]
DestroyedBudgetLimit           → ProcessesBudgetLimits
UpdatedBudgetLimit             → ProcessesBudgetLimits
CreatedCurrencyExchangeRate    → ProcessesExchangeRates  [recalc primary amounts]
DestroyedCurrencyExchangeRate  → ProcessesExchangeRates
UpdatedCurrencyExchangeRate    → ProcessesExchangeRates
UserGroupChangedPrimaryCurrency → RecalculatesPrimaryCurrencyAmounts  [recalc ALL amounts]
CreatedNewAccount              → UpdatesAccountInformation  [credit + virtual balance]
UpdatedExistingAccount         → UpdatesAccountInformation  [credit + balance + rule rename]
NewUserRegistered              → HandlesNewUserRegistration  [create group + membership + seed rates]
```

### Phase 2 (notifications only)

```
UserLoggedInFromNewIpAddress  → StoresNewIpAddress + NotifiesUserAboutNewIpAddress
UserFailedLoginAttempt        → NotifiesUserAboutFailedLogin
AccessTokenCreated            → NotifiesUserAboutNewAccessToken
UserHasEnabledMFA             → NotifiesUserAboutEnabledMFA
UserHasDisabledMFA            → NotifiesUserAboutDisabledMFA
UserHasFewMFABackupCodesLeft  → NotifiesUserAboutFewCodesLeft
UserHasNoMFABackupCodesLeft   → NotifiesUserAboutNoCodesLeft
UserKeepsFailingMFA           → NotifiesUserAboutRepeatedMFAFailures
UserHasUsedBackupCode         → NotifiesUserAboutUsedBackupCode
UserHasGeneratedNewBackupCodes → NotifiesUserAboutNewBackupCodes
UserRequestedNewPassword      → SendsUserNewPassword
RuleActionFailed*             → NotifiesUserAboutFailedRuleAction
SubscriptionNeedsExtension    → NotifiesAboutExtensionOrRenewal
SubscriptionsAreOverdue       → NotifiesAboutOverdueSubscriptions
NewInvitationCreated          → NotifiesAboutNewInvitation
NewUserRegistered             → (welcome email in HandlesNewUserRegistration)
UnknownUserTriedLogin         → NotifiesOwnerAboutUnknownUser
SystemFoundNewVersionOnline   → NotifiesOwnerAboutNewVersion
```

## 6. Go Implementation Notes

### Notification Service Interface

```go
type NotificationService interface {
    // Send notification to user through their active channels
    Send(ctx context.Context, userID int64, notificationType string, data interface{}) error

    // Get active channels for a user
    GetActiveChannels(ctx context.Context, userID int64) []string

    // Check if user has notification enabled
    IsEnabled(ctx context.Context, userID int64, notificationType string) bool
}

// Channel implementations
type EmailChannel struct { mailer *mail.Mail }
type SlackChannel struct { webhookURL string }
type PushoverChannel struct { appToken, userToken string }
```

### Event Bus Pattern

```go
// Use Go channels or event library (e.g., asynq, watermill)
type EventBus interface {
    Publish(ctx context.Context, event Event) error
    Subscribe(eventType string, handler EventHandler)
}

// Critical: transaction events must be processed AFTER commit
// to avoid race conditions with webhook/audit reads
```

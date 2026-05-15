# RBAC & Permissions

Gofin uses a hierarchical role-based access control system with two levels: **group roles** (broad permissions) and **wallet roles** (per-wallet permissions).

## Two-Level Permission Model

```
┌─────────────────────────────────────┐
│         Group Role (21 levels)       │
│  Controls what you can do in a group │
│  e.g., manage budgets, view reports │
└──────────────┬──────────────────────┘
               │ applies within
               ▼
┌─────────────────────────────────────┐
│       Wallet Role (3 levels)         │
│  Controls access to specific wallets │
│  e.g., owner, editor, viewer         │
└─────────────────────────────────────┘
```

A user must pass **both** checks to access a resource:
1. Their **group role** must meet the minimum required for the endpoint
2. Their **wallet role** (if the endpoint is wallet-scoped) must meet the required level

## Group Roles

21 hierarchical roles. Each role implicitly includes all roles below it.

| Level | Role | Description |
|-------|------|-------------|
| 21 | `owner` | Full control. Manage group settings, members, and all resources |
| 20 | `full` | Full access to all resources except group settings |
| 19 | `view_memberships` | View group member list and roles |
| 18 | `view_reports` | Access analytics and reports |
| 17 | `manage_currencies` | Add, edit, delete currencies and exchange rates |
| 16 | `read_currencies` | View currencies and exchange rates |
| 15 | `manage_webhooks` | Create and manage webhook endpoints |
| 14 | `read_webhooks` | View webhook configurations |
| 13 | `manage_recurring` | Create, edit, delete recurring transactions |
| 12 | `read_recurring` | View recurring transactions |
| 11 | `manage_rules` | Create and manage automation rules |
| 10 | `read_rules` | View automation rules |
| 9 | `manage_subscriptions` | Manage bill subscriptions |
| 8 | `read_subscriptions` | View bill subscriptions |
| 7 | `manage_piggy_banks` | Create, edit, delete piggy banks |
| 6 | `read_piggy_banks` | View piggy banks |
| 5 | `manage_budgets` | Create, edit, delete budgets |
| 4 | `read_budgets` | View budgets |
| 3 | `manage_meta` | Manage categories, tags, object groups, notes |
| 2 | `manage_transactions` | Create, edit, delete transactions |
| 1 | `read_only` | Read-only access to all data |

### How Hierarchy Works

```go
// Full access — can do everything members with role level 1-20 can do
HasPermission("owner", "read_only")     // ✅ true (21 >= 1)
HasPermission("full", "manage_budgets") // ✅ true (20 >= 5)
HasPermission("manage_meta", "manage_budgets") // ❌ false (3 < 5)
HasPermission("read_only", "manage_transactions") // ❌ false (1 < 2)
```

### Common Role Assignments

| Scenario | Recommended Role |
|----------|-----------------|
| Family member (view only) | `read_only` |
| Personal accountant | `manage_transactions` |
| Finance manager | `manage_budgets` |
| Business partner | `view_reports` |
| Co-owner | `full` |
| Group creator | `owner` |

## Wallet Roles

3 roles that control access to individual wallets.

| Role | Permissions |
|------|------------|
| **owner** | Full access. Share wallet, manage members, delete wallet |
| **editor** | Create and modify transactions on this wallet |
| **viewer** | Read-only access to transactions and balance |

### Wallet Role Hierarchy

```
owner > editor > viewer
```

- **Owner** is always the wallet creator
- **Owner** can invite others as editors or viewers
- **Editor** can create/modify transactions but cannot manage members
- **Viewer** can only view data

### Permission Matrix

| Action | Owner | Editor | Viewer |
|--------|-------|--------|--------|
| View transactions | ✅ | ✅ | ✅ |
| Create transaction | ✅ | ✅ | ❌ |
| Edit transaction | ✅ | ✅ | ❌ |
| Delete transaction | ✅ | ✅ | ❌ |
| View balance | ✅ | ✅ | ✅ |
| View members | ✅ | ✅ | ✅ |
| Invite member | ✅ | ❌ | ❌ |
| Remove member | ✅ | ❌ | ❌ |
| Change member role | ✅ | ❌ | ❌ |
| Edit wallet settings | ✅ | ❌ | ❌ |
| Delete wallet | ✅ | ❌ | ❌ |

## API-Level Permission Enforcement

Permissions are enforced at two middleware levels:

### Group RBAC Middleware

Applied to all protected endpoints. Checks the user's group role against the minimum required:

```go
// Only users with manage_budgets role or higher can access
app.Get("/api/v1/budgets", rbac.RBACMiddleware(auth.RoleManageBudgets), handler.Index)
```

### Wallet RBAC Middleware

Applied to wallet-scoped endpoints. Checks wallet membership:

```go
// Only wallet owners can delete a wallet
app.Delete("/api/v1/wallets/:wallet_id",
    middleware.WalletRBAC(memberRepo, "owner"),
    handler.Delete,
)

// Editors and owners can create transactions
app.Post("/api/v1/transactions",
    middleware.WalletRBAC(memberRepo, "editor"),
    handler.Create,
)

// All members (viewers, editors, owners) can list transactions
app.Get("/api/v1/transactions",
    middleware.WalletRBAC(memberRepo, "viewer"),
    handler.Index,
)
```

## Admin Access

The `owner` group role also grants **global admin** privileges:

- Access to `/api/v1/admin/*` endpoints
- Create system-wide users
- View audit logs
- Toggle feature flags

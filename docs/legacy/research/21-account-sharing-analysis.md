# Account Sharing & Wallet Model — Analysis & Redesign

---

> Analisis mendalam model account Firefly III dan proposal redesign untuk Go API.
> Menjawab pertanyaan: rename "account" → "wallet/kantong", sharing ke multiple users.

## 1. Current State — Firefly III

### 1.1 Apakah 1 User Bisa Punya Beberapa Account?

**Ya.** User bisa punya unlimited accounts. Relasi: `User ──(1:N)──> Account` via `accounts.user_id`.

### 1.2 Current Account Model

```
accounts table:
  id              BIGINT PK
  user_id         BIGINT FK → users        ← OWNER (exclusive)
  user_group_id   BIGINT FK → user_groups  ← GROUP (organizational)
  account_type_id BIGINT FK → account_types
  name            VARCHAR(255)
  active          BOOLEAN
  virtual_balance DECIMAL  (string)
  iban            VARCHAR(255)
  native_virtual_balance DECIMAL (string)
```

### 1.3 Apakah Ada Sharing? **TIDAK ADA.**

- Tidak ada `account_user` pivot table
- Tidak ada `account_members` table
- Tidak ada Laravel Policies
- `Account::routeBinder()` cek `user_id` saja → kalau bukan pemilik → 404
- `account_role` di meta hanyalah label UI (`sharedAsset`, `savingAsset`, dll) — **bukan permission**

### 1.4 UserGroup Scoping (Partial)

UserGroup punya 22 group-level permissions, tapi **tidak mengontrol akses ke account individual**:

```
UserGroup
  ├── User A (role: OWNER)
  │     └── accounts: [Wallet A1, Wallet A2]
  ├── User B (role: FULL)
  │     └── accounts: [Wallet B1]
  └── User B TIDAK BISA akses Wallet A1 (karena user_id berbeda)
```

### 1.5 Problem Statement

| Problem | Impact |
|---------|--------|
| Account single-owner | Pasangan/keluarga tidak bisa share wallet |
| Group RBAC tidak apply ke account | User FULL di group tetap tidak bisa lihat account user lain |
| `sharedAsset` role hanya label | Menyesatkan — tidak ada logic sharing |
| Tidak ada read-only access | Tidak bisa share view-only ke anggota keluarga |

## 2. Proposal: Rename "Account" → "Wallet"

### 2.1 Terminology Mapping

| Firefly III (Lama) | Go API (Baru) | Keterangan |
|--------------------|---------------|------------|
| `Account` | `Wallet` | Model utama — dompet/kantong |
| `account_type` | `wallet_type` | Jenis wallet |
| `Asset account` | `Wallet` | Dompet utama (bank, e-wallet) |
| `Expense account` | `Payee` | Penerima pembayaran (auto-created) |
| `Revenue account` | `Income Source` | Sumber income (auto-created) |
| `Initial balance account` | `System: Opening Balance` | Auto-created system |
| `Reconciliation account` | `System: Reconciliation` | Auto-created system |
| `Import account` | `System: Import` | Auto-created system |
| `Cash account` | `Cash Wallet` | Tunai |
| `Credit card` | `Credit Wallet` | Kartu kredit |
| `Loan` | `Loan Wallet` | Pinjaman |
| `Debt` | `Debt Wallet` | Hutang |
| `Mortgage` | `Mortgage Wallet` | KPR |
| `Liability credit` | `System: Liability Credit` | Auto-created system |
| `Beneficiary account` | (deprecated) | Alias untuk Expense |
| `Default account` | (deprecated) | Alias untuk Asset |

### 2.2 Wallet Types (Go Enum)

```go
type WalletType string

const (
    WalletTypeAsset           WalletType = "asset"           // User wallet
    WalletTypeCash            WalletType = "cash"            // Cash wallet
    WalletTypeCreditCard      WalletType = "credit_card"     // Credit card
    WalletTypeLoan            WalletType = "loan"            // Loan
    WalletTypeDebt            WalletType = "debt"            // Debt
    WalletTypeMortgage        WalletType = "mortgage"        // Mortgage
    // System types (auto-created, hidden from user)
    WalletTypeExpense         WalletType = "expense"         // Payee
    WalletTypeRevenue         WalletType = "revenue"         // Income source
    WalletTypeInitialBalance  WalletType = "initial_balance"  // Opening balance
    WalletTypeReconciliation  WalletType = "reconciliation"   // Reconciliation
    WalletTypeImport          WalletType = "import"          // Import
    WalletTypeLiabilityCredit WalletType = "liability_credit" // Liability credit
)
```

### 2.3 User-Facing vs System Wallets

| Category | Types | User Can Create? | User Can See? |
|----------|-------|-----------------|---------------|
| **User Wallets** | asset, cash, credit_card, loan, debt, mortgage | Yes | Yes |
| **System Wallets** | expense, revenue, initial_balance, reconciliation, import, liability_credit | No (auto) | No (hidden) |

## 3. Proposal: Wallet Sharing

### 3.1 New Table: `wallet_members`

```sql
CREATE TABLE wallet_members (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    wallet_id       BIGINT NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role            VARCHAR(20) NOT NULL DEFAULT 'viewer',
    invited_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY (wallet_id, user_id)
);

CREATE INDEX idx_wallet_members_user ON wallet_members(user_id);
CREATE INDEX idx_wallet_members_wallet ON wallet_members(wallet_id);
```

### 3.2 Wallet Member Roles

| Role | Value | Permissions |
|------|-------|------------|
| **Owner** | `owner` | Full CRUD + manage members + delete wallet |
| **Editor** | `editor` | Create/edit/delete transactions |
| **Viewer** | `viewer` | View only (read transactions, balances) |

### 3.3 Access Control Logic

```go
func CanAccessWallet(ctx context.Context, userID int64, walletID int64, requiredRole string) bool {
    // 1. Check ownership (wallets.user_id = userID)
    if isOwner(ctx, userID, walletID) {
        return true // owner has full access
    }

    // 2. Check membership (wallet_members)
    member := getWalletMember(ctx, walletID, userID)
    if member == nil {
        return false
    }

    // 3. Role hierarchy: owner > editor > viewer
    roleHierarchy := map[string]int{
        "owner":  3,
        "editor": 2,
        "viewer": 1,
    }

    return roleHierarchy[member.Role] >= roleHierarchy[requiredRole]
}

// Usage:
// Create transaction: requiredRole = "editor"
// View transactions:  requiredRole = "viewer"
// Delete wallet:      requiredRole = "owner"
// Manage members:     requiredRole = "owner"
```

### 3.4 Wallet Sharing API

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| `GET` | `/api/v1/wallets/{id}/members` | List wallet members | owner/editor/viewer |
| `POST` | `/api/v1/wallets/{id}/members` | Add member | owner |
| `PUT` | `/api/v1/wallets/{id}/members/{userId}` | Change member role | owner |
| `DELETE` | `/api/v1/wallets/{id}/members/{userId}` | Remove member | owner |
| `GET` | `/api/v1/wallets/shared-with-me` | List wallets shared with me | any |

### 3.5 Share Wallet Request/Response

**POST /api/v1/wallets/{id}/members**

```json
// Request
{
  "user_id": 5,
  "role": "viewer"
}

// Response
{
  "wallet_id": 1,
  "user_id": 5,
  "email": "partner@example.com",
  "role": "viewer",
  "invited_by": 1,
  "created_at": "2024-01-01T00:00:00Z"
}
```

### 3.6 Transaction Access with Sharing

Saat wallet di-share, member dengan role `editor` atau `viewer` bisa:

```go
// Modified account access check
func GetAccessibleWallets(ctx context.Context, userID int64) []int64 {
    // 1. Own wallets
    ownWallets := getWalletIDsByOwner(ctx, userID)

    // 2. Shared wallets
    sharedWallets := getWalletIDsByMembership(ctx, userID)

    return append(ownWallets, sharedWallets...)
}

// Transaction queries scope to accessible wallets
func GetUserTransactions(ctx context.Context, userID int64, db *sqlx.DB) []Transaction {
    walletIDs := GetAccessibleWallets(ctx, userID)
    return queryTransactionsByWalletIDs(ctx, db, walletIDs)
}
```

### 3.7 Sharing Rules

| Rule | Description |
|------|-------------|
| Owner selalu `wallets.user_id` | Owner tidak bisa di-kick, role tidak bisa diubah |
| System wallets tidak bisa di-share | Expense, revenue, initial balance, dll |
| Editor tidak bisa manage members | Hanya owner |
| Viewer tidak bisa create/edit/delete transactions | Read-only |
| Member bisa leave sendiri | DELETE /api/v1/wallets/{id}/members/me |
| Sharing tidak mengubah UserGroup scope | Wallet sharing independent dari group membership |

## 4. Impact on Existing Features

### 4.1 Transaction Creation

| Scenario | Source Account | Destination Account |
|----------|---------------|---------------------|
| User creates from own wallet | Must be owner or editor | Can be any (including payee) |
| User creates to shared wallet | Must be owner or editor of dest | Must be owner or editor of source |

### 4.2 Search & Filter

Transaction search HARUS scope ke accessible wallets:

```go
// Before (Firefly III): WHERE source.account_id IN (user's own accounts)
// After (Go):           WHERE source.account_id IN (own + shared wallets)
```

### 4.3 Reports & Charts

Reports HARUS hanya include accessible wallets.

### 4.4 Webhooks

Webhooks fire untuk transaction yang ter-create di wallet yang di-share — semua member (owner, editor, viewer) yang punya webhook aktif akan menerima.

### 4.5 Rules

Rules berjalan pada wallet yang di-share — semua member melihat hasilnya.

## 5. Database Schema Changes Summary

### New Tables

```sql
-- Wallet members (sharing)
CREATE TABLE wallet_members (...);

-- Rename accounts → wallets (optional, bisa pakai view)
CREATE VIEW wallets AS SELECT * FROM accounts;
```

### Modified Tables

```sql
-- accounts table: no schema change needed
-- (user_id stays as owner, wallet_members handles sharing)
```

### Migration Path (from Firefly III)

```
1. Create wallet_members table (empty)
2. Existing accounts: user_id = owner, no members
3. Sharing is opt-in: users manually add members
4. No data migration needed for existing accounts
```

## 6. Recommendation

### Phase 1: Core (No Sharing)

- Rename `Account` → `Wallet` di API level
- Keep `user_id` as owner
- Implement wallet CRUD with all type constraints from doc 17
- Keep existing UserGroup scoping

### Phase 2: Wallet Sharing

- Add `wallet_members` table
- Add sharing API endpoints
- Modify transaction queries to include shared wallets
- Add role-based access (owner/editor/viewer)
- Update search, reports, charts to respect sharing

### Phase 3: Advanced Sharing

- Sharing notifications (invite, accept, leave)
- Activity log per wallet
- Pending invitations
- Share via email link

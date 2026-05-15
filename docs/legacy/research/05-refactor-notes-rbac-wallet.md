# Refactor Notes: RBAC per Account/Wallet

## Analisis Kondisi Saat Ini

### Apa yang Sudah Ada

1. **UserGroup (Administration)** — shared workspace. Semua data milik group.
2. **Two-tier RBAC** — Global roles (`owner`, `demo`) + Group-level roles (22 permissions).
3. **Multi-account** — User bisa punya banyak financial accounts (Asset, Cash, Credit card, dll) dalam satu group.
4. **Account types** — 14 tipe: Asset, Cash, Credit card, Loan, Mortgage, Debt, Expense, Revenue, dll.

### Apa yang TIDAK Ada (Gap Analysis)

| Gap | Detail |
|-----|--------|
| **Tidak ada RBAC per account** | Role hanya di level Group. Tidak bisa memberi akses ke account tertentu saja. |
| **Tidak ada account sharing** | Account hanya punya satu `user_id`. Tidak ada konsep "member" di level account. |
| **Tidak ada Laravel Policies** | Authorization dilakukan manual via traits dan middleware. |
| **Tidak ada API scopes** | Token OAuth punya akses penuh ke semua data di group. |
| **Tidak ada granular ownership** | Account langsung milik user (`user_id`). Bukan shared. |
| **Binder returns 404** | `routeBinder()` throw `NotFoundHttpException` untuk unauthorized access — ini sengaja untuk hide resources, tapi kurang informatif untuk debugging. |
| **Group switching di DB** | `user_group_id` disimpan di kolom user, bukan session. Bisa bikin race condition. |

### Model Ownership Saat Ini

```
User ──(1:N)──> Account
     │
     └──(1:N)──> GroupMembership <──(N:1)──> UserGroup ──(1:N)──> Account
```

Account punya **dua ownership**: `user_id` (langsung) dan `user_group_id` (via group). Tapi akses ke account di-control **hanya via group**, bukan per-account.

---

## Proposal: RBAC per Account/Wallet

### Visi

1 user bisa punya beberapa **wallet** (account).
1 wallet bisa diakses oleh **beberapa user** dengan **role berbeda**.
Contoh: User A punya wallet pribadi (full access) dan wallet keluarga (shared dengan User B yang hanya READ_ONLY).

### Arsitektur yang Diusulkan

#### 1. Tabel Baru: `account_members`

```sql
CREATE TABLE account_members (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    account_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    account_role VARCHAR(50) NOT NULL DEFAULT 'viewer',
    invited_by BIGINT UNSIGNED NULL,
    joined_at TIMESTAMP NULL,
    UNIQUE KEY unique_account_user (account_id, user_id),
    FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (invited_by) REFERENCES users (id) ON DELETE SET NULL
);
```

#### 2. Enum Baru: `AccountRoleEnum`

```php
enum AccountRoleEnum: string
{
    case VIEWER    = 'viewer';      // Hanya lihat saldo & transaksi
    case REPORTER  = 'reporter';    // Viewer + generate reports
    case CONTRIBUTOR = 'contributor'; // Reporter + create/edit transactions
    case MANAGER   = 'manager';     // Contributor + manage budgets, categories, piggy banks
    case CO_OWNER  = 'co_owner';    // Manager + manage members, edit account settings
    case OWNER     = 'owner';       // Full control + delete account
}
```

**Hierarchy**: `OWNER > CO_OWNER > MANAGER > CONTRIBUTOR > REPORTER > VIEWER`

#### 3. Perubahan Model

**Account** — tambah relationships:

```php
// app/Models/Account.php

public function members(): HasMany
{
    return $this->hasMany(AccountMember::class);
}

public function memberUsers(): BelongsToMany
{
    return $this->belongsToMany(User::class, 'account_members')
        ->withPivot(['account_role', 'invited_by', 'joined_at'])
        ->withTimestamps();
}
```

**User** — tambah relationships:

```php
// app/User.php

public function accountMemberships(): HasMany
{
    return $this->hasMany(AccountMember::class);
}

public function sharedAccounts(): BelongsToMany
{
    return $this->belongsToMany(Account::class, 'account_members')
        ->withPivot(['account_role', 'invited_by', 'joined_at'])
        ->withTimestamps();
}
```

#### 4. Model Baru: `AccountMember`

```php
class AccountMember extends Model
{
    use SoftDeletes;

    protected $fillable = [
        'account_id', 'user_id', 'account_role', 'invited_by', 'joined_at',
    ];

    protected $casts = [
        'joined_at' => 'datetime',
    ];

    // Relationships
    public function account(): BelongsTo { /* ... */ }
    public function user(): BelongsTo { /* ... */ }
    public function inviter(): BelongsTo { /* ... */ }

    // Authorization helpers
    public function hasRoleOrHigher(AccountRoleEnum $role): bool { /* ... */ }
}
```

#### 5. Authorization Helpers di User

```php
// app/User.php

public function getAccountRole(Account $account): ?AccountRoleEnum
{
    $membership = AccountMember::where('account_id', $account->id)
        ->where('user_id', $this->id)
        ->first();

    return $membership ? AccountRoleEnum::from($membership->account_role) : null;
}

public function hasAccountRoleOrHigher(Account $account, AccountRoleEnum $role): bool
{
    $userRole = $this->getAccountRole($account);
    if (!$userRole) return false;

    // OWNER dan CO_OWNER dari UserGroup tetap override
    if ($this->hasRoleInGroupOrOwner($account->userGroup, UserRoleEnum::FULL)) {
        return true;
    }

    return $userRole->rank() >= $role->rank();
}
```

---

## Strategi Implementasi (Phased)

### Phase 1: Foundation (Non-Breaking)

1. **Buat migration** `account_members` table
2. **Buat model** `AccountMember`
3. **Buat enum** `AccountRoleEnum`
4. **Auto-seed**: Saat create Account, otomatis buat `AccountMember` dengan role OWNER untuk user yang membuatnya
5. **Buat relationships** di Account dan User model
6. **Buat API endpoints**:
   - `GET /api/v1/accounts/{account}/members` — list members
   - `POST /api/v1/accounts/{account}/members` — invite member
   - `PUT /api/v1/accounts/{account}/members/{member}` — update role
   - `DELETE /api/v1/accounts/{account}/members/{member}` — remove member

### Phase 2: Authorization Layer

1. **Buat Account Policy** (Laravel Policy):
   ```php
   class AccountPolicy
   {
       public function view(User $user, Account $account): bool
       {
           return $user->hasAccountRoleOrHigher($account, AccountRoleEnum::VIEWER);
       }

       public function update(User $user, Account $account): bool
       {
           return $user->hasAccountRoleOrHigher($account, AccountRoleEnum::CO_OWNER);
       }

       public function delete(User $user, Account $account): bool
       {
           return $user->hasAccountRoleOrHigher($account, AccountRoleEnum::OWNER);
       }

       public function createTransaction(User $user, Account $account): bool
       {
           return $user->hasAccountRoleOrHigher($account, AccountRoleEnum::CONTRIBUTOR);
       }

       public function manageBudgets(User $user, Account $account): bool
       {
           return $user->hasAccountRoleOrHigher($account, AccountRoleEnum::MANAGER);
       }

       public function manageMembers(User $user, Account $account): bool
       {
           return $user->hasAccountRoleOrHigher($account, AccountRoleEnum::CO_OWNER);
       }
   }
   ```

2. **Buat Transaction Policy**:
   ```php
   class TransactionPolicy
   {
       public function view(User $user, Transaction $transaction): bool
       {
           return $user->hasAccountRoleOrHigher($transaction->account, AccountRoleEnum::VIEWER);
       }

       public function create(User $user, Account $source, Account $dest): bool
       {
           return $user->hasAccountRoleOrHigher($source, AccountRoleEnum::CONTRIBUTOR)
               && $user->hasAccountRoleOrHigher($dest, AccountRoleEnum::CONTRIBUTOR);
       }
   }
   ```

3. **Integrate ke Binder** — update `Account::routeBinder()` untuk cek `account_members`:
   ```php
   public static function routeBinder($value): Account
   {
       $account = Account::where('id', (int)$value)
           ->where('user_group_id', $user->user_group_id)
           ->first();

       // Check via account_members
       $member = AccountMember::where('account_id', $account->id)
           ->where('user_id', $user->id)
           ->first();

       if (!$member && !$user->hasRoleInGroupOrOwner($account->userGroup, UserRoleEnum::FULL)) {
           throw new NotFoundHttpException;
       }

       return $account;
   }
   ```

### Phase 3: Scoping Queries

1. **Scope transactions** — filter berdasarkan account membership:
   ```php
   // app/Repositories/Transaction/TransactionRepository.php
   public function setUser(User $user): void
   {
       $this->user = $user;
       $accessibleAccountIds = AccountMember::where('user_id', $user->id)
           ->pluck('account_id')
           ->toArray();

       // Jika FULL/OWNER di group, skip scoping
       if ($user->hasRoleInGroupOrOwner($user->userGroup, UserRoleEnum::FULL)) {
           return; // akses semua account di group
       }

       $this->accessibleAccountIds = $accessibleAccountIds;
   }
   ```

2. **Scope API list endpoints** — hanya tampilkan account yang user punya akses

### Phase 4: API Scopes (Optional)

1. **Definisikan Passport scopes**:
   ```php
   // AppServiceProvider
   Passport::tokensCan([
       'accounts:read'    => 'Read account balances',
       'accounts:write'   => 'Create/edit accounts',
       'transactions:read'  => 'Read transactions',
       'transactions:write' => 'Create/edit transactions',
       'reports'          => 'Generate reports',
   ]);
   ```

2. **Apply scope middleware** ke API routes

### Phase 5: Notification & Invitation

1. **Email invitation** saat add member ke account
2. **Accept/decline flow** via email link
3. **Notification** saat role di-change atau di-remove

---

## Backward Compatibility

### Yang TIDAK Berubah

- UserGroup tetap sebagai level tertinggi scoping
- Global roles (`owner`, `demo`) tetap
- Group-level roles (`UserRoleEnum`) tetap
- Semua existing API endpoints tetap (tambah parameter optional)
- Existing `Account.user_id` tetap (sebagai "primary owner")

### Yang Berubah

- `Account::routeBinder()` — tambah cek `account_members`
- Transaction queries — scope ke accessible accounts
- Account list — filter by membership
- Tambah API endpoints untuk member management

### Migration Strategy

```php
// Step 1: Buat tabel baru
Schema::create('account_members', function (Blueprint $table) { /* ... */ });

// Step 2: Seed existing accounts
$accounts = Account::all();
foreach ($accounts as $account) {
    AccountMember::firstOrCreate([
        'account_id' => $account->user_id,
        'user_id'    => $account->user_id,
        'account_role' => 'owner',
        'joined_at'  => $account->created_at,
    ]);
}

// Step 3: Tambah index untuk performance
Schema::table('account_members', function (Blueprint $table) {
    $table->index(['user_id', 'account_role']);
    $table->index(['account_id', 'account_role']);
});
```

---

## Risk & Pertimbangan

| Risk | Mitigasi |
|------|----------|
| **Breaking change** | Phase 1 & 2 tidak break existing behavior. Account owner otomatis di-seed. |
| **Performance** | `account_members` query per-request. Mitigasi: eager load, cache membership per session. |
| **Complex authorization** | Dua level scoping (group + account). Mitigasi: group FULL/OWNER override account-level checks. |
| **Data leakage** | Reporter bisa lihat semua transaksi. Mitigasi: dokumentasi jelas per role. |
| **Invitation abuse** | User invite siapa saja. Mitigasi: hanya CO_OWNER+ yang bisa invite. |
| **Race condition group switch** | `user_group_id` di DB. Mitigasi: pertimbangkan pindah ke session. |

---

## Pertanyaan Terbuka

1. **Apakah perlu role per-account atau per-account-type?** (misalnya: user bisa CONTRIBUTOR di semua Asset accounts tapi VIEWER di Credit Card)
2. **Apakah budget/piggy bank perlu di-share terpisah dari account?** (saat ini budget/piggy bank milik group)
3. **Apakah notification system perlu di-enhance untuk invitation?**
4. **Apakah perlu audit log untuk member changes?** (role assigned, removed, invitation sent)
5. **Apakah perlu soft-delete untuk membership?** (supaya bisa di-restore)
6. **Bagaimana handle split transactions antara account yang berbeda permission?** (source account accessible tapi destination tidak)

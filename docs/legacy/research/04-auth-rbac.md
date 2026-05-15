# Authentication & Authorization System

## Authentication Guards

| Guard | Driver | Provider | Purpose |
|-------|--------|----------|---------|
| `web` (default) | `session` | `users` (Eloquent) | Web login, session cookies. Remember-me: 364 hari |
| `remote_user_guard` | `remote_user_guard` | `remote_user_provider` | SSO / reverse-proxy auth (REMOTE_USER header) |
| `api` | `passport` | `users` (Eloquent) | OAuth2 API token auth via Laravel Passport |

Default guard dikontrol via `AUTHENTICATION_GUARD` env variable.

---

## Login Flow

**File**: `app/Http/Controllers/Auth/LoginController.php`

1. Login field: `email` (hardcoded)
2. Throttled login attempts
3. On success: fire `UserSuccessfullyLoggedIn` event
4. On failure: fire `UnknownUserTriedLogin` atau `UserFailedLoginAttempt` event
5. Route: `POST /login`

---

## Registration Flow

**File**: `app/Http/Controllers/Auth/RegisterController.php`

1. Registration bisa di-disable jika `single_user_mode=true` dan sudah ada user
2. Support invitation codes (`InvitedUser` records, 2-day expiry)
3. On register: fire `NewUserRegistered` event

### Side Effects (New User Registration)

**File**: `app/Listeners/Security/System/HandlesNewUserRegistration.php`

```
1. Attach global role: Jika user pertama, attach role 'owner'
2. Create UserGroup: title = user's email
3. Create GroupMembership: user + group + OWNER role
4. Set user.user_group_id = new group
5. Send notification emails
```

---

## Password Reset

- Reset tokens expire: 60 menit
- Throttle: 1 request per 300 detik
- Password confirmation timeout: 3 jam (10800 detik)
- Min password length: 16 karakter

---

## Two-Factor Authentication (2FA/MFA)

**Package**: `pragmarx/google2fa-laravel`

- MFA codes disimpan di kolom `mfa_secret`
- **Replay protection**: Codes di-track 5 menit history
- **Backup codes**: Disimpan di preferences (`mfa_recovery`), warning saat <=3 codes
- **Failure counter**: Tracked di preferences (`mfa_failure_count`), warning di 3 dan 10 attempts

---

## Remote User Authentication (SSO)

**File**: `app/Support/Authentication/RemoteUserProvider.php`

- Custom `UserProvider` untuk reverse-proxy SSO
- Auto-create User dengan random password saat pertama kali
- User pertama otomatis dapat role `owner`

---

## Two-Tier Role System

### Tier 1: Global Roles (`roles` table)

**Table**: `roles` (name, display_name, description)
**Pivot**: `role_user` (User M:N Role)

| Role | Purpose | Check Method |
|------|---------|--------------|
| `owner` | System administrator | `IsAdmin` middleware, `$user->hasRole('owner')` |
| `demo` | Read-only restricted | `IsDemoUser` middleware |

- User pertama otomatis dapat `owner` role
- Admin bisa grant/revoke `owner` role di `/settings/users`

### Tier 2: Group-Level Roles (`user_roles` table)

**Table**: `user_roles` (title)
**Pivot**: `group_memberships` (user_id, user_group_id, user_role_id)

**Enum**: `app/Enums/UserRoleEnum.php` — 22 permission values

| Enum Value | Title | Description |
|------------|-------|-------------|
| `READ_ONLY` | `ro` | Bisa baca semua kecuali members |
| `MANAGE_TRANSACTIONS` | `mng_trx` | CRUD transactions (required untuk pakai group) |
| `MANAGE_META` | `mng_meta` | Edit categories/tags/object-groups |
| `READ_BUDGETS` | `read_budgets` | Baca budgets |
| `READ_PIGGY_BANKS` | `read_piggies` | Baca piggy banks |
| `READ_SUBSCRIPTIONS` | `read_subscriptions` | Baca subscriptions/bills |
| `READ_RULES` | `read_rules` | Baca rules |
| `READ_RECURRING` | `read_recurring` | Baca recurring |
| `READ_WEBHOOKS` | `read_webhooks` | Baca webhooks |
| `READ_CURRENCIES` | `read_currencies` | Baca currencies |
| `MANAGE_BUDGETS` | `mng_budgets` | Manage budgets |
| `MANAGE_PIGGY_BANKS` | `mng_piggies` | Manage piggy banks |
| `MANAGE_SUBSCRIPTIONS` | `mng_subscriptions` | Manage subscriptions |
| `MANAGE_RULES` | `mng_rules` | Manage rules |
| `MANAGE_RECURRING` | `mng_recurring` | Manage recurring |
| `MANAGE_WEBHOOKS` | `mng_webhooks` | Manage webhooks |
| `MANAGE_CURRENCIES` | `mng_currencies` | Manage currencies |
| `VIEW_REPORTS` | `view_reports` | View/generate reports |
| `VIEW_MEMBERSHIPS` | `view_memberships` | Lihat group members |
| `FULL` | `full` | Semua kecuali hapus creator & hapus group |
| `OWNER` | `owner` | Reserved untuk original creator |

**Permission Hierarchy**:
```
OWNER
  └── FULL
        └── VIEW_MEMBERSHIPS
              └── VIEW_REPORTS
                    └── MANAGE_* (per domain)
                          └── READ_* (per domain)
                                └── MANAGE_META
                                      └── MANAGE_TRANSACTIONS
                                            └── READ_ONLY
```

`FULL` dan `OWNER` **cascade down** — jika user punya role FULL/OWNER di sebuah group, mereka otomatis punya semua role di bawahnya.

---

## Authorization Enforcement

### Tidak ada Laravel Gates atau Policies

Authorization diimplementasikan melalui:

### 1. Middleware

| Middleware | Check |
|-----------|-------|
| `Authenticate` | `auth()->check()` + not blocked |
| `IsAdmin` / `IsAdminApi` | global `owner` role |
| `IsDemoUser` / `ApiDemoUser` | global `demo` role (block write) |

### 2. FormRequest `authorize()` via `ChecksLogin` trait

```php
// app/Support/Request/ChecksLogin.php
public function authorize(): bool
{
    $user = auth()->user();
    $userGroup = $this->getUserGroup(); // from user_group_id param or user default
    foreach ($this->acceptedRoles as $role) {
        if ($user->hasRoleInGroupOrOwner($userGroup, $role)) return true;
    }
    return false;
}
```

### 3. API Controller `ValidatesUserGroupTrait`

```php
// app/Support/Http/Api/ValidatesUserGroupTrait.php
protected function validateUserGroup(Request $request): void
{
    // 1. Check authenticated
    // 2. Check membership in group
    // 3. Check has acceptedRoles in group
}
```

### 4. Route Model Binding with Authorization

```php
// app/Models/UserGroup.php
public static function routeBinder($value): UserGroup
{
    $group = self::find((int)$value);
    // Check user has READ_ONLY+ role or global owner
    if (!$user->hasRoleInGroupOrOwner($group, UserRoleEnum::READ_ONLY)) {
        throw new NotFoundHttpException; // silently hides unauthorized resources
    }
    return $group;
}
```

### 5. Repository `checkUserGroupAccess()`

```php
// app/Support/Repositories/UserGroup/UserGroupTrait.php
protected function checkUserGroupAccess(UserRoleEnum $role): void
{
    if (!$this->user->hasRoleInGroupOrOwner($this->userGroup, $role)) {
        throw new AuthorizationException;
    }
}
```

---

## User Groups ("Administrations")

### Konsep

UserGroup adalah **shared workspace**. Semua financial data (accounts, budgets, bills, tags, transactions, dll) **milik UserGroup**, bukan langsung milik User.

### Database Schema

```
user_groups: id, title, timestamps, soft_deletes
user_roles:  id, title, timestamps, soft_deletes
group_memberships: id, user_id, user_group_id, user_role_id (unique composite)
```

### User Group Switching

**File**: `app/Repositories/UserGroup/UserGroupRepository.php`

```php
public function useUserGroup(UserGroup $userGroup): void
{
    $this->user->user_group_id = $userGroup->id;
    $this->user->save(); // persisted to DB, NOT session
}
```

> Switching disimpan di database (kolom `users.user_group_id`), bukan di session. Tidak ada session-based switching.

### Administration UI Routes

```
GET  /administrations          → list all groups
GET  /administrations/create   → create group
GET  /administrations/edit/{id} → edit group
```

### API Group Switching

```
GET /api/v1/accounts?user_group_id=2  → switch context to group 2
GET /api/v1/accounts                   → use user's default group
```

---

## Data Ownership Flow

```
1. User login → user.user_group_id = default group
2. API request → Binder middleware resolves UserGroup from param or default
3. Controller → validatesUserGroup() checks membership + role
4. Repository → setUser() scopes all queries to user + group
5. Model binders → routeBinder() checks ownership (returns 404 if no access)
```

**Semua query di-scope ke level UserGroup**. Account, Transaction, Budget, dll semuanya punya `user_group_id`. User hanya bisa akses data di group yang mereka menjadi member.

# OAuth Token Compatibility — Laravel Passport ke Golang

---

> Analisis lengkap struktur token Firefly III (Laravel Passport) dan strategi kompatibilitas untuk implementasi di Go.

## 1. Token Format: JWT (Bukan Opaque)

Sebelumnya diduga opaque token, tapi setelah analisis mendalam:

**Token adalah JWT** yang di-sign dengan RSA 4096-bit, dan **juga** disimpan di database untuk revocation capability.

### Bukti:

1. **RSA key pair** di-generate saat setup (`CreateStuff::createOAuthKeys()`):
   ```php
   $key = RSA::createKey(4096);
   file_put_contents('storage/oauth-public.key', $key->getPublicKey());
   file_put_contents('storage/oauth-private.key', $key->toString('PKCS1'));
   ```

2. **Key management** via `OAuthKeys` class — generate, encrypt ke DB, restore, verify.

3. **Laravel Passport behavior**: Saat RSA key pair ada, Passport otomatis issue JWT access tokens.

4. **`CreateFreshApiToken` middleware** — membuat encrypted JWT cookie untuk SPA auth.

### Mode: "Database-backed JWT"

Token adalah JWT **dan** disimpan server-side di `oauth_access_tokens` table. Ini memberikan:
- Signature verification tanpa DB hit (JWT)
- Revocation capability (DB `revoked` flag)
- Best of both worlds

---

## 2. Database Schema — OAuth Tables

### oauth_access_tokens

| Column | Type | Nullable | Notes |
|--------|------|----------|-------|
| `id` | string(100) | NO | **PK** — JWT token string |
| `user_id` | integer | YES | **Index** |
| `client_id` | integer | NO | FK ke oauth_clients |
| `name` | string | YES | Human-readable token name |
| `scopes` | text | YES | JSON-encoded scope list |
| `revoked` | boolean | NO | |
| `created_at` | timestamp | NO | |
| `updated_at` | timestamp | NO | |
| `expires_at` | dateTime | YES | Token expiry |

### oauth_refresh_tokens

| Column | Type | Nullable | Notes |
|--------|------|----------|-------|
| `id` | string(100) | NO | **PK** |
| `access_token_id` | string(100) | NO | **Index** FK ke oauth_access_tokens |
| `revoked` | boolean | NO | |
| `expires_at` | dateTime | YES | |

### oauth_clients

| Column | Type | Nullable | Notes |
|--------|------|----------|-------|
| `id` | increments | NO | **PK** |
| `user_id` | integer | YES | **Index** |
| `name` | string | NO | |
| `secret` | string(100) | NO | Plaintext client secret |
| `redirect` | text | NO | Redirect URIs |
| `personal_access_client` | boolean | NO | PAT flag |
| `password_client` | boolean | NO | Password grant flag |
| `revoked` | boolean | NO | |
| `created_at` | timestamp | NO | |
| `updated_at` | timestamp | NO | |

### oauth_personal_access_clients

| Column | Type | Nullable | Notes |
|--------|------|----------|-------|
| `id` | increments | NO | **PK** |
| `client_id` | integer | NO | **Index** FK ke oauth_clients |
| `created_at` | timestamp | NO | |
| `updated_at` | timestamp | NO | |

### oauth_auth_codes

| Column | Type | Nullable | Notes |
|--------|------|----------|-------|
| `id` | string(100) | NO | **PK** |
| `user_id` | integer | NO | |
| `client_id` | integer | NO | |
| `scopes` | text | YES | |
| `revoked` | boolean | NO | |
| `expires_at` | dateTime | YES | |

### personal_access_tokens (Firefly III-specific, non-Passport)

| Column | Type | Nullable | Notes |
|--------|------|----------|-------|
| `id` | bigIncrements | NO | **PK** |
| `tokenable_type` | string | NO | Polymorphic type |
| `tokenable_id` | bigint | NO | User ID |
| `name` | string | NO | |
| `token` | string(64) | NO | **Unique** — hashed prefix |
| `abilities` | text | YES | |
| `last_used_at` | timestamp | YES | |
| `created_at` | timestamp | NO | |
| `updated_at` | timestamp | NO | |

---

## 3. Token Expiry

| Token Type | Expiry | Config |
|------------|--------|--------|
| **Access token** | 14 hari | `Passport::tokensExpireIn(now()->addDays(14))` |
| **Refresh token** | 30 hari (default Passport) | Tidak explicit di code |
| **Personal access token** | **Never expire** | Tidak explicit di code |

---

## 4. Guard Configuration

```php
// config/auth.php
'guards' => [
    'web' => ['driver' => 'session', 'provider' => 'users', 'remember' => 364*24*60],
    'remote_user_guard' => ['driver' => 'remote_user_guard', 'provider' => 'remote_user_provider'],
    'api' => ['driver' => 'passport', 'provider' => 'users'],
]
```

Semua API route (`/api/*`) menggunakan middleware `auth:api` (Passport driver).

---

## 5. Grant Types

| Grant | Endpoint | Keterangan |
|-------|----------|------------|
| Password | `POST /oauth/token` | `grant_type=password` |
| Client Credentials | `POST /oauth/token` | `grant_type=client_credentials` |
| Authorization Code | `GET /oauth/authorize` → `POST /oauth/token` | |
| Refresh Token | `POST /oauth/token` | `grant_type=refresh_token` |
| Personal Access | Via authenticated UI | Never expire |

### Password Grant Request

```
POST /oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=password&client_id=1&client_secret=xxx
&username=user@example.com&password=xxx
```

### Response

```json
{
    "token_type": "Bearer",
    "expires_in": 31536000,
    "access_token": "eyJ0eXAiOiJKV1...",
    "refresh_token": "def502..."
}
```

---

## 6. Scopes

**Tidak ada scopes yang defined di production code.**

`Passport::tokensCan()` tidak dipanggil di `AuthServiceProvider::boot()`. Route `GET /oauth/scopes` akan return empty list.

> Proposal scopes ada di `05-refactor-notes-rbac-wallet.md` tapi belum diimplementasi.

---

## 7. Token Revocation

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/oauth/tokens/{token_id}` | DELETE | Revoke access token |
| `/oauth/personal-access-tokens/{token_id}` | DELETE | Revoke PAT |

Semua token management routes dilindungi `auth:web`.

Validasi token: League OAuth2 Server cek `revoked` flag dan `expires_at` pada setiap request.

---

## 8. Client Auto-Configuration

PAT client auto-created saat user pertama kali buka profile page (`ProfileController::index()`):

```php
if (null === PersonalAccessClient::where('user_id', null)->first()) {
    // create PAT client
}
```

Environment variables:
- `PASSPORT_PERSONAL_ACCESS_CLIENT_ID`
- `PASSPORT_PERSONAL_ACCESS_CLIENT_SECRET`

---

## 9. Go Compatibility Strategy

### Opsi A: Query DB Langsung (Recommended untuk Strangler Fig)

Karena token juga disimpan di DB (`oauth_access_tokens.id` = token string):

```go
// 1. Extract token dari Authorization header
token := extractBearerToken(r.Header.Get("Authorization"))

// 2. Query DB
var accessToken OAuthAccessToken
err := db.Get(&accessToken,
    "SELECT * FROM oauth_access_tokens WHERE id = ? AND revoked = 0 AND expires_at > NOW()",
    token)

// 3. Get user
var user User
err = db.Get(&user, "SELECT * FROM users WHERE id = ?", accessToken.UserID)

// 4. Verify JWT signature (opsional, untuk keamanan tambahan)
// publicPem, _ := os.ReadFile("storage/oauth-public.key")
// publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicPem)
// token.Parse(publicKey)
```

**Keuntungan:**
- Shared DB dengan Laravel — token yang dikeluarkan Laravel bisa di-verify Go
- Tidak perlu implementasi full OAuth2 server
- Revocation langsung dari DB

### Opsi B: Verify JWT Signature Only

```go
publicKey := loadRSAPublicKey("storage/oauth-public.key")
token, err := jwt.Parse(bearerToken, func(t *jwt.Token) (interface{}, error) {
    return publicKey, nil
})
```

**Kelemahan:** Tidak bisa cek revocation tanpa DB query juga.

### Opsi C: Full OAuth2 Server di Go (Recommended untuk Production)

Implementasi full OAuth2 server menggunakan library seperti:
- `golang.org/x/oauth2`
- `github.com/go-oauth2/oauth2/v4`
- `github.com/ory/fosite`

**Strategi migrasi:**
1. Phase 1: Query DB langsung (compatibility)
2. Phase 2: Implementasi OAuth2 server di Go
3. Phase 3: Matikan Laravel Passport, Go handle semua auth

---

## 10. Key Takeaways

| Aspek | Detail |
|-------|--------|
| **Token format** | JWT (RS256, RSA 4096-bit) + DB storage |
| **Access token expiry** | 14 hari |
| **Refresh token expiry** | 30 hari (default) |
| **PAT expiry** | Never |
| **Scopes** | Tidak ada yang defined |
| **Grant types** | Password, client_credentials, authorization_code, refresh_token |
| **Key storage** | File (`storage/oauth-*.key`) + encrypted DB backup |
| **API auth** | `auth:api` (Passport) di semua `/api/*` routes |
| **User model trait** | `HasApiTokens` |
| **Client secrets** | Plaintext di DB |
| **Strangler Fig approach** | Query `oauth_access_tokens` table langsung dari Go |

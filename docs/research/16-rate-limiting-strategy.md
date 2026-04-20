# Rate Limiting Strategy — Go API

---

> Analisis rate limiting Firefly III (existing) dan rekomendasi untuk Go API baru.
> Firefly III **hampir tidak punya rate limiting** — ini vulnerability yang harus diperbaiki.

## 1. Existing Firefly III Rate Limiting

### Yang Ada (sangat Terbatas)

| Endpoint | Limit | Mechanism |
|----------|-------|-----------|
| `POST /oauth/token` | 60/min | Laravel default throttle middleware |
| `POST /login` (web) | 5/min, lockout 1min | `ThrottlesLogins` trait |
| `POST /password/email` | 1 per 5 min per email | Config throttle |

### Yang TIDAK Ada (Vulnerabilities)

| Endpoint | Vulnerability |
|----------|--------------|
| `POST /register` | Bot account creation — unlimited |
| `POST /oauth/token` (password grant) | 60/min terlalu longgar untuk brute force |
| Semua `/api/v1/*` | Zero rate limiting — DoS, scraping, abuse |
| `/api/v1/search/*` | Expensive DB queries — unlimited |
| `POST /export/*` | CSV generation DoS — unlimited |
| `POST /webhooks/*/submit` | Webhook message flooding — unlimited |
| `GET /api/v1/cron/{token}` | Cron abuse — unlimited |

---

## 2. Recommended Go Library

**Primary:** `github.com/ulule/limiter` + Redis

```go
import limiter "github.com/ulule/limiter/v3"
import "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
import "github.com/ulule/limiter/v3/drivers/store/redis"
```

**Fitur:**
- Sliding window, fixed window, token bucket strategies
- Redis backend (distributed)
- Multiple rate limit per route
- Easy Fiber middleware integration

---

## 3. HTTP Response Specification

### Headers (semua response)

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 97
X-RateLimit-Reset: 1718234567
```

### 429 Response

```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests. Please retry after 42 seconds.",
    "details": {
      "limit": 100,
      "remaining": 0,
      "reset_at": "2024-06-13T12:30:00Z"
    }
  }
}
```

```
HTTP/1.1 429 Too Many Requests
Retry-After: 42
```

---

## 4. Rate Limits by Category

### Authentication

| Endpoint | Limit | Window | Strategy | Scope |
|----------|-------|--------|----------|-------|
| `POST /auth/login` | **5** | 1 min | Fixed | IP + email |
| `POST /auth/login` | **20** | 1 hour | Sliding | IP |
| `POST /auth/login` | **10** | 15 min | Fixed | Email |
| `POST /auth/token` (password) | **10** | 1 min | Fixed | IP + client |
| `POST /auth/token` (client_credentials) | **60** | 1 min | Token bucket | Client ID |
| `POST /auth/token/refresh` | **30** | 1 min | Token bucket | Client ID |
| `POST /auth/register` | **3** | 1 hour | Fixed | IP |
| `POST /auth/register` | **5** | 24 hours | Fixed | IP |
| `POST /auth/password/email` | **3** | 1 hour | Fixed | IP |
| `POST /auth/password/email` | **1** | 5 min | Fixed | Email |
| `POST /auth/password/reset` | **5** | 1 hour | Fixed | IP |
| `GET /auth/2fa/verify` | **10** | 1 min | Fixed | Session |

### Progressive Login Escalation

| Failed Attempts | Action |
|----------------|--------|
| 3 per email | 5-second response delay |
| 5 per email | Lock 1 minute |
| 10 per email | Lock 15 minutes |
| 20 per email | Lock 1 hour + notify user |
| 50 per IP | Block IP 24 hours |

### General API CRUD

| Pattern | Limit | Window | Strategy | Scope |
|---------|-------|--------|----------|-------|
| GET (list) | **120** | 1 min | Token bucket | API key |
| GET (single) | **120** | 1 min | Token bucket | API key |
| POST (create) | **60** | 1 min | Token bucket | API key |
| PUT/PATCH (update) | **60** | 1 min | Token bucket | API key |
| DELETE | **30** | 1 min | Token bucket | API key |
| Bulk operations | **10** | 1 min | Token bucket | API key |
| Mass operations | **5** | 1 min | Token bucket | API key |

### Search (Expensive Queries)

| Endpoint | Limit | Window | Strategy | Scope |
|----------|-------|--------|----------|-------|
| `POST /search/transactions` | **20** | 1 min | Sliding | API key |
| `POST /search/transactions/count` | **30** | 1 min | Sliding | API key |
| `GET /search/accounts` | **30** | 1 min | Sliding | API key |
| `GET /autocomplete/*` | **60** | 1 min | Token bucket | API key |

### Reports & Charts (CPU Intensive)

| Endpoint | Limit | Window | Strategy | Scope |
|----------|-------|--------|----------|-------|
| `GET /charts/*` | **20** | 1 min | Sliding | API key |
| `GET /reports/*` | **10** | 1 min | Sliding | API key |
| `GET /summary/*` | **30** | 1 min | Sliding | API key |

### Webhooks

| Endpoint | Limit | Window | Strategy | Scope |
|----------|-------|--------|----------|-------|
| `POST /webhooks` (create) | **10** | 1 min | Token bucket | API key |
| `POST /webhooks/{id}/submit` | **30** | 1 min | Token bucket | API key |
| `GET /webhooks/{id}/messages` | **60** | 1 min | Token bucket | API key |
| Other CRUD | **60** | 1 min | Token bucket | API key |

### File Upload

| Endpoint | Limit | Window | Strategy | Scope |
|----------|-------|--------|----------|-------|
| `POST /attachments` (upload) | **10** | 1 min | Token bucket | API key |
| `POST /attachments` (upload) | **50** | 1 hour | Sliding | API key |
| `GET /attachments/*/download` | **60** | 1 min | Token bucket | API key |

### Export

| Endpoint | Limit | Window | Strategy | Scope |
|----------|-------|--------|----------|-------|
| `POST /export` | **3** | 1 min | Fixed | API key |
| `POST /export` | **10** | 1 hour | Sliding | API key |
| `GET /export/*/download` | **10** | 1 min | Token bucket | API key |

### System

| Endpoint | Limit | Window | Strategy | Scope |
|----------|-------|--------|----------|-------|
| `GET /health` | **30** | 1 min | Fixed | IP |
| `GET /cron/{token}` | **1** | 1 min | Fixed | Token |
| Admin endpoints | **60** | 1 min | Token bucket | User |

### Global IP Rate Limit (Layer 1)

| Limit | Window | Strategy | Scope |
|-------|--------|----------|-------|
| **1000** | 1 min | Token bucket | IP |
| **10000** | 1 hour | Sliding | IP |

---

## 5. Strategy Selection Guide

| Strategy | Best For | Accuracy | Memory |
|----------|----------|----------|--------|
| **Fixed Window** | Auth, exports | Low (boundary burst) | Low |
| **Sliding Window** | Search, reports, hourly caps | High | Medium |
| **Token Bucket** | General CRUD, API calls | Medium | Low |

---

## 6. Rate Limit Key Design (Redis)

```
rl:ip:{ip_address}:{category}           # Per IP (unauthenticated)
rl:apikey:{api_key_id}:{category}       # Per API key (authenticated)
rl:user:{user_id}:{category}             # Per user (authenticated)
rl:login:ip:{ip}:email:{email_hash}      # Per IP + email (login)
rl:oauth:client:{client_id}              # Per OAuth client
```

---

## 7. Middleware Stack Order

```
1. Panic recovery
2. Request ID / correlation ID
3. CORS
4. Global IP rate limit (Layer 1)
5. Authentication
6. Per-endpoint-category rate limit (Layer 2)
7. Per-user/per-API-key rate limit (Layer 3)
8. Route handler
```

---

## 8. Configuration (Environment Variables)

```env
RATE_LIMIT_AUTH_LOGIN_PER_MIN=5
RATE_LIMIT_AUTH_REGISTER_PER_HOUR=3
RATE_LIMIT_CRUD_READ_PER_MIN=120
RATE_LIMIT_CRUD_WRITE_PER_MIN=60
RATE_LIMIT_SEARCH_PER_MIN=20
RATE_LIMIT_EXPORT_PER_MIN=3
RATE_LIMIT_WEBHOOK_SUBMIT_PER_MIN=30
RATE_LIMIT_FILE_UPLOAD_PER_MIN=10
RATE_LIMIT_GLOBAL_PER_MIN=1000
RATE_LIMIT_STRATEGY=sliding_window
RATE_LIMIT_BACKEND=redis
RATE_LIMIT_REDIS_URL=redis://localhost:6379/2
```

---

## 9. Go Implementation Example

```go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/ulule/limiter/v3"
    "github.com/ulule/limiter/v3/drivers/middleware/fiber"
    "github.com/ulule/limiter/v3/drivers/store/redis"
)

func RateLimit(store *redis.Store, rate limiter.Rate, keyPrefix string) fiber.Handler {
    return fiber.New(func(c *fiber.Ctx) error {
        // Extract scope: IP, API key, or user ID
        scope := extractScope(c)

        limiterInstance := limiter.New(store, rate).
            WithKeyPrefix(keyPrefix).
            WithKey(scope)

        return limiterInstance(c)
    })
}

// Usage in route setup
api.Post("/transactions", middleware.RateLimit(
    redisStore,
    limiter.Rate{Period: time.Minute, Limit: 60},
    "rl:crud:write",
))
```

---

## 10. Monitoring

Setiap rate limit check emit metrics:
- `rate_limit_check_total{endpoint, scope, result="allowed|denied"}`
- `rate_limit_remaining{endpoint, scope}` (gauge)
- `rate_limit_denied_total{endpoint, scope}` (counter)

Memungkinkan:
- Deteksi legitimate user yang ke-block
- Identifikasi attack patterns (spike denied)
- Monitor endpoint-specific load

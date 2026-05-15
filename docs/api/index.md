# API Reference

Gofin provides a RESTful JSON API. The full OpenAPI 3.0 specification is available below.

## Interactive Documentation

<iframe src="https://redocly.github.io/redoc/?url=https://raw.githubusercontent.com/ajianaz/gofin-full/main/docs/openapi.yaml" frameborder="0" width="100%" height="800px" loading="lazy" style="border: 1px solid var(--vp-c-divider); border-radius: 8px;"></iframe>

## Base URL

| Environment | URL |
|------------|-----|
| Production | `https://your-domain/api/v1` |
| Development | `http://localhost:8080/api/v1` |

## Authentication

All protected endpoints require a JWT bearer token:

```http
Authorization: Bearer <access_token>
```

### Obtain a Token

```bash
# Login
curl -X POST https://your-domain/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "your-password"}'

# Response
{
  "access_token": "eyJhbG...",
  "refresh_token": "eyJhbG...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

### Refresh a Token

```bash
curl -X POST https://your-domain/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "eyJhbG..."}'
```

## Rate Limiting

API requests are rate-limited. Default: 100 requests per 60 seconds.

Rate limit headers are included in every response:

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1715731200
```

## Error Format

All errors follow a consistent format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Human-readable error description",
    "details": {}
  }
}
```

### HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 201 | Created |
| 204 | No Content (deleted successfully) |
| 400 | Bad Request (validation error) |
| 401 | Unauthorized (missing or invalid token) |
| 403 | Forbidden (insufficient permissions) |
| 404 | Not Found |
| 409 | Conflict (duplicate resource) |
| 422 | Unprocessable Entity (business logic error) |
| 429 | Too Many Requests (rate limited) |
| 500 | Internal Server Error |

## API Resources

### Core Resources

| Resource | Endpoints | Description |
|----------|-----------|-------------|
| **Auth** | `/api/v1/auth/*` | Login, register, refresh, OAuth |
| **Users** | `/api/v1/users/*` | User management (admin) |
| **Wallets** | `/api/v1/wallets/*` | Financial accounts |
| **Transactions** | `/api/v1/transactions/*` | Double-entry transactions |
| **Categories** | `/api/v1/categories/*` | Transaction categories |
| **Tags** | `/api/v1/tags/*` | Transaction tags |
| **Budgets** | `/api/v1/budgets/*` | Spending budgets |
| **Piggy Banks** | `/api/v1/piggy-banks/*` | Savings goals |
| **Bills** | `/api/v1/bills/*` | Recurring bills |
| **Recurring** | `/api/v1/recurring/*` | Recurring transactions |

### Advanced Resources

| Resource | Endpoints | Description |
|----------|-----------|-------------|
| **Analytics** | `/api/v1/analytics/*` | Spending analysis |
| **Reports** | `/api/v1/reports/*` | Financial reports |
| **Export** | `/api/v1/export/*` | CSV/OFX download |
| **Rules** | `/api/v1/rules/*` | Automation rules |
| **Webhooks** | `/api/v1/webhooks/*` | Webhook management |
| **Groups** | `/api/v1/groups/*` | User groups |
| **Currencies** | `/api/v1/currencies/*` | Currency management |
| **Audit Logs** | `/api/v1/audit-logs/*` | Activity audit trail |
| **Notifications** | `/api/v1/notifications/*` | User notifications |
| **Admin** | `/api/v1/admin/*` | System administration |

## Common Patterns

### Pagination

List endpoints support pagination via query parameters:

```
GET /api/v1/transactions?page=1&per_page=20
```

| Parameter | Default | Description |
|-----------|---------|-------------|
| `page` | 1 | Page number |
| `per_page` | 20 | Items per page (max: 100) |

### Sorting

```
GET /api/v1/transactions?sort=created_at&order=desc
```

| Parameter | Default | Description |
|-----------|---------|-------------|
| `sort` | `created_at` | Sort field |
| `order` | `desc` | Sort order: `asc` or `desc` |

### Filtering

Most list endpoints support filtering:

```
GET /api/v1/transactions?wallet_id=xxx&category_id=yyy&start_date=2026-01-01&end_date=2026-01-31
```

### Double-Entry Transactions

Transactions always have a source and destination wallet:

```json
{
  "description": "Grocery shopping",
  "amount": 50000,
  "source_wallet_id": "wallet-uuid-1",
  "destination_wallet_id": "wallet-uuid-2",
  "category_id": "category-uuid",
  "tags": ["food", "weekly"],
  "notes": "Weekly grocery run"
}
```

## Download OpenAPI Spec

The full specification is available in the repository:

- [openapi.yaml](https://raw.githubusercontent.com/ajianaz/gofin-full/main/docs/openapi.yaml)

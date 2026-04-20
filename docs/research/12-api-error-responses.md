# API Error Responses — Status Codes & Messages

---

> Lengkap error response codes dan messages per endpoint. Digunakan untuk Go error handling middleware.

## JSON Response Body Shapes

### Shape 1 — Standard Error (umum)

```go
type ErrorResponse struct {
    Message   string `json:"message"`
    Exception string `json:"exception"`
}
```
Digunakan: 400, 401, 404, 405, 406, 410, 415, 500

### Shape 2 — Validation Error

```go
type ValidationErrorResponse struct {
    Message string              `json:"message"`
    Errors  map[string][]string `json:"errors"`
}
```
Digunakan: 422 (Laravel validation — format paling umum)

### Shape 3 — Rate Limit Error

```go
type RateLimitErrorResponse struct {
    Error struct {
        Code    string `json:"code"`
        Message string `json:"message"`
        Details struct {
            Limit     int    `json:"limit"`
            Remaining int    `json:"remaining"`
            ResetAt   string `json:"reset_at"`
        } `json:"details"`
    } `json:"error"`
}
```
Digunakan: 429 (Go API baru, tidak ada di Firefly III)

---

## HTTP Status Code Reference

| Code | When |
|------|------|
| 200 | Successful GET, PUT, DELETE |
| 201 | Successful POST (create) |
| 204 | Successful delete/destroy (empty body) |
| 400 | Bad request (invalid headers) |
| 401 | Unauthenticated / unauthorized |
| 403 | Demo user write attempt |
| 404 | Resource not found |
| 405 | Method not allowed |
| 406 | Bad Accept header |
| 409 | Conflict (currency in use) |
| 410 | Transaction deleted by rule after creation |
| 415 | Bad/missing Content-Type |
| 422 | Validation failure |
| 429 | Rate limit exceeded (Go API baru) |
| 500 | Internal error |

---

## 400 Bad Request

| Trigger | Message | Source |
|---------|---------|--------|
| Invalid `X-Trace-Id` header | `Bad X-Trace-Id header.` | AcceptHeaders middleware |
| General bad request | `<exception message>` | HttpException |

## 401 Unauthorized

| Trigger | Message | Source |
|---------|---------|--------|
| Missing/expired OAuth token | `The user is not logged in but must be.` | Authenticate middleware |
| Blocked account (`blocked=1`) | `Blocked account.` | Authenticate middleware |
| Email changed block | `Blocked account.` | Authenticate middleware |
| OAuth server error | `<OAuth error message>` | OAuthServerException |
| Authorization failure | `<message>` | AuthorizationException |
| Admin API guest user | `"Unauthorized."` (plain text) | IsAdminApi middleware |

## 403 Forbidden

| Trigger | Message | Body | Source |
|---------|---------|------|--------|
| Demo user write attempt | — | Empty `""` | ApiDemoUser middleware |

## 404 Not Found

| Trigger | Message | Source |
|---------|---------|--------|
| Any model not found (route binding) | `Resource not found` | Global handler |
| Webhooks disabled | `Resource not found` | Webhook controllers |
| Exchange rate not found | `Resource not found` | CurrencyExchangeRate\ShowController |
| Attachment not found | `Resource not found` | Attachment controllers |
| Transaction group not found | `Resource not found` | Transaction\UpdateController |
| Transaction journal not found | `Resource not found` | Transaction\ShowController |

## 405 Method Not Allowed

| Trigger | Message | Source |
|---------|---------|--------|
| Wrong HTTP method | Debug: full exception. Production: `Internal Firefly III Exception: <msg>` | Global handler |

## 406 Not Acceptable

| Trigger | Message | Source |
|---------|---------|--------|
| Bad Accept header (global) | `Accept header "<value>" is not something this server can provide.` | AcceptHeaders middleware |
| Bad Accept header (controller) | `Sorry, Accept header "<value>" is not something this endpoint can provide.` | Controller base |

**Valid Accept headers:** `application/json`, `application/vnd.api+json`, `application/x-www-form-urlencoded`, `application/octet-stream`, `*/*`

## 409 Conflict

| Trigger | Message | Body | Source |
|---------|---------|------|--------|
| Disable currency in use | — | Empty `[]` | Currency\UpdateController |
| Disable only remaining currency | — | Empty `[]` | Currency\UpdateController |

## 410 Gone

| Trigger | Message | Source |
|---------|---------|--------|
| Transaction deleted by rule after creation | `200032: Cannot find transaction. Possibly, a rule deleted this transaction after its creation.` | Transaction\StoreController |

## 415 Unsupported Media Type

| Trigger | Message | Source |
|---------|---------|--------|
| Missing Content-Type (POST/PUT) | `Content-Type header cannot be empty.` | AcceptHeaders middleware |
| Invalid Content-Type (POST/PUT) | `Content-Type cannot be "<value>"` | AcceptHeaders middleware |

**Valid Content-Type:** `application/json`, `application/vnd.api+json`, `application/x-www-form-urlencoded`, `application/octet-stream`

> Exempt: bulk transactions, attachment upload

## 422 Unprocessable Entity

### Format A — Laravel Validation (umum)

```json
{
  "message": "The given data was invalid.",
  "errors": {
    "field_name": ["Error message 1"],
    "another_field": ["Error message"]
  }
}
```

### Format B — Custom ValidationException

```json
{
  "message": "Validation exception: <exception message>",
  "errors": {"field": "Field is invalid"}
}
```

### Format C — Empty Body

```json
[]
```
Trigger: attachment upload failure

### Format D — Duplicate Transaction

Format A dengan error message: `"Duplicate of transaction #<group_id>."` di field `transactions.0.description`

## 500 Internal Server Error

| Trigger | Message | Body | Source |
|---------|---------|------|--------|
| Delete self (admin) | — | Empty `[]` | UserController |
| Unhandled exception (production) | `Internal Firefly III Exception: <msg>` | Standard error shape | Global handler |
| Unhandled exception (debug) | Full exception + trace | Debug error shape | Global handler |

## 204 No Content (Success)

| Endpoint | When |
|----------|------|
| `DELETE /data/destroy` | Successful data destruction |
| `DELETE /data/purge` | Successful data purge |
| `POST /data/bulk/transactions` | Successful bulk update |
| `POST /attachments/{id}/upload` | Successful upload |
| `DELETE /users/{id}` | Successful user deletion |

---

## Go Error Handling Pattern

```go
func ErrorHandler(ctx *fiber.Ctx, err error) error {
    code := fiber.StatusInternalServerError

    var e *fiber.Error
    if errors.As(err, &e) {
        code = e.Code
    }

    var ve *ValidationError
    if errors.As(err, &ve) {
        return ctx.Status(422).JSON(ValidationErrorResponse{
            Message: "The given data was invalid.",
            Errors:  ve.FieldErrors,
        })
    }

    return ctx.Status(code).JSON(ErrorResponse{
        Message:   http.StatusText(code),
        Exception: reflect.TypeOf(err).String(),
    })
}
```

---

## Demo User Protection

Demo user (`role=demo`) dilindungi oleh middleware `ApiDemoUser`:
- **Store/Update/Delete**: return 403 dengan empty body
- **Attachment endpoints**: return 404 (disembunyikan)
- Berlaku di: transactions, accounts, budgets, bills, categories, tags, piggy banks, rules, rule groups, recurrences, webhooks, available budgets, object groups

## Admin-Only Endpoints

Membutuhkan middleware `api-admin` (`IsAdminApi`):
- `DELETE /currencies/{code}`
- `POST /currencies`
- `PUT/DELETE /link-types/*`
- `PUT /configuration/{key}`
- Semua `/users/*` CRUD

## Accept Header Enforcement

Dua layer:
1. **Global middleware** (`AcceptHeaders`): cek terhadap list yang luas
2. **Controller level**: cek terhadap `$accepts` array controller (default: `application/json`, `application/vnd.api+json`)

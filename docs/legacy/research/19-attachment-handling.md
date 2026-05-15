# Attachment Handling — Complete Specification

---

> File upload, storage, MIME validation, API endpoints, cascade behavior.

## 1. Attachment Model

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | bigint PK | Auto-increment |
| `attachable_id` | bigint | Parent model ID |
| `attachable_type` | string | Parent model class name |
| `user_id` | bigint FK | Owner user |
| `user_group_id` | bigint FK | User group |
| `md5` | string | MD5 hash of file content |
| `filename` | string | Original filename (max 255) |
| `mime` | string | MIME type |
| `title` | string | Display title (max 255) |
| `description` | string | Description text (max 32768) |
| `size` | int | File size in bytes |
| `uploaded` | boolean | Whether file content is uploaded |
| `created_at` | timestamp | |
| `updated_at` | timestamp | |
| `deleted_at` | timestamp | Soft delete |

### Polymorphic Types (valid)

`Account`, `Bill`, `Budget`, `Category`, `PiggyBank`, `Tag`, `TransactionJournal`, `Recurrence`

> `Transaction` type auto-remapped ke `TransactionJournal` saat create.

## 2. File Storage

| Aspect | Value |
|--------|-------|
| Disk name | `upload` |
| Root path | `storage/upload/` |
| File name pattern | `at-{id}.data` (selalu `.data` extension) |
| Example | `storage/upload/at-42.data` |
| Encryption | Tidak dienkripsi saat simpan, toleransi decrypt error saat baca |
| Soft delete | Ya — `deleted_at` column |

## 3. MIME Type Whitelist

### Images
`image/jpeg`, `image/svg+xml`, `image/png`, `image/heic`, `image/heic-sequence`, `image/webp`, `image/gif`, `image/tiff`, `image/bmp`, `image/x-icon`, `image/vnd.microsoft.icon`

### Documents
`application/pdf`, `text/plain`, `text/html`, `text/xml`, `application/xml`, `application/json`, `message/rfc822`

### Office
`application/msword`, `application/vnd.openxmlformats-officedocument.wordprocessingml.*`, `application/vnd.ms-excel`, `application/vnd.openxmlformats-officedocument.spreadsheetml.*`, `application/vnd.ms-powerpoint`, `application/vnd.openxmlformats-officedocument.presentationml.*`

### OpenOffice / ODF
`application/vnd.oasis.opendocument.*`, `application/vnd.sun.xml.*`, `application/x-iwork-pages-sffpages`

### Generic
`application/octet-stream` (accepts anything)

### Max Upload Size
**1 GB** (`1073741824` bytes)

## 4. API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/attachments` | List all user's attachments (paginated) |
| `POST` | `/api/v1/attachments` | Create attachment metadata (no file) |
| `GET` | `/api/v1/attachments/{id}` | Get attachment metadata |
| `GET` | `/api/v1/attachments/{id}/download` | Download file content |
| `POST` | `/api/v1/attachments/{id}/upload` | Upload file content (raw body) |
| `PUT` | `/api/v1/attachments/{id}` | Update metadata |
| `DELETE` | `/api/v1/attachments/{id}` | Delete attachment (soft) |

### Sub-resource Endpoints

| Parent | Path |
|--------|------|
| Account | `/api/v1/accounts/{id}/attachments` |
| Bill | `/api/v1/bills/{id}/attachments` |
| Budget | `/api/v1/budgets/{id}/attachments` |
| Category | `/api/v1/categories/{id}/attachments` |
| Piggy Bank | `/api/v1/piggy-banks/{id}/attachments` |
| Tag | `/api/v1/tags/{id}/attachments` |
| Transaction Group | `/api/v1/transaction-groups/{id}/attachments` |

## 5. Request Validation

### StoreRequest

| Field | Rules |
|-------|-------|
| `filename` | required, min:1, max:255 |
| `title` | min:1, max:255 |
| `notes` | min:1, max:32768 |
| `attachable_type` | required, in:Account,Bill,Budget,Category,PiggyBank,Tag,Transaction,TransactionJournal,Recurrence |
| `attachable_id` | required, numeric, IsValidAttachmentModel |

### UpdateRequest

Same as StoreRequest but `filename`, `title`, `attachable_type` are optional.

### Upload Request

- Content-Type: raw binary (not multipart)
- Body: file content
- MIME detection via `finfo_file()` (not from Content-Type header)
- Validation: MIME must be in whitelist, size must be ≤ 1GB

## 6. Response Format

```json
{
  "id": "42",
  "created_at": "2024-01-01T00:00:00+00:00",
  "updated_at": "2024-01-01T00:00:00+00:00",
  "attachable_id": "15",
  "attachable_type": "TransactionJournal",
  "hash": "d41d8cd98f00b204e9800998ecf8427e",
  "filename": "receipt.pdf",
  "download_url": "https://example.com/api/v1/attachments/42/download",
  "upload_url": "https://example.com/api/v1/attachments/42/upload",
  "title": "Store receipt",
  "notes": "Note text or null",
  "mime": "application/pdf",
  "size": 12345,
  "links": [{"rel": "self", "uri": "/attachment/42"}]
}
```

> `hash` = MD5 of file content. `attachable_type` prefix `FireflyIII\Models\` di-strip.

## 7. Upload Flow

### API Flow (2-step)

```
1. POST /api/v1/attachments
   Body: {"filename": "receipt.pdf", "attachable_type": "TransactionJournal", "attachable_id": 15}
   → Creates record: uploaded=false, size=0, md5=""
   → Returns attachment with upload_url

2. POST /api/v1/attachments/{id}/upload
   Body: <raw file bytes>
   → finfo_file() detect MIME
   → Validate MIME in whitelist
   → Write to storage/upload/at-{id}.data
   → Update record: md5, mime, size, uploaded=true
   → Returns 204 No Content
```

## 8. Download Flow

```
1. Validate: uploaded=true AND size > 0
2. Read file: storage/upload/at-{id}.data
3. Attempt Crypt::decrypt() — if fails, return raw content
4. Set Content-Disposition: attachment; filename="{original_filename}"
5. Set Content-Type: application/octet-stream
```

## 9. Cascade Delete Behavior

Saat parent model dihapus, semua attachments ikut terhapus (file + DB):

| Parent Model | Cascade Trigger |
|-------------|-----------------|
| Account | `DeletedAccountObserver` → destroy attachments |
| Category | `DeletedCategoryObserver` → destroy attachments |
| Tag | `DeletedTagObserver` → destroy attachments |
| PiggyBank | `PiggyBankObserver` → destroy attachments |
| Recurrence | `DeletedRecurrenceObserver` → destroy attachments |
| TransactionJournal | `DeletedTransactionJournalObserver` → destroy attachments |

File deletion path: `storage/upload/at-{id}.data`

## 10. Go Implementation

```go
func UploadAttachment(c *fiber.Ctx, db *sqlx.DB) error {
    attachmentID := c.Params("id")

    // Validate ownership
    var attachment Attachment
    if err := db.Get(&attachment,
        "SELECT * FROM attachments WHERE id = ? AND user_id = ? AND deleted_at IS NULL",
        attachmentID, c.Locals("userID")); err != nil {
        return c.Status(404).JSON(ErrorResponse{Message: "Resource not found"})
    }

    // Read body
    body := c.Body()

    // Detect MIME
    mime := detectMIME(body) // using mime.TypeByExtension + fallback to http.DetectContentType

    // Validate MIME
    if !isValidMIME(mime) {
        return c.Status(422).JSON(ValidationErrorResponse{
            Message: "The given data was invalid.",
            Errors: map[string][]string{"file": {"Invalid MIME type: " + mime}},
        })
    }

    // Validate size
    if len(body) > 1073741824 {
        return c.Status(422).JSON(ValidationErrorResponse{
            Message: "The given data was invalid.",
            Errors: map[string][]string{"file": {"File too large. Maximum size is 1GB."}},
        })
    }

    // Store file
    path := filepath.Join(uploadDir, fmt.Sprintf("at-%s.data", attachmentID))
    if err := os.WriteFile(path, body, 0644); err != nil {
        return c.Status(500).JSON(ErrorResponse{Message: "Failed to store file"})
    }

    // Update record
    md5 := md5.Sum(body)
    _, err := db.Exec(`
        UPDATE attachments SET md5 = ?, mime = ?, size = ?, uploaded = 1, updated_at = NOW()
        WHERE id = ?
    `, hex.EncodeToString(md5[:]), mime, len(body), attachmentID)

    return c.SendStatus(204)
}
```

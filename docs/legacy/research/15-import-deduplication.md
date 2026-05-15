# Import Deduplication — Hash Algorithm

---

> Algoritma deduplication transaction menggunakan SHA-256 hash. Digunakan saat import untuk mencegah duplikat.

## Algorithm Overview

**Exact-match SHA-256 hash** stored as journal metadata (`journal_meta` table).

## 1. Hash Computation

**File:** `TransactionJournalFactory.php`, method `hashArray()`

```
1. Ambil full transaction data array (NullArrayObject)
2. Remove field: import_hash_v2, original_source
3. JSON-encode seluruh array
4. Compute: hash('sha256', json_string)
5. Return: 64-character hex digest
```

### Pseudocode

```go
func hashTransactionData(row map[string]interface{}) string {
    delete(row, "import_hash_v2")
    delete(row, "original_source")

    jsonBytes, err := json.Marshal(row)
    if err != nil {
        // fallback: use timestamp (effectively disables dedup)
        jsonBytes = []byte(fmt.Sprintf("%f", float64(time.Now().UnixNano())))
    }

    hash := sha256.Sum256(jsonBytes)
    return hex.EncodeToString(hash[:])
}
```

## 2. Fields Included in Hash

Hash dihitung dari **seluruh field** per transaction journal:

| Category | Fields |
|----------|--------|
| **Type & ordering** | `type`, `order` |
| **Date** | `date` (serialized Carbon) |
| **Description** | `description` |
| **Currency** | `currency_id`, `currency_code` |
| **Foreign currency** | `foreign_currency_id`, `foreign_currency_code` |
| **Amounts** | `amount`, `foreign_amount` |
| **Location** | `latitude`, `longitude`, `zoom_level` |
| **Source account** | `source_id`, `source_name`, `source_iban`, `source_number`, `source_bic` |
| **Destination account** | `destination_id`, `destination_name`, `destination_iban`, `destination_number`, `destination_bic` |
| **Budget** | `budget_id`, `budget_name` |
| **Category** | `category_id`, `category_name` |
| **Bill** | `bill_id`, `bill_name` |
| **Piggy bank** | `piggy_bank_id`, `piggy_bank_name` |
| **Flags** | `reconciled` |
| **Notes** | `notes` |
| **Tags** | `tags` (array) |
| **Custom fields** | `internal_reference`, `external_id`, `recurrence_id`, `bunq_payment_id`, `external_url` |
| **SEPA fields** | `sepa_cc`, `sepa_ct_op`, `sepa_ct_id`, `sepa_db`, `sepa_country`, `sepa_ep`, `sepa_ci`, `sepa_batch_id` |
| **Custom dates** | `interest_date`, `book_date`, `process_date`, `due_date`, `payment_date`, `invoice_date` |

**Excluded:** `import_hash_v2` (hash sendiri), `original_source` (e.g. `"ff3-v6.5.9"`)

### Hash Sensitivity

Hash **exact match** — perbedaan apapun di field manapun menghasilkan hash berbeda:
- Order of keys dalam JSON
- Exact values (resolved IDs, names, IBANs)
- String representation of date objects (includes timezone)
- `null` vs empty string

## 3. Hash Storage

| Aspect | Detail |
|--------|--------|
| **Table** | `journal_meta` (EAV pattern) |
| **name** | `'import_hash_v2'` |
| **data** | JSON-encoded hash string: `"\"abc123...\""` |
| **Auto-stored** | Via `TransactionJournalMetaFactory::storeMetaFields()` |

## 4. Duplicate Detection Flow

```
1. TransactionJournalFactory::createJournal(row)
2. hash = hashArray(row) → row['import_hash_v2'] = hash
3. IF errorOnHash == true:
   a. Query journal_meta WHERE name='import_hash_v2' AND data=json_encode(hash)
   b. JOIN transaction_journals → scope to current user_id
   c. Include soft-deleted (withTrashed)
   d. IF found → throw DuplicateTransactionException
4. Create journal
5. Store import_hash_v2 in journal_meta
```

## 5. Control Parameter

| `error_if_duplicate_hash` | Behavior |
|---------------------------|----------|
| `false` (default) | Hash computed & stored, tapi TIDAK ada duplicate check |
| `true` | Hash computed, duplicate check dijalankan |

### Flow Through System

```
StoreRequest::getAll()
  → extract error_if_duplicate_hash
  → TransactionGroupRepository::store()
    → TransactionGroupFactory::create()
      → TransactionJournalFactory::setErrorOnHash(value)
      → createJournal() → hashArray() → errorIfDuplicate()
```

## 6. On Duplicate Found

| Step | Action |
|------|--------|
| 1 | `DuplicateTransactionException` thrown |
| 2 | Jika batch (split transaction) → semua journal 1..N-1 **force-deleted** |
| 3 | Exception caught di StoreController |
| 4 | Converted to 422 Validation Error |
| 5 | Message: `"Duplicate of transaction #<group_id>."` |
| 6 | Field: `transactions.0.description` |

### 422 Response Format

```json
{
  "message": "The given data was invalid.",
  "errors": {
    "transactions.0.description": ["Duplicate of transaction #42."]
  }
}
```

## 7. Atomicity

- Split transactions: jika duplicate di journal N, journal 1..N-1 di-**rollback** via `forceDeleteOnError()`
- Seluruh group creation bersifat **atomic**: semua journal tercipta atau tidak sama sekali

## 8. Separate: getCompareHash() (NOT Import Dedup)

Hash terpisah di `TransactionGroupRepository::getCompareHash()` untuk deteksi "meaningful changes" saat update:

```go
func getCompareHash(group TransactionGroup) string {
    sum := "0"
    names := ""
    for _, journal := range group.Journals {
        names += journal.Date.Format("2006-01-02-15:04:05")
        for _, tx := range journal.Transactions {
            if tx.Amount.IsNegative() {
                sum = sum.Add(tx.Amount)
                names += tx.Account.Name
            }
        }
    }
    return sha256(sum + "-" + names)
}
```

**Usage:** Cek apakah update transaction perlu re-run rules/webhooks (bandingkan old vs new hash).

**Only considers:** dates, negative-amount sums, source account names.

## 9. Legacy import_hash (v1)

| Aspect | Detail |
|--------|--------|
| **Name** | `import_hash` |
| **Status** | Legacy, tidak lagi di-compute di codebase saat ini |
| **Migration** | Dipertahankan saat upgrade ke group model |
| **Config** | Tetap listed di `journal_meta_fields` |

## 10. Go Implementation

```go
package dedup

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
)

func ComputeImportHashV2(data map[string]interface{}) string {
    // Remove self-referential fields
    delete(data, "import_hash_v2")
    delete(data, "original_source")

    jsonBytes, err := json.Marshal(data)
    if err != nil {
        return ""
    }

    hash := sha256.Sum256(jsonBytes)
    return hex.EncodeToString(hash[:])
}

func CheckDuplicate(db *sqlx.DB, userID int64, hash string) (int64, bool) {
    var journalID int64
    err := db.Get(&journalID, `
        SELECT tj.id
        FROM journal_meta jm
        INNER JOIN transaction_journals tj ON tj.id = jm.transaction_journal_id
        WHERE jm.name = 'import_hash_v2'
          AND jm.data = ?
          AND tj.user_id = ?
    `, "\""+hash+"\"", userID)

    if err != nil {
        return 0, false
    }
    return journalID, true
}
```

## 11. Summary Table

| Aspect | Detail |
|--------|--------|
| **Hash name** | `import_hash_v2` |
| **Algorithm** | SHA-256 |
| **Input** | JSON-encoded array semua transaction fields (minus import_hash_v2, original_source) |
| **Detection type** | Exact match (binary hash comparison) |
| **Scope** | Per-user |
| **Includes soft-deleted** | Yes |
| **Opt-in** | `error_if_duplicate_hash = true` |
| **Default** | Disabled (hash stored but not checked) |
| **On duplicate** | 422: `"Duplicate of transaction #<group_id>."` |
| **Atomicity** | Yes — entire group rolled back |
| **Storage** | `journal_meta` (EAV: name='import_hash_v2') |

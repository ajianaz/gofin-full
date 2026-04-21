package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"time"

	google "github.com/google/uuid"
)

// NewV7 generates a UUID v7 (RFC 9562) — time-sortable, globally unique.
// 48-bit timestamp ms | 4-bit version (7) | 12-bit rand | 2-bit variant (1) | 62-bit rand
func NewV7() google.UUID {
	var u google.UUID
	now := time.Now().UTC()

	millis := uint64(now.UnixMilli())
	binary.BigEndian.PutUint64(u[0:], millis)

	u[6] = (u[6] & 0x0F) | 0x70 // version 7
	u[8] = (u[8] & 0x3F) | 0x80 // variant 1

	_, _ = rand.Read(u[6:])
	u[6] = (u[6] & 0x0F) | 0x70 // re-apply version after rand
	u[8] = (u[8] & 0x3F) | 0x80 // re-apply variant after rand

	return u
}

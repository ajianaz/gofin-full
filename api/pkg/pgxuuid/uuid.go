package pgxuuid

import (
	"database/sql/driver"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Codec struct{}

func (Codec) FormatSupported(int16) bool { return true }
func (Codec) PreferredFormat() int16     { return pgtype.BinaryFormatCode }

func (c Codec) PlanEncode(m *pgtype.Map, oid uint32, format int16, value any) pgtype.EncodePlan {
	return encodePlan{}
}

type encodePlan struct{}

func (encodePlan) Encode(value any, buf []byte) ([]byte, error) {
	u, ok := value.(uuid.UUID)
	if !ok {
		return nil, nil
	}
	return append(buf, u[:]...), nil
}

func (c Codec) PlanScan(m *pgtype.Map, oid uint32, format int16, target any) pgtype.ScanPlan {
	return scanPlan{}
}

type scanPlan struct{}

func (scanPlan) Scan(src []byte, dst any) error {
	if src == nil {
		return nil
	}
	u, ok := dst.(*uuid.UUID)
	if !ok {
		return nil
	}
	copied := make([]byte, 16)
	copy(copied, src)
	*u = uuid.UUID(copied)
	return nil
}

func (c Codec) DecodeDatabaseSQLValue(m *pgtype.Map, oid uint32, format int16, src []byte) (driver.Value, error) {
	if src == nil {
		return nil, nil
	}
	u := uuid.UUID{}
	copy(u[:], src)
	return u[:], nil
}

func (c Codec) DecodeValue(m *pgtype.Map, oid uint32, format int16, src []byte) (any, error) {
	if src == nil {
		return nil, nil
	}
	u := uuid.UUID{}
	copy(u[:], src)
	return u, nil
}

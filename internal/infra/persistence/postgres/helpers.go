package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/sqlc-dev/pqtype"
)

func IsNoRowsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func NullString(s *string) sql.NullString {
	if s == nil || *s == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func NullInt32(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{
		Int32: *i,
		Valid: true,
	}
}

func NullBool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{
		Bool:  *b,
		Valid: true,
	}
}

func NullRawMessage(m *map[string]string) pqtype.NullRawMessage {
	if m == nil {
		return pqtype.NullRawMessage{}
	}
	bytes, err := json.Marshal(*m)
	if err != nil {
		return pqtype.NullRawMessage{}
	}
	return pqtype.NullRawMessage{
		RawMessage: bytes,
		Valid:      true,
	}
}

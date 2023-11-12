package myutil

import (
	"database/sql"
	"time"
)

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullInt(s int64) sql.NullInt64 {
	if s == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: s,
		Valid: true,
	}
}

func NewNullTime(s time.Time) sql.NullTime {
	if s.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  s,
		Valid: true,
	}
}

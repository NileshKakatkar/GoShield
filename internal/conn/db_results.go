package conn

import (
	"database/sql"
	"goshield/internal/gbl"
)

type ExecResult struct {
	InsertId     int64
	RowsAffected int64
	Error        error
}

type QueryResult struct {
	Rows  *sql.Rows
	Error error
}

func GetFieldValue(field sql.NullString) string {
	if field.Valid {
		return field.String
	}
	return gbl.EmptyString
}

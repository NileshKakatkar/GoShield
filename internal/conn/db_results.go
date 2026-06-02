package conn

import "database/sql"

type ExecResult struct {
	InsertId     int64
	RowsAffected int64
	Error        error
}

type QueryResult struct {
	Rows  *sql.Rows
	Error error
}

package conn

import (
	"context"
	"fmt"
	"goshield/internal/config"
	"goshield/internal/log"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

var sysDb *sqlx.DB = nil

type SqlDb struct {
	SysDb *sqlx.DB
}

func InitDbConn() error {
	dbConfig := config.GetDbConfig()
	var err error
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable\u0026connect_timeout=5", dbConfig.UserName, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	sysDb, err = sqlx.ConnectContext(context.Background(), "postgres", dbUrl)
	if err != nil {
		log.Logger.Errorf("Error while creating connection: %+v\n", err)
		return err
	}
	log.Logger.Info("Database connection successful")
	return nil
}

func NewConnection() *SqlDb {
	return &SqlDb{SysDb: sysDb}
}

func (sqlDb *SqlDb) QuerySql(query string, params ...any) QueryResult {
	return sqlDb.querySqlData(sqlDb.SysDb, query, params...)
}

func (sqlDb *SqlDb) querySqlData(db *sqlx.DB, query string, params ...any) QueryResult {
	rows, err := db.Query(query, params...)
	if err != nil {
		log.Logger.Errorf("Error while executing query: %v", err)
	}
	return QueryResult{Rows: rows, Error: err}
}

func (sqlDb *SqlDb) ExecSql(query string, params ...any) ExecResult {
	return sqlDb.execSqlQry(sqlDb.SysDb, query, params...)
}

func (sqlDb *SqlDb) execSqlQry(db *sqlx.DB, query string, params ...any) ExecResult {
	result, err := db.Exec(query, params...)
	execResult := ExecResult{Error: err}
	if err != nil {
		log.Logger.Errorf("Error while executing query: %v", err)
	} else {
		execResult.RowsAffected, _ = result.RowsAffected()
		execResult.InsertId, _ = result.LastInsertId()
		log.Logger.Debugf("ExecSql rowsAffected: %d, lastInsertId: %d,",
			execResult.RowsAffected, execResult.InsertId)
	}
	return execResult
}

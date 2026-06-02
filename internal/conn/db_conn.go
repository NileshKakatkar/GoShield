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

func InitDbConn() {
	dbConfig := config.GetDbConfig()
	var err error
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable\u0026connect_timeout=5", dbConfig.UserName, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	sysDb, err = sqlx.ConnectContext(context.Background(), "postgres", dbUrl)
	if err != nil {
		log.Logger.Errorf("Error while creating connection: %+v\n", err)
		return
	}
	log.Logger.Info("Database connection successful")
}

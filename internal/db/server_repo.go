package db

import (
	"database/sql"
	"goshield/internal/conn"
	"goshield/internal/log"
	"goshield/internal/model"
)

type ServerRepo interface {
	GetAllServers() ([]model.ServerInfo, error)
	AddServer(server model.ServerInfo) conn.ExecResult
	IsServerExist(server model.ServerInfo) (bool, error)
	GetServerByUnqiqId(uniqId int) (*model.ServerInfo, error)
	UpdateServerStatus(server model.ServerInfo) error
}

type serverRepo struct{}

func NewServerRepo() ServerRepo {
	return &serverRepo{}
}

func (s *serverRepo) GetAllServers() ([]model.ServerInfo, error) {
	var servers []model.ServerInfo

	qryRes := conn.NewConnection().QuerySql("select uniq_id, host_address, port, method, status from serverinfo")

	if qryRes.Error != nil {
		log.Logger.Errorf("Error while getting data from database: %+v\n", qryRes.Error)
		return servers, qryRes.Error
	}

	rows := qryRes.Rows
	defer rows.Close()

	if rows != nil {

		for rows.Next() {
			var server model.ServerInfo
			err := rows.Scan(&server.UniqId, &server.HostAddress, &server.Port, &server.Method, &server.Status)

			if err != nil {
				log.Logger.Errorf("Error while extracting data: %+v\n", err)
				return servers, err
			} else {
				servers = append(servers, server)
			}

		}
	}

	return servers, nil
}

func (s *serverRepo) AddServer(server model.ServerInfo) conn.ExecResult {
	execRes := conn.NewConnection().ExecSql("insert into serverinfo (host_address, port, method, status) values ($1, $2, $3, $4)", server.HostAddress, server.Port, server.Method, server.Status)

	if execRes.Error != nil {
		log.Logger.Errorf("Error while inserting data: %+v\n", execRes.Error)
	} else {
		log.Logger.Info("Data inserted successfully")
	}

	return execRes
}

func (s *serverRepo) IsServerExist(server model.ServerInfo) (bool, error) {
	isExist := false
	qryRes := conn.NewConnection().QuerySql("select host_address, port from serverinfo where host_address=$1 and port=$2", server.HostAddress, server.Port)

	if qryRes.Error != nil {
		log.Logger.Errorf("Error in query: %+v\n", qryRes.Error)
		return false, qryRes.Error
	}

	rows := qryRes.Rows
	defer rows.Close()

	if rows != nil {

		var hostAddress sql.NullString
		var port int

		for rows.Next() {

			err := rows.Scan(&hostAddress, &port)

			if err != nil {
				log.Logger.Errorf("Error while scanning data: %+v\n", err)
				return false, err
			} else {
				isExist = server.HostAddress == conn.GetFieldValue(hostAddress) && server.Port == port
			}
		}
	}
	return isExist, nil
}

func (s *serverRepo) GetServerByUnqiqId(uniqId int) (*model.ServerInfo, error) {
	qryRes := conn.NewConnection().QuerySql("select host_address, port, method, status from serverinfo where uniq_id=$1", uniqId)

	if qryRes.Error != nil {
		log.Logger.Errorf("Error while query execution: %+v\n", qryRes.Error)
		return nil, qryRes.Error
	}

	var server model.ServerInfo

	rows := qryRes.Rows
	defer rows.Close()

	if rows != nil {
		for rows.Next() {
			err := rows.Scan(&server.HostAddress, &server.Port, &server.Method, &server.Status)
			if err != nil {
				log.Logger.Errorf("Error while extracting data: %+v\n", err)
				return nil, err
			}
		}
	}
	return &server, nil
}

func (s *serverRepo) UpdateServerStatus(server model.ServerInfo) error {
	execRes := conn.NewConnection().ExecSql("update serverinfo set status=$1 where uniq_id=$2", server.Status, server.UniqId)

	if execRes.Error != nil {
		log.Logger.Errorf("Error while executing request: %+v\n", execRes.Error)
		return execRes.Error
	}

	return nil
}

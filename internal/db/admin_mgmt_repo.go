package db

import (
	"goshield/internal/conn"
	"goshield/internal/log"
	"goshield/internal/model"
)

type AdminRepo interface {
	GetAdminDetails() (model.Admin, error)
}

type adminRepo struct{}

func NewAdminRepo() AdminRepo {
	return &adminRepo{}
}

func (a *adminRepo) GetAdminDetails() (model.Admin, error) {
	var admin model.Admin

	qryRes := conn.NewConnection().QuerySql("select uniq_id, admin_name, admin_id, password from admin")

	if qryRes.Error != nil {
		log.Logger.Errorf("Error while executing query: %+v\n", qryRes.Error)
		return admin, qryRes.Error
	}

	rows := qryRes.Rows
	defer rows.Close()

	if rows != nil {
		for rows.Next() {
			err := rows.Scan(&admin.UniqId, &admin.AdminName, &admin.AdminId, &admin.Password)
			if err != nil {
				log.Logger.Errorf("Error while extracting data: %+v\n", err)
			}
		}
	}

	return admin, nil
}

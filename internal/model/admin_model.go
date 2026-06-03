package model

type Admin struct {
	UniqId    int    `json:"uniq_id"`
	AdminName string `json:"admin_name"`
	AdminId   string `json:"admin_id"`
	Password  string `json:"password"`
}

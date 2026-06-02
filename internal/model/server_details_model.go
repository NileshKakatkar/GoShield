package model

type ServerInfo struct {
	UniqId      int    `json:"uniq_id"`
	HostAddress string `json:"host_address"`
	Port        int    `json:"port"`
	Method      string `json:"method"`
	Status      string `json:"status"`
}

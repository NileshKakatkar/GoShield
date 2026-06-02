package model

type ServerResp struct {
	ServersList []ServerInfo `json:"server_info_list"`
	Server      ServerInfo   `json:"server_info"`
	Message     string       `json:"message"`
	Status      int          `json:"status"`
}

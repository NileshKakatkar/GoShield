package model

type Server struct {
	ServerAddress string `json:"server_address"`
	ServerPort    int    `json:"server_port"`
	HttpMethod    string `json:"http_method"`
}

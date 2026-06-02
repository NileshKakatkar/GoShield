package model

type Config struct {
	Proxy    ProxyConfig     `json:"proxy"`
	DbConfig DbConfig        `json:"dbconfig"`
}

type ProxyConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type DbConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"username"`
	Database string `json:"database"`
	Password string `json:"password"`
}

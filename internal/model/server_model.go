package model

type Config struct {
	Proxy    ProxyConfig     `json:"proxy"`
	Backends []BackendConfig `json:"backends"`
	DbConfig DbConfig        `json:"dbconfig"`
}

type ProxyConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type BackendConfig struct {
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

package model

type Config struct {
	Proxy    ProxyConfig     `json:"proxy"`
	Backends []BackendConfig `json:"backends"`
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

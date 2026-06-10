package model

type Config struct {
	Proxy        ProxyConfig  `json:"proxy"`
	DbConfig     DbConfig     `json:"dbconfig"`
	Certificates Certificates `json:"certificates"`
	RedisConfig  RedisConfig  `json:"redisConfig"`
}

type ProxyConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	TlsPort  int    `json:"tls_port"`
	Protocol string `json:"protocol"`
}

type DbConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"username"`
	Database string `json:"database"`
	Password string `json:"password"`
}

type Certificates struct {
	PublicKey  string `json:"pub_key"`
	PrivateKey string `json:"pri_key"`
}

type RedisConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

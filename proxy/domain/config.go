package domain

type ServerConfig interface {
	GetScheme() string
	GetPort() string
	GetAddress() string
	GetCertFilePath() string
	GetKeyFilePath() string
}

type ConfigModel interface {
	NextServer() interface{}
}

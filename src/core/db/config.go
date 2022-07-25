package db

type Config struct {
	Host               string
	Port               string
	Username           string
	Password           string
	Database           string
	MaxOpenConnections int
	MaxIdleConnections int
}

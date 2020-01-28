package fw

import "github.com/jinzhu/gorm"

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

type DBConnector interface {
	Connect(config DBConfig) (*gorm.DB, error)
}

package mddb

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sentadmedia/elf/fw"
)

var _ fw.DBConnector = (*GormPstgresConnector)(nil)

type GormPstgresConnector struct {
}

func (p GormPstgresConnector) Connect(dbConfig fw.DBConfig) (*gorm.DB, error) {
	dataSource := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DbName,
	)

	db, err := gorm.Open("postgres", dataSource)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewPostgresConnector() GormPstgresConnector {
	return GormPstgresConnector{}
}

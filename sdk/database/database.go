package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Interface interface {
	GetInstance() *gorm.DB
}

type Config struct {
	Host     string
	Port     string
	Password string
	Username string
	Database string
}

type sql struct {
	Db *gorm.DB
}

func InitDB(config Config) (Interface, error) {
	sql := sql{}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	sql.Db = db

	return &sql, nil
}

func (s *sql) GetInstance() *gorm.DB {
	return s.Db
}

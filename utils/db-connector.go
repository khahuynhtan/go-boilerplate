package utils

import (
	"fmt"
	"net/url"
	"sync"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

type DBConnector struct {
	DB_NAME     string
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
}

func NewConnector(dbName, dbHost, dbPort, dbUser, dbPassword string) *DBConnector {
	return &DBConnector{
		DB_NAME:     dbName,
		DB_HOST:     dbHost,
		DB_PORT:     dbPort,
		DB_USER:     dbUser,
		DB_PASSWORD: dbPassword,
	}
}

func (dbConnector DBConnector) Connect() *gorm.DB {
	once.Do(func() {

		params := url.Values{}
		params.Set("host", dbConnector.DB_HOST)
		params.Set("user", dbConnector.DB_USER)
		params.Set("password", dbConnector.DB_PASSWORD)
		params.Set("dbname", dbConnector.DB_NAME)
		params.Set("port", dbConnector.DB_PORT)
		params.Set("sslmode", "disable")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			dbConnector.DB_HOST, dbConnector.DB_USER, dbConnector.DB_PASSWORD, dbConnector.DB_NAME, dbConnector.DB_PORT)

		db, err := gorm.Open(driver.Open(dsn), &gorm.Config{})
		if err != nil {
			Logger(ErrorLevel, "Failed to connect to the database: %v")
			return
		}
		Logger(InfoLevel, "Database connection is successful")
		instance = db
	})
	if instance == nil {
		Logger(ErrorLevel, "Database connection instance is nil")
	}
	return instance
}

func GetDB() *gorm.DB {
	if instance == nil {
		Logger(ErrorLevel, "Database connection instance is nil")
		return nil
	}
	return instance
}

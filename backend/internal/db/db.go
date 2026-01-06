package db

import (
	"backend/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	// Example DSN (Data Source Name):
	/*
		host=localhost
		user=postgres
		password=postgres
		dbname=mydb port=5432
		sslmode=disable
		TimeZone=Asia/Ho_Chi_Minh
	*/
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// config DB
	sqlDB, err := database.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// Create tables if not exists
	if err := database.AutoMigrate(
		&models.Account{},
		&models.LogHead{},
		&models.LogContent{},
	); err != nil {
		return nil, err
	}

	return database, nil
}

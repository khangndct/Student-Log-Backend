package db

import (
	"backend/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(sqlitePath string) (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{}) // default config
	if err != nil {
		return nil, err
	}

	// This create the tables if not exists
	if err := database.AutoMigrate(
		&models.Account{},
		&models.LogHead{},
		&models.LogContent{},
	); err != nil {
		return nil, err
	}

	return database, nil
}

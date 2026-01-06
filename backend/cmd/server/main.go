package main

import (
	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/routes"
	"backend/internal/utils"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

const jwtSecret = "dev_secret_change_me"

// This is a function to seed an admin when currently there is no admin
func seedAdmin(database *gorm.DB) {
	var count int64
	database.Model(&models.Account{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		return
	}

	hash, _ := utils.HashPassword("admin123")
	_ = database.Create(&models.Account{
		Username: "admin",
		Email:    "admin@example.com",
		Phone:    0,
		Password: hash,
	}).Error

	log.Println("Seeded admin account: admin / admin123")
}

func main() {
	database, err := db.Connect("app.db")
	if err != nil {
		log.Fatal(err)
	}
	seedAdmin(database)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routes.Register(e, database, jwtSecret)

	log.Println("Server running at :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/routes"
	"backend/internal/utils"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func loadEnvFiles(paths ...string) {
	for _, path := range paths {
		if loadEnvFile(path) {
			return
		}
	}
}

func loadEnvFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"'")

		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}

	return true
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func buildPostgresDSN() (string, error) {
	user := strings.TrimSpace(os.Getenv("POSTGRES_USER"))
	password := strings.TrimSpace(os.Getenv("POSTGRES_PASSWORD"))
	if user == "" || password == "" {
		return "", fmt.Errorf("missing POSTGRES_USER or POSTGRES_PASSWORD")
	}

	host := envOrDefault("POSTGRES_HOST", "localhost")
	port := envOrDefault("POSTGRES_PORT", "5432")
	dbName := envOrDefault("POSTGRES_DB", "postgres")
	sslMode := envOrDefault("POSTGRES_SSLMODE", "disable")
	timeZone := envOrDefault("POSTGRES_TIMEZONE", "Asia/Ho_Chi_Minh")

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host,
		user,
		password,
		dbName,
		port,
		sslMode,
		timeZone,
	), nil
}

func loadJWTSecret() (string, error) {
	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		return "", fmt.Errorf("missing JWT_SECRET")
	}
	return secret, nil
}

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

// This is a function to seed a member account for testing purposes
func seedMember(database *gorm.DB) {
	var count int64
	database.Model(&models.Account{}).Where("username = ?", "member").Count(&count)
	if count > 0 {
		return
	}

	hash, _ := utils.HashPassword("member123")
	_ = database.Create(&models.Account{
		Username: "member",
		Email:    "member@example.com",
		Phone:    1234567890,
		Password: hash,
	}).Error

	log.Println("Seeded member account: member / member123")
}

// This is a function to seed a log head with admin as owner and admin+member as writers, plus a sample log content
func seedLogHeadWithContent(database *gorm.DB) {
	// Check if log head already exists (by checking for a specific subject)
	var count int64
	database.Model(&models.LogHead{}).Where("subject = ?", "Sample Student Diary Log").Count(&count)
	if count > 0 {
		return
	}

	// Get admin and member account IDs
	var admin models.Account
	if err := database.Where("username = ?", "admin").First(&admin).Error; err != nil {
		log.Println("Warning: Admin account not found, skipping log head seeding")
		return
	}

	var member models.Account
	if err := database.Where("username = ?", "member").First(&member).Error; err != nil {
		log.Println("Warning: Member account not found, skipping log head seeding")
		return
	}

	// Create log head with admin as owner and both admin and member as writers
	now := time.Now()
	logHead := models.LogHead{
		Subject:      "Sample Student Diary Log",
		StartDate:    now.AddDate(0, 0, -7), // 7 days ago
		EndDate:      now.AddDate(0, 0, 30),  // 30 days from now
		WriterIDList: pq.Int64Array{int64(admin.ID), int64(member.ID)},
		OwnerID:      uint(admin.ID),
	}

	if err := database.Create(&logHead).Error; err != nil {
		log.Printf("Warning: Failed to create log head: %v", err)
		return
	}

	// Create a sample log content (event) for this log head, written by admin
	logContent := models.LogContent{
		LogHeadID: logHead.ID,
		WriterID:  uint(admin.ID),
		Content:   "Today we had a great class discussion about project management. All students participated actively and shared their ideas.",
		Date:      now,
	}

	if err := database.Create(&logContent).Error; err != nil {
		log.Printf("Warning: Failed to create log content: %v", err)
		return
	}

	log.Printf("Seeded log head (ID: %d) with owner: admin (ID: %d), writers: admin (ID: %d) and member (ID: %d)", 
		logHead.ID, admin.ID, admin.ID, member.ID)
	log.Printf("Seeded log content (ID: %d) for log head (ID: %d) written by admin", 
		logContent.ID, logHead.ID)
}

func main() {
	loadEnvFiles(".env", filepath.Join("..", ".env"))
	dsn, err := buildPostgresDSN()
	if err != nil {
		log.Fatal(err)
	}

	jwtSecret, err := loadJWTSecret()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.Connect(dsn)
	if err != nil {
		log.Fatal(err)
	}
	seedAdmin(database)
	seedMember(database)
	seedLogHeadWithContent(database)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routes.Register(e, database, jwtSecret)

	log.Println("Server running at :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

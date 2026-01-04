package models

import "time"

type Account struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"not null"` // "admin" | "member"
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

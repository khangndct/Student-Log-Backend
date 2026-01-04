package models

import "time"

type LogHead struct {
	ID         uint   `gorm:"primaryKey"`
	Title      string `gorm:"not null"`
	WriteScope string `gorm:"not null"` // "all" | "owner" | "admin"
	OwnerID    uint   `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	LogContents []LogContent `gorm:"constraint:OnDelete:CASCADE;"`
}

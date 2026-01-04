package models

import "time"

type LogContent struct {
	ID        uint      `gorm:"primaryKey"`
	LogHeadID uint      `gorm:"index;not null"`
	UserID    uint      `gorm:"index;not null"`
	Content   string    `gorm:"not null"`
	LogTime   time.Time `gorm:"not null"` // timestamp of input
	CreatedAt time.Time
	UpdatedAt time.Time
}

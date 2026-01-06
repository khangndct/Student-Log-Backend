package models

import "time"

type LogContent struct {
	ID        uint      `gorm:"primaryKey"`
	LogHeadID uint      `gorm:"index;not null"`
	WriterID  uint      `gorm:"index;not null"`
	Content   string    `gorm:"not null"`
	Date      time.Time `gorm:"not null"` // timestamp of input
}

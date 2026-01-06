package models

import "time"

type LogContent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	LogHeadID uint      `gorm:"index;not null" json:"log_head_id"`
	WriterID  uint      `gorm:"index;not null" json:"writer_id"`
	Content   string    `gorm:"not null" json:"content"`
	Date      time.Time `gorm:"not null" json:"date"` // timestamp of input
}

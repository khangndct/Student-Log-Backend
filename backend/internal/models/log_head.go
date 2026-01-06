package models

import (
	"time"

	"github.com/lib/pq"
)

type LogHead struct {
	ID           uint   `gorm:"primaryKey"`
	Subject      string `gorm:"not null"`
	StartDate    time.Time
	EndDate      time.Time
	WriterIDList pq.Int64Array `gorm:"not null"`
	OwnerID      uint          `gorm:"not null"`

	LogContents []LogContent `gorm:"constraint:OnDelete:CASCADE;"`
}

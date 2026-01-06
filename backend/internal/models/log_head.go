package models

import (
	"time"

	"github.com/lib/pq"
)

type LogHead struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Subject      string        `gorm:"not null" json:"subject"`
	StartDate    time.Time     `json:"start_date"`
	EndDate      time.Time     `json:"end_date"`
	WriterIDList pq.Int64Array `gorm:"type:bigint[];not null" json:"writer_id_list"`
	OwnerID      uint          `gorm:"not null" json:"owner_id"`

	LogContents []LogContent `gorm:"constraint:OnDelete:CASCADE;" json:"log_contents"`
}

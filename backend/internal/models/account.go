package models

type Account struct {
	ID       int64  `gorm:"primaryKey;type:bigint"`
	Username string `gorm:"type:varchar(255);not null"`
	Phone    int64  `gorm:"type:bigint;not null"`
	Email    string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Password string `gorm:"type:varchar(255);not null"` // Hashed
}

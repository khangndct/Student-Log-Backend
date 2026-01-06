package models

type Account struct {
	ID       int64  `gorm:"primaryKey;type:bigint" json:"id"`
	Username string `gorm:"type:varchar(255);not null" json:"username"`
	Phone    int64  `gorm:"type:bigint;not null" json:"phone"`
	Email    string `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"password"` // Hashed
}

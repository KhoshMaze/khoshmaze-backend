package types

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName        string `gorm:"type:varchar(255);not null"`
	LastName         string `gorm:"type:varchar(255);not null"`
	Phone            string `gorm:"type:varchar(255);not null;unique"`
	SubscribtionType uint   `gorm:"not null;default:0"`
}

type TokenBlacklist struct {
	ExpiresAt time.Time
	Value     string `gorm:"type:text;not null;unique"`
	UserID    uint
}

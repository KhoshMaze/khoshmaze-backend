package types

import "gorm.io/gorm"

type Restaurant struct {
	gorm.Model
	Name    string `gorm:"not null"`
	URL     string `gorm:"uniqueIndex, not null"`
	OwnerID uint `gorm:"not null;foreignKey:UserID"`
}

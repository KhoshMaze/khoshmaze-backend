package types

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	PrimaryColor   string `gorm:"varchar(255);not null"`
	SecondaryColor string `gorm:"varchar(255);not null"`
	Foods          []Food
	BranchID       uint `gorm:"not null;index:,unique"`
}

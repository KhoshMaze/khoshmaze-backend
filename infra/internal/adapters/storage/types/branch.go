package types

import "gorm.io/gorm"

type Branch struct {
	gorm.Model
	RestaurantID   uint   `gorm:"not null;"`
	Name           string `gorm:"type:varchar(255)"`
	Address        string
	Phone          string `gorm:"type:varchar(255)"`
	PrimaryColor   string `gorm:"type:varchar(255);not null;default:#000000"`
	SecondaryColor string `gorm:"type:varchar(255);not null;default:#FFFFFF"`
	Foods          []Food
}

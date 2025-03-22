package types

import "gorm.io/gorm"

type Branch struct {
	gorm.Model
	RestaurantID uint `gorm:"not null;foreignKey:RestaurantID"`
	Address      string
	Phone        string `gorm:"not null"`
}

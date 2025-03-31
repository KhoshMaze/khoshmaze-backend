package types

import "gorm.io/gorm"

type Restaurant struct {
	gorm.Model
	Name    string `gorm:"not null;type:varchar(255)"`
	URL     string `gorm:"uniqueIndex;type:varchar(255); not null"`
	OwnerID uint   `gorm:"not null;"`
	Branches []*Branch `gorm:"foreignKey:RestaurantID"`
}

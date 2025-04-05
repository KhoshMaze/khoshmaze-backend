package types

import "time"

type Food struct {
	ID          uint      `gorm:"primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;not null"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text;"`
	Type        string    `gorm:"type:varchar(255);not null"`
	IsAvailable bool      `gorm:"boolean;not null;default:true"`
	Price       float64   
	FoodPrices  []FoodPrice
	Images      []FoodImage
	MenuID      uint      `gorm:"not null"`
}

type FoodImage struct {
	ID     uint   `gorm:"primaryKey"`
	Image  []byte `gorm:"type:bytea;not null"`
	FoodID uint   `gorm:"not null"`
}

type FoodPrice struct {
	ID    uint    `gorm:"primaryKey"`
	Price float64 `gorm:"not null"`
	FoodID  uint  `gorm:"not null"`
}

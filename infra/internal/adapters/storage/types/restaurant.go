package types

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	ID                  uint      `gorm:"primaryKey"`
	CreatedAt           time.Time `gorm:"not null;autoCreateTime"`
	MaxBranchCount      uint      `gorm:"not null;default:1"`
	ExpiresAt           time.Time `gorm:"not null"`
	SubscriptionPriceID uint
	SubscriptionPrice   SubscriptionPrice `gorm:"foreignKey:SubscriptionPriceID"`
}

type SubscriptionPrice struct {
	gorm.Model
	Price         uint           `gorm:"not null;index:,sort:desc;"`
	Type          uint8          `gorm:"not null;"`
	Subscriptions []Subscription `gorm:"foreignKey:SubscriptionPriceID"`
}

type Restaurant struct {
	gorm.Model
	Name           string    `gorm:"not null;type:varchar(255)"`
	URL            string    `gorm:"uniqueIndex;type:varchar(255); not null;"`
	OwnerID        uint      `gorm:"not null;"`
	Branches       []*Branch `gorm:"foreignKey:RestaurantID"`
	SubscriptionID uint      `gorm:"constraint:OnDelete:CASCADE;unique;"`
	Subscription   Subscription
}

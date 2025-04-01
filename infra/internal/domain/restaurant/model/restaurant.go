package model

import (
	"errors"
	"time"
)

type SubscriptionType uint8

const (
	SubscriptionTypeNormal SubscriptionType = iota
	SubscriptionTypePremium
)

type Subscription struct {
	ID             uint
	MaxBranchCount uint
	CreatedAt      time.Time
	ExpiresAt      time.Time
	Price          SubscriptionPrice
}

type SubscriptionPrice struct {
	ID    uint
	Price uint
	Type  SubscriptionType
}

type Restaurant struct {
	ID           uint
	Name         string
	URL          string
	OwnerID      uint
	Branches     []*Branch
	Subscription Subscription
}

type RestaurantFilter struct {
	ID  uint
	URL string
}

func (r *Restaurant) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}

	if r.URL == "" {
		return errors.New("url is required")
	}

	if r.OwnerID == 0 {
		return errors.New("owner id is required")
	}

	return nil
}

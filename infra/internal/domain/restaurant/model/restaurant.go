package model

import (
	"errors"
)

type Restaurant struct {
	ID       uint
	Name     string
	URL      string
	OwnerID  uint
	Branches []*Branch
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

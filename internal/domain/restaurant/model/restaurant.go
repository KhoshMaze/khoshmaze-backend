package model

import (
	"errors"

	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"
)

type Restaurant struct {
	ID      uint
	Name    string
	URL     string
	OwnerID model.UserID
}

type RestaurantFilter struct {
	ID  model.UserID
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
package model

import "github.com/KhoshMaze/khoshmaze-backend/internal/domain/menu/model"

type Branch struct {
	ID           uint
	RestaurantID uint
	Name         string
	Address      string
	Phone        string
	Menu         model.Menu
}

type BranchFilter struct {
	ID             uint
	RestaurantName string
}

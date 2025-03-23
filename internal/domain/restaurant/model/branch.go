package model

type Branch struct {
	ID           uint
	RestaurantID uint
	Name         string
	Address      string
	Phone        string
}

type BranchFilter struct {
	ID uint
}

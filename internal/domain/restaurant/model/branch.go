package model

type Branch struct {
	ID           uint
	RestaurantID uint
	Address      string
	Phone        string
}

type BranchFilter struct {
	ID uint
}

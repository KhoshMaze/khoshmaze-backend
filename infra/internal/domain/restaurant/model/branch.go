package model

type Branch struct {
	ID             uint
	RestaurantID   uint
	Name           string
	Address        string
	Phone          string
	PrimaryColor   string
	SecondaryColor string
}

type BranchFilter struct {
	ID             uint
	RestaurantName string
}

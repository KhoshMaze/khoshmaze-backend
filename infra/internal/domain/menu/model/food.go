package model

type Food struct {
	ID          uint
	Name        string
	Description string
	Type        string
	IsAvailable bool
	MenuID      uint
	Price       float64
	Images      []FoodImage
}

type FoodPrice struct {
	ID uint
	Price float64
	FoodID uint
}

type FoodImage struct {
	ID uint 
	Image []byte 
	FoodID uint
}
package model

type Food struct {
	ID          uint
	Name        string
	Description string
	Type        string
	IsAvailable bool
	Price       *float64
	Images      []FoodImage
}

type FoodPrice struct {
	ID uint
	Price float64
}

type FoodImage struct {
	ID uint 
	Image []byte 
	FoodID uint
}
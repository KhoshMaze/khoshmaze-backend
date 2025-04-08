package mapper

import (
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/fp"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/menu/model"
)

func FoodStorageToDomain(food types.Food) *model.Food {
	return &model.Food{
		ID:          food.ID,
		Name:        food.Name,
		Description: food.Description,
		Type:        food.Type,
		IsAvailable: food.IsAvailable,
		Price:       food.Price,
		Images: fp.Map(food.Images, func(image types.FoodImage) model.FoodImage {
			return *FoodImageStorageToDomain(&image)
		}),
		BranchID: food.BranchID,
	}
}

func FoodDomainToStorage(food model.Food) *types.Food {
	return &types.Food{
		ID:          food.ID,
		Name:        food.Name,
		Description: food.Description,
		Type:        food.Type,
		IsAvailable: food.IsAvailable,
		Price:       food.Price,
		Images: fp.Map(food.Images, func(image model.FoodImage) types.FoodImage {
			return *FoodImageDomainToStorage(&image)
		}),
		BranchID: food.BranchID,
	}
}

func FoodImageStorageToDomain(image *types.FoodImage) *model.FoodImage {
	return &model.FoodImage{
		ID:     image.ID,
		Image:  image.Image,
		FoodID: image.FoodID,
	}
}

func FoodImageDomainToStorage(image *model.FoodImage) *types.FoodImage {
	return &types.FoodImage{
		ID:     image.ID,
		Image:  image.Image,
		FoodID: image.FoodID,
	}
}

package mapper

import (
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/model"
	userModel "github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"
	"gorm.io/gorm"
)

func RestaurantDomainToStorage(restaurantDomain *model.Restaurant) *types.Restaurant {
	return &types.Restaurant{
		Model: gorm.Model{
			ID: restaurantDomain.ID,
		},
		Name:    restaurantDomain.Name,
		URL:     restaurantDomain.URL,
		OwnerID: uint(restaurantDomain.OwnerID),
	}
}

func RestaurantStorageToDomain(restaurantStorage *types.Restaurant) *model.Restaurant {
	return &model.Restaurant{
		ID:      restaurantStorage.ID,
		Name:    restaurantStorage.Name,
		URL:     restaurantStorage.URL,
		OwnerID: userModel.UserID(restaurantStorage.OwnerID),
	}
}

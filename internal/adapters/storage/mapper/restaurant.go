package mapper

import (
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/model"
	"gorm.io/gorm"
)

func RestaurantDomainToStorage(restaurantDomain *model.Restaurant) *types.Restaurant {
	restaurantStorage := &types.Restaurant{
		Model: gorm.Model{
			ID: restaurantDomain.ID,
		},
		Name:     restaurantDomain.Name,
		URL:      restaurantDomain.URL,
		OwnerID:  uint(restaurantDomain.OwnerID),
		Branches: make([]*types.Branch, 0),
	}

	for _, branch := range restaurantDomain.Branches {
		restaurantStorage.Branches = append(restaurantStorage.Branches, BranchDomainToStorage(branch))
	}

	return restaurantStorage
}

func RestaurantStorageToDomain(restaurantStorage *types.Restaurant) *model.Restaurant {
	restaurantDomain := &model.Restaurant{
		ID:       restaurantStorage.ID,
		Name:     restaurantStorage.Name,
		URL:      restaurantStorage.URL,
		OwnerID:  restaurantStorage.OwnerID,
		Branches: make([]*model.Branch, 0),
	}

	for _, branch := range restaurantStorage.Branches {
		restaurantDomain.Branches = append(restaurantDomain.Branches, BranchStorageToDomain(branch))
	}

	return restaurantDomain
}

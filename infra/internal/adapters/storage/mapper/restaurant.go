package mapper

import (
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/model"
	"gorm.io/gorm"
)

func SubscriptionDomainToStorage(subscriptionDomain *model.Subscription) types.Subscription {
	return types.Subscription{
		ID: subscriptionDomain.ID,
		CreatedAt: subscriptionDomain.CreatedAt,
		ExpiresAt: subscriptionDomain.ExpiresAt,
		MaxBranchCount: subscriptionDomain.MaxBranchCount,
		SubscriptionPrice: types.SubscriptionPrice{
			Model: gorm.Model{
				ID: subscriptionDomain.Price.ID,
			},
			Price: subscriptionDomain.Price.Price,
			Type: uint8(subscriptionDomain.Price.Type),
		},
	}
}

func SubscriptionStorageToDomain(subscriptionStorage *types.Subscription) model.Subscription {
	return model.Subscription{
		ID: subscriptionStorage.ID,
		CreatedAt: subscriptionStorage.CreatedAt,
		MaxBranchCount: subscriptionStorage.MaxBranchCount,
		ExpiresAt: subscriptionStorage.ExpiresAt,
		Price: model.SubscriptionPrice{
			ID: subscriptionStorage.SubscriptionPrice.ID,
			Price: subscriptionStorage.SubscriptionPrice.Price,
			Type: model.SubscriptionType(subscriptionStorage.SubscriptionPrice.Type),
		},
	}
}


func RestaurantDomainToStorage(restaurantDomain *model.Restaurant) *types.Restaurant {
	restaurantStorage := &types.Restaurant{
		Model: gorm.Model{
			ID: restaurantDomain.ID,
		},
		Name:     restaurantDomain.Name,
		URL:      restaurantDomain.URL,
		OwnerID:  uint(restaurantDomain.OwnerID),
		Branches: make([]*types.Branch, 0),
		Subscription: SubscriptionDomainToStorage(&restaurantDomain.Subscription),
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
		Subscription: SubscriptionStorageToDomain(&restaurantStorage.Subscription),
	}

	for _, branch := range restaurantStorage.Branches {
		restaurantDomain.Branches = append(restaurantDomain.Branches, BranchStorageToDomain(branch))
	}

	return restaurantDomain
}

package mapper

import (
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	permModel "github.com/KhoshMaze/khoshmaze-backend/internal/domain/permission/model"
	restaurantModel "github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/model"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"
	"gorm.io/gorm"
)

func UserDomainToStorage(userDomain model.User) *types.User {
	user := types.User{
		Model: gorm.Model{
			ID:        uint(userDomain.ID),
			CreatedAt: userDomain.CreatedAt,
			UpdatedAt: userDomain.UpdatedAt,
		},
		FirstName:        userDomain.FirstName,
		LastName:         userDomain.LastName,
		Phone:            string(userDomain.Phone),
		SubscriptionType: uint(userDomain.SubscriptionType),
		Permissions:      uint64(userDomain.Permissions),
		Restaurants:      make([]*types.Restaurant, 0),
	}

	for _, restaurant := range userDomain.Restaurants {
		user.Restaurants = append(user.Restaurants, RestaurantDomainToStorage(restaurant))
	}

	return &user
}

func UserStorageToDomain(user types.User) *model.User {
	userDomain := &model.User{
		ID:               model.UserID(user.ID),
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		Phone:            model.Phone(user.Phone),
		SubscriptionType: IntegerToSubscriptionType(user.SubscriptionType),
		Permissions:      permModel.Permission(user.Permissions),
		Restaurants:      make([]*restaurantModel.Restaurant, 0),
	}

	for _, restaurant := range user.Restaurants {
		userDomain.Restaurants = append(userDomain.Restaurants, RestaurantStorageToDomain(restaurant))
	}

	return userDomain
}

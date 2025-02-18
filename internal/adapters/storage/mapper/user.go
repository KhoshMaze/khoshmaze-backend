package mapper

import (
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"
	"gorm.io/gorm"
)

func UserDomainToStorage(userDomain model.User) *types.User {
	return &types.User{
		Model: gorm.Model{
			ID:        uint(userDomain.UserID),
			CreatedAt: userDomain.CreatedAt,
			UpdatedAt: userDomain.UpdatedAt,
		},
		FirstName:        userDomain.FirstName,
		LastName:         userDomain.LastName,
		Phone:            string(userDomain.Phone),
		SubscribtionType: uint(userDomain.SubscribtionType),
	}

}

func UserStorageToDomain(user types.User) *model.User {
	return &model.User{
		UserID:           model.UserID(user.ID),
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		Phone:            model.Phone(user.Phone),
		SubscribtionType: IntegerToSubscribtionType(user.SubscribtionType),
	}
}

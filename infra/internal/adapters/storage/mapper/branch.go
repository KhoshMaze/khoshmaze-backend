package mapper

import (
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/model"
	"gorm.io/gorm"
)

func BranchDomainToStorage(branchDomain *model.Branch) *types.Branch {
	return &types.Branch{
		Model: gorm.Model{
			ID: branchDomain.ID,
		},
		Name:           branchDomain.Name,
		RestaurantID:   uint(branchDomain.RestaurantID),
		Address:        branchDomain.Address,
		Phone:          branchDomain.Phone,
		PrimaryColor:   branchDomain.PrimaryColor,
		SecondaryColor: branchDomain.SecondaryColor,
	}
}

func BranchStorageToDomain(branchStorage *types.Branch) *model.Branch {
	return &model.Branch{
		ID:           branchStorage.ID,
		Name:         branchStorage.Name,
		RestaurantID: branchStorage.RestaurantID,
		Address:      branchStorage.Address,
		Phone:        branchStorage.Phone,
		PrimaryColor: branchStorage.PrimaryColor,
		SecondaryColor: branchStorage.SecondaryColor,
	}
}

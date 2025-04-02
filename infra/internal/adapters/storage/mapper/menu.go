package mapper

import (
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/fp"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/menu/model"
	"gorm.io/gorm"
)

func MenuStorageToDomain(menu types.Menu) *model.Menu {
	return &model.Menu{
		ID:             menu.ID,
		PrimaryColor:   menu.PrimaryColor,
		SecondaryColor: menu.SecondaryColor,
		Foods: fp.Map(menu.Foods, func(food types.Food) model.Food {
			return *FoodStorageToDomain(food)
		}),
	}
}

func MenuDomainToStorage(menu model.Menu) *types.Menu {
	return &types.Menu{
		Model: gorm.Model{
			ID: menu.ID,
		},	
		PrimaryColor:   menu.PrimaryColor,
		SecondaryColor: menu.SecondaryColor,
		Foods: fp.Map(menu.Foods, func(food model.Food) types.Food {
			return *FoodDomainToStorage(food)
		}),
	}
}

package port

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/menu/model"
)

type Service interface {
	AddFoodToMenu(ctx context.Context, food model.Food) (uint, error)
	UpdateFoodInMenu(ctx context.Context, food model.Food) error
	DeleteFoodFromMenu(ctx context.Context, foodID uint) error
	GetAllFoods(ctx context.Context, pagination *common.Pagination, branchID uint) (*common.PaginatedResponse[*model.Food], error)
	GetFoodByID(ctx context.Context, id uint) (*model.Food, error)
	GetImagesByFoodID(ctx context.Context, foodID uint, pagination *common.Pagination) (*common.PaginatedResponse[*model.FoodImage], error)
	AddFoodImageToFood(ctx context.Context, image *model.FoodImage) error
	DeleteFoodImageFromFood(ctx context.Context,imageID uint) error
}

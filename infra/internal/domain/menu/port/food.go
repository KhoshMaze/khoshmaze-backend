package port

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/menu/model"
)

type FoodRepository interface {
	Create(ctx context.Context, food model.Food) (uint, error)
	GetAll(ctx context.Context, pagination *common.Pagination, branchID uint) (*common.PaginatedResponse[*model.Food], error)
	GetByID(ctx context.Context, id uint) (*model.Food, error)
	GetImagesByFoodID(ctx context.Context, foodID uint, pagination *common.Pagination) (*common.PaginatedResponse[*model.FoodImage], error)
	Update(ctx context.Context, food model.Food) error
	Delete(ctx context.Context, id uint) error
	CreateImage(ctx context.Context, image *model.FoodImage) error
	DeleteImage(ctx context.Context, id uint) error
}

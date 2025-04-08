package menu

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/menu/model"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/menu/port"
)

type service struct {
	foodRepo port.FoodRepository
}

func NewService(foodRepo port.FoodRepository) port.Service {
	return &service{
		foodRepo: foodRepo,
	}
}

func (s *service) AddFoodToMenu(ctx context.Context, food model.Food) (uint, error) {
	return s.foodRepo.Create(ctx, food)
}

func (s *service) UpdateFoodInMenu(ctx context.Context, food model.Food) error {
	return s.foodRepo.Update(ctx, food)
}

func (s *service) DeleteFoodFromMenu(ctx context.Context, foodID uint) error {
	return s.foodRepo.Delete(ctx, foodID)
}

func (s *service) GetAllFoods(ctx context.Context, pagination *common.Pagination, branchID uint) (*common.PaginatedResponse[*model.Food], error) {
	return s.foodRepo.GetAll(ctx, pagination, branchID)
}

func (s *service) GetFoodByID(ctx context.Context, id uint) (*model.Food, error) {
	return s.foodRepo.GetByID(ctx, id)
}

func (s *service) GetImagesByFoodID(ctx context.Context, foodID uint, pagination *common.Pagination) (*common.PaginatedResponse[*model.FoodImage], error) {
	return s.foodRepo.GetImagesByFoodID(ctx, foodID, pagination)
}

func (s *service) AddFoodImageToFood(ctx context.Context, image *model.FoodImage) error {
	return s.foodRepo.CreateImage(ctx, image)
}

func (s *service) DeleteFoodImageFromFood(ctx context.Context, imageID uint) error {
	return s.foodRepo.DeleteImage(ctx, imageID)
}







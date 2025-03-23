package restaurant

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/model"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/port"
)

type service struct {
	restaurantRepo port.RestaurantRepo
	branchRepo     port.BranchRepo
}

func NewService(restaurantRepo port.RestaurantRepo, branchRepo port.BranchRepo) port.Service {
	return &service{
		restaurantRepo: restaurantRepo,
		branchRepo:     branchRepo}
}

func (s *service) CreateRestaurant(ctx context.Context, restaurant model.Restaurant) (uint, error) {
	if err := restaurant.Validate(); err != nil {
		return 0, err
	}
	return s.restaurantRepo.Create(ctx, restaurant)
}

func (s *service) GetRestaurantByFilter(ctx context.Context, filter *model.RestaurantFilter) (*model.Restaurant, error) {
	return s.restaurantRepo.GetByFilter(ctx, filter)
}

func (s *service) CreateBranch(ctx context.Context, branch model.Branch) (uint, error) {
	return s.branchRepo.Create(ctx, branch)
}

func (s *service) GetBranchByFilter(ctx context.Context, filter *model.BranchFilter) (*model.Branch, error) {
	return s.branchRepo.GetByFilter(ctx, filter)
}

func (s *service) GetAllRestaurants(ctx context.Context, pagination *common.Pagination) (*common.PaginatedResponse[*model.Restaurant], error) {
	return s.restaurantRepo.GetAll(ctx, pagination)
}

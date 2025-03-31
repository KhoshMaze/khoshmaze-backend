package port

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/model"
)

type Service interface {
	CreateRestaurant(ctx context.Context, restaurant model.Restaurant) (uint, error)
	GetRestaurantByFilter(ctx context.Context, filter *model.RestaurantFilter) (*model.Restaurant, error)
	GetAllRestaurants(ctx context.Context, pagination *common.Pagination) (*common.PaginatedResponse[*model.Restaurant], error)
	CreateBranch(ctx context.Context, branch model.Branch) (uint, error)
	GetBranchByFilter(ctx context.Context, filter *model.BranchFilter) (*model.Branch, error)
}

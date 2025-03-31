package port

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/model"
)

type RestaurantRepo interface {
	Create(ctx context.Context, restaurant model.Restaurant) (uint, error)
	GetByFilter(ctx context.Context, filter *model.RestaurantFilter) (*model.Restaurant, error)
	GetAll(ctx context.Context, pagination *common.Pagination) (*common.PaginatedResponse[*model.Restaurant], error)
}

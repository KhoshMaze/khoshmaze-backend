package storage

import (
	"context"
	"errors"

	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/cache"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/fp"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/mapper"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/model"
	"gorm.io/gorm"
)

type restaurantRepo struct {
	db    *gorm.DB
	cache cache.Provider
}

func NewRestaurantRepo(db *gorm.DB, cache cache.Provider) *restaurantRepo {
	return &restaurantRepo{db: db, cache: cache}
}

func (r *restaurantRepo) Create(ctx context.Context, restaurantDomain model.Restaurant) (uint, error) {
	restaurant := mapper.RestaurantDomainToStorage(&restaurantDomain)
	return restaurant.ID, r.db.Table("restaurants").WithContext(ctx).Create(restaurant).Error
}

func (r *restaurantRepo) GetByFilter(ctx context.Context, filter *model.RestaurantFilter) (*model.Restaurant, error) {
	var restaurant types.Restaurant

	q := r.db.Table("restaurants").WithContext(ctx)

	if filter.ID > 0 {
		q = q.Where("id = ?", filter.ID)
	}

	if len(filter.URL) > 0 {
		q = q.Where("url = ?", filter.URL)
	}

	err := q.First(&restaurant).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return mapper.RestaurantStorageToDomain(&restaurant), nil
}

func (r *restaurantRepo) GetAll(ctx context.Context, pagination *common.Pagination) (*common.PaginatedResponse[*model.Restaurant], error) {
	var restaurants []types.Restaurant
	var totalItems int64

	if err := r.db.Table("restaurants").WithContext(ctx).Count(&totalItems).Error; err != nil {
		return nil, err
	}

	err := r.db.Table("restaurants").
		WithContext(ctx).
		Offset(pagination.Offset()).
		Limit(pagination.PageSize).
		Find(&restaurants).Error

	if err != nil {
		return nil, err
	}

	restaurantsDomain := fp.Map(restaurants, func(restaurant types.Restaurant) *model.Restaurant {
		return mapper.RestaurantStorageToDomain(&restaurant)
	})

	return common.NewPaginatedResponse(restaurantsDomain, totalItems, pagination.Page, pagination.PageSize), nil
}

// func (r *restaurantRepo) getSubscriptionPrice(ctx context.Context, subType uint8) (*types.SubscriptionPrice, error) {
// 	var subscriptionPrice types.SubscriptionPrice
// 	err := r.db.Table("subscription_prices").
// 		WithContext(ctx).
// 		Where("type = ?", subType).
// 		Order("created_at DESC").
// 		First(&subscriptionPrice).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &subscriptionPrice, nil
// }

package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/cache"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/fp"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/mapper"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/menu/model"
	"gorm.io/gorm"
)

type foodRepo struct {
	db    *gorm.DB
	cache cache.Provider
}

func NewFoodRepo(db *gorm.DB, cache cache.Provider) *foodRepo {
	return &foodRepo{db: db, cache: cache}
}

func (r *foodRepo) Create(ctx context.Context, foodDomain model.Food) (uint, error) {
	food := mapper.FoodDomainToStorage(foodDomain)
	return food.ID, r.db.Table("foods").WithContext(ctx).Create(food).Error
}

func (r *foodRepo) Update(ctx context.Context, foodDomain model.Food) error {
	var f types.Food

	food := mapper.FoodDomainToStorage(foodDomain)
	q := r.db.Table("foods").WithContext(ctx).Where("id = ?", food.ID)

	if err := q.First(&f).Error; err != nil {
		return err
	}

	if err := q.Updates(food).Error; err != nil {
		return err
	}
	if f.Price == 0 || food.Price == 0 {
		return nil
	}
	if f.Price != food.Price {
		r.db.Table("food_prices").WithContext(ctx).Create(&types.FoodPrice{
			FoodID: food.ID,
			Price:  food.Price,
		})
	}

	return nil
}

func (r *foodRepo) Delete(ctx context.Context, id uint) error {
	return r.db.Table("foods").WithContext(ctx).Delete(&types.Food{}, id).Error
}

func (r *foodRepo) GetAll(ctx context.Context, pagination *common.Pagination, menuID uint) (*common.PaginatedResponse[*model.Food], error) {
	var foods []types.Food

	var totalItems int64
	oc := cache.NewObjectCacher[*common.PaginatedResponse[*model.Food]](r.cache, cache.SerializationTypeGob)
	if cached, err := oc.Get(ctx, fmt.Sprintf("menu:%d:foods:page:%d:size:%d", menuID, pagination.Page, pagination.PageSize)); cached != nil && err == nil {
		return cached, nil
	}

	q := r.db.Table("foods").WithContext(ctx)

	if menuID != 0 {
		q = q.Where("menu_id = ?", menuID)
	}

	if err := q.Count(&totalItems).Error; err != nil {
		return nil, err
	}

	if err := q.Offset(pagination.Offset()).
		Limit(pagination.PageSize).
		Find(&foods).Error; err != nil {
		return nil, err
	}

	foodsDomain := fp.Map(foods, func(food types.Food) *model.Food {
		return mapper.FoodStorageToDomain(food)
	})

	response := common.NewPaginatedResponse(foodsDomain, totalItems, pagination.Page, pagination.PageSize)

	oc.Set(ctx, fmt.Sprintf("menu:%d:foods:page:%d:size:%d", menuID, pagination.Page, pagination.PageSize), time.Minute*30, response)

	return response, nil
}

func (r *foodRepo) GetByID(ctx context.Context, id uint) (*model.Food, error) {
	var food types.Food

	oc := cache.NewObjectCacher[model.Food](r.cache, cache.SerializationTypeGob)
	if food, err := oc.Get(ctx, fmt.Sprintf("foods:%d", id)); err == nil {
		return &food, nil
	}

	if err := r.db.Table("foods").WithContext(ctx).Where("id = ?", id).First(&food).Error; err != nil {
		return &model.Food{}, err
	}

	result := mapper.FoodStorageToDomain(food)

	oc.Set(ctx, fmt.Sprintf("foods:%d", id), time.Minute*10, *result)

	return result, nil
}

func (r *foodRepo) GetImagesByFoodID(ctx context.Context, foodID uint, pagination *common.Pagination) (*common.PaginatedResponse[*model.FoodImage], error) {
	var images []types.FoodImage

	var totalItems int64
	if err := r.db.Table("food_images").WithContext(ctx).Where("food_id = ?", foodID).Count(&totalItems).Error; err != nil {
		return nil, err
	}

	if err := r.db.Table("food_images").WithContext(ctx).Where("food_id = ?", foodID).Offset(pagination.Offset()).Limit(pagination.PageSize).Find(&images).Error; err != nil {
		return nil, err
	}

	imagesDomain := fp.Map(images, func(image types.FoodImage) *model.FoodImage {
		return mapper.FoodImageStorageToDomain(&image)
	})

	return common.NewPaginatedResponse(imagesDomain, totalItems, pagination.Page, pagination.PageSize), nil
}

func (r *foodRepo) CreateImage(ctx context.Context, image *model.FoodImage) error {
	return r.db.Table("food_images").WithContext(ctx).Create(image).Error
}

func (r *foodRepo) DeleteImage(ctx context.Context, id uint) error {
	return r.db.Table("food_images").WithContext(ctx).Delete(&types.FoodImage{}, id).Error
}

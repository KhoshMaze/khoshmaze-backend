package storage

import (
	"context"
	"errors"

	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/cache"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/mapper"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	mnuDomain "github.com/KhoshMaze/khoshmaze-backend/internal/domain/menu/model"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/model"
	"gorm.io/gorm"
)

type branchRepo struct {
	db    *gorm.DB
	cache cache.Provider
}

func NewBranchRepo(db *gorm.DB, cache cache.Provider) *branchRepo {
	return &branchRepo{db: db, cache: cache}
}

func (r *branchRepo) Create(ctx context.Context, branchDomain model.Branch) (uint, error) {
	branch := mapper.BranchDomainToStorage(&branchDomain)

	err := r.db.Table("branches").WithContext(ctx).Create(branch).Error

	mnu := NewMenuRepo(r.db)
	menu := &mnuDomain.Menu{
		BranchID: branch.ID,
	}

	if err := mnu.create(ctx, menu); err != nil {
		return 0, err
	}
	return branch.ID, err
}

func (r *branchRepo) GetByFilter(ctx context.Context, filter *model.BranchFilter) (*model.Branch, error) {
	var branch types.Branch
	q := r.db.WithContext(ctx)

	if filter.RestaurantName != "" {
		// q = q.Joins("JOIN restaurants ON restaurants.id = branches.restaurant_id").
		// 	Where("restaurants.url = ?", filter.RestaurantName).Offset(int(filter.ID) - 1).Limit(1)

		// somehow window function is faster than offset limits (even though there's not much data in the table)
		q = q.Raw(`
		WITH ranked_branches AS (
			SELECT b.*, 
				   ROW_NUMBER() OVER (PARTITION BY r.id ORDER BY b.id) as rn
			FROM branches b
			JOIN restaurants r ON r.id = b.restaurant_id
			WHERE r.url = ?
		)
		SELECT * FROM ranked_branches WHERE rn = ?
	`, filter.RestaurantName, filter.ID)

	} else if filter.ID > 0 {
		q.Where("id = ?", filter.ID)
	}

	err := q.First(&branch).Error

	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return mapper.BranchStorageToDomain(&branch), nil
}

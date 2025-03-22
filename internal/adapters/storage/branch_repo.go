package storage

import (
	"context"
	"errors"

	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/cache"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/mapper"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
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
	return branch.ID, r.db.Table("branches").WithContext(ctx).Create(branch).Error
}

func (r *branchRepo) GetByFilter(ctx context.Context, filter *model.BranchFilter) (*model.Branch, error) {
	var branch types.Branch

	q := r.db.Table("branches").WithContext(ctx)

	if filter.ID > 0 {
		q = q.Where("id = ?", filter.ID)
	}

	err := q.First(&branch).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return mapper.BranchStorageToDomain(&branch), nil
}

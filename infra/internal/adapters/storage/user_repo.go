package storage

import (
	"context"
	"errors"

	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/cache"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/mapper"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/port"
	"gorm.io/gorm"
)

type userRepo struct {
	db    *gorm.DB
	cache cache.Provider
}

func NewUserRepo(db *gorm.DB, cacheProvider cache.Provider) port.Repo {
	return &userRepo{db: db, cache: cacheProvider}
}

func (r *userRepo) Create(ctx context.Context, userDomain model.User) (model.UserID, error) {
	user := mapper.UserDomainToStorage(userDomain)
	return model.UserID(user.ID), r.db.Table("users").WithContext(ctx).Create(user).Error
}

func (r *userRepo) IsBannedToken(ctx context.Context, value string) error {
	var token types.TokenBlacklist
	q := r.db.Table("token_blacklists").Where("value = ?", value)
	return q.First(&token).Error

}

func (r *userRepo) CreateBannedToken(ctx context.Context, tokenDomain model.TokenBlacklist) error {
	token := &types.TokenBlacklist{
		ExpiresAt: tokenDomain.ExpiresAt,
		Value:     tokenDomain.Value,
		UserID:    uint(tokenDomain.UserID),
	}

	return r.db.Table("token_blacklists").WithContext(ctx).Create(token).Error
}

func (r *userRepo) GetByFilter(ctx context.Context, filter *model.UserFilter) (*model.User, error) {
	var user types.User

	q := r.db.Table("users").WithContext(ctx)

	if filter.ID > 0 {
		q = q.Where("id = ?", filter.ID)
	}

	if len(filter.Phone) > 0 {
		q = q.Where("phone = ?", filter.Phone)
	}

	err := q.First(&user).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if user.ID == 0 {
		return nil, nil
	}

	return mapper.UserStorageToDomain(user), nil
}

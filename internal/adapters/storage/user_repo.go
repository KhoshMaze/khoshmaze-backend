package storage

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/mapper"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/port"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) port.Repo {
	return &userRepo{db}
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

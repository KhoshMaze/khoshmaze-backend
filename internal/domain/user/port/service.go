package port

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"
)

type Service interface {
	CreateUser(ctx context.Context, user model.User) (model.UserID, error)
	GetUserByFilter(ctx context.Context, filter *model.UserFilter) (*model.User, error)
	IsBannedToken(ctx context.Context, value string) bool
	CreateBannedToken(ctx context.Context, token model.TokenBlacklist) error
}

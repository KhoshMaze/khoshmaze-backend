package port

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"
)

type Repo interface {
	Create(ctx context.Context, user model.User) (model.UserID, error)
	GetByFilter(ctx context.Context, filter *model.UserFilter) (*model.User, error)
	IsBannedToken(ctx context.Context, value string) error
	CreateBannedToken(ctx context.Context, token model.TokenBlacklist) error
}

package handlers

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/api/service"
	"github.com/KhoshMaze/khoshmaze-backend/config"
	"github.com/KhoshMaze/khoshmaze-backend/internal/app"
)

type ServiceGetter[T any] func(context.Context) T

func UserServiceGetter(appContainer app.App, cfg config.ServerConfig) ServiceGetter[*service.UserService] {
	return func(ctx context.Context) *service.UserService {
		return service.NewUserService(appContainer.UserService(ctx),
			cfg.AuthSecret, cfg.RefreshSecret, cfg.AuthExpMinute, cfg.AuthRefreshMinute)
	}
}

package handlers

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/api/service"
	"github.com/KhoshMaze/khoshmaze-backend/config"
	"github.com/KhoshMaze/khoshmaze-backend/internal/app"
)

type ServiceGetter[T any] func(context.Context) T

func UserServiceGetter(appContainer app.App, cfgServer config.ServerConfig) ServiceGetter[*service.UserService] {
	return func(ctx context.Context) *service.UserService {
		return service.NewUserService(appContainer.UserService(ctx),
			cfgServer.AuthSecret, cfgServer.RefreshSecret, cfgServer.AESSecret, cfgServer.AuthExpMinute, cfgServer.AuthRefreshMinute, appContainer.NotificationService(ctx))
	}
}

func RestaurantServiceGetter(appContainer app.App) ServiceGetter[*service.RestaurantService] {
	return func(ctx context.Context) *service.RestaurantService {
		return service.NewRestaurantService(appContainer.RestaurantService(ctx))
	}
}

func MenuServiceGetter(appContainer app.App) ServiceGetter[*service.MenuService] {
	return func(ctx context.Context) *service.MenuService {
		return service.NewMenuService(appContainer.MenuService(ctx))
	}
}

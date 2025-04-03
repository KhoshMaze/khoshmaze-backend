package app

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/api/middlewares"
	"github.com/KhoshMaze/khoshmaze-backend/config"
	notifPort "github.com/KhoshMaze/khoshmaze-backend/internal/domain/notification/port"
	restaurantPort "github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/port"
	userPort "github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/port"
	"gorm.io/gorm"
)

type App interface {
	DB() *gorm.DB
	Config() config.Config
	UserService(ctx context.Context) userPort.Service
	NotificationService(ctx context.Context) notifPort.Service
	RestaurantService(ctx context.Context) restaurantPort.Service
	AnomalyDetectionService() *middlewares.GeoAnomalyService
}

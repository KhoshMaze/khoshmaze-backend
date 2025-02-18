package app

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/config"
	userPort "github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/port"
	"gorm.io/gorm"
)

type App interface {
	DB() *gorm.DB
	Config() config.Config
	UserService(ctx context.Context) userPort.Service
}

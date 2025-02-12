package app

import (
	"github.com/KhoshMaze/khoshmaze-backend/config"
	"gorm.io/gorm"
)

type App interface {
	DB() *gorm.DB
	Config() config.Config
}
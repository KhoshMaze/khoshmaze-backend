package app

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/config"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/postgres"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user"
	userPort "github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/port"
	"gorm.io/gorm"

	appCtx "github.com/KhoshMaze/khoshmaze-backend/internal/adapters/context"
)

type app struct {
	db          *gorm.DB
	cfg         config.Config
	userService userPort.Service
}

func (a *app) DB() *gorm.DB {
	return a.db
}

func (a *app) Config() config.Config {
	return a.cfg
}

func (a *app) setDB() error {
	db, err := postgres.NewPsqlGormConnection(postgres.DBConnOptions{
		User:   a.cfg.DB.User,
		Pass:   a.cfg.DB.Password,
		Host:   a.cfg.DB.Host,
		Port:   a.cfg.DB.Port,
		DBName: a.cfg.DB.Database,
	})

	if err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&types.User{},
		&types.TokenBlacklist{},
	); err != nil {
		return err
	}

	a.db = db
	return nil
}

func (a *app) UserService(ctx context.Context) userPort.Service {
	db := appCtx.GetDB(ctx)

	if db == nil {
		if a.userService == nil {
			a.userService = a.userServiceWithDB(a.db)
		}
		return a.userService
	}

	return a.userServiceWithDB(db)

}

func (a *app) userServiceWithDB(db *gorm.DB) userPort.Service {
	return user.NewService(storage.NewUserRepo(db))
}

func MustNewApp(cfg config.Config) App {
	app, err := newApp(cfg)

	if err != nil {
		panic(err)
	}

	return app
}

func newApp(cfg config.Config) (App, error) {
	a := &app{
		cfg: cfg,
	}

	if err := a.setDB(); err != nil {
		return nil, err
	}

	return a, nil

}

package app

import (
	"context"
	"fmt"

	"github.com/KhoshMaze/khoshmaze-backend/config"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/cache"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/postgres"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/redis"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/notification"
	notifPort "github.com/KhoshMaze/khoshmaze-backend/internal/domain/notification/port"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user"
	userPort "github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/port"
	"github.com/go-co-op/gocron/v2"
	"gorm.io/gorm"

	appCtx "github.com/KhoshMaze/khoshmaze-backend/internal/adapters/context"
)

type app struct {
	db                  *gorm.DB
	cfg                 config.Config
	userService         userPort.Service
	notificationService notifPort.Service
	redisProvider       cache.Provider
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
		&types.Outbox{},
		&types.Notification{},
	); err != nil {
		return err
	}

	a.db = db
	return nil
}

func (a *app) setRedis() {
	a.redisProvider = redis.NewRedisProvider(fmt.Sprintf("%s:%d", a.cfg.Redis.Host, a.cfg.Redis.Port), a.cfg.Redis.Password, 0, "khoshmaze-cache")
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
	return user.NewService(storage.NewUserRepo(db, a.redisProvider))
}

func (a *app) notifServiceWithDB(db *gorm.DB) notifPort.Service {
	return notification.NewService(storage.NewNotificationRepo(db, a.redisProvider),
		user.NewService(storage.NewUserRepo(db, a.redisProvider)), storage.NewOutboxRepo(db), a.cfg.SMS)
}

func (a *app) NotificationService(ctx context.Context) notifPort.Service {
	db := appCtx.GetDB(ctx)
	if db == nil {
		if a.notificationService == nil {
			a.notificationService = a.notifServiceWithDB(a.db)
		}
		return a.notificationService
	}

	return a.notifServiceWithDB(db)
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

	a.setRedis()
	a.registerOutboxHandlers()
	return a, nil

}

func (a *app) registerOutboxHandlers() error {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	common.RegisterOutboxRunner(a.notifServiceWithDB(a.db), scheduler)

	scheduler.Start()

	return nil
}

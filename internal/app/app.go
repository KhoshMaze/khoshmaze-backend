package app

import (
	"github.com/KhoshMaze/khoshmaze-backend/config"
	"gorm.io/gorm"
)

type app struct {
	db  *gorm.DB
	cfg config.Config
}

func (a *app) DB() *gorm.DB {
	return a.db
}

func (a *app) Config() config.Config {
	return a.cfg
}

func MustNewApp(cfg config.Config) App {
	app, err := newApp(cfg)

	if err != nil {
		panic(err)
	}

	return app
}

func newApp(cfg config.Config) (App, error) {
	app := &app{
		cfg: cfg,
	}

	return app, nil

}

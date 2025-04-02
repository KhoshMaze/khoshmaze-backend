package context

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/logger"
	"gorm.io/gorm"
)

var defaultLogger *logger.CustomLogger

func init() {
	defaultLogger = logger.NewLogger()
}

type appContext struct {
	context.Context
	db           *gorm.DB
	shouldCommit bool
	logger       *logger.CustomLogger
}

type AppContextOpt func(*appContext) *appContext // option pattern

func WithDB(db *gorm.DB, shouldCommit bool) AppContextOpt {
	return func(ac *appContext) *appContext {
		ac.db = db
		ac.shouldCommit = shouldCommit
		return ac
	}
}

func WithLogger(logger *logger.CustomLogger) AppContextOpt {
	return func(ac *appContext) *appContext {
		ac.logger = logger
		return ac
	}
}

func NewAppContext(parent context.Context, opts ...AppContextOpt) context.Context {
	ctx := &appContext{Context: parent}
	for _, opt := range opts {
		ctx = opt(ctx)
	}

	return ctx
}

func SetDB(ctx context.Context, db *gorm.DB, shouldCommit bool) {
	appCtx, ok := ctx.(*appContext)
	if !ok {
		return
	}

	appCtx.db = db
	appCtx.shouldCommit = shouldCommit
}

func GetDB(ctx context.Context) *gorm.DB {
	appCtx, ok := ctx.(*appContext)
	if !ok {
		return nil
	}

	return appCtx.db
}

func Commit(ctx context.Context) error {
	appCtx, ok := ctx.(*appContext)
	if !ok || !appCtx.shouldCommit {
		return nil
	}

	return appCtx.db.Commit().Error
}

func Rollback(ctx context.Context) error {
	appCtx, ok := ctx.(*appContext)
	if !ok || !appCtx.shouldCommit {
		return nil
	}

	return appCtx.db.Rollback().Error
}

func CommitOrRollback(ctx context.Context, shouldLog bool) error {
	commitErr := Commit(ctx)
	if commitErr == nil {
		return nil
	}

	if shouldLog {
		logger := GetLogger(ctx)
		logger.Error("error on committing transaction", "err", commitErr.Error())
	}

	if err := Rollback(ctx); err != nil {
		logger := GetLogger(ctx)
		logger.Error("error on rollback transaction", "err", err.Error())
	}

	return commitErr
}

func SetLogger(ctx context.Context, logger *logger.CustomLogger) {
	if appCtx, ok := ctx.(*appContext); ok {
		appCtx.logger = logger
	}
}

func GetLogger(ctx context.Context) *logger.CustomLogger {
	appCtx, ok := ctx.(*appContext)
	if !ok || appCtx.logger == nil {
		return defaultLogger
	}

	return appCtx.logger
}

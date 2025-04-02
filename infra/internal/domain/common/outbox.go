package common

import (
	"context"
	"time"

	appCtx "github.com/KhoshMaze/khoshmaze-backend/internal/adapters/context"
	"github.com/go-co-op/gocron/v2"
)

type OutboxStatus uint8

const (
	OutboxStatusCreated OutboxStatus = iota + 1
	OutboxStatusPicked
	OutboxStatusDone
	OutboxStatusFailed
)

type OutboxType uint8

const (
	OutboxTypeNotif OutboxType = iota + 1
)

type OutboxHandler[T any] interface {
	Handle(ctx context.Context, outboxes []T) error
	Cleanup(ctx context.Context, status OutboxStatus, outboxes []T) error
	Query(ctx context.Context, amount uint, status OutboxStatus) ([]T, error)
	Interval() time.Duration
}

type OutboxRunner[T any] struct {
	handler   OutboxHandler[T]
	scheduler gocron.Scheduler
}

func RegisterOutboxRunner[T any](handler OutboxHandler[T], scheduler gocron.Scheduler) {
	runner := &OutboxRunner[T]{
		handler:   handler,
		scheduler: scheduler,
	}
	runner.register()
}

func (o *OutboxRunner[T]) register() {
	o.scheduler.NewJob(
		gocron.DurationJob(o.handler.Interval()),
		gocron.NewTask(func() { // poller logic
			ctx := context.Background()
			outboxes, err := o.handler.Query(context.Background(), 100, OutboxStatusCreated)
			ctx = appCtx.NewAppContext(ctx)
			logger := appCtx.GetLogger(ctx)

			if err != nil {
				logger.Error("failed to fetch outboxes", "err", err.Error())
				return
			}

			if err := o.handler.Handle(ctx, outboxes); err != nil {
				logger.Error("failed to handle outboxes", "err", err.Error())
			}
		}),
	)

	o.scheduler.NewJob(gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(23,30,0))), gocron.NewTask(func() {

		cleaner := func(status OutboxStatus) {
			ctx := context.Background()
			ctx = appCtx.NewAppContext(ctx)
			logger := appCtx.GetLogger(ctx)
			logger.Info("outbox cleanup started", "status", status)
			for {

				outboxes, err := o.handler.Query(ctx, 1000, status)

				if len(outboxes) <= 0 {
					break
				}

				if err != nil {
					logger.Error("failed to fetch outboxes", "err", err.Error())
					return
				}

				if err := o.handler.Cleanup(ctx, status, outboxes); err != nil {
					logger.Error("failed to cleanup outboxes", "err", err.Error())
					return
				}
			}
			logger.Info("outbox cleanup finished", "status", status)
		}

		cleaner(OutboxStatusDone)
		cleaner(OutboxStatusFailed)
	}))
}

type OutboxID uint

type OutboxRepo interface {
	UpdateStatus(ctx context.Context, status OutboxStatus, id OutboxID) error
	UpdateBulkStatuses(ctx context.Context, status OutboxStatus, ids ...OutboxID) error
	DeleteBulk(ctx context.Context, status OutboxStatus, ids ...OutboxID) error
}

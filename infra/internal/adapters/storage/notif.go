package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/cache"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/mapper"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/notification/model"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/notification/port"
	userDomain "github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"
	"gorm.io/gorm"
)

type notifRepo struct {
	db            *gorm.DB
	cacheProvider cache.Provider
}

func NewNotificationRepo(db *gorm.DB, cacheProvider cache.Provider) port.Repo {
	return &notifRepo{
		db:            db,
		cacheProvider: cacheProvider,
	}
}

func (r *notifRepo) Create(ctx context.Context, notif *model.Notification) (model.NotifID, error) {
	oc := cache.NewObjectCacher[string](r.cacheProvider, cache.SerializationTypeJSON)
	if notif.ForAuthorization {
		if err := oc.Set(ctx, fmt.Sprintf("notifs:%s", notif.Phone), notif.TTL, notif.Content); err != nil {
			return 0, err
		}

		return 0, nil
	}

	no := mapper.Notification2Storage(notif)
	if err := r.db.WithContext(ctx).Table("notifications").Create(no).Error; err != nil {
		return 0, err
	}
	return model.NotifID(no.ID), nil
}

func (r *notifRepo) CreateOutbox(ctx context.Context, no *model.NotificationOutbox) error {
	outbox, err := mapper.NotifOutbox2Storage(no)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Table("outboxes").Create(outbox).Error
}

func (r *notifRepo) QueryOutboxes(ctx context.Context, limit uint, status common.OutboxStatus) ([]model.NotificationOutbox, error) {
	var outboxes []types.Outbox

	err := r.db.WithContext(ctx).Table("outboxes").
		Where(`"type" = ?`, common.OutboxTypeNotif).
		Where("status = ?", status).
		Limit(int(limit)).Scan(&outboxes).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	result := make([]model.NotificationOutbox, len(outboxes))

	for i := range outboxes {
		v, err := mapper.OutboxStorage2Notif(outboxes[i])
		if err != nil {
			return nil, err
		}
		result[i] = v
	}

	return result, nil
}

func (r *notifRepo) GetUserNotifValue(ctx context.Context, phone userDomain.Phone) (string, error) {
	oc := cache.NewObjectCacher[string](r.cacheProvider, cache.SerializationTypeJSON)
	v, err := oc.Get(ctx, fmt.Sprintf("notifs:%s", phone))
	return v, err
}

func (r *notifRepo) DeleteUserNotifValue(ctx context.Context, phone userDomain.Phone) error {
	oc := cache.NewObjectCacher[string](r.cacheProvider, cache.SerializationTypeJSON)
	return oc.Del(ctx, fmt.Sprintf("notifs:%s", phone))
}

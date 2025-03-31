package port

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/notification/model"
	userDomain "github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"
)

type Repo interface {
	Create(ctx context.Context, notif *model.Notification) (model.NotifID, error)
	CreateOutbox(ctx context.Context, outbox *model.NotificationOutbox) error
	QueryOutboxes(ctx context.Context, limit uint, status common.OutboxStatus) ([]model.NotificationOutbox, error)
	GetUserNotifValue(ctx context.Context, phone userDomain.Phone) (string, error)
	DeleteUserNotifValue(ctx context.Context, phone userDomain.Phone) error
}

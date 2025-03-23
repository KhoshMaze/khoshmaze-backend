package port

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/notification/model"
	userDomain "github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"
)

type Service interface {
	Send(ctx context.Context, notif *model.Notification) error
	CheckUserNotifValue(ctx context.Context, phone userDomain.Phone, val string) (bool, error)
	DeleteUserNotifValue(ctx context.Context, phone userDomain.Phone) error
	common.OutboxHandler[model.NotificationOutbox]
}

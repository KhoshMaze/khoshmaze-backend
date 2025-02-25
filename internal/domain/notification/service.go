package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/KhoshMaze/khoshmaze-backend/config"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/fp"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/notification"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/notification/model"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/notification/port"
	userDomain "github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"
	userPort "github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/port"
)

type service struct {
	repo       port.Repo
	outboxRepo common.OutboxRepo
	userPort   userPort.Service
	cfg        config.SMSConfig
}

func NewService(repo port.Repo, userPort userPort.Service, outboxRepo common.OutboxRepo, cfg config.SMSConfig) port.Service {
	return &service{
		repo:       repo,
		userPort:   userPort,
		outboxRepo: outboxRepo,
		cfg:        cfg,
	}
}

func (s *service) Send(ctx context.Context, notif *model.Notification) error {
	// var (
	// 	user *userDomain.User
	// 	err  error
	// )
	// if !notif.ForAuthorization {
	// 	user, err = s.userPort.GetUserByFilter(ctx, &userDomain.UserFilter{
	// 		ID: notif.UserID,
	// 	})

	// 	if err != nil {
	// 		return err
	// 	}

	// }
	notifID, err := s.repo.Create(ctx, notif)
	if err != nil {
		return err
	}
	return s.repo.CreateOutbox(ctx, &model.NotificationOutbox{
		NotifID: notifID,
		Data: model.OutboxData{
			Dest: func() string {
				switch notif.Type {
				case model.NotifTypeSMS:
					return string(notif.Phone)
				default:
					return ""
				}
			}(),
			Content: notif.Content,
			Type:    notif.Type,
		},
		Status: common.OutboxStatusCreated,
		Type:   common.OutboxTypeNotif,
	})
}

func (s *service) Handle(ctx context.Context, outboxes []model.NotificationOutbox) error {
	outBoxIDs := fp.Map(outboxes, func(o model.NotificationOutbox) common.OutboxID {
		return o.OutboxID
	})

	if err := s.outboxRepo.UpdateBulkStatuses(ctx, common.OutboxStatusPicked, outBoxIDs...); err != nil {
		return fmt.Errorf("failed to update notif outbox statuses to picked %w", err)
	}

	for _, outbox := range outboxes {
		fmt.Printf("dest : %s, content : %s\n", outbox.Data.Dest, outbox.Data.Content)
		go notification.SendSMS(&outbox.Data, s.cfg)
	}

	if err := s.outboxRepo.UpdateBulkStatuses(ctx, common.OutboxStatusDone, outBoxIDs...); err != nil {
		return fmt.Errorf("failed to update notif outbox statuses to done %w", err)
	}

	return nil
}

func (s *service) Interval() time.Duration {
	return time.Second * 10
}

func (s *service) Query(ctx context.Context) ([]model.NotificationOutbox, error) {
	return s.repo.QueryOutboxes(ctx, 100, common.OutboxStatusCreated)
}

func (s *service) CheckUserNotifValue(ctx context.Context, phone userDomain.Phone, val string) (bool, error) {
	expected, err := s.repo.GetUserNotifValue(ctx, phone)
	if err != nil {
		return false, err
	}

	return expected == val, nil
}

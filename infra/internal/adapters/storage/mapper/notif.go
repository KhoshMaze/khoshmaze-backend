package mapper

import (
	"encoding/json"

	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/notification/model"
)

func NotifOutbox2Storage(no *model.NotificationOutbox) (*types.Outbox, error) {
	data, err := json.Marshal(&no.Data)
	if err != nil {
		return nil, err
	}

	return &types.Outbox{
		Data:   data,
		RefID:  uint(no.NotifID),
		Type:   uint8(no.Type),
		Status: uint8(no.Status),
	}, nil
}

func Notification2Storage(no *model.Notification) *types.Notification {
	return &types.Notification{
		Content: no.Content,
		To:      uint(no.UserID),
		Type:    uint8(no.Type),
	}
}

func OutboxStorage2Notif(outbox types.Outbox) (model.NotificationOutbox, error) {
	var outboxData model.OutboxData
	err := json.Unmarshal([]byte(outbox.Data), &outboxData)
	if err != nil {
		return model.NotificationOutbox{}, err
	}

	return model.NotificationOutbox{
		OutboxID: common.OutboxID(outbox.ID),
		NotifID:  model.NotifID(outbox.RefID),
		Data:     outboxData,
		Status:   common.OutboxStatus(outbox.Status),
		Type:     common.OutboxType(outbox.Type),
	}, nil
}

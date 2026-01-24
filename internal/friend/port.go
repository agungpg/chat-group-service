package friend

import (
	"context"

	notif "github.com/agungpg/group-chat-service/internal/notification"
)

type NotificationService interface {
	SendPushNotification(ctx context.Context, payload notif.PushNotificationPayload) error
}

type UserRepository interface {
	GetActiveUserDeviceTokens(ctx context.Context, userID string) (string, error)
}

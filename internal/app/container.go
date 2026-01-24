package app

import (
	"github.com/agungpg/group-chat-service/internal/auth"
	"github.com/agungpg/group-chat-service/internal/friend"
	"github.com/agungpg/group-chat-service/internal/notification"
	"github.com/agungpg/group-chat-service/internal/profile"
	"github.com/uptrace/bun"
)

type Container struct {
	AuthHandler         *auth.Handler
	ProfileHandler      *profile.Handler
	FriendHandler       *friend.Handler
	NotificationHandler *notification.Handler
	// ConversationHandler *conversation.Handler
	// InviteHandler       *invite.Handler
	// ChatWsHandler fiber.Handler // or *chat.WsHandler
	AuthMiddlerware *auth.MiddleWare
}

func NewContainer(db *bun.DB) *Container {
	authM := auth.NewMiddleware()
	authRepo := auth.NewRepository(db)
	profileRepo := profile.NewRepository(db)
	friendRepo := friend.NewRepository(db)

	authSvc := auth.NewService(authRepo, profileRepo)
	profileSvc := profile.NewService(profileRepo)
	notificationSvc := notification.NewService()
	friendSvc := friend.NewService(friendRepo, authRepo, notificationSvc)

	authH := auth.NewHandler(authSvc)
	profileH := profile.NewHandler(profileSvc)
	friendH := friend.NewHandler(friendSvc)
	notificationH := notification.NewHandler(notificationSvc)

	return &Container{
		AuthHandler:         authH,
		ProfileHandler:      profileH,
		FriendHandler:       friendH,
		AuthMiddlerware:     authM,
		NotificationHandler: notificationH,
	}

}

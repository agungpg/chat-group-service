package app

import (
	"github.com/agungpg/group-chat-service/config"
	"github.com/agungpg/group-chat-service/internal/auth"
	"github.com/agungpg/group-chat-service/internal/files"
	"github.com/agungpg/group-chat-service/internal/friend"
	"github.com/agungpg/group-chat-service/internal/notification"
	"github.com/agungpg/group-chat-service/internal/profile"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/uptrace/bun"
)

type Container struct {
	AuthHandler         *auth.Handler
	ProfileHandler      *profile.Handler
	FriendHandler       *friend.Handler
	NotificationHandler *notification.Handler
	FilesHandler        *files.Handler
	// ConversationHandler *conversation.Handler
	// InviteHandler       *invite.Handler
	// ChatWsHandler fiber.Handler // or *chat.WsHandler
	AuthMiddlerware *auth.MiddleWare
}

func NewContainer(db *bun.DB, s3Client *s3.Client, strgConf *config.StorageConfig) *Container {
	authM := auth.NewMiddleware()
	authRepo := auth.NewRepository(db)
	profileRepo := profile.NewRepository(db)
	friendRepo := friend.NewRepository(db)
	filesRepo := files.NewRepository(db)

	authSvc := auth.NewService(authRepo, profileRepo)
	notificationSvc := notification.NewService()
	friendSvc := friend.NewService(friendRepo, authRepo, notificationSvc)
	filesAdapter := files.NewS3Adapter(s3Client)
	filesSvc := files.NewService(filesRepo, filesAdapter, strgConf)
	profileSvc := profile.NewService(profileRepo, filesSvc)

	authH := auth.NewHandler(authSvc)
	profileH := profile.NewHandler(profileSvc)
	friendH := friend.NewHandler(friendSvc)
	notificationH := notification.NewHandler(notificationSvc)
	filesH := files.NewHandler(filesSvc)

	return &Container{
		AuthHandler:         authH,
		ProfileHandler:      profileH,
		FriendHandler:       friendH,
		AuthMiddlerware:     authM,
		NotificationHandler: notificationH,
		FilesHandler:        filesH,
	}

}

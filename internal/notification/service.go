package notification

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) SendPushNotification(ctx context.Context, payload PushNotificationPayload) error {

	fileName := os.Getenv("FIREBASE_CREDENTIALS")
	projectName := os.Getenv("FIREBASE_PROJECT_NAME")
	wd, _ := os.Getwd()
	credPath := filepath.Join(wd, "internal", "config", "firebase", fileName)

	if _, err := os.Stat(credPath); err != nil {
		return fmt.Errorf("credentials file not found at %s: %w", credPath, err)
	}

	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: projectName,
	}, option.WithCredentialsFile(credPath))
	if err != nil {
		return err
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		fmt.Println("clienterr : ", err)
		return err
	}

	msg := &messaging.Message{
		Token: payload.Token,
		Notification: &messaging.Notification{
			Title: payload.Title,
			Body:  payload.Body,
		},
		Data: map[string]string{
			"type": payload.Data.Type,
			"id":   payload.Data.Id,
		},
	}

	id, err := client.Send(ctx, msg)
	if err != nil {
		return fmt.Errorf("send push failed: %w", err)
	}

	return nil
}

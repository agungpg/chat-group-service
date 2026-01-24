package friend

import (
	"context"
	"fmt"

	notif "github.com/agungpg/group-chat-service/internal/notification"
	"github.com/google/uuid"
)

type Service struct {
	repo     *Repository
	userRepo UserRepository
	notif    NotificationService
}

func NewService(repo *Repository, userRepo UserRepository, notif NotificationService) *Service {
	return &Service{repo: repo, userRepo: userRepo, notif: notif}
}

func (s *Service) SendFriendRequest(ctx context.Context, requesterID, addresseeID string) error {
	fr := FriendRequest{
		ID:          uuid.New().String(),
		RequesterID: requesterID,
		AddresseeID: addresseeID,
		Status:      "pending",
	}

	err := s.repo.CreateFriendRequest(ctx, fr)
	if err != nil {
		return err
	}

	token, err := s.userRepo.GetActiveUserDeviceTokens(ctx, addresseeID)
	if err != nil || token == "" {
		return err
	}

	payload := notif.PushNotificationPayload{
		Token: token,
		Title: "New friend request",
		Body:  "You have a new friend request.",
		Data: notif.PushNotificationData{
			Type: "friend_request",
			Id:   fr.ID,
		},
	}

	s.notif.SendPushNotification(ctx, payload)
	return err
}

func (s *Service) GetFriendRequestList(ctx context.Context, userId, requestType string, limit, offset int) ([]FriendRequestDTO, int, error) {

	var (
		friendRequestList []FriendRequest
		total             int
		err               error
	)

	if requestType == "incoming" {
		friendRequestList, total, err = s.repo.GetIncomingRequestList(ctx, userId, limit, offset)
	} else {
		friendRequestList, total, err = s.repo.GetOutgoingRequestsList(ctx, userId, limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	userIds := getFriendRequestUserId(friendRequestList, requestType)
	userSummaries, err := s.repo.GetUserSummaries(ctx, userIds)
	if err != nil {
		return nil, 0, err
	}

	dtos := make([]FriendRequestDTO, 0, len(friendRequestList))
	for _, fr := range friendRequestList {
		var id string

		if requestType == "incoming" {
			id = fr.RequesterID
		} else {
			id = fr.AddresseeID
		}

		usr := userSummaries[id]
		dtos = append(dtos, FriendRequestDTO{
			ID:          fr.ID,
			UserID:      id,
			UserName:    usr.Username,
			DisplayName: usr.DisplayName,
			AvatarUrl:   usr.AvatarURL,
			Status:      string(fr.Status),
		})
	}

	return dtos, total, nil
}

func getFriendRequestUserId(frList []FriendRequest, requestType string) []string {
	ids := make([]string, 0, len(frList))
	seen := make(map[string]struct{}, len(frList))
	for _, fr := range frList {

		if requestType == "incoming" {
			if _, ok := seen[fr.RequesterID]; ok {
				continue
			}
			seen[fr.RequesterID] = struct{}{}
			ids = append(ids, fr.RequesterID)

		} else {
			if _, ok := seen[fr.AddresseeID]; ok {
				continue
			}
			seen[fr.AddresseeID] = struct{}{}
			ids = append(ids, fr.AddresseeID)
		}
	}

	fmt.Println("fr ids: ", ids)
	return ids
}

package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/agungpg/group-chat-service/internal/profile"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo    *Repository
	profile ProfileCreator
}

func NewService(repo *Repository, profile ProfileCreator) *Service {
	return &Service{repo, profile}
}

func (s *Service) Register(ctx context.Context, payload RegistrationPayload) error {

	fmt.Println("Register service is running")
	tx, err := s.repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if rErr := tx.Rollback(); rErr != nil && rErr != sql.ErrTxDone {
			fmt.Printf("rollback error: %v", rErr)
		}
	}()

	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	user := &User{
		ID:        uuid.New().String(),
		Username:  payload.Username,
		Email:     payload.Email,
		Password:  string(hashedPwd),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repo.CreateUser(ctx, user, &tx)
	if err != nil {
		return err
	}

	profile := &profile.UserProfile{
		UserID:      user.ID,
		DisplayName: payload.DisplayName,
		AvatarURL:   payload.AvatarURL,
		Bio:         payload.Bio,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	err = s.profile.CreateProfile(ctx, profile, &tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Service) Login(ctx context.Context, payload LoginPayload) (string, error) {
	user, err := s.repo.FindByUsername(ctx, payload.Username)
	if err != nil {
		return "", errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return "", errors.New("invalid password")
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}

func (s *Service) RegisterUserDevice(ctx context.Context, payload DeviceRegistrationPayload) error {
	userDevice := &UserDevice{
		ID:          uuid.New().String(),
		UserID:      payload.UserID,
		DeviceID:    payload.DeviceID,
		Platform:    payload.Platform,
		Provider:    payload.Provider,
		Token:       payload.Token,
		AppVersion:  payload.AppVersion,
		DeviceModel: payload.DeviceModel,
		LastSeenAt:  time.Now(),
	}

	err := s.repo.RegiserUserDevice(ctx, userDevice)

	return err
}

func (s *Service) UnRegisterUserDevice(ctx context.Context, deviceId string) error {
	err := s.repo.UnRegisterUserDevice(ctx, deviceId)
	return err
}

func (s *Service) CheckIsUsernameOrEmailTaken(ctx context.Context, payload CheckUsernameAndEmailDTO) (interface{}, error) {
	users, err := s.repo.GetUserByEmailOrUsername(ctx, payload.Username, payload.Email)

	availability := map[string]bool{
		"isEmailToken":    false,
		"isUsernameTaken": false,
	}

	for i := 0; i < len(users); i++ {
		if users[i].Email == payload.Email {
			availability["isEmailToken"] = true
		}
		if users[i].Username == payload.Username {
			availability["isUsernameTaken"] = true
		}
	}

	return availability, err
}

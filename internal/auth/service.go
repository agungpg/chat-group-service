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
	fmt.Println("CreateUser success")
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

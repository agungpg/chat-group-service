package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) Register(ctx context.Context, payload RegistrationPayload) error {
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	user := &User{
		ID:        uuid.New().String(),
		Username:  payload.Username,
		Email:     payload.Email,
		Password:  string(hashedPwd),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.CreateUser(ctx, user)
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

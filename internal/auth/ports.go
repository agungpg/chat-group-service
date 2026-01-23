package auth

import (
	"context"

	"github.com/agungpg/group-chat-service/internal/profile"
	"github.com/uptrace/bun"
)

type ProfileCreator interface {
	CreateProfile(ctx context.Context, profile *profile.UserProfile, tx *bun.Tx) error
}

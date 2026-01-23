package profile

import (
	"time"

	"github.com/uptrace/bun"
)

// UserProfile maps to the user_profiles table.
type UserProfile struct {
	bun.BaseModel `bun:"table:user_profiles"`

	UserID      string    `bun:"user_id,pk,notnull"`
	DisplayName string    `bun:"display_name,notnull"`
	AvatarURL   string    `bun:"avatar_url,nullzero"`
	Bio         string    `bun:"bio,nullzero"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

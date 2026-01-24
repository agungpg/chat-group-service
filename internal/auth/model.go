package auth

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            string    `bun:"id,pk,unique,notnull"`
	Username      string    `bun:"username,unique,notnull"`
	Password      string    `bun:"password,notnull"`
	Email         string    `bun:"email,unique,notnull"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

type UserDevice struct {
	bun.BaseModel `bun:"table:user_devices"`

	ID          string    `bun:"id,pk,notnull"`
	UserID      string    `bun:"user_id,nullzero"`
	DeviceID    string    `bun:"device_id,nullzero"`
	Platform    string    `bun:"platform,notnull"` // ios | android
	Provider    string    `bun:"provider,notnull"` // fcm | apns
	Token       string    `bun:"token,notnull"`
	AppVersion  string    `bun:"app_version,nullzero"`
	DeviceModel string    `bun:"device_model,nullzero"`
	LastSeenAt  time.Time `bun:"last_seen_at,nullzero"`
	RevokedAt   time.Time `bun:"revoked_at,nullzero"`
	IsActive    bool      `bun:"is_active,notnull,default:true"`
	IsDeleted   bool      `bun:"is_deleted,notnull,default:false"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

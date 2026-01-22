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

package friend

import (
	"time"

	"github.com/uptrace/bun"
)

// FriendRequestStatus mirrors the public.friend_request_status enum.
type FriendRequestStatus string

const (
	FriendRequestStatusPending   FriendRequestStatus = "pending"
	FriendRequestStatusAccepted  FriendRequestStatus = "accepted"
	FriendRequestStatusDeclined  FriendRequestStatus = "declined"
	FriendRequestStatusCancelled FriendRequestStatus = "cancelled"
)

// FriendRequest maps to the public.friend_requests table.
type FriendRequest struct {
	bun.BaseModel `bun:"table:friend_requests"`

	ID          string              `bun:"id,pk,default:gen_random_uuid()"`
	RequesterID string              `bun:"requester_id,notnull"`
	AddresseeID string              `bun:"addressee_id,notnull"`
	Status      FriendRequestStatus `bun:"status,notnull,default:pending"`
	CreatedAt   time.Time           `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   time.Time           `bun:"updated_at,notnull,default:current_timestamp"`
}

// Friendship maps to the public.friendships table.
type Friendship struct {
	bun.BaseModel `bun:"table:friendships"`

	UserLowID  string    `bun:"user_low_id,pk,notnull"`
	UserHighID string    `bun:"user_high_id,pk,notnull"`
	CreatedAt  time.Time `bun:"created_at,notnull,default:current_timestamp"`
}

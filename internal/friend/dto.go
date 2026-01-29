package friend

type FriendRequestDTO struct {
	ID          string `json:"id" bun:"id"`
	UserID      string `json:"userId" bun:"user_id"`
	UserName    string `json:"userName" bun:"user_name"`
	DisplayName string `json:"displayName" bun:"user_name"`
	AvatarUrl   string `json:"avatarUrl" bun:"avatar_url"`
	Status      string `json:"status" bun:"status"`
}

type SendFriendRequestPayload struct {
	AddresseeID string `json:"userId"`
}

type ActionFriendRequestPayload struct {
	ID string `json:"id"`
}

type RemoveFriendPayload struct {
	UserID string `json:"userId"`
}

type FriendUser struct {
	ID          string `json:"id"`
	UserName    string `json:"userName"`
	DisplayName string `json:"displayName"`
	AvatarUrl   string `json:"avatarUrl"`
}

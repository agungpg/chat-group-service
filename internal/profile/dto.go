package profile

type ProfileResponse struct {
	UserID      string `json:"id"`
	DisplayName string `json:"name"`
	AvatarURL   string `json:"avatar"`
	Bio         string `json:"bio"`
}

package profile

type ProfileResponse struct {
	UserID      string `json:"id"`
	DisplayName string `json:"name"`
	AvatarURL   string `json:"avatar"`
	Bio         string `json:"bio"`
}

type UpdateProfileRequest struct {
	UserID      string `json:"id"`
	DisplayName string `json:"name,omitempty"`
	AvatarURL   string `json:"avatar,omitempty"`
	Bio         string `json:"bio,omitempty"`
}

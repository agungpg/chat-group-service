package auth

type RegistrationPayload struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName,omitempty"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
	Bio         string `json:"bio,omitempty"`
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

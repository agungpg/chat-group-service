package auth

type RegistrationPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

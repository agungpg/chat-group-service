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

type DeviceRegistrationPayload struct {
	UserID      string `json:"user_id,omitempty"`
	DeviceID    string `json:"device_id"`
	Platform    string `json:"platform"` // e.g., "ios", "android"
	Provider    string `json:"provider"` // e.g., "fcm", "apns"
	Token       string `json:"token"`    // The push notification token
	AppVersion  string `json:"app_version,omitempty"`
	DeviceModel string `json:"device_model,omitempty"` // e.g., "iPhone 15 Pro"
}

type UnRegisterDevicePayload struct {
	DeviceID string `json:"device_id"`
}

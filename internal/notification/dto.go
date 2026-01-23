package notification

type PushNotificationData struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type PushNotificationPayload struct {
	Token string `json:"token"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Data  PushNotificationData
}

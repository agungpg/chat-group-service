package app

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, c *Container) {
	app.Post("/auth/register", c.AuthHandler.Register)
	app.Post("/auth/login", c.AuthHandler.Login)
	app.Post("/auth/devices/register", c.AuthMiddlerware.JWTMiddleware(), c.AuthHandler.RegisterUserDevice)
	app.Post("/auth/devices/unregister", c.AuthMiddlerware.JWTMiddleware(), c.AuthHandler.UnRegisterUserDevice)

	app.Get("/profile", c.AuthMiddlerware.JWTMiddleware(), c.ProfileHandler.GetProfile)
	app.Put("/profile", c.AuthMiddlerware.JWTMiddleware(), c.ProfileHandler.UpdateProfile)

	app.Post("/friend/send-request", c.AuthMiddlerware.JWTMiddleware(), c.FriendHandler.SendFriendRequest)
	app.Get("/friend/requests", c.AuthMiddlerware.JWTMiddleware(), c.FriendHandler.GetIncomingRequestList)
	app.Post("/friend/accept-request", c.AuthMiddlerware.JWTMiddleware(), c.FriendHandler.AcceptRequest)

	app.Post("/notification/push", c.AuthMiddlerware.JWTMiddleware(), c.NotificationHandler.PushNotification)
}

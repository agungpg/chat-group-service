package app

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, c *Container) {
	app.Post("/auth/register", c.AuthHandler.Register)
	app.Post("/auth/login", c.AuthHandler.Login)
	app.Post("/auth/devices/register", c.AuthMiddlerware.JWTMiddleware(), c.AuthHandler.RegisterUserDevice)
	app.Post("/auth/devices/unregister", c.AuthMiddlerware.JWTMiddleware(), c.AuthHandler.UnRegisterUserDevice)

	app.Get("/profile", c.AuthMiddlerware.JWTMiddleware(), c.ProfileHandler.GetProfile)
	app.Put("/profile", c.AuthMiddlerware.JWTMiddleware(), c.ProfileHandler.UpdateProfile)

	app.Post("/notification/push", c.AuthMiddlerware.JWTMiddleware(), c.NotificationHandler.PushNotification)
}

package app

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, c *Container) {
	app.Post("/auth/register", c.AuthHandler.Register)
	app.Post("/auth/login", c.AuthHandler.Login)

	app.Get("/profile", c.AuthMiddlerware.JWTMiddleware(), c.ProfileHandler.GetProfile)
	app.Put("/profile", c.AuthMiddlerware.JWTMiddleware(), c.ProfileHandler.UpdateProfile)
}

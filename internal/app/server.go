package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewServer(c *Container) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "fiber-chat-app",
	})

	// --- global middlewares ---
	app.Use(recover.New())

	// Log requests (good for dev; tune for prod)
	app.Use(logger.New())

	// CORS (adjust origins for production)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // change to your frontend domain later
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PATCH,DELETE,OPTIONS",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})

	// --- routes ---
	RegisterRoutes(app, c)

	return app
}

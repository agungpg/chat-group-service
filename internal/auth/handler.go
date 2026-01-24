package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Change parameter type to fiber.Router to match project handler pattern
func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/register", h.Register)
	router.Post("/login", h.Login)
}

func (h *Handler) Register(c *fiber.Ctx) error {
	fmt.Printf("Register handler is running")
	var payload RegistrationPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload",
		})
	}
	err := h.service.Register(c.Context(), payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var payload LoginPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
		})
	}

	token, err := h.service.Login(c.Context(), payload)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":   token,
		"message": "success",
	})
}

func (h *Handler) RegisterUserDevice(c *fiber.Ctx) error {
	var payload DeviceRegistrationPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}
	claims := c.Locals("userClaims").(jwt.MapClaims)
	userID := claims["id"].(string)

	payload.UserID = userID
	err := h.service.RegisterUserDevice(c.Context(), payload)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "user device registered!",
	})
}

func (h *Handler) UnRegisterUserDevice(c *fiber.Ctx) error {
	var payload UnRegisterDevicePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}
	err := h.service.UnRegisterUserDevice(c.Context(), payload.DeviceID)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "user device registered!",
	})
}

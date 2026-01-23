package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetProfile(c *fiber.Ctx) error {

	claims := c.Locals("userClaims").(jwt.MapClaims)
	userID := claims["id"].(string)
	profile, err := h.service.GetProfileByUserId(c.Context(), userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":    profile,
		"status":  200,
		"message": nil,
	})
}

func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	var payload UpdateProfileRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
		})
	}

	err := h.service.UpdateProfile(c.Context(), payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Update profile is success!",
	})
}

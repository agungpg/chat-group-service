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
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":   profile,
		"status": 200,
		"error":  nil,
	})
}

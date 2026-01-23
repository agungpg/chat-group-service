package notification

import "github.com/gofiber/fiber/v2"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) PushNotification(c *fiber.Ctx) error {
	var payload PushNotificationPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload",
		})
	}
	err := h.service.SendPushNotification(c.Context(), payload)
	if err != nil {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Failed to send notification",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Notification sent",
	})
}

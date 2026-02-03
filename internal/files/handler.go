package files

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

func (h *Handler) PresignUploadUrl(c *fiber.Ctx) error {
	var body UploadPresignRequestDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body payload",
		})
	}

	claims := c.Locals("userClaims").(jwt.MapClaims)
	userId := claims["id"].(string)

	res, err := h.service.CreatePresignUpload(c.Context(), body, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    res,
		"message": "ok",
	})
}

func (h *Handler) PresignViewUrl(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := h.service.GetPresignView(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    res,
		"message": "ok",
	})
}

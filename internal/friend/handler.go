package friend

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SendFriendRequest(c *fiber.Ctx) error {
	var payload SendFriendRequestPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
		})
	}
	claims := c.Locals("userClaims").(jwt.MapClaims)
	requesterID := claims["id"].(string)
	err := h.service.SendFriendRequest(c.Context(), requesterID, payload.AddresseeID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "request sent succesfully",
	})
}

func (h *Handler) GetIncomingRequestList(c *fiber.Ctx) error {
	page := 1
	if pageParam := c.Query("page"); pageParam != "" {
		if parsed, err := strconv.Atoi(pageParam); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := 20
	limitParam := c.Query("size")
	if limitParam != "" {
		if parsed, err := strconv.Atoi(limitParam); err == nil && parsed > 0 {
			// cap limit to avoid fetching excessively large pages
			if parsed > 100 {
				limit = 100
			} else {
				limit = parsed
			}
		}
	}

	offset := (page - 1) * limit

	requestType := c.Query("type")
	if requestType == "" {
		requestType = "incoming"
	}

	claims := c.Locals("userClaims").(jwt.MapClaims)
	userId := claims["id"].(string)
	incomingList, total, err := h.service.GetFriendRequestList(c.Context(), userId, requestType, limit, offset)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    incomingList,
		"message": "ok",
		"pagination": fiber.Map{
			"pageNumber": page,
			"size":       limit,
			"totalData":  total,
		},
	})
}

func (h *Handler) AcceptRequest(c *fiber.Ctx) error {
	var payload ActionFriendRequestPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
		})
	}

	claims := c.Locals("userClaims").(jwt.MapClaims)
	userId := claims["id"].(string)

	err := h.service.AcceptRequest(c.Context(), userId, payload.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "friend request accepted!",
	})
}

func (h *Handler) DeclineRequest(c *fiber.Ctx) error {
	var payload ActionFriendRequestPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
		})
	}

	err := h.service.DeclineRequest(c.Context(), payload.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "friend request declined!",
	})
}

func (h *Handler) CancelRequest(c *fiber.Ctx) error {
	var payload ActionFriendRequestPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
		})
	}

	claims := c.Locals("userClaims").(jwt.MapClaims)
	userId := claims["id"].(string)

	if err := h.service.CancelRequest(c.Context(), userId, payload.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "friend request cancelled!",
	})
}

func (h *Handler) RemoveFriend(c *fiber.Ctx) error {
	var payload RemoveFriendPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
		})
	}

	claims := c.Locals("userClaims").(jwt.MapClaims)
	userId := claims["id"].(string)

	if err := h.service.RemoveFriend(c.Context(), userId, payload.UserID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "friend removed!",
	})
}

func (h *Handler) ListFriends(c *fiber.Ctx) error {
	page := 1
	if pageParam := c.Query("page"); pageParam != "" {
		if parsed, err := strconv.Atoi(pageParam); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := 20
	if limitParam := c.Query("size"); limitParam != "" {
		if parsed, err := strconv.Atoi(limitParam); err == nil && parsed > 0 {
			if parsed > 100 {
				limit = 100
			} else {
				limit = parsed
			}
		}
	}

	offset := (page - 1) * limit

	claims := c.Locals("userClaims").(jwt.MapClaims)
	userId := claims["id"].(string)

	friends, total, err := h.service.ListFriends(c.Context(), userId, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    friends,
		"message": "ok",
		"pagination": fiber.Map{
			"pageNumber": page,
			"size":       limit,
			"totalData":  total,
		},
	})
}

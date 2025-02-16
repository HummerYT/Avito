package info

import (
	"github.com/gofiber/fiber/v2"

	"AvitoTask/internal/models"
)

type Handler struct {
	uc infoUser
}

func NewHandler(uc infoUser) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Handle(c *fiber.Ctx) error {
	userID, ok := c.Locals("UserID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": models.ErrAuthUser.Error(),
		})
	}

	username, infoResp, err := h.uc.GetInfo(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(ConvertInfoResponse(infoResp, userID, username))
}

package send_coin

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"AvitoTask/internal/models"
	"AvitoTask/internal/usecase/send_coin"
)

type Handler struct {
	sender sender
}

func NewHandler(s sender) *Handler {
	return &Handler{
		sender: s,
	}
}

func (h *Handler) Handle(ctx *fiber.Ctx) error {
	fromUser, ok := ctx.Context().Value("UserID").(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errors": models.ErrAuthUser.Error(),
		})
	}

	var req request
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}

	if err := validate(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}

	err := h.sender.SendCoin(ctx.Context(), fromUser, req.ToUser, req.Amount)
	if errors.Is(err, send_coin.ErrNotEnoughCoins) || errors.Is(err, send_coin.ErrSameUser) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}

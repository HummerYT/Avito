package buy_item

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"AvitoTask/internal/models"
	"AvitoTask/internal/usecase/buy_item"
)

type Handler struct {
	buyer buyer
}

func NewHandler(b buyer) *Handler {
	return &Handler{
		buyer: b,
	}
}

func (h *Handler) Handle(ctx *fiber.Ctx) error {
	userID, ok := ctx.Context().Value("UserID").(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errors": models.ErrAuthUser.Error(),
		})
	}

	item := ctx.Params("item")
	if item == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": "item is required field",
		})
	}

	costItem, ok := models.PriceItem[item]
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": fmt.Sprintf("item %s is not exist", item),
		})
	}

	err := h.buyer.BuyItem(ctx.Context(), userID, item, costItem)
	if errors.Is(err, buy_item.ErrNotEnoughCoins) {
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

package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"AvitoTask/internal/models"
	"AvitoTask/internal/usecase/auth"
)

type Handler struct {
	Auth Auth
}

func NewHandler(a Auth) *Handler {
	return &Handler{
		Auth: a,
	}
}

func (c *Handler) Handle(ctx *fiber.Ctx) error {
	var request userAuthIn

	err := ctx.BodyParser(&request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": err.Error()})
	}

	username, password, err := validate(request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": err.Error()})
	}

	userID, err := c.Auth.RegisterUser(ctx.Context(), models.User{
		ID:       uuid.New().String(),
		Username: username,
		Password: password,
	})
	if errors.Is(err, auth.ErrIncorrectPassword) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": err.Error()})
	}
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"errors": err.Error()})
	}

	ctx.Locals("UserID", userID)

	return ctx.Next()
}

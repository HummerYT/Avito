package auth

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	passwordValidator "github.com/wagslane/go-password-validator"

	"AvitoTask/internal/models"
)

type userAuthIn struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,max=255"`
}

func validate(u userAuthIn) (string, string, error) {
	validate := validator.New()
	if err := validate.Struct(u); err != nil {
		return "", "", fmt.Errorf("%s: %w", models.ErrValidation, err)
	}

	if err := passwordValidator.Validate(u.Password, float64(models.MinEntropyBits)); err != nil {
		return "", "", fmt.Errorf("password is too simple: %w", err)
	}

	return u.Username, u.Password, nil
}

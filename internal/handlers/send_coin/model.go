package send_coin

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"AvitoTask/internal/models"
)

type request struct {
	ToUser string `json:"toUser" validate:"required"`
	Amount int64  `json:"amount" validate:"required,min=1"`
}

func validate(r request) error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return fmt.Errorf("%s: %w", models.ErrValidation, err)
	}

	return nil
}

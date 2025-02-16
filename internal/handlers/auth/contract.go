package auth

import (
	"context"

	"AvitoTask/internal/models"
)

type Auth interface {
	RegisterUser(ctx context.Context, user models.User) (string, error)
}

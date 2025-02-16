//go:generate mockgen -source=contract.go -destination=mocks/mock.go -package=mocks $GOPACKAGE
package auth

import (
	"AvitoTask/internal/models"
	"context"
)

type insert interface {
	IsUserExists(ctx context.Context, user models.User) (bool, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	InsertUser(ctx context.Context, user models.User) (string, error)
}

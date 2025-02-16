package info

import (
	"context"

	"AvitoTask/internal/models"
)

type infoUser interface {
	GetInfo(ctx context.Context, userID string) (username string, res models.InfoResponse, err error)
}

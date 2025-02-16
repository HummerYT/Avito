//go:generate mockgen -source=contract.go -destination=mocks/mock.go -package=mocks $GOPACKAGE
//go:generate mockgen -destination=mocks/mock_tx.go -package=mocks github.com/jackc/pgx/v5 Tx
package buy_item

import (
	"context"

	"github.com/jackc/pgx/v5"

	"AvitoTask/internal/models"
)

type user interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
	GetUserById(ctx context.Context, tx pgx.Tx, userID string) (models.User, error)
	GetUserByLoginWithTx(ctx context.Context, tx pgx.Tx, login string) (models.User, error)
	IsUserExists(ctx context.Context, user models.User) (bool, error)
	UpdateUserCoins(ctx context.Context, tx pgx.Tx, userID string, newCoins int64) error
	GetUserCoins(ctx context.Context, tx pgx.Tx, userID string) (int64, error)
}

type inventory interface {
	GetInventoryItem(ctx context.Context, tx pgx.Tx, userID, itemType string) (int64, error)
	InsertInventoryItem(ctx context.Context, tx pgx.Tx, id, userID, itemType string) error
	UpdateInventoryItem(ctx context.Context, tx pgx.Tx, userID, itemType string, newQuantity int64) error
}

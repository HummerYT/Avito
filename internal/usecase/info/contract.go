//go:generate mockgen -source=contract.go -destination=mocks/mock.go -package=mocks $GOPACKAGE
//go:generate mockgen -destination=mocks/mock_tx.go -package=mocks github.com/jackc/pgx/v5 Tx
package info

import (
	"context"

	"github.com/jackc/pgx/v5"

	"AvitoTask/internal/models"
)

type user interface {
	GetUserCoins(ctx context.Context, tx pgx.Tx, userID string) (int64, error)
	GetUserById(ctx context.Context, tx pgx.Tx, userID string) (models.User, error)
}

type inventory interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
	GetInventoryItem(ctx context.Context, tx pgx.Tx, userID, itemType string) (int64, error)
	GetUserInventory(ctx context.Context, tx pgx.Tx, userID string) ([]models.InventoryItem, error)
}

type transaction interface {
	GetUserTransactions(ctx context.Context, tx pgx.Tx, userID string) ([]models.TransactionItem, error)
}

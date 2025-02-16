//go:generate mockgen -source=contract.go -destination=mocks/mock.go -package=mocks $GOPACKAGE
//go:generate mockgen -destination=mocks/mock_tx.go -package=mocks github.com/jackc/pgx/v5 Tx
package send_coin

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
}

type transaction interface {
	InsertTransaction(ctx context.Context, tx pgx.Tx, id, fromUserID, toUserID string, amount int64) error
}

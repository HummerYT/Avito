//go:generate mockgen -source=contract.go -destination=mocks/mock.go -package=mocks $GOPACKAGE
//go:generate mockgen -destination=mocks/mock_tx.go -package=mocks github.com/jackc/pgx/v5 Tx
package auth

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type pool interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

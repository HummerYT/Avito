package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"AvitoTask/internal/models"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) InsertTransaction(ctx context.Context, tx pgx.Tx, id, fromUserID, toUserID string, amount int64) error {
	query := `
        INSERT INTO transactions (id, from_user_id, to_user_id, amount)
        VALUES ($1, $2, $3, $4)
    `
	_, err := tx.Exec(ctx, query, id, fromUserID, toUserID, amount)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}
	return nil
}

func (r *Repository) GetUserTransactions(ctx context.Context, tx pgx.Tx, userID string) ([]models.TransactionItem, error) {
	query := `
        SELECT from_user_id, to_user_id, amount, created_at, users.username
        FROM transactions
        LEFT JOIN users ON users.id = transactions.to_user_id
        WHERE from_user_id = $1
           OR to_user_id = $1
        ORDER BY created_at DESC
    `
	rows, err := tx.Query(ctx, query, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var result []models.TransactionItem
	for rows.Next() {
		var t models.TransactionItem
		if err := rows.Scan(&t.FromUserID, &t.ToUserID, &t.Amount, &t.CreatedAt, &t.Username); err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %w", err)
		}
		result = append(result, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return result, nil
}

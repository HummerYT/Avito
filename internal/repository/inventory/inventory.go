package inventory

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

func NewInsertRepo(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *Repository) GetInventoryItem(ctx context.Context, tx pgx.Tx, userID, itemType string) (int64, error) {
	var quantity int64
	query := `SELECT quantity 
              FROM inventory 
              WHERE user_id = $1 AND item_type = $2
              LIMIT 1`
	err := tx.QueryRow(ctx, query, userID, itemType).Scan(&quantity)
	if err != nil {
		return 0, err
	}
	return quantity, nil
}

func (r *Repository) InsertInventoryItem(ctx context.Context, tx pgx.Tx, id, userID, itemType string) error {
	query := `
        INSERT INTO inventory (id, user_id, item_type, quantity)
        VALUES ($1, $2, $3, 1)
    `
	_, err := tx.Exec(ctx, query, id, userID, itemType)
	if err != nil {
		return fmt.Errorf("failed to insert new item '%s' for user %s: %w", itemType, userID, err)
	}
	return nil
}

func (r *Repository) UpdateInventoryItem(ctx context.Context, tx pgx.Tx, userID, itemType string, newQuantity int64) error {
	query := `
        UPDATE inventory 
        SET quantity = $1
        WHERE user_id = $2 AND item_type = $3
    `
	_, err := tx.Exec(ctx, query, newQuantity, userID, itemType)
	if err != nil {
		return fmt.Errorf("failed to update item '%s' for user %s: %w", itemType, userID, err)
	}
	return nil
}

func (r *Repository) GetUserInventory(ctx context.Context, tx pgx.Tx, userID string) ([]models.InventoryItem, error) {
	query := `
        SELECT item_type, quantity
        FROM inventory
        WHERE user_id = $1
    `
	rows, err := tx.Query(ctx, query, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query inventory: %w", err)
	}
	defer rows.Close()

	var result []models.InventoryItem
	for rows.Next() {
		var it models.InventoryItem
		if err := rows.Scan(&it.ItemType, &it.Quantity); err != nil {
			return nil, fmt.Errorf("failed to scan inventory row: %w", err)
		}
		result = append(result, it)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return result, nil
}

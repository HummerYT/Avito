package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"AvitoTask/internal/models"
)

var (
	ErrNoUserExist = errors.New("no user exist")
)

type Repository struct {
	pool pool
}

func NewInsertRepo(pool pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) IsUserExists(ctx context.Context, user models.User) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = $1`
	err := r.pool.QueryRow(ctx, query, user.Username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}

func (r *Repository) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	var dbUser models.User
	query := `SELECT id, username, password, coins FROM users WHERE username = $1 LIMIT 1`

	row := r.pool.QueryRow(ctx, query, login)
	err := row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password, &dbUser.Coins)
	if err != nil {
		return dbUser, fmt.Errorf("failed to scan user: %w", err)
	}

	return dbUser, nil
}

func (r *Repository) GetUserByLoginWithTx(ctx context.Context, tx pgx.Tx, login string) (models.User, error) {
	var dbUser models.User
	query := `SELECT id, username, password, coins FROM users WHERE username = $1 LIMIT 1`

	row := tx.QueryRow(ctx, query, login)
	err := row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password, &dbUser.Coins)
	if err != nil {
		return dbUser, fmt.Errorf("failed to scan user: %w", err)
	}

	return dbUser, nil
}

func (r *Repository) InsertUser(ctx context.Context, user models.User) (string, error) {
	var userID string
	query := `
        INSERT INTO users (id, username, password)
        VALUES ($1, $2, $3)
        RETURNING id
    `
	err := r.pool.QueryRow(ctx, query, user.ID, user.Username, user.Password).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("failed to insert user: %w", err)
	}

	return userID, nil
}

func (r *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *Repository) GetUserById(ctx context.Context, tx pgx.Tx, userID string) (models.User, error) {
	var user models.User
	query := `SELECT id, username, password, coins 
              FROM users
              WHERE id = $1
              LIMIT 1`

	err := tx.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Coins,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrNoUserExist
	}
	if err != nil {
		return user, fmt.Errorf("cannot find user '%s': %w", userID, err)
	}

	return user, nil
}

func (r *Repository) UpdateUserCoins(ctx context.Context, tx pgx.Tx, userID string, newCoins int64) error {
	query := `UPDATE users
              SET coins = $1
              WHERE id = $2`

	_, err := tx.Exec(ctx, query, newCoins, userID)
	if err != nil {
		return fmt.Errorf("failed to update coins for userID=%s: %w", userID, err)
	}

	return nil
}

func (r *Repository) GetUserCoins(ctx context.Context, tx pgx.Tx, userID string) (int64, error) {
	var coins int64
	query := `SELECT coins FROM users WHERE id = $1 LIMIT 1`
	err := tx.QueryRow(ctx, query, userID).Scan(&coins)
	if err != nil {
		return 0, fmt.Errorf("failed to get user coins (userID=%s): %w", userID, err)
	}
	return coins, nil
}

package info

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"AvitoTask/internal/models"
)

type Usecase struct {
	repoUser        user
	repoInfo        inventory
	repoTransaction transaction
	TX              func(ctx context.Context) (pgx.Tx, error)
}

func New(repoUser user, repo inventory, t transaction) *Usecase {
	return &Usecase{
		repoUser:        repoUser,
		repoInfo:        repo,
		repoTransaction: t,
		TX:              repo.BeginTx,
	}
}

func (uc *Usecase) GetInfo(ctx context.Context, userID string) (username string, res models.InfoResponse, err error) {
	var tx pgx.Tx
	tx, err = uc.TX(ctx)
	if err != nil {
		return "", res, fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	userFrom, err := uc.repoUser.GetUserById(ctx, tx, userID)
	if err != nil {
		return "", res, err
	}
	res.Coins = userFrom.Coins

	items, err := uc.repoInfo.GetUserInventory(ctx, tx, userID)
	if err != nil {
		return "", res, err
	}
	for _, it := range items {
		res.Inventory = append(res.Inventory, models.InventoryItem{
			ItemType: it.ItemType,
			Quantity: it.Quantity,
		})
	}

	txs, err := uc.repoTransaction.GetUserTransactions(ctx, tx, userID)
	if err != nil {
		return "", res, err
	}
	for _, t := range txs {
		res.Transactions = append(res.Transactions, models.TransactionItem{
			FromUserID:   t.FromUserID,
			ToUserID:     t.ToUserID,
			Amount:       t.Amount,
			CreatedAt:    t.CreatedAt,
			FromUsername: t.FromUsername,
			ToUserName:   t.ToUserName,
		})
	}

	return userFrom.Username, res, nil
}

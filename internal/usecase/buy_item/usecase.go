package buy_item

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrNotEnoughCoins = errors.New("not enough coins to buy this item")

type Usecase struct {
	repoUser      user
	repoInventory inventory
}

func NewUsecase(u user, i inventory) *Usecase {
	return &Usecase{
		repoUser:      u,
		repoInventory: i,
	}
}

func (u *Usecase) BuyItem(ctx context.Context, userID string, item string, cost int64) (err error) {
	tx, err := u.repoUser.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	currentCoins, err := u.repoUser.GetUserCoins(ctx, tx, userID)
	if err != nil {
		return err
	}

	if currentCoins < cost {
		err = ErrNotEnoughCoins
		return err
	}

	newCoins := currentCoins - cost
	if err = u.repoUser.UpdateUserCoins(ctx, tx, userID, newCoins); err != nil {
		return err
	}

	quantity, err := u.repoInventory.GetInventoryItem(ctx, tx, userID, item)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if insertErr := u.repoInventory.InsertInventoryItem(ctx, tx, uuid.New().String(), userID, item); insertErr != nil {
				err = insertErr
				return err
			}
		} else {
			return err
		}
	}

	newQuantity := quantity + 1
	if err = u.repoInventory.UpdateInventoryItem(ctx, tx, userID, item, newQuantity); err != nil {
		return err
	}

	return nil
}

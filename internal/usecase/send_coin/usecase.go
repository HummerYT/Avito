package send_coin

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrSameUser       = errors.New("cannot send coins to the same user")
	ErrNotEnoughCoins = errors.New("user does not have enough coins to send")
)

type Usecase struct {
	repoUser        user
	repoTransaction transaction
}

func NewUsecase(repoUser user, repoTransaction transaction) *Usecase {
	return &Usecase{
		repoUser:        repoUser,
		repoTransaction: repoTransaction,
	}
}

func (u *Usecase) SendCoin(ctx context.Context, fromUser, toUser string, amount int64) (err error) {
	tx, err := u.repoUser.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	fromData, err := u.repoUser.GetUserById(ctx, tx, fromUser)
	if err != nil {
		return fmt.Errorf("failed to get user by id: %w", err)
	}

	if fromData.Username == toUser {
		return ErrSameUser
	}

	toData, err := u.repoUser.GetUserByLoginWithTx(ctx, tx, toUser)
	if err != nil {
		return fmt.Errorf("failed to get user by id: %w", err)
	}

	if fromData.Coins < amount {
		err = ErrNotEnoughCoins
		return ErrNotEnoughCoins
	}

	newFromCoins := fromData.Coins - amount
	newToCoins := toData.Coins + amount

	if err = u.repoUser.UpdateUserCoins(ctx, tx, fromData.ID, newFromCoins); err != nil {
		return fmt.Errorf("failed to update user coins: %w", err)
	}
	if err = u.repoUser.UpdateUserCoins(ctx, tx, toData.ID, newToCoins); err != nil {
		return fmt.Errorf("failed to update user coins: %w", err)
	}

	if err = u.repoTransaction.InsertTransaction(ctx, tx, uuid.New().String(), fromData.ID, toData.ID, amount); err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	return nil
}

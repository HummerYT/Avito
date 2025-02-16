package send_coin_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"AvitoTask/internal/models"
	"AvitoTask/internal/usecase/send_coin"
	"AvitoTask/internal/usecase/send_coin/mocks"
)

func TestSendCoin_SameUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockUser := mocks.NewMockuser(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockTransaction := mocks.NewMocktransaction(ctrl)

	mockUser.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(ctx).Return(nil)

	fromData := models.User{ID: "user123", Username: "user123", Coins: 100}
	mockUser.EXPECT().GetUserById(gomock.Any(), gomock.Any(), gomock.Any()).Return(fromData, nil)

	uc := send_coin.NewUsecase(mockUser, mockTransaction)
	err := uc.SendCoin(ctx, "user123", "user123", 100)
	if !errors.Is(err, send_coin.ErrSameUser) {
		t.Errorf("expected error %v, got %v", send_coin.ErrSameUser, err)
	}
}

func TestSendCoin_BeginTxError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockUser := mocks.NewMockuser(ctrl)
	mockTransaction := mocks.NewMocktransaction(ctrl)

	beginErr := errors.New("begin tx error")
	mockUser.EXPECT().BeginTx(ctx).Return(nil, beginErr)

	uc := send_coin.NewUsecase(mockUser, mockTransaction)
	err := uc.SendCoin(ctx, "user123", "user456", 100)
	expectedMsg := fmt.Sprintf("failed to begin transaction: %v", beginErr)
	if err == nil || err.Error() != expectedMsg {
		t.Errorf("expected error %q, got %v", expectedMsg, err)
	}
}

func TestSendCoin_GetUserByIdError_FromUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockUser := mocks.NewMockuser(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	mockTransaction := mocks.NewMocktransaction(ctrl)

	mockUser.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	getUserErr := errors.New("get user error")
	mockUser.EXPECT().GetUserById(ctx, mockTx, "user123").
		Return(models.User{}, getUserErr)
	mockTx.EXPECT().Rollback(ctx).Return(nil)

	uc := send_coin.NewUsecase(mockUser, mockTransaction)
	err := uc.SendCoin(ctx, "user123", "user456", 100)
	expectedMsg := fmt.Sprintf("failed to get user by id: %v", getUserErr)
	if err == nil || err.Error() != expectedMsg {
		t.Errorf("expected error %q, got %v", expectedMsg, err)
	}
}

func TestSendCoin_GetUserByIdError_ToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockUser := mocks.NewMockuser(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockTransaction := mocks.NewMocktransaction(ctrl)

	mockUser.EXPECT().BeginTx(ctx).Return(mockTx, nil)

	fromData := models.User{ID: "user123", Coins: 200}
	mockUser.EXPECT().GetUserById(ctx, mockTx, "user123").Return(fromData, nil)

	getUserErr := errors.New("get user error")
	mockUser.EXPECT().GetUserByLoginWithTx(ctx, mockTx, "user456").Return(models.User{}, getUserErr)
	mockTx.EXPECT().Rollback(ctx).Return(nil)

	uc := send_coin.NewUsecase(mockUser, mockTransaction)
	err := uc.SendCoin(ctx, "user123", "user456", 100)
	expectedMsg := fmt.Sprintf("failed to get user by id: %v", getUserErr)
	if err == nil || err.Error() != expectedMsg {
		t.Errorf("expected error %q, got %v", expectedMsg, err)
	}
}

func TestSendCoin_NotEnoughCoins(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockUser := mocks.NewMockuser(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockTransaction := mocks.NewMocktransaction(ctrl)
	mockUser.EXPECT().BeginTx(ctx).Return(mockTx, nil)

	fromData := models.User{ID: "user123", Coins: 50}
	toData := models.User{ID: "user456", Coins: 100}
	mockUser.EXPECT().GetUserById(ctx, mockTx, "user123").Return(fromData, nil)
	mockUser.EXPECT().GetUserByLoginWithTx(ctx, mockTx, "user456").Return(toData, nil)
	mockTx.EXPECT().Rollback(ctx).Return(nil)

	uc := send_coin.NewUsecase(mockUser, mockTransaction)
	err := uc.SendCoin(ctx, "user123", "user456", 100)
	if err == nil || !errors.Is(err, send_coin.ErrNotEnoughCoins) {
		t.Errorf("expected error %v, got %v", send_coin.ErrNotEnoughCoins, err)
	}
}

func TestSendCoin_UpdateUserCoinsError_FromUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockUser := mocks.NewMockuser(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockTransaction := mocks.NewMocktransaction(ctrl)

	mockUser.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	fromData := models.User{ID: "user123", Coins: 200}
	toData := models.User{ID: "user456", Coins: 100}
	mockUser.EXPECT().GetUserById(ctx, mockTx, "user123").Return(fromData, nil)
	mockUser.EXPECT().GetUserByLoginWithTx(ctx, mockTx, "user456").Return(toData, nil)

	newFromCoins := fromData.Coins - 100
	updateErr := errors.New("update coins error")
	mockUser.EXPECT().UpdateUserCoins(ctx, mockTx, "user123", newFromCoins).Return(updateErr)
	mockTx.EXPECT().Rollback(ctx).Return(nil)

	uc := send_coin.NewUsecase(mockUser, mockTransaction)
	err := uc.SendCoin(ctx, "user123", "user456", 100)
	expectedMsg := fmt.Sprintf("failed to update user coins: %v", updateErr)
	if err == nil || err.Error() != expectedMsg {
		t.Errorf("expected error %q, got %v", expectedMsg, err)
	}
}

func TestSendCoin_UpdateUserCoinsError_ToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockUser := mocks.NewMockuser(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockTransaction := mocks.NewMocktransaction(ctrl)

	mockUser.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	fromData := models.User{ID: "user123", Coins: 200}
	toData := models.User{ID: "user456", Coins: 100}
	mockUser.EXPECT().GetUserById(ctx, mockTx, "user123").Return(fromData, nil)
	mockUser.EXPECT().GetUserByLoginWithTx(ctx, mockTx, "user456").Return(toData, nil)

	newFromCoins := fromData.Coins - 100
	newToCoins := toData.Coins + 100

	mockUser.EXPECT().UpdateUserCoins(ctx, mockTx, "user123", newFromCoins).Return(nil)

	updateErr := errors.New("update to coins error")
	mockUser.EXPECT().UpdateUserCoins(ctx, mockTx, "user456", newToCoins).Return(updateErr)
	mockTx.EXPECT().Rollback(ctx).Return(nil)

	uc := send_coin.NewUsecase(mockUser, mockTransaction)
	err := uc.SendCoin(ctx, "user123", "user456", 100)
	expectedMsg := fmt.Sprintf("failed to update user coins: %v", updateErr)
	if err == nil || err.Error() != expectedMsg {
		t.Errorf("expected error %q, got %v", expectedMsg, err)
	}
}

func TestSendCoin_InsertTransactionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockUser := mocks.NewMockuser(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockTransaction := mocks.NewMocktransaction(ctrl)

	mockUser.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	fromData := models.User{ID: "user123", Coins: 200}
	toData := models.User{ID: "user456", Coins: 100}
	mockUser.EXPECT().GetUserById(ctx, mockTx, "user123").Return(fromData, nil)
	mockUser.EXPECT().GetUserByLoginWithTx(ctx, mockTx, "user456").Return(toData, nil)
	newFromCoins := fromData.Coins - 100
	newToCoins := toData.Coins + 100
	mockUser.EXPECT().UpdateUserCoins(ctx, mockTx, "user123", newFromCoins).Return(nil)
	mockUser.EXPECT().UpdateUserCoins(ctx, mockTx, "user456", newToCoins).Return(nil)

	insertErr := errors.New("insert transaction error")
	mockTransaction.EXPECT().
		InsertTransaction(ctx, mockTx, gomock.Any(), "user123", "user456", gomock.Any()).
		Return(insertErr)
	mockTx.EXPECT().Rollback(ctx).Return(nil)

	uc := send_coin.NewUsecase(mockUser, mockTransaction)
	err := uc.SendCoin(ctx, "user123", "user456", 100)
	expectedMsg := fmt.Sprintf("failed to insert transaction: %v", insertErr)
	if err == nil || err.Error() != expectedMsg {
		t.Errorf("expected error %q, got %v", expectedMsg, err)
	}
}

func TestSendCoin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockUser := mocks.NewMockuser(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockTransaction := mocks.NewMocktransaction(ctrl)

	mockUser.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	fromData := models.User{ID: "user123", Coins: 200}
	toData := models.User{ID: "user456", Coins: 100}
	mockUser.EXPECT().GetUserById(ctx, mockTx, "user123").Return(fromData, nil)
	mockUser.EXPECT().GetUserByLoginWithTx(ctx, mockTx, "user456").Return(toData, nil)

	newFromCoins := fromData.Coins - 100
	newToCoins := toData.Coins + 100
	mockUser.EXPECT().UpdateUserCoins(ctx, mockTx, "user123", newFromCoins).Return(nil)
	mockUser.EXPECT().UpdateUserCoins(ctx, mockTx, "user456", newToCoins).Return(nil)
	mockTransaction.EXPECT().
		InsertTransaction(ctx, mockTx, gomock.Any(), "user123", "user456", gomock.Any()).
		Return(nil)
	mockTx.EXPECT().Commit(ctx).Return(nil)

	uc := send_coin.NewUsecase(mockUser, mockTransaction)
	err := uc.SendCoin(ctx, "user123", "user456", 100)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

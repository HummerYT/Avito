package info_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"

	"AvitoTask/internal/models"
	"AvitoTask/internal/usecase/info"
	"AvitoTask/internal/usecase/info/mocks"
)

func TestGetInfo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "user123"
	user := models.User{
		ID:       userID,
		Username: userID,
		Coins:    100,
	}

	mockUser := mocks.NewMockuser(ctrl)
	mockInventory := mocks.NewMockinventory(ctrl)
	mockTransaction := mocks.NewMocktransaction(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	uc := info.New(mockUser, mockInventory, mockTransaction)
	uc.TX = func(ctx context.Context) (pgx.Tx, error) {
		return mockTx, nil
	}

	expectedCoins := int64(100)
	expectedInventory := []models.InventoryItem{
		{ItemType: "sword", Quantity: 1},
		{ItemType: "shield", Quantity: 2},
	}
	expectedTransactions := []models.TransactionItem{
		{
			FromUserID: "user123",
			ToUserID:   "user456",
			Amount:     50,
			CreatedAt:  time.Now(),
		},
	}

	mockUser.
		EXPECT().
		GetUserById(ctx, mockTx, userID).
		Return(user, nil)
	mockInventory.
		EXPECT().
		GetUserInventory(ctx, mockTx, userID).
		Return(expectedInventory, nil)
	mockTransaction.
		EXPECT().
		GetUserTransactions(ctx, mockTx, userID).
		Return(expectedTransactions, nil)

	mockTx.
		EXPECT().
		Commit(ctx).
		Return(nil)

	_, res, err := uc.GetInfo(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Coins != expectedCoins {
		t.Errorf("expected coins %d, got %d", expectedCoins, res.Coins)
	}
	if len(res.Inventory) != len(expectedInventory) {
		t.Errorf("expected inventory length %d, got %d", len(expectedInventory), len(res.Inventory))
	}
	if len(res.Transactions) != len(expectedTransactions) {
		t.Errorf("expected transactions length %d, got %d", len(expectedTransactions), len(res.Transactions))
	}
}

func TestGetInfo_TXError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "user123"

	mockUser := mocks.NewMockuser(ctrl)
	mockInventory := mocks.NewMockinventory(ctrl)
	mockTransaction := mocks.NewMocktransaction(ctrl)

	uc := info.New(mockUser, mockInventory, mockTransaction)
	expectedErr := errors.New("begin tx error")
	uc.TX = func(ctx context.Context) (pgx.Tx, error) {
		return nil, expectedErr
	}

	_, _, err := uc.GetInfo(ctx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != ("begin tx: " + expectedErr.Error()) {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGetInfo_GetUserCoinsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "user123"
	user := models.User{
		ID:       userID,
		Username: userID,
		Coins:    100,
	}

	mockUser := mocks.NewMockuser(ctrl)
	mockInventory := mocks.NewMockinventory(ctrl)
	mockTransaction := mocks.NewMocktransaction(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	uc := info.New(mockUser, mockInventory, mockTransaction)
	uc.TX = func(ctx context.Context) (pgx.Tx, error) {
		return mockTx, nil
	}
	expectedErr := errors.New("get user coins error")
	mockUser.
		EXPECT().
		GetUserById(ctx, mockTx, userID).
		Return(user, expectedErr)
	mockTx.
		EXPECT().
		Rollback(ctx).
		Return(nil)

	_, _, err := uc.GetInfo(ctx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestGetInfo_GetUserInventoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "user123"
	user := models.User{
		ID:       userID,
		Username: userID,
		Coins:    100,
	}

	mockUser := mocks.NewMockuser(ctrl)
	mockInventory := mocks.NewMockinventory(ctrl)
	mockTransaction := mocks.NewMocktransaction(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	uc := info.New(mockUser, mockInventory, mockTransaction)
	uc.TX = func(ctx context.Context) (pgx.Tx, error) {
		return mockTx, nil
	}

	expectedCoins := int64(100)
	expectedErr := errors.New("get inventory error")

	mockUser.
		EXPECT().
		GetUserById(ctx, mockTx, userID).
		Return(user, nil)
	mockInventory.
		EXPECT().
		GetUserInventory(ctx, mockTx, userID).
		Return(nil, expectedErr)

	mockTx.
		EXPECT().
		Rollback(ctx).
		Return(nil)

	_, res, err := uc.GetInfo(ctx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
	if res.Coins != expectedCoins {
		t.Errorf("expected coins %d, got %d", expectedCoins, res.Coins)
	}
	if len(res.Inventory) != 0 {
		t.Errorf("expected empty inventory, got %v", res.Inventory)
	}
}

func TestGetInfo_GetUserTransactionsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "user123"
	user := models.User{
		ID:       userID,
		Username: userID,
		Coins:    100,
	}

	mockUser := mocks.NewMockuser(ctrl)
	mockInventory := mocks.NewMockinventory(ctrl)
	mockTransaction := mocks.NewMocktransaction(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	uc := info.New(mockUser, mockInventory, mockTransaction)
	uc.TX = func(ctx context.Context) (pgx.Tx, error) {
		return mockTx, nil
	}

	expectedCoins := int64(100)
	expectedInventory := []models.InventoryItem{
		{ItemType: "sword", Quantity: 1},
	}
	expectedErr := errors.New("get transactions error")

	mockUser.
		EXPECT().
		GetUserById(ctx, mockTx, userID).
		Return(user, nil)
	mockInventory.
		EXPECT().
		GetUserInventory(ctx, mockTx, userID).
		Return(expectedInventory, nil)
	mockTransaction.
		EXPECT().
		GetUserTransactions(ctx, mockTx, userID).
		Return(nil, expectedErr)
	mockTx.
		EXPECT().
		Rollback(ctx).
		Return(nil)

	_, res, err := uc.GetInfo(ctx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
	if res.Coins != expectedCoins {
		t.Errorf("expected coins %d, got %d", expectedCoins, res.Coins)
	}
	if len(res.Inventory) != len(expectedInventory) {
		t.Errorf("expected inventory length %d, got %d", len(expectedInventory), len(res.Inventory))
	}
	if len(res.Transactions) != 0 {
		t.Errorf("expected empty transactions, got %v", res.Transactions)
	}
}

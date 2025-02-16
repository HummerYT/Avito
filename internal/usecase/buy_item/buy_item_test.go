package buy_item_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"

	"AvitoTask/internal/usecase/buy_item"
	"AvitoTask/internal/usecase/buy_item/mocks"
)

func TestBuyItem_BeginTxError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "user123"
	item := "sword"
	cost := int64(100)

	mockUser := mocks.NewMockuser(ctrl)
	mockInventory := mocks.NewMockinventory(ctrl)

	beginErr := errors.New("begin tx error")
	mockUser.EXPECT().BeginTx(ctx).Return(nil, beginErr)

	uc := buy_item.NewUsecase(mockUser, mockInventory)
	err := uc.BuyItem(ctx, userID, item, cost)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	expectedMsg := fmt.Sprintf("failed to begin tx: %v", beginErr)
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}
}

func TestBuyItem_GetUserCoinsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "user123"
	item := "sword"
	cost := int64(100)

	mockUser := mocks.NewMockuser(ctrl)
	mockInventory := mocks.NewMockinventory(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	mockUser.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	getCoinsErr := errors.New("failed to get coins")
	mockUser.EXPECT().GetUserCoins(ctx, mockTx, userID).Return(int64(0), getCoinsErr)
	mockTx.EXPECT().Rollback(ctx).Return(nil)

	uc := buy_item.NewUsecase(mockUser, mockInventory)
	err := uc.BuyItem(ctx, userID, item, cost)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, getCoinsErr) {
		t.Errorf("expected error %v, got %v", getCoinsErr, err)
	}
}

func TestBuyItem_NotEnoughCoins(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "user123"
	item := "sword"
	cost := int64(100)

	mockUser := mocks.NewMockuser(ctrl)
	mockInventory := mocks.NewMockinventory(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	mockUser.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	mockUser.EXPECT().GetUserCoins(ctx, mockTx, userID).Return(int64(50), nil)
	mockTx.EXPECT().Rollback(ctx).Return(nil)

	uc := buy_item.NewUsecase(mockUser, mockInventory)
	err := uc.BuyItem(ctx, userID, item, cost)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, buy_item.ErrNotEnoughCoins) {
		t.Errorf("expected error %v, got %v", buy_item.ErrNotEnoughCoins, err)
	}
}

func TestBuyItem_UpdateUserCoinsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "user123"
	item := "sword"
	cost := int64(100)
	startingCoins := int64(150)

	mockUser := mocks.NewMockuser(ctrl)
	mockInventory := mocks.NewMockinventory(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	mockUser.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	mockUser.EXPECT().GetUserCoins(ctx, mockTx, userID).Return(startingCoins, nil)
	newCoins := startingCoins - cost
	updateErr := errors.New("failed to update coins")
	mockUser.EXPECT().UpdateUserCoins(ctx, mockTx, userID, newCoins).Return(updateErr)
	mockTx.EXPECT().Rollback(ctx).Return(nil)

	uc := buy_item.NewUsecase(mockUser, mockInventory)
	err := uc.BuyItem(ctx, userID, item, cost)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, updateErr) {
		t.Errorf("expected error %v, got %v", updateErr, err)
	}
}

func TestBuyItem_GetInventoryItemError_NotNoRows(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "user123"
	item := "sword"
	cost := int64(100)
	startingCoins := int64(200)

	mockUser := mocks.NewMockuser(ctrl)
	mockInventory := mocks.NewMockinventory(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	mockUser.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	mockUser.EXPECT().GetUserCoins(ctx, mockTx, userID).Return(startingCoins, nil)
	newCoins := startingCoins - cost
	mockUser.EXPECT().UpdateUserCoins(ctx, mockTx, userID, newCoins).Return(nil)

	invErr := errors.New("inventory error")
	mockInventory.EXPECT().GetInventoryItem(ctx, mockTx, userID, item).Return(int64(0), invErr)
	mockTx.EXPECT().Rollback(ctx).Return(nil)

	uc := buy_item.NewUsecase(mockUser, mockInventory)
	err := uc.BuyItem(ctx, userID, item, cost)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, invErr) {
		t.Errorf("expected error %v, got %v", invErr, err)
	}
}

func TestBuyItem_InsertInventoryItemError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "user123"
	item := "sword"
	cost := int64(100)
	startingCoins := int64(200)

	mockUser := mocks.NewMockuser(ctrl)
	mockInventory := mocks.NewMockinventory(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	mockUser.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	mockUser.EXPECT().GetUserCoins(ctx, mockTx, userID).Return(startingCoins, nil)
	newCoins := startingCoins - cost
	mockUser.EXPECT().UpdateUserCoins(ctx, mockTx, userID, newCoins).Return(nil)

	mockInventory.EXPECT().GetInventoryItem(ctx, mockTx, userID, item).Return(int64(0), pgx.ErrNoRows)
	insertErr := errors.New("failed to insert inventory")
	mockInventory.EXPECT().InsertInventoryItem(ctx, mockTx, gomock.Any(), userID, item).Return(insertErr)
	mockTx.EXPECT().Rollback(ctx).Return(nil)

	uc := buy_item.NewUsecase(mockUser, mockInventory)
	err := uc.BuyItem(ctx, userID, item, cost)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, insertErr) {
		t.Errorf("expected error %v, got %v", insertErr, err)
	}
}

func TestBuyItem_UpdateInventoryItemError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "user123"
	item := "sword"
	cost := int64(100)
	startingCoins := int64(200)

	mockUser := mocks.NewMockuser(ctrl)
	mockInventory := mocks.NewMockinventory(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	mockUser.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	mockUser.EXPECT().GetUserCoins(ctx, mockTx, userID).Return(startingCoins, nil)
	newCoins := startingCoins - cost
	mockUser.EXPECT().UpdateUserCoins(ctx, mockTx, userID, newCoins).Return(nil)

	existingQuantity := int64(2)
	mockInventory.EXPECT().GetInventoryItem(ctx, mockTx, userID, item).Return(existingQuantity, nil)
	updateInvErr := errors.New("failed to update inventory")
	newQuantity := existingQuantity + 1
	mockInventory.EXPECT().UpdateInventoryItem(ctx, mockTx, userID, item, newQuantity).Return(updateInvErr)
	mockTx.EXPECT().Rollback(ctx).Return(nil)

	uc := buy_item.NewUsecase(mockUser, mockInventory)
	err := uc.BuyItem(ctx, userID, item, cost)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, updateInvErr) {
		t.Errorf("expected error %v, got %v", updateInvErr, err)
	}
}

func TestBuyItem_Success_ExistingItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "user123"
	item := "sword"
	cost := int64(100)
	startingCoins := int64(200)

	mockUser := mocks.NewMockuser(ctrl)
	mockInventory := mocks.NewMockinventory(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	mockUser.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	mockUser.EXPECT().GetUserCoins(ctx, mockTx, userID).Return(startingCoins, nil)
	newCoins := startingCoins - cost
	mockUser.EXPECT().UpdateUserCoins(ctx, mockTx, userID, newCoins).Return(nil)

	existingQuantity := int64(3)
	mockInventory.EXPECT().GetInventoryItem(ctx, mockTx, userID, item).Return(existingQuantity, nil)
	newQuantity := existingQuantity + 1
	mockInventory.EXPECT().UpdateInventoryItem(ctx, mockTx, userID, item, newQuantity).Return(nil)

	mockTx.EXPECT().Commit(ctx).Return(nil)

	uc := buy_item.NewUsecase(mockUser, mockInventory)
	err := uc.BuyItem(ctx, userID, item, cost)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuyItem_Success_NewItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := "user123"
	item := "sword"
	cost := int64(100)
	startingCoins := int64(200)

	mockUser := mocks.NewMockuser(ctrl)
	mockInventory := mocks.NewMockinventory(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	mockUser.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	mockUser.EXPECT().GetUserCoins(ctx, mockTx, userID).Return(startingCoins, nil)
	newCoins := startingCoins - cost
	mockUser.EXPECT().UpdateUserCoins(ctx, mockTx, userID, newCoins).Return(nil)

	mockInventory.EXPECT().GetInventoryItem(ctx, mockTx, userID, item).Return(int64(0), pgx.ErrNoRows)

	mockInventory.EXPECT().InsertInventoryItem(ctx, mockTx, gomock.Any(), userID, item).Return(nil)

	mockInventory.EXPECT().UpdateInventoryItem(ctx, mockTx, userID, item, int64(1)).Return(nil)
	mockTx.EXPECT().Commit(ctx).Return(nil)

	uc := buy_item.NewUsecase(mockUser, mockInventory)
	err := uc.BuyItem(ctx, userID, item, cost)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

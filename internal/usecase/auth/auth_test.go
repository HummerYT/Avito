package auth_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"AvitoTask/internal/models"
	"AvitoTask/internal/usecase/auth"
	"AvitoTask/internal/usecase/auth/mocks"
)

func TestRegisterUser_IsUserExistsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	testUser := models.User{Username: "testuser", Password: "password123"}

	mockInsert := mocks.NewMockinsert(ctrl)
	expectedErr := errors.New("db error")
	mockInsert.EXPECT().
		IsUserExists(ctx, testUser).
		Return(false, expectedErr)

	client := auth.New(mockInsert)
	id, err := client.RegisterUser(ctx, testUser)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if id != "" {
		t.Errorf("expected empty id, got %s", id)
	}
	expectedMsg := fmt.Sprintf("failed user exists: %v", expectedErr)
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}
}

func TestRegisterUser_UserExists_GetUserByLoginError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	testUser := models.User{Username: "testuser", Password: "password123"}

	mockInsert := mocks.NewMockinsert(ctrl)
	mockInsert.EXPECT().
		IsUserExists(ctx, testUser).
		Return(true, nil)
	getUserErr := errors.New("get user error")
	mockInsert.EXPECT().
		GetUserByLogin(ctx, testUser.Username).
		Return(models.User{}, getUserErr)

	client := auth.New(mockInsert)
	id, err := client.RegisterUser(ctx, testUser)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if id != "" {
		t.Errorf("expected empty id, got %s", id)
	}
	expectedMsg := fmt.Sprintf("failed get user by username: %v", getUserErr)
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}
}

func TestRegisterUser_UserExists_IncorrectPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	plainPassword := "password123"
	testUser := models.User{Username: "testuser", Password: plainPassword}

	dbUser := models.User{
		ID:       "user-id-123",
		Username: "testuser",
		Password: "hashed_password",
	}

	mockInsert := mocks.NewMockinsert(ctrl)
	mockInsert.EXPECT().
		IsUserExists(ctx, testUser).
		Return(true, nil)
	mockInsert.EXPECT().
		GetUserByLogin(ctx, testUser.Username).
		Return(dbUser, nil)

	client := auth.New(mockInsert)
	client.CompareHashAndPassword = func(hash, password string) (bool, error) {
		return false, errors.New("password mismatch")
	}
	id, err := client.RegisterUser(ctx, testUser)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, auth.ErrIncorrectPassword) {
		t.Errorf("expected error %v, got %v", auth.ErrIncorrectPassword, err)
	}
	if id != "" {
		t.Errorf("expected empty id, got %s", id)
	}
}

func TestRegisterUser_UserExists_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	plainPassword := "password123"
	testUser := models.User{Username: "testuser", Password: plainPassword}
	dbUser := models.User{
		ID:       "user-id-123",
		Username: "testuser",
		Password: "hashed_password",
	}

	mockInsert := mocks.NewMockinsert(ctrl)
	mockInsert.EXPECT().
		IsUserExists(ctx, testUser).
		Return(true, nil)
	mockInsert.EXPECT().
		GetUserByLogin(ctx, testUser.Username).
		Return(dbUser, nil)

	client := auth.New(mockInsert)
	client.CompareHashAndPassword = func(hash, password string) (bool, error) {
		if hash == "hashed_password" && password == plainPassword {
			return true, nil
		}
		return false, errors.New("password mismatch")
	}
	id, err := client.RegisterUser(ctx, testUser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != dbUser.ID {
		t.Errorf("expected id %s, got %s", dbUser.ID, id)
	}
}

func TestRegisterUser_UserDoesNotExist_CreateHashPasswordError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	testUser := models.User{Username: "testuser", Password: "password123"}

	mockInsert := mocks.NewMockinsert(ctrl)
	mockInsert.EXPECT().
		IsUserExists(ctx, testUser).
		Return(false, nil)

	client := auth.New(mockInsert)
	client.CreateHashPassword = func(password string) (string, error) {
		return "", errors.New("hash error")
	}

	id, err := client.RegisterUser(ctx, testUser)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expectedMsg := "failed generate password: hash error"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}
	if id != "" {
		t.Errorf("expected empty id, got %s", id)
	}
}

func TestRegisterUser_UserDoesNotExist_InsertUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	testUser := models.User{Username: "testuser", Password: "password123"}

	mockInsert := mocks.NewMockinsert(ctrl)
	mockInsert.EXPECT().
		IsUserExists(ctx, testUser).
		Return(false, nil)

	client := auth.New(mockInsert)
	hashed := "hashed_password_new"
	client.CreateHashPassword = func(password string) (string, error) {
		return hashed, nil
	}

	expectedUser := testUser
	expectedUser.Password = hashed
	insertErr := errors.New("insert error")
	mockInsert.EXPECT().
		InsertUser(ctx, expectedUser).
		Return("", insertErr)

	id, err := client.RegisterUser(ctx, testUser)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expectedMsg := fmt.Sprintf("failed of create user: %v", insertErr)
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}
	if id != "" {
		t.Errorf("expected empty id, got %s", id)
	}
}

func TestRegisterUser_UserDoesNotExist_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	testUser := models.User{Username: "testuser", Password: "password123"}

	mockInsert := mocks.NewMockinsert(ctrl)
	mockInsert.EXPECT().
		IsUserExists(ctx, testUser).
		Return(false, nil)
	client := auth.New(mockInsert)

	hashed := "hashed_password_new"
	client.CreateHashPassword = func(password string) (string, error) {
		return hashed, nil
	}

	expectedUser := testUser
	expectedUser.Password = hashed
	expectedUserID := uuid.New().String()
	mockInsert.EXPECT().
		InsertUser(ctx, expectedUser).
		Return(expectedUserID, nil)

	id, err := client.RegisterUser(ctx, testUser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != expectedUserID {
		t.Errorf("expected userID %s, got %s", expectedUserID, id)
	}
}

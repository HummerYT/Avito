package auth

import (
	"context"
	"errors"
	"fmt"

	"AvitoTask/internal/models"
	"AvitoTask/internal/utils"
)

var ErrIncorrectPassword = errors.New("incorrect password")

type Client struct {
	insert                 insert
	CreateHashPassword     func(password string) (string, error)
	CompareHashAndPassword func(hash string, password string) (bool, error)
}

func New(insert insert) *Client {
	return &Client{
		insert:                 insert,
		CreateHashPassword:     utils.CreateHashPassword,
		CompareHashAndPassword: utils.CompareHashAndPassword,
	}
}

func (c *Client) RegisterUser(ctx context.Context, user models.User) (string, error) {
	isExists, err := c.insert.IsUserExists(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed user exists: %w", err)
	}

	if isExists {
		dbUser, err := c.insert.GetUserByLogin(ctx, user.Username)
		if err != nil {
			return "", fmt.Errorf("failed get user by username: %w", err)
		}

		_, err = c.CompareHashAndPassword(dbUser.Password, user.Password)
		if err != nil {
			return "", ErrIncorrectPassword
		}

		return dbUser.ID, nil
	}

	hashPassword, err := c.CreateHashPassword(user.Password)
	if err != nil {
		return "", fmt.Errorf("failed generate password: %w", err)
	}

	user.Password = hashPassword

	userID, err := c.insert.InsertUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed of create user: %w", err)
	}

	return userID, nil
}

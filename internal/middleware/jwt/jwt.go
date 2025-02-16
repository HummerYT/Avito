package jwt

import (
	"AvitoTask/internal/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

type Middleware struct {
	SecretKey string
}

func NewMiddleware(secretKey string) *Middleware {
	return &Middleware{
		SecretKey: secretKey,
	}
}

// SignedToken - подписание JWT для авторизированного пользователя токена
func (m *Middleware) SignedToken(ctx *fiber.Ctx) error {
	userID, ok := ctx.Context().Value("UserID").(string)
	if !ok {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": "header user is empty",
		})
	}

	payload := jwt.MapClaims{
		"ExpiresAt": jwt.NewNumericDate(time.Now().UTC().Add(models.DurationJwtToken)),
		"IssuedAt":  jwt.NewNumericDate(time.Now().UTC()),
		"UserID":    userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token.Header["kid"] = userID

	secretKey := []byte(m.SecretKey)

	jwtToken, err := token.SignedString(secretKey)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}

	ctx.Set(models.AuthorizationToken, jwtToken)

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"token": jwtToken,
	})
}

func (m *Middleware) CompareToken(c *fiber.Ctx) error {
	tokenStr := c.Get(models.AuthorizationToken, "")
	if tokenStr == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token is empty",
		})
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	secretKey := []byte(m.SecretKey)

	jwtToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": fmt.Sprintf("JWT token is not valid: %v", err),
		})
	}

	payload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	expires, ok := payload["ExpiresAt"].(float64)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing or invalid 'ExpiresAt' field",
		})
	}

	expiresAt := time.Unix(int64(expires), 0)
	if time.Now().After(expiresAt) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "JWT token is expired",
		})
	}

	c.Locals("UserID", jwtToken.Header["kid"])

	return c.Next()
}

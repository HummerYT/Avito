package main

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"AvitoTask/internal/config"
	"AvitoTask/internal/handlers/auth"
	"AvitoTask/internal/handlers/buy_item"
	"AvitoTask/internal/handlers/info"
	"AvitoTask/internal/handlers/send_coin"
	"AvitoTask/internal/middleware/jwt"
	authRepository "AvitoTask/internal/repository/auth"
	"AvitoTask/internal/repository/inventory"
	"AvitoTask/internal/repository/transaction"
	authUsecase "AvitoTask/internal/usecase/auth"
	buyItemUsecase "AvitoTask/internal/usecase/buy_item"
	infoUsecase "AvitoTask/internal/usecase/info"
	sendCoinUseCase "AvitoTask/internal/usecase/send_coin"
)

func main() {
	ctx := context.Background()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(
		cors.New(
			cors.Config{
				Next:             nil,
				AllowOriginsFunc: nil,
				AllowOrigins:     "*",
				AllowMethods: strings.Join([]string{
					fiber.MethodGet,
					fiber.MethodPost,
					fiber.MethodHead,
					fiber.MethodPut,
					fiber.MethodDelete,
					fiber.MethodPatch,
				}, ","),
				AllowCredentials: false,
				MaxAge:           0,
				AllowHeaders:     "Authorization, Reset",
				ExposeHeaders:    "Authorization, Reset",
			},
		),
	)

	cfg := config.MustConfig(nil)

	if err := cfg.Postgres.MigrationsUp(); err != nil {
		panic(err)
	}

	pool := config.NewPostgres(ctx, cfg.Postgres)
	defer pool.Close()

	// repository group
	authPool := authRepository.NewInsertRepo(pool)
	transactionPool := transaction.NewRepository(pool)
	buyItemPool := inventory.NewInsertRepo(pool)

	// usecase group
	authUC := authUsecase.New(authPool)
	sendCoinUC := sendCoinUseCase.NewUsecase(authPool, transactionPool)
	buyItemUC := buyItemUsecase.NewUsecase(authPool, buyItemPool)
	infoUC := infoUsecase.New(authPool, buyItemPool, transactionPool)

	// handlers group
	authHandler := auth.NewHandler(authUC)
	sendCoinHandler := send_coin.NewHandler(sendCoinUC)
	buyItemHandler := buy_item.NewHandler(buyItemUC)
	infoHandler := info.NewHandler(infoUC)

	// middleware group
	jwtToken := jwt.NewMiddleware(cfg.JWT.Secret)

	api := app.Group("/api")
	api.Post("/auth", authHandler.Handle, jwtToken.SignedToken)
	api.Post("/sendCoin", jwtToken.CompareToken, sendCoinHandler.Handle)
	api.Get("/buy/:item", jwtToken.CompareToken, buyItemHandler.Handle)
	api.Get("/info", jwtToken.CompareToken, infoHandler.Handle)

	log.Println(cfg.App.String())
	if err := app.Listen(cfg.App.String()); err != nil {
		panic("app not start")
	}
}

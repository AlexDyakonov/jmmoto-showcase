package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shampsdev/go-telegram-template/pkg/config"
	"github.com/shampsdev/go-telegram-template/pkg/gateways/rest"
	"github.com/shampsdev/go-telegram-template/pkg/usecase"
	"github.com/shampsdev/go-telegram-template/pkg/utils/slogx"
)

// @title           Cats server
// @version         1.0
// @description     Manage chats, users
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Token
func main() {
	cfg := config.Load(".env")
	log := cfg.Logger()
	log.Info("Hello from Cats server!")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	ctx = slogx.NewCtx(ctx, log)

	pgConfig := cfg.PGXConfig()
	pool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		slogx.Fatal(log, "failed to connect to database", slogx.Err(err))
	}
	defer pool.Close()

	s := rest.NewServer(ctx, cfg, usecase.Setup(ctx, cfg, pool))
	if err := s.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slogx.Fatal(log, "failed to run app", slogx.Err(err))
	}
}

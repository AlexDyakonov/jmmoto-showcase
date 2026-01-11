package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shampsdev/go-telegram-template/pkg/config"
	"github.com/shampsdev/go-telegram-template/pkg/gateways/tg"
	"github.com/shampsdev/go-telegram-template/pkg/utils/slogx"
)

func main() {
	cfg := config.Load(".env")
	log := cfg.Logger()
	log.Info("Hello from Motorcycle Showcase tgbot!")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	ctx = slogx.NewCtx(ctx, log)

	pgConfig := cfg.PGXConfig()
	pool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		slogx.Fatal(log, "failed to connect to database", slogx.Err(err))
	}
	defer pool.Close()

	b, err := tg.NewBot(ctx, cfg, pool)
	if err != nil {
		slogx.Fatal(log, "failed to create telegram bot", slogx.Err(err))
	}

	b.Run(ctx)
}

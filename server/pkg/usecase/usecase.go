package usecase

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shampsdev/go-telegram-template/pkg/config"
	"github.com/shampsdev/go-telegram-template/pkg/parser"
	"github.com/shampsdev/go-telegram-template/pkg/repo/pg"
	"github.com/shampsdev/go-telegram-template/pkg/repo/s3"
)

type Cases struct {
	User       *User
	Motorcycle *Motorcycle
	Analytics  *Analytics
}

func Setup(ctx context.Context, cfg *config.Config, db *pgxpool.Pool) Cases {
	userRepo := pg.NewUserRepo(db)
	motorcycleRepo := pg.NewMotorcycleRepo(db)
	analyticsRepo := pg.NewAnalyticsRepo(db)

	storage, err := s3.NewStorage(cfg.S3)
	if err != nil {
		panic(err)
	}

	motorcycleParser := parser.NewMotorcycleParser()

	userCase := NewUser(ctx, userRepo, storage)
	motorcycleCase := NewMotorcycle(motorcycleRepo, storage, motorcycleParser)
	analyticsCase := NewAnalytics(analyticsRepo)

	return Cases{
		User:       userCase,
		Motorcycle: motorcycleCase,
		Analytics:  analyticsCase,
	}
}

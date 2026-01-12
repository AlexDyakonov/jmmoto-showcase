package repo

import (
	"context"
	"errors"

	"github.com/shampsdev/go-telegram-template/pkg/domain"
)

var (
	ErrNotFound = errors.New("not found")
)

type User interface {
	Create(ctx context.Context, user *domain.CreateUser) (string, error)
	Patch(ctx context.Context, id string, user *domain.PatchUser) error
	Filter(ctx context.Context, filter *domain.FilterUser) ([]*domain.User, error)
	Delete(ctx context.Context, id string) error
}

type Motorcycle interface {
	Create(ctx context.Context, motorcycle *domain.CreateMotorcycle) (string, error)
	Patch(ctx context.Context, id string, motorcycle *domain.PatchMotorcycle) error
	Filter(ctx context.Context, filter *domain.FilterMotorcycle) ([]*domain.Motorcycle, error)
	Delete(ctx context.Context, id string) error
	AddPhotos(ctx context.Context, motorcycleID string, photoURLs []string) error
}

type ImageStorage interface {
	SaveImageByURL(ctx context.Context, url, key string) (string, error)
	SaveImageByBytes(ctx context.Context, bytes []byte, key string) (string, error)
}

type Analytics interface {
	RecordVisit(ctx context.Context, visit *domain.CreateUserVisit) error
	GetUserStats(ctx context.Context, userID string) (*domain.UserVisitStats, error)
	GetAllUserStats(ctx context.Context) ([]*domain.UserVisitStats, error)
}
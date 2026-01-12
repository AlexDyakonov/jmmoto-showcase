package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/shampsdev/go-telegram-template/pkg/domain"
	"github.com/shampsdev/go-telegram-template/pkg/usecase"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

var BotToken string

// contextKey - тип для ключей контекста
type contextKey string

const (
	UserContextKey contextKey = "user"
)

// ExtractUserTGDataFromToken извлекает данные пользователя из Telegram токена
func ExtractUserTGDataFromToken(token string) (*domain.UserTGData, error) {
	if token == "" {
		return nil, fmt.Errorf("missing X-API-Token header")
	}

	expIn := 2 * time.Hour
	err := initdata.Validate(token, BotToken, expIn)
	if err != nil {
		return nil, fmt.Errorf("failed to validate initdata: %w", err)
	}

	parsed, err := initdata.Parse(token)
	if err != nil {
		return nil, fmt.Errorf("failed to parse initdata: %w", err)
	}

	return &domain.UserTGData{
		TelegramID:       parsed.User.ID,
		FirstName:        parsed.User.FirstName,
		LastName:         parsed.User.LastName,
		TelegramUsername: parsed.User.Username,
		Avatar:           parsed.User.PhotoURL,
	}, nil
}

// AuthenticateUserFromToken аутентифицирует пользователя по токену
func AuthenticateUserFromToken(ctx context.Context, token string, userCase *usecase.User) (*domain.User, error) {
	tgUser, err := ExtractUserTGDataFromToken(token)
	if err != nil {
		return nil, err
	}

	return userCase.GetByTGData(ctx, tgUser)
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*domain.User)
	return user, ok
}

// SetUserInContext добавляет пользователя в контекст
func SetUserInContext(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

// AuthenticateFromInput - функция для аутентификации из input структуры (для совместимости)
func AuthenticateFromInput(ctx context.Context, token string, userCase *usecase.User) (*domain.User, error) {
	return AuthenticateUserFromToken(ctx, token, userCase)
}

// RequireAdminFromInput - проверяет аутентификацию и админские права из input
func RequireAdminFromInput(ctx context.Context, token string, userCase *usecase.User) (*domain.User, error) {
	user, err := AuthenticateUserFromToken(ctx, token, userCase)
	if err != nil {
		return nil, err
	}

	if !user.IsAdmin {
		return nil, fmt.Errorf("admin access required: user is not an admin")
	}

	return user, nil
}

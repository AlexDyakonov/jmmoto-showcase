package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/shampsdev/go-telegram-template/pkg/domain"
	"github.com/shampsdev/go-telegram-template/pkg/usecase"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

var BotToken string

func extractUserTGDataFromToken(token string) (*domain.UserTGData, error) {
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

func RequireAdmin(userCase *usecase.User) func(ctx context.Context, input interface{}) error {
	return func(ctx context.Context, input interface{}) error {
		// Получаем токен из заголовка
		var token string
		if req, ok := ctx.Value("request").(*http.Request); ok {
			token = req.Header.Get("X-API-Token")
		}

		// Извлекаем данные пользователя из токена
		tgUser, err := extractUserTGDataFromToken(token)
		if err != nil {
			return huma.Error401Unauthorized("failed to authenticate", err)
		}

		// Получаем пользователя из БД
		user, err := userCase.GetByTGData(ctx, tgUser)
		if err != nil {
			return huma.Error401Unauthorized("failed to get user", err)
		}

		// Проверяем, что пользователь - админ
		if !user.IsAdmin {
			return huma.Error403Forbidden("admin access required", fmt.Errorf("user is not an admin"))
		}

		return nil
	}
}


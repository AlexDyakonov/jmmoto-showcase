package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/shampsdev/go-telegram-template/pkg/gateways/rest/auth"
	"github.com/shampsdev/go-telegram-template/pkg/usecase"
)

// RequireAdmin - middleware для проверки админских прав (используется в операциях Huma)
func RequireAdmin(userCase *usecase.User) func(ctx context.Context, input interface{}) error {
	return func(ctx context.Context, input interface{}) error {
		// Получаем токен из заголовка через контекст запроса
		var token string
		if req, ok := ctx.Value("request").(*http.Request); ok {
			token = req.Header.Get("X-API-Token")
		}

		// Аутентифицируем пользователя
		user, err := auth.AuthenticateUserFromToken(ctx, token, userCase)
		if err != nil {
			return huma.Error401Unauthorized("failed to authenticate", err)
		}

		// Проверяем админские права
		if !user.IsAdmin {
			return huma.Error403Forbidden("admin access required", fmt.Errorf("user is not an admin"))
		}

		return nil
	}
}


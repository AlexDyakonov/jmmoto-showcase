package analytics

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

func authenticateUserFromToken(ctx context.Context, token string, userCase *usecase.User) (*domain.User, error) {
	tgUser, err := extractUserTGDataFromToken(token)
	if err != nil {
		return nil, err
	}

	return userCase.GetByTGData(ctx, tgUser)
}

type RecordVisitInput struct {
	XAPIToken string  `header:"X-API-Token" required:"true" doc:"Telegram init data"`
	Body      struct {
		Source *string `json:"source,omitempty" doc:"Visit source (e.g., channel_post, direct_bot)"`
	} `json:"body"`
}

type RecordVisitOutput struct {
	Body struct {
		Success bool `json:"success"`
	} `json:"body"`
}

func RecordVisitHandler(analyticsCase *usecase.Analytics, userCase *usecase.User) func(ctx context.Context, input *RecordVisitInput) (*RecordVisitOutput, error) {
	return func(ctx context.Context, input *RecordVisitInput) (*RecordVisitOutput, error) {
		user, err := authenticateUserFromToken(ctx, input.XAPIToken, userCase)
		if err != nil {
			return nil, huma.Error401Unauthorized("authentication required", err)
		}

		err = analyticsCase.RecordUserVisit(ctx, user.ID, input.Body.Source)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to record visit", err)
		}

		return &RecordVisitOutput{
			Body: struct {
				Success bool `json:"success"`
			}{Success: true},
		}, nil
	}
}

type GetUserStatsInput struct {
	XAPIToken string `header:"X-API-Token" required:"true" doc:"Telegram init data"`
}

type GetUserStatsOutput struct {
	Body *domain.UserVisitStats `json:"body"`
}

func GetUserStatsHandler(analyticsCase *usecase.Analytics, userCase *usecase.User) func(ctx context.Context, input *GetUserStatsInput) (*GetUserStatsOutput, error) {
	return func(ctx context.Context, input *GetUserStatsInput) (*GetUserStatsOutput, error) {
		user, err := authenticateUserFromToken(ctx, input.XAPIToken, userCase)
		if err != nil {
			return nil, huma.Error401Unauthorized("authentication required", err)
		}

		stats, err := analyticsCase.GetUserStats(ctx, user.ID)
		if err != nil {
			return nil, huma.Error404NotFound("user stats not found", err)
		}

		return &GetUserStatsOutput{Body: stats}, nil
	}
}

type GetAllStatsInput struct {
	XAPIToken string `header:"X-API-Token" required:"true" doc:"Telegram init data"`
}

type GetAllStatsOutput struct {
	Body []*domain.UserVisitStats `json:"body"`
}

func GetAllStatsHandler(analyticsCase *usecase.Analytics, userCase *usecase.User) func(ctx context.Context, input *GetAllStatsInput) (*GetAllStatsOutput, error) {
	return func(ctx context.Context, input *GetAllStatsInput) (*GetAllStatsOutput, error) {
		user, err := authenticateUserFromToken(ctx, input.XAPIToken, userCase)
		if err != nil {
			return nil, huma.Error401Unauthorized("authentication required", err)
		}

		// Проверяем, что пользователь - админ
		if !user.IsAdmin {
			return nil, huma.Error403Forbidden("admin access required", fmt.Errorf("user is not an admin"))
		}

		stats, err := analyticsCase.GetAllUserStats(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to get stats", err)
		}

		return &GetAllStatsOutput{Body: stats}, nil
	}
}

func SetupHuma(api huma.API, cases usecase.Cases) {
	// Записать заход пользователя
	huma.Register(api, huma.Operation{
		OperationID: "record-visit",
		Method:      http.MethodPost,
		Path:        "/analytics/visit",
		Summary:     "Record user visit",
		Tags:        []string{"analytics"},
		Security: []map[string][]string{
			{"ApiKeyAuth": {}},
		},
	}, RecordVisitHandler(cases.Analytics, cases.User))

	// Получить статистику текущего пользователя
	huma.Register(api, huma.Operation{
		OperationID: "get-user-stats",
		Method:      http.MethodGet,
		Path:        "/analytics/my-stats",
		Summary:     "Get current user visit statistics",
		Tags:        []string{"analytics"},
		Security: []map[string][]string{
			{"ApiKeyAuth": {}},
		},
	}, GetUserStatsHandler(cases.Analytics, cases.User))

	// Получить статистику всех пользователей (только для админов)
	huma.Register(api, huma.Operation{
		OperationID: "get-all-stats",
		Method:      http.MethodGet,
		Path:        "/analytics/all-stats",
		Summary:     "Get all users visit statistics (admin only)",
		Tags:        []string{"analytics"},
		Security: []map[string][]string{
			{"ApiKeyAuth": {}},
		},
	}, GetAllStatsHandler(cases.Analytics, cases.User))
}

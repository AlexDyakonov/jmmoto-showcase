package analytics

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/shampsdev/go-telegram-template/pkg/domain"
	"github.com/shampsdev/go-telegram-template/pkg/gateways/rest/auth"
	"github.com/shampsdev/go-telegram-template/pkg/usecase"
)


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
		user, err := auth.AuthenticateFromInput(ctx, input.XAPIToken, userCase)
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
		user, err := auth.AuthenticateFromInput(ctx, input.XAPIToken, userCase)
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
		// Проверяем аутентификацию и админские права
		_, err := auth.RequireAdminFromInput(ctx, input.XAPIToken, userCase)
		if err != nil {
			return nil, huma.Error403Forbidden("admin access required", err)
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

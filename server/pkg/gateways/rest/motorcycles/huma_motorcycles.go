package motorcycles

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

// Public handlers

type GetMotorcyclesInput struct {
	XAPIToken string  `header:"X-API-Token" required:"true" doc:"Telegram init data"`
	Status    string  `query:"status" doc:"Filter by status (available, reserved, sold)"`
	Title     string  `query:"title" doc:"Filter by title (partial match)"`
	MinPrice  float64 `query:"minPrice" doc:"Minimum price"`
	MaxPrice  float64 `query:"maxPrice" doc:"Maximum price"`
}

type GetMotorcyclesOutput struct {
	Body []domain.Motorcycle `json:"body"`
}

func GetMotorcyclesHandler(motorcycleCase *usecase.Motorcycle, userCase *usecase.User) func(ctx context.Context, input *GetMotorcyclesInput) (*GetMotorcyclesOutput, error) {
	return func(ctx context.Context, input *GetMotorcyclesInput) (*GetMotorcyclesOutput, error) {
		_, err := authenticateUserFromToken(ctx, input.XAPIToken, userCase)
		if err != nil {
			return nil, huma.Error401Unauthorized("authentication required", err)
		}

		filter := &domain.FilterMotorcycle{
			IncludePhotos: true,
		}

		if input.Status != "" {
			status := domain.MotorcycleStatus(input.Status)
			filter.Status = &status
		}
		if input.Title != "" {
			filter.Title = &input.Title
		}
		if input.MinPrice > 0 {
			filter.MinPrice = &input.MinPrice
		}
		if input.MaxPrice > 0 {
			filter.MaxPrice = &input.MaxPrice
		}

		motorcycles, err := motorcycleCase.ListMotorcycles(ctx, filter)
		if err != nil {
			return nil, huma.Error400BadRequest("failed to get motorcycles", err)
		}

		result := make([]domain.Motorcycle, len(motorcycles))
		for i, m := range motorcycles {
			result[i] = *m
		}

		return &GetMotorcyclesOutput{Body: result}, nil
	}
}

type GetMotorcycleInput struct {
	XAPIToken string `header:"X-API-Token" required:"true" doc:"Telegram init data"`
	ID        string `path:"id" doc:"Motorcycle ID"`
}

type GetMotorcycleOutput struct {
	Body domain.Motorcycle `json:"body"`
}

func GetMotorcycleHandler(motorcycleCase *usecase.Motorcycle, userCase *usecase.User) func(ctx context.Context, input *GetMotorcycleInput) (*GetMotorcycleOutput, error) {
	return func(ctx context.Context, input *GetMotorcycleInput) (*GetMotorcycleOutput, error) {
		// Проверяем аутентификацию пользователя
		_, err := authenticateUserFromToken(ctx, input.XAPIToken, userCase)
		if err != nil {
			return nil, huma.Error401Unauthorized("authentication required", err)
		}

		motorcycle, err := motorcycleCase.GetMotorcycle(ctx, input.ID)
		if err != nil {
			return nil, huma.Error404NotFound("motorcycle not found", err)
		}

		return &GetMotorcycleOutput{Body: *motorcycle}, nil
	}
}

// Admin handlers

type CreateMotorcycleFromURLInput struct {
	XAPIToken string                      `header:"X-API-Token" required:"true" doc:"Telegram init data"`
	Body      domain.CreateMotorcycleFromURL `json:"body"`
}

type CreateMotorcycleFromURLOutput struct {
	Body domain.Motorcycle `json:"body"`
}

func CreateMotorcycleFromURLHandler(motorcycleCase *usecase.Motorcycle, userCase *usecase.User) func(ctx context.Context, input *CreateMotorcycleFromURLInput) (*CreateMotorcycleFromURLOutput, error) {
	return func(ctx context.Context, input *CreateMotorcycleFromURLInput) (*CreateMotorcycleFromURLOutput, error) {
		user, err := authenticateUserFromToken(ctx, input.XAPIToken, userCase)
		if err != nil {
			return nil, huma.Error401Unauthorized("failed to authenticate user", err)
		}

		if !user.IsAdmin {
			return nil, huma.Error403Forbidden("admin access required", fmt.Errorf("user is not an admin"))
		}

		motorcycle, err := motorcycleCase.CreateMotorcycleFromURL(usecase.NewContext(ctx, user), input.Body.URL)
		if err != nil {
			return nil, huma.Error400BadRequest("failed to create motorcycle from URL", err)
		}

		return &CreateMotorcycleFromURLOutput{Body: *motorcycle}, nil
	}
}

type PatchMotorcycleInput struct {
	XAPIToken string                `header:"X-API-Token" required:"true" doc:"Telegram init data"`
	ID        string                 `path:"id" doc:"Motorcycle ID"`
	Body      domain.PatchMotorcycle `json:"body"`
}

type PatchMotorcycleOutput struct {
	Body domain.Motorcycle `json:"body"`
}

func PatchMotorcycleHandler(motorcycleCase *usecase.Motorcycle, userCase *usecase.User) func(ctx context.Context, input *PatchMotorcycleInput) (*PatchMotorcycleOutput, error) {
	return func(ctx context.Context, input *PatchMotorcycleInput) (*PatchMotorcycleOutput, error) {
		user, err := authenticateUserFromToken(ctx, input.XAPIToken, userCase)
		if err != nil {
			return nil, huma.Error401Unauthorized("failed to authenticate user", err)
		}

		if !user.IsAdmin {
			return nil, huma.Error403Forbidden("admin access required", fmt.Errorf("user is not an admin"))
		}

		motorcycle, err := motorcycleCase.PatchMotorcycle(usecase.NewContext(ctx, user), input.ID, &input.Body)
		if err != nil {
			return nil, huma.Error400BadRequest("failed to update motorcycle", err)
		}

		return &PatchMotorcycleOutput{Body: *motorcycle}, nil
	}
}

type UpdateMotorcycleStatusInput struct {
	XAPIToken string                `header:"X-API-Token" required:"true" doc:"Telegram init data"`
	ID        string                 `path:"id" doc:"Motorcycle ID"`
	Body      struct {
		Status domain.MotorcycleStatus `json:"status"`
	} `json:"body"`
}

type UpdateMotorcycleStatusOutput struct {
	Body domain.Motorcycle `json:"body"`
}

func UpdateMotorcycleStatusHandler(motorcycleCase *usecase.Motorcycle, userCase *usecase.User) func(ctx context.Context, input *UpdateMotorcycleStatusInput) (*UpdateMotorcycleStatusOutput, error) {
	return func(ctx context.Context, input *UpdateMotorcycleStatusInput) (*UpdateMotorcycleStatusOutput, error) {
		user, err := authenticateUserFromToken(ctx, input.XAPIToken, userCase)
		if err != nil {
			return nil, huma.Error401Unauthorized("failed to authenticate user", err)
		}

		if !user.IsAdmin {
			return nil, huma.Error403Forbidden("admin access required", fmt.Errorf("user is not an admin"))
		}

		motorcycle, err := motorcycleCase.UpdateMotorcycleStatus(usecase.NewContext(ctx, user), input.ID, input.Body.Status)
		if err != nil {
			return nil, huma.Error400BadRequest("failed to update motorcycle status", err)
		}

		return &UpdateMotorcycleStatusOutput{Body: *motorcycle}, nil
	}
}

func SetupHuma(api huma.API, cases usecase.Cases) {
	// Public endpoints (require authentication)
	huma.Register(api, huma.Operation{
		OperationID: "get-motorcycles",
		Method:      http.MethodGet,
		Path:        "/motorcycles",
		Summary:     "Get all motorcycles (authenticated users only)",
		Tags:        []string{"motorcycles"},
		Security: []map[string][]string{
			{"ApiKeyAuth": {}},
		},
	}, GetMotorcyclesHandler(cases.Motorcycle, cases.User))

	huma.Register(api, huma.Operation{
		OperationID: "get-motorcycle",
		Method:      http.MethodGet,
		Path:        "/motorcycles/{id}",
		Summary:     "Get motorcycle by ID (authenticated users only)",
		Tags:        []string{"motorcycles"},
		Security: []map[string][]string{
			{"ApiKeyAuth": {}},
		},
	}, GetMotorcycleHandler(cases.Motorcycle, cases.User))

	// Admin endpoints
	huma.Register(api, huma.Operation{
		OperationID: "create-motorcycle-from-url",
		Method:      http.MethodPost,
		Path:        "/admin/motorcycle/from-url",
		Summary:     "Create motorcycle from URL (admin only)",
		Tags:        []string{"admin", "motorcycles"},
		Security: []map[string][]string{
			{"ApiKeyAuth": {}},
		},
	}, CreateMotorcycleFromURLHandler(cases.Motorcycle, cases.User))

	huma.Register(api, huma.Operation{
		OperationID: "patch-motorcycle",
		Method:      http.MethodPatch,
		Path:        "/admin/motorcycle/{id}",
		Summary:     "Update motorcycle (admin only)",
		Tags:        []string{"admin", "motorcycles"},
		Security: []map[string][]string{
			{"ApiKeyAuth": {}},
		},
	}, PatchMotorcycleHandler(cases.Motorcycle, cases.User))

	huma.Register(api, huma.Operation{
		OperationID: "update-motorcycle-status",
		Method:      http.MethodPatch,
		Path:        "/admin/motorcycle/{id}/status",
		Summary:     "Update motorcycle status (admin only)",
		Tags:        []string{"admin", "motorcycles"},
		Security: []map[string][]string{
			{"ApiKeyAuth": {}},
		},
	}, UpdateMotorcycleStatusHandler(cases.Motorcycle, cases.User))
}


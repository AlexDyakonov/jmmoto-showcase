package cats

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

type GetAllCatsOutput struct {
	Body []domain.Cat `json:"body"`
}

type CreateCatInput struct {
	XAPIToken string            `header:"X-API-Token" required:"true" doc:"Telegram init data"`
	Body      domain.CreateCat  `json:"body"`
}

type CreateCatOutput struct {
	Body domain.Cat `json:"body"`
}

type MeCatsInput struct {
	XAPIToken string `header:"X-API-Token" required:"true" doc:"Telegram init data"`
}

type MeCatsOutput struct {
	Body []domain.Cat `json:"body"`
}

type PatchCatInput struct {
	XAPIToken string          `header:"X-API-Token" required:"true" doc:"Telegram init data"`
	ID        string          `path:"id" doc:"Cat ID"`
	Body      domain.PatchCat `json:"body"`
}

type PatchCatOutput struct {
	Body domain.Cat `json:"body"`
}

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

func GetAllCatsHandler(catCase *usecase.Cat) func(ctx context.Context, input *struct{}) (*GetAllCatsOutput, error) {
	return func(ctx context.Context, input *struct{}) (*GetAllCatsOutput, error) {
		cats, err := catCase.ListAllCats(ctx)
		if err != nil {
			return nil, huma.Error400BadRequest("failed to get cats", err)
		}

		// Convert []*domain.Cat to []domain.Cat
		result := make([]domain.Cat, len(cats))
		for i, cat := range cats {
			result[i] = *cat
		}
		return &GetAllCatsOutput{Body: result}, nil
	}
}

func CreateCatHandler(catCase *usecase.Cat, userCase *usecase.User) func(ctx context.Context, input *CreateCatInput) (*CreateCatOutput, error) {
	return func(ctx context.Context, input *CreateCatInput) (*CreateCatOutput, error) {
		user, err := authenticateUserFromToken(ctx, input.XAPIToken, userCase)
		if err != nil {
			return nil, huma.Error401Unauthorized("failed to authenticate user", err)
		}

		cat, err := catCase.CreateCat(usecase.NewContext(ctx, user), &input.Body)
		if err != nil {
			return nil, huma.Error400BadRequest("failed to create cat", err)
		}

		return &CreateCatOutput{Body: *cat}, nil
	}
}

func MeCatsHandler(catCase *usecase.Cat, userCase *usecase.User) func(ctx context.Context, input *MeCatsInput) (*MeCatsOutput, error) {
	return func(ctx context.Context, input *MeCatsInput) (*MeCatsOutput, error) {
		user, err := authenticateUserFromToken(ctx, input.XAPIToken, userCase)
		if err != nil {
			return nil, huma.Error401Unauthorized("failed to authenticate user", err)
		}

		cats, err := catCase.MeCats(usecase.NewContext(ctx, user))
		if err != nil {
			return nil, huma.Error400BadRequest("failed to get user cats", err)
		}

		// Convert []*domain.Cat to []domain.Cat
		result := make([]domain.Cat, len(cats))
		for i, cat := range cats {
			result[i] = *cat
		}
		return &MeCatsOutput{Body: result}, nil
	}
}

func PatchCatHandler(catCase *usecase.Cat, userCase *usecase.User) func(ctx context.Context, input *PatchCatInput) (*PatchCatOutput, error) {
	return func(ctx context.Context, input *PatchCatInput) (*PatchCatOutput, error) {
		user, err := authenticateUserFromToken(ctx, input.XAPIToken, userCase)
		if err != nil {
			return nil, huma.Error401Unauthorized("failed to authenticate user", err)
		}

		cat, err := catCase.PatchCat(usecase.NewContext(ctx, user), input.ID, &input.Body)
		if err != nil {
			return nil, huma.Error400BadRequest("failed to update cat", err)
		}

		return &PatchCatOutput{Body: *cat}, nil
	}
}

func SetupHuma(api huma.API, cases usecase.Cases) {
	// GET /cats - Get all cats
	huma.Register(api, huma.Operation{
		OperationID: "get-all-cats",
		Method:      http.MethodGet,
		Path:        "/cats",
		Summary:     "Get all cats",
		Tags:        []string{"cats"},
	}, GetAllCatsHandler(cases.Cat))

	// POST /cats - Create cat
	huma.Register(api, huma.Operation{
		OperationID: "create-cat",
		Method:      http.MethodPost,
		Path:        "/cats",
		Summary:     "Create cat",
		Tags:        []string{"cats"},
		Security: []map[string][]string{
			{"ApiKeyAuth": {}},
		},
	}, CreateCatHandler(cases.Cat, cases.User))

	// GET /cats/me - Get my cats
	huma.Register(api, huma.Operation{
		OperationID: "get-me-cats",
		Method:      http.MethodGet,
		Path:        "/cats/me",
		Summary:     "Get my cats",
		Tags:        []string{"cats"},
		Security: []map[string][]string{
			{"ApiKeyAuth": {}},
		},
	}, MeCatsHandler(cases.Cat, cases.User))

	// PATCH /cats/id/{id} - Update cat
	huma.Register(api, huma.Operation{
		OperationID: "patch-cat",
		Method:      http.MethodPatch,
		Path:        "/cats/id/{id}",
		Summary:     "Update cat",
		Tags:        []string{"cats"},
		Security: []map[string][]string{
			{"ApiKeyAuth": {}},
		},
	}, PatchCatHandler(cases.Cat, cases.User))
}

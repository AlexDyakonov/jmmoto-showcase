package user

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

type CreateMeInput struct {
	XAPIToken string             `header:"X-API-Token" required:"true" doc:"Telegram init data"`
	Body      domain.CreateUser  `json:"body"`
}

type CreateMeOutput struct {
	Body domain.User `json:"body"`
}

type GetMeInput struct {
	XAPIToken string `header:"X-API-Token" required:"true" doc:"Telegram init data"`
}

type GetMeOutput struct {
	Body domain.User `json:"body"`
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

func CreateMeHandler(userCase *usecase.User) func(ctx context.Context, input *CreateMeInput) (*CreateMeOutput, error) {
	return func(ctx context.Context, input *CreateMeInput) (*CreateMeOutput, error) {
		tgUser, err := extractUserTGDataFromToken(input.XAPIToken)
		if err != nil {
			return nil, huma.Error401Unauthorized("failed to validate user", err)
		}
		
		user, err := userCase.CreateMe(usecase.NewContextWithTGData(ctx, tgUser), &input.Body)
		if err != nil {
			return nil, huma.Error400BadRequest("failed to create user", err)
		}

		return &CreateMeOutput{Body: *user}, nil
	}
}

func GetMeHandler(userCase *usecase.User) func(ctx context.Context, input *GetMeInput) (*GetMeOutput, error) {
	return func(ctx context.Context, input *GetMeInput) (*GetMeOutput, error) {
		user, err := authenticateUserFromToken(ctx, input.XAPIToken, userCase)
		if err != nil {
			return nil, huma.Error401Unauthorized("failed to authenticate user", err)
		}
		
		userProfile, err := userCase.GetMe(usecase.NewContext(ctx, user))
		if err != nil {
			return nil, huma.Error400BadRequest("failed to get user", err)
		}

		return &GetMeOutput{Body: *userProfile}, nil
	}
}

func SetupHuma(api huma.API, cases usecase.Cases) {
	// POST /users/me - Create user
	huma.Register(api, huma.Operation{
		OperationID: "create-me",
		Method:      http.MethodPost,
		Path:        "/users/me",
		Summary:     "Create me",
		Tags:        []string{"users"},
		Security: []map[string][]string{
			{"ApiKeyAuth": {}},
		},
	}, CreateMeHandler(cases.User))

	// GET /users/me - Get current user
	huma.Register(api, huma.Operation{
		OperationID: "get-me",
		Method:      http.MethodGet,
		Path:        "/users/me",
		Summary:     "Get me",
		Tags:        []string{"users"},
		Security: []map[string][]string{
			{"ApiKeyAuth": {}},
		},
	}, GetMeHandler(cases.User))
}

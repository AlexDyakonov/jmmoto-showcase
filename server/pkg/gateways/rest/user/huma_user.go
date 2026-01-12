package user

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/shampsdev/go-telegram-template/pkg/domain"
	"github.com/shampsdev/go-telegram-template/pkg/gateways/rest/auth"
	"github.com/shampsdev/go-telegram-template/pkg/usecase"
)

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


func CreateMeHandler(userCase *usecase.User) func(ctx context.Context, input *CreateMeInput) (*CreateMeOutput, error) {
	return func(ctx context.Context, input *CreateMeInput) (*CreateMeOutput, error) {
		tgUser, err := auth.ExtractUserTGDataFromToken(input.XAPIToken)
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
		user, err := auth.AuthenticateFromInput(ctx, input.XAPIToken, userCase)
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

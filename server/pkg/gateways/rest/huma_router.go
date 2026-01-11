package rest

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/shampsdev/go-telegram-template/pkg/gateways/rest/motorcycles"
	"github.com/shampsdev/go-telegram-template/pkg/gateways/rest/user"
	"github.com/shampsdev/go-telegram-template/pkg/usecase"
)

func setupHumaRouter(api huma.API, useCases usecase.Cases) {
	user.SetupHuma(api, useCases)
	motorcycles.SetupHuma(api, useCases)
}

func NewHumaAPI(ctx context.Context, useCases usecase.Cases) (huma.API, *chi.Mux) {
	router := chi.NewMux()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Api-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	apiRouter := chi.NewRouter()

	config := huma.DefaultConfig("Motorcycle Showcase API", "1.0.0")
	config.Info.Description = "Manage motorcycles showcase"
	
	config.Servers = []*huma.Server{
		{URL: "/api/v1"},
	}

	config.OpenAPIPath = "/openapi"
	config.DocsPath = "/docs"
	config.SchemasPath = "/schemas"

	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"ApiKeyAuth": {
			Type: "apiKey",
			In:   "header",
			Name: "X-API-Token",
		},
	}

	api := humachi.New(apiRouter, config)
	
	setupHumaRouter(api, useCases)

	router.Mount("/api/v1", apiRouter)

	return api, router
}

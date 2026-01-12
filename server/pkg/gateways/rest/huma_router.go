package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/shampsdev/go-telegram-template/pkg/gateways/rest/analytics"
	"github.com/shampsdev/go-telegram-template/pkg/gateways/rest/motorcycles"
	"github.com/shampsdev/go-telegram-template/pkg/gateways/rest/user"
	"github.com/shampsdev/go-telegram-template/pkg/usecase"
	"github.com/shampsdev/go-telegram-template/pkg/utils/slogx"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)


func extractUserInfoFromToken(token string) string {
	if token == "" {
		return "anonymous"
	}
	
	parsed, err := initdata.Parse(token)
	if err != nil {
		return "invalid_token"
	}
	
	if parsed.User.ID == 0 {
		return "no_user_data"
	}
	
	userInfo := fmt.Sprintf("tg_id:%d", parsed.User.ID)
	if parsed.User.Username != "" {
		userInfo += fmt.Sprintf(" @%s", parsed.User.Username)
	}
	if parsed.User.FirstName != "" {
		userInfo += fmt.Sprintf(" (%s", parsed.User.FirstName)
		if parsed.User.LastName != "" {
			userInfo += fmt.Sprintf(" %s", parsed.User.LastName)
		}
		userInfo += ")"
	}
	
	return userInfo
}

func setupHumaRouter(api huma.API, useCases usecase.Cases) {
	user.SetupHuma(api, useCases)
	motorcycles.SetupHuma(api, useCases)
	analytics.SetupHuma(api, useCases)
}

func NewHumaAPI(ctx context.Context, useCases usecase.Cases) (huma.API, *chi.Mux) {
	router := chi.NewMux()
	log := slogx.FromCtx(ctx)

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			
			next.ServeHTTP(ww, r)
			
			duration := time.Since(start)
			
			// Извлекаем информацию о пользователе из Telegram токена
			userInfo := "anonymous"
			if token := r.Header.Get("X-API-Token"); token != "" {
				userInfo = extractUserInfoFromToken(token)
			}
			
			log.Info("HTTP Request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("query", r.URL.RawQuery),
				slog.Int("status", ww.Status()),
				slog.Duration("duration", duration),
				slog.String("user", userInfo),
				slog.String("user_agent", r.UserAgent()),
				slog.String("remote_addr", r.RemoteAddr),
			)
		})
	})

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

package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"github.com/shampsdev/go-telegram-template/pkg/config"
	"github.com/shampsdev/go-telegram-template/pkg/gateways/rest/auth"
	"github.com/shampsdev/go-telegram-template/pkg/usecase"
	"github.com/tj/go-spin"
	"golang.org/x/sync/errgroup"
)

const shutdownDuration = 1500 * time.Millisecond

type Server struct {
	HttpServer http.Server
	API        huma.API
	Router     *chi.Mux
}

func NewServer(ctx context.Context, cfg *config.Config, useCases usecase.Cases) *Server {
	// Set bot token for authentication
	auth.BotToken = cfg.TG.BotToken
	
	api, router := NewHumaAPI(ctx, useCases)

	s := &Server{
		API:    api,
		Router: router,
		HttpServer: http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
			Handler: router,
		},
	}

	return s
}

func (s *Server) Run(ctx context.Context) error {
	eg := errgroup.Group{}
	eg.Go(func() error {
		return s.HttpServer.ListenAndServe()
	})

	<-ctx.Done()
	err := s.HttpServer.Shutdown(ctx)
	shutdownWait()
	return err
}

func shutdownWait() {
	spinner := spin.New()
	const spinIterations = 20
	for range spinIterations {
		fmt.Printf("\rgraceful shutdown %s ", spinner.Next())
		time.Sleep(shutdownDuration / spinIterations)
	}
	fmt.Println()
}

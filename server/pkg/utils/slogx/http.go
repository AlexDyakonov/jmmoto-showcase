package slogx

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func InjectHTTP(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.New().String()
			logger := log.With(slog.String("request_id", requestID))
			ctx := NewCtx(r.Context(), logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

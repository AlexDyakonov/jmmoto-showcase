package usecase

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/shampsdev/go-telegram-template/pkg/domain"
	"github.com/shampsdev/go-telegram-template/pkg/repo"
)

type Analytics struct {
	analyticsRepo repo.Analytics
}

func NewAnalytics(analyticsRepo repo.Analytics) *Analytics {
	return &Analytics{
		analyticsRepo: analyticsRepo,
	}
}

func (a *Analytics) RecordUserVisit(ctx context.Context, userID string, source *string) error {
	sessionID := generateSessionID()
	
	visit := &domain.CreateUserVisit{
		UserID:    userID,
		SessionID: sessionID,
		Source:    source,
	}

	return a.analyticsRepo.RecordVisit(ctx, visit)
}

func (a *Analytics) GetUserStats(ctx context.Context, userID string) (*domain.UserVisitStats, error) {
	return a.analyticsRepo.GetUserStats(ctx, userID)
}

func (a *Analytics) GetAllUserStats(ctx context.Context) ([]*domain.UserVisitStats, error) {
	return a.analyticsRepo.GetAllUserStats(ctx)
}

func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

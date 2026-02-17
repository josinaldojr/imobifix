package service

import (
	"context"
	"time"

	"github.com/josinaldojr/imobifix-api/internal/domain"
	"github.com/josinaldojr/imobifix-api/internal/repo"
)

type QuotesService struct {
	db *repo.DB
}

func NewQuotesService(db *repo.DB) *QuotesService {
	return &QuotesService{db: db}
}

func (s *QuotesService) Create(ctx context.Context, brlToUsd float64, effectiveAt *time.Time) (domain.Quote, error) {
	t := time.Now().UTC()
	if effectiveAt != nil {
		t = effectiveAt.UTC()
	}
	return s.db.CreateQuote(ctx, brlToUsd, t)
}

func (s *QuotesService) Current(ctx context.Context) (*domain.Quote, error) {
	return s.db.GetCurrentQuote(ctx)
}

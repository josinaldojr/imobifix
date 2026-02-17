package service

import (
	"context"
	"time"

	"github.com/josinaldojr/imobifx-api/internal/domain"
)

type QuotesService struct {
	db QuotesRepository
}

func NewQuotesService(db QuotesRepository) *QuotesService {
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

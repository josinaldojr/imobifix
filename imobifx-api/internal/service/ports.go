package service

import (
	"context"
	"time"

	"github.com/josinaldojr/imobifix-api/internal/domain"
	"github.com/josinaldojr/imobifix-api/internal/repo"
)

type AdsRepository interface {
	CreateAd(ctx context.Context, ad domain.Ad) (domain.Ad, error)
	ListAds(ctx context.Context, f repo.AdsFilter, page, pageSize int) ([]domain.Ad, int, error)
	GetCurrentQuote(ctx context.Context) (*domain.Quote, error)
}

type QuotesRepository interface {
	CreateQuote(ctx context.Context, brlToUsd float64, effectiveAt time.Time) (domain.Quote, error)
	GetCurrentQuote(ctx context.Context) (*domain.Quote, error)
}

type ViaCEPClient interface {
	Lookup(ctx context.Context, cep8digits string) (domain.Address, error)
}

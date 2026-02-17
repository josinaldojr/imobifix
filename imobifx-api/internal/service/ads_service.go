package service

import (
	"context"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/josinaldojr/imobifix-api/internal/domain"
	"github.com/josinaldojr/imobifix-api/internal/repo"
)

type AdsService struct {
	db        *repo.DB
	imagesDir string
}

func NewAdsService(db *repo.DB, imagesDir string) *AdsService {
	_ = os.MkdirAll(imagesDir, 0o755)
	return &AdsService{db: db, imagesDir: imagesDir}
}

func (s *AdsService) Create(ctx context.Context, ad domain.Ad) (domain.Ad, error) {
	return s.db.CreateAd(ctx, ad)
}

func (s *AdsService) List(ctx context.Context, f repo.AdsFilter, page, pageSize int, quote *domain.Quote) (domain.AdsListResponse, error) {
	ads, total, err := s.db.ListAds(ctx, f, page, pageSize)
	if err != nil {
		return domain.AdsListResponse{}, err
	}

	resp := domain.AdsListResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Items:    make([]domain.AdItem, 0, len(ads)),
	}

	if quote != nil {
		resp.QuoteUsed = &struct {
			BrlToUsd    float64   `json:"brl_to_usd"`
			EffectiveAt time.Time `json:"effective_at"`
		}{
			BrlToUsd:    quote.BrlToUsd,
			EffectiveAt: quote.EffectiveAt,
		}
	}

	for _, a := range ads {
		item := domain.AdItem{
			ID:        a.ID,
			Type:      a.Type,
			PriceBRL:  a.PriceBRL,
			CreatedAt: a.CreatedAt,
		}

		item.Address.CEP = a.CEP
		item.Address.Street = a.Street
		item.Address.Number = a.Number
		item.Address.Complement = a.Complement
		item.Address.Neighborhood = a.Neighborhood
		item.Address.City = a.City
		item.Address.State = a.State

		if a.ImagePath != nil && strings.TrimSpace(*a.ImagePath) != "" {
			u := "/static/images/" + *a.ImagePath
			item.ImageURL = &u
		}

		if quote != nil {
			v := round2(a.PriceBRL * quote.BrlToUsd)
			item.PriceUSD = &v
		}

		resp.Items = append(resp.Items, item)
	}

	return resp, nil
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func (s *AdsService) SaveImage(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = ".bin"
	}
	return uuid.NewString() + ext
}

func (s *AdsService) ImagePath(name string) string {
	return filepath.Join(s.imagesDir, name)
}

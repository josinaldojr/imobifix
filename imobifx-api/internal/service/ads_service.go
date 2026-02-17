package service

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/josinaldojr/imobifix-api/internal/domain"
	"github.com/josinaldojr/imobifix-api/internal/repo"
	"github.com/josinaldojr/imobifix-api/internal/usecase"
	"github.com/josinaldojr/imobifix-api/internal/validation"
)

type AdsService struct {
	db           AdsRepository
	imagesDir    string
	maxImageSize int64
}

func NewAdsService(db AdsRepository, imagesDir string, maxImageSize int64) *AdsService {
	_ = os.MkdirAll(imagesDir, 0o755)
	return &AdsService{db: db, imagesDir: imagesDir, maxImageSize: maxImageSize}
}

func (s *AdsService) Create(ctx context.Context, in usecase.CreateAdInput, image *multipart.FileHeader) (domain.Ad, error) {
	if err := validation.ValidateCreateAdInput(&in); err != nil {
		return domain.Ad{}, err
	}
	if err := validation.ValidateImage(image, s.maxImageSize); err != nil {
		return domain.Ad{}, err
	}

	var imageName *string
	if image != nil {
		name, err := s.saveImage(image)
		if err != nil {
			return domain.Ad{}, err
		}
		imageName = &name
	}

	ad := domain.Ad{
		Type:         in.Type,
		PriceBRL:     in.PriceBRL,
		ImagePath:    imageName,
		CEP:          in.CEP,
		Street:       in.Street,
		Number:       in.Number,
		Complement:   in.Complement,
		Neighborhood: in.Neighborhood,
		City:         in.City,
		State:        in.State,
	}

	return s.db.CreateAd(ctx, ad)
}

func (s *AdsService) List(ctx context.Context, in usecase.ListAdsInput) (domain.AdsListResponse, error) {
	if err := validation.ValidateListAdsInput(in); err != nil {
		return domain.AdsListResponse{}, err
	}

	var f repo.AdsFilter
	f.Type = in.Type
	f.City = in.City
	f.State = in.State
	f.MinPrice = in.MinPrice
	f.MaxPrice = in.MaxPrice

	quote, err := s.db.GetCurrentQuote(ctx)
	if err != nil {
		return domain.AdsListResponse{}, err
	}

	ads, total, err := s.db.ListAds(ctx, f, in.Page, in.PageSize)
	if err != nil {
		return domain.AdsListResponse{}, err
	}

	resp := domain.AdsListResponse{
		Page:     in.Page,
		PageSize: in.PageSize,
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
		resp.Items = append(resp.Items, domain.ToAdItemWithQuote(a, quote))
	}

	return resp, nil
}

func (s *AdsService) saveImage(file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		ext = ".bin"
	}
	name := uuid.NewString() + ext
	dstPath := filepath.Join(s.imagesDir, name)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}
	return name, nil
}

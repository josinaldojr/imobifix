package service

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/josinaldojr/imobifix-api/internal/domain"
	"github.com/josinaldojr/imobifix-api/internal/repo"
	"github.com/josinaldojr/imobifix-api/internal/validation"
)

type AdsService struct {
	db           *repo.DB
	imagesDir    string
	maxImageSize int64
}

func NewAdsService(db *repo.DB, imagesDir string, maxImageSize int64) *AdsService {
	_ = os.MkdirAll(imagesDir, 0o755)
	return &AdsService{db: db, imagesDir: imagesDir, maxImageSize: maxImageSize}
}

func (s *AdsService) Create(ctx context.Context, in domain.CreateAdInput, image *multipart.FileHeader) (domain.Ad, error) {
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

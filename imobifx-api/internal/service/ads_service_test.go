package service_test

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/josinaldojr/imobifix-api/internal/domain"
	"github.com/josinaldojr/imobifix-api/internal/repo"
	"github.com/josinaldojr/imobifix-api/internal/service"
	"github.com/josinaldojr/imobifix-api/internal/usecase"
)

type fakeAdsRepo struct {
	lastCreated domain.Ad

	createCalled bool
	listCalled   bool
	quoteCalled  bool

	createFn func(ctx context.Context, ad domain.Ad) (domain.Ad, error)
	listFn   func(ctx context.Context, f repo.AdsFilter, page, pageSize int) ([]domain.Ad, int, error)
	quoteFn  func(ctx context.Context) (*domain.Quote, error)
}

func (f *fakeAdsRepo) CreateAd(ctx context.Context, ad domain.Ad) (domain.Ad, error) {
	f.createCalled = true
	f.lastCreated = ad
	if f.createFn != nil {
		return f.createFn(ctx, ad)
	}
	ad.ID = "ad-1"
	ad.CreatedAt = time.Now().UTC()
	return ad, nil
}

func (f *fakeAdsRepo) ListAds(ctx context.Context, flt repo.AdsFilter, page, pageSize int) ([]domain.Ad, int, error) {
	f.listCalled = true
	if f.listFn != nil {
		return f.listFn(ctx, flt, page, pageSize)
	}
	return nil, 0, nil
}

func (f *fakeAdsRepo) GetCurrentQuote(ctx context.Context) (*domain.Quote, error) {
	f.quoteCalled = true
	if f.quoteFn != nil {
		return f.quoteFn(ctx)
	}
	return nil, nil
}

func TestAdsService_Create_OK_NoImage_FormatsCEP_AndCallsRepo(t *testing.T) {
	db := &fakeAdsRepo{}
	tmp := t.TempDir()

	svc := service.NewAdsService(db, tmp, 5*1024*1024)

	in := usecase.CreateAdInput{
		Type:         "SALE",
		PriceBRL:     250000,
		CEP:          "58000000",
		Street:       "Rua A",
		Neighborhood: "Centro",
		City:         "João Pessoa",
		State:        "PB",
	}

	_, err := svc.Create(context.Background(), in, nil)
	require.NoError(t, err)
	require.True(t, db.createCalled)

	require.Equal(t, "58000-000", db.lastCreated.CEP)
	require.Nil(t, db.lastCreated.ImagePath)
}

func TestAdsService_Create_InvalidInput_DoesNotCallRepo(t *testing.T) {
	db := &fakeAdsRepo{}
	svc := service.NewAdsService(db, t.TempDir(), 5*1024*1024)

	in := usecase.CreateAdInput{
		Type:         "X",
		PriceBRL:     -1,
		CEP:          "123",
		Street:       "",
		Neighborhood: "",
		City:         "",
		State:        "P",
	}

	_, err := svc.Create(context.Background(), in, nil)
	require.Error(t, err)
	require.False(t, db.createCalled)
}

func TestAdsService_Create_WithImage_SavesFile_AndSetsImagePath(t *testing.T) {
	db := &fakeAdsRepo{}
	tmp := t.TempDir()

	svc := service.NewAdsService(db, tmp, 5*1024*1024)

	in := usecase.CreateAdInput{
		Type:         "RENT",
		PriceBRL:     1500,
		CEP:          "58000-000",
		Street:       "Rua B",
		Neighborhood: "Bairro",
		City:         "João Pessoa",
		State:        "PB",
	}

	fh := makeMultipartFileHeader(t, "image", "house.jpg", "image/jpeg", []byte("fake-jpeg-bytes"))

	_, err := svc.Create(context.Background(), in, fh)
	require.NoError(t, err)
	require.True(t, db.createCalled)

	require.NotNil(t, db.lastCreated.ImagePath)
	require.NotEmpty(t, *db.lastCreated.ImagePath)

	_, statErr := os.Stat(filepath.Join(tmp, *db.lastCreated.ImagePath))
	require.NoError(t, statErr)
}

func TestAdsService_List_SetsQuoteUsed_AndReturnsItems(t *testing.T) {
	db := &fakeAdsRepo{
		quoteFn: func(ctx context.Context) (*domain.Quote, error) {
			return &domain.Quote{
				ID:          "q1",
				BrlToUsd:    0.2,
				EffectiveAt: time.Date(2026, 2, 16, 10, 0, 0, 0, time.UTC),
				CreatedAt:   time.Now().UTC(),
			}, nil
		},
		listFn: func(ctx context.Context, f repo.AdsFilter, page, pageSize int) ([]domain.Ad, int, error) {
			return []domain.Ad{
				{
					ID:           "ad-1",
					Type:         "SALE",
					PriceBRL:     100,
					CEP:          "58000-000",
					Street:       "Rua A",
					Neighborhood: "Centro",
					City:         "João Pessoa",
					State:        "PB",
					CreatedAt:    time.Now().UTC(),
				},
			}, 1, nil
		},
	}

	svc := service.NewAdsService(db, t.TempDir(), 5*1024*1024)

	typ := "SALE"
	in := usecase.ListAdsInput{
		Page:     1,
		PageSize: 10,
		Type:     &typ,
	}

	resp, err := svc.List(context.Background(), in)
	require.NoError(t, err)

	require.True(t, db.quoteCalled)
	require.True(t, db.listCalled)

	require.Equal(t, 1, resp.Page)
	require.Equal(t, 10, resp.PageSize)
	require.Equal(t, 1, resp.Total)
	require.Len(t, resp.Items, 1)

	require.NotNil(t, resp.QuoteUsed)
	require.Equal(t, 0.2, resp.QuoteUsed.BrlToUsd)
	require.True(t, resp.QuoteUsed.EffectiveAt.Equal(time.Date(2026, 2, 16, 10, 0, 0, 0, time.UTC)))
}

func makeMultipartFileHeader(t *testing.T, field, filename, contentType string, data []byte) *multipart.FileHeader {
	t.Helper()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, field, filename))
	h.Set("Content-Type", contentType)

	part, err := w.CreatePart(h)
	require.NoError(t, err)

	_, err = part.Write(data)
	require.NoError(t, err)

	require.NoError(t, w.Close())

	req := httptest.NewRequest("POST", "/", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	require.NoError(t, req.ParseMultipartForm(10<<20))

	fhs := req.MultipartForm.File[field]
	require.Len(t, fhs, 1)

	require.Equal(t, contentType, fhs[0].Header.Get("Content-Type"))

	return fhs[0]
}

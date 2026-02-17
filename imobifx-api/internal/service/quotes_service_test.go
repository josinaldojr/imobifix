package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/josinaldojr/imobifix-api/internal/domain"
	"github.com/josinaldojr/imobifix-api/internal/service"
)

type fakeQuotesRepo struct {
	createCalled  bool
	currentCalled bool
	lastRate      float64
	lastEff       time.Time

	createFn  func(ctx context.Context, rate float64, eff time.Time) (domain.Quote, error)
	currentFn func(ctx context.Context) (*domain.Quote, error)
}

func (f *fakeQuotesRepo) CreateQuote(ctx context.Context, rate float64, eff time.Time) (domain.Quote, error) {
	f.createCalled = true
	f.lastRate = rate
	f.lastEff = eff
	if f.createFn != nil {
		return f.createFn(ctx, rate, eff)
	}
	return domain.Quote{ID: "q1", BrlToUsd: rate, EffectiveAt: eff, CreatedAt: time.Now().UTC()}, nil
}

func (f *fakeQuotesRepo) GetCurrentQuote(ctx context.Context) (*domain.Quote, error) {
	f.currentCalled = true
	if f.currentFn != nil {
		return f.currentFn(ctx)
	}
	return nil, nil
}

func TestQuotesService_Create_UsesProvidedEffectiveAtUTC(t *testing.T) {
	db := &fakeQuotesRepo{}
	svc := service.NewQuotesService(db)

	eff := time.Date(2026, 2, 16, 10, 0, 0, 0, time.FixedZone("X", -3*3600))
	_, err := svc.Create(context.Background(), 0.19, &eff)
	require.NoError(t, err)

	require.True(t, db.createCalled)
	require.Equal(t, 0.19, db.lastRate)
	require.True(t, db.lastEff.Equal(eff.UTC()))
}

func TestQuotesService_Current(t *testing.T) {
	expected := &domain.Quote{ID: "q1", BrlToUsd: 0.2, EffectiveAt: time.Now().UTC(), CreatedAt: time.Now().UTC()}
	db := &fakeQuotesRepo{
		currentFn: func(ctx context.Context) (*domain.Quote, error) { return expected, nil },
	}
	svc := service.NewQuotesService(db)

	got, err := svc.Current(context.Background())
	require.NoError(t, err)
	require.True(t, db.currentCalled)
	require.Equal(t, expected.ID, got.ID)
}

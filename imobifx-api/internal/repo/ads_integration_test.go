//go:build integration

package repo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/josinaldojr/imobifx-api/internal/domain"
	"github.com/josinaldojr/imobifx-api/internal/repo"
)

func TestAds_CreateAndList_WithFilters(t *testing.T) {
	dsn := testDSN()
	if dsn == "" {
		t.Skip("TEST_DB_DSN/DB_DSN not set")
	}

	db, err := repo.NewPostgres(dsn)
	require.NoError(t, err)
	defer db.Close()

	_, _ = db.Pool.Exec(context.Background(), "TRUNCATE TABLE ads RESTART IDENTITY")
	_, _ = db.Pool.Exec(context.Background(), "TRUNCATE TABLE quotes RESTART IDENTITY")

	_, err = db.CreateAd(context.Background(), domain.Ad{
		Type:         "SALE",
		PriceBRL:     250000,
		CEP:          "58000-000",
		Street:       "Rua A",
		Neighborhood: "Centro",
		City:         "Joao Pessoa",
		State:        "PB",
	})
	require.NoError(t, err)

	_, err = db.CreateAd(context.Background(), domain.Ad{
		Type:         "RENT",
		PriceBRL:     1800,
		CEP:          "58000-000",
		Street:       "Rua B",
		Neighborhood: "Bairro",
		City:         "Joao Pessoa",
		State:        "PB",
	})
	require.NoError(t, err)

	typ := "SALE"
	filter := repo.AdsFilter{Type: &typ}
	items, total, err := db.ListAds(context.Background(), filter, 1, 10)
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, items, 1)
	require.Equal(t, "SALE", items[0].Type)
}

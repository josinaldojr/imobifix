//go:build integration

package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/josinaldojr/imobifx-api/internal/repo"
)

func TestQuotes_CreateAndGetCurrent(t *testing.T) {
	dsn := testDSN()
	if dsn == "" {
		t.Skip("TEST_DB_DSN/DB_DSN not set")
	}

	db, err := repo.NewPostgres(dsn)
	require.NoError(t, err)
	defer db.Close()

	_, _ = db.Pool.Exec(context.Background(), "TRUNCATE TABLE quotes RESTART IDENTITY")

	q1, err := db.CreateQuote(context.Background(), 0.19, time.Date(2026, 2, 16, 10, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.NotEmpty(t, q1.ID)

	q2, err := db.CreateQuote(context.Background(), 0.20, time.Date(2026, 2, 17, 10, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	cur, err := db.GetCurrentQuote(context.Background())
	require.NoError(t, err)
	require.NotNil(t, cur)
	require.Equal(t, q2.ID, cur.ID)
	require.Equal(t, 0.20, cur.BrlToUsd)
}

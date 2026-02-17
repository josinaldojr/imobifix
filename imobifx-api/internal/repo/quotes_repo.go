package repo

import (
	"context"
	"time"

	"github.com/josinaldojr/imobifx-api/internal/domain"
)

func (d *DB) CreateQuote(ctx context.Context, brlToUsd float64, effectiveAt time.Time) (domain.Quote, error) {
	row := d.Pool.QueryRow(ctx, `
		INSERT INTO quotes (brl_to_usd, effective_at)
		VALUES ($1, $2)
		RETURNING id, brl_to_usd, effective_at, created_at
	`, brlToUsd, effectiveAt)

	var q domain.Quote
	err := row.Scan(&q.ID, &q.BrlToUsd, &q.EffectiveAt, &q.CreatedAt)
	return q, err
}

func (d *DB) GetCurrentQuote(ctx context.Context) (*domain.Quote, error) {
	rows, err := d.Pool.Query(ctx, `
		SELECT id, brl_to_usd, effective_at, created_at
		FROM quotes
		ORDER BY effective_at DESC
		LIMIT 1
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var q domain.Quote
	if err := rows.Scan(&q.ID, &q.BrlToUsd, &q.EffectiveAt, &q.CreatedAt); err != nil {
		return nil, err
	}
	return &q, nil
}

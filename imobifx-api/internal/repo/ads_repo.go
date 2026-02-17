package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/josinaldojr/imobifx-api/internal/domain"
)

type AdsFilter struct {
	Type     *string
	City     *string
	State    *string
	MinPrice *float64
	MaxPrice *float64
}

func (d *DB) CreateAd(ctx context.Context, ad domain.Ad) (domain.Ad, error) {
	row := d.Pool.QueryRow(ctx, `
		INSERT INTO ads (
			type, price_brl, image_path,
			cep, street, number, complement, neighborhood, city, state
		) VALUES (
			$1,$2,$3,
			$4,$5,$6,$7,$8,$9,$10
		)
		RETURNING id, type, price_brl, image_path,
		          cep, street, number, complement, neighborhood, city, state, created_at
	`, ad.Type, ad.PriceBRL, ad.ImagePath,
		ad.CEP, ad.Street, ad.Number, ad.Complement, ad.Neighborhood, ad.City, ad.State)

	var out domain.Ad
	err := row.Scan(&out.ID, &out.Type, &out.PriceBRL, &out.ImagePath,
		&out.CEP, &out.Street, &out.Number, &out.Complement, &out.Neighborhood, &out.City, &out.State, &out.CreatedAt)
	return out, err
}

func (d *DB) ListAds(ctx context.Context, f AdsFilter, page, pageSize int) ([]domain.Ad, int, error) {
	where, args := buildAdsWhere(f)
	offset := (page - 1) * pageSize

	countSQL := "SELECT count(*) FROM ads " + where
	var total int
	if err := d.Pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listSQL := fmt.Sprintf(`
		SELECT id, type, price_brl, image_path,
		       cep, street, number, complement, neighborhood, city, state, created_at
		FROM ads
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)

	rows, err := d.Pool.Query(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]domain.Ad, 0, pageSize)
	for rows.Next() {
		var a domain.Ad
		if err := rows.Scan(&a.ID, &a.Type, &a.PriceBRL, &a.ImagePath,
			&a.CEP, &a.Street, &a.Number, &a.Complement, &a.Neighborhood, &a.City, &a.State, &a.CreatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, a)
	}
	return out, total, nil
}

func buildAdsWhere(f AdsFilter) (string, []interface{}) {
	clauses := []string{}
	args := []interface{}{}

	add := func(expr string, val interface{}) {
		args = append(args, val)
		clauses = append(clauses, fmt.Sprintf(expr, len(args)))
	}

	if f.Type != nil {
		add("type = $%d", *f.Type)
	}
	if f.City != nil {
		add("city = $%d", *f.City)
	}
	if f.State != nil {
		add("state = $%d", *f.State)
	}
	if f.MinPrice != nil {
		add("price_brl >= $%d", *f.MinPrice)
	}
	if f.MaxPrice != nil {
		add("price_brl <= $%d", *f.MaxPrice)
	}

	if len(clauses) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(clauses, " AND "), args
}

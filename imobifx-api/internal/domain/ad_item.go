package domain

import (
	"strings"
	"time"
)

type AdItem struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	PriceBRL float64  `json:"price_brl"`
	PriceUSD *float64 `json:"price_usd"`
	ImageURL *string  `json:"image_url"`
	Address  struct {
		CEP          string  `json:"cep"`
		Street       string  `json:"street"`
		Number       *string `json:"number,omitempty"`
		Complement   *string `json:"complement,omitempty"`
		Neighborhood string  `json:"neighborhood"`
		City         string  `json:"city"`
		State        string  `json:"state"`
	} `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

func ToAdItem(a Ad) AdItem {
	item := AdItem{
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

	if a.ImagePath != nil && *a.ImagePath != "" {
		u := "/static/images/" + *a.ImagePath
		item.ImageURL = &u
	}
	return item
}

func round2(val float64) float64 {
	return float64(int(val*100+0.5)) / 100
}

func ToAdItemWithQuote(a Ad, quote *Quote) AdItem {
	item := AdItem{
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
	return item
}

package domain

import "time"

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

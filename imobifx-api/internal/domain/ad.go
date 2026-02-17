package domain

import "time"

type Ad struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	PriceBRL     float64   `json:"price_brl"`
	ImagePath    *string   `json:"-"`
	CEP          string    `json:"cep"`
	Street       string    `json:"street"`
	Number       *string   `json:"number,omitempty"`
	Complement   *string   `json:"complement,omitempty"`
	Neighborhood string    `json:"neighborhood"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
}

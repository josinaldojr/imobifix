package domain

import "time"

type AdsListResponse struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`

	QuoteUsed *struct {
		BrlToUsd    float64   `json:"brl_to_usd"`
		EffectiveAt time.Time `json:"effective_at"`
	} `json:"quote_used"`

	Items []AdItem `json:"items"`
}

package domain

import "time"

type Quote struct {
	ID          string    `json:"id"`
	BrlToUsd    float64   `json:"brl_to_usd"`
	EffectiveAt time.Time `json:"effective_at"`
	CreatedAt   time.Time `json:"created_at"`
}

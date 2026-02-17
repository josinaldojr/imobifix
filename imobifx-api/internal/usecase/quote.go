package usecase

type CreateQuoteInput struct {
	BrlToUsd    float64 `json:"brl_to_usd"`
	EffectiveAt string  `json:"effective_at"` 
}

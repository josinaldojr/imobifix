package usecase

type CreateAdInput struct {
	Type         string
	PriceBRL     float64
	CEP          string
	Street       string
	Number       *string
	Complement   *string
	Neighborhood string
	City         string
	State        string
}

type ListAdsInput struct {
	Page     int
	PageSize int

	Type     *string
	City     *string
	State    *string
	MinPrice *float64
	MaxPrice *float64
}
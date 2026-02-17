package domain

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

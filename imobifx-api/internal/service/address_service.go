package service

import (
	"context"
	"net/http"

	"github.com/josinaldojr/imobifx-api/internal/domain"
	"github.com/josinaldojr/imobifx-api/internal/errors"
	"github.com/josinaldojr/imobifx-api/internal/integrations/viacep"
	"github.com/josinaldojr/imobifx-api/internal/validation"
)

type AddressService struct {
	client ViaCEPClient
}

func NewAddressService(client ViaCEPClient) *AddressService {
	return &AddressService{client: client}
}

func (s *AddressService) Lookup(ctx context.Context, rawCEP string) (domain.Address, error) {
	cep8, ok := validation.NormalizeCEP(rawCEP)
	if !ok {
		return domain.Address{}, errors.New(http.StatusBadRequest, "CEP_INVALID", "CEP inválido. Use 8 dígitos.", map[string]string{"cep": rawCEP})
	}

	addr, err := s.client.Lookup(ctx, cep8)
	if err == viacep.ErrNotFound {
		return domain.Address{}, errors.New(http.StatusNotFound, "CEP_NOT_FOUND", "CEP não encontrado.", map[string]string{"cep": rawCEP})
	}
	if err == viacep.ErrUnavailable || err == viacep.ErrInvalidAnswer {
		return domain.Address{}, errors.New(http.StatusServiceUnavailable, "VIA_CEP_UNAVAILABLE", "Não foi possível consultar o CEP agora. Preencha o endereço manualmente.", nil)
	}
	if err != nil {
		return domain.Address{}, err
	}
	return addr, nil
}

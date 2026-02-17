package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/josinaldojr/imobifx-api/internal/domain"
	"github.com/josinaldojr/imobifx-api/internal/errors"
	"github.com/josinaldojr/imobifx-api/internal/integrations/viacep"
	"github.com/josinaldojr/imobifx-api/internal/service"
)

type fakeViaCEP struct {
	fn func(ctx context.Context, cep8 string) (domain.Address, error)
}

func (f fakeViaCEP) Lookup(ctx context.Context, cep8 string) (domain.Address, error) {
	return f.fn(ctx, cep8)
}

func TestAddressService_Lookup_InvalidCEP(t *testing.T) {
	svc := service.NewAddressService(fakeViaCEP{
		fn: func(ctx context.Context, cep8 string) (domain.Address, error) {
			t.Fatal("should not call viacep on invalid cep")
			return domain.Address{}, nil
		},
	})

	_, err := svc.Lookup(context.Background(), "123")
	require.Error(t, err)

	var appErr *errors.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, 400, appErr.Status)
	require.Equal(t, "CEP_INVALID", appErr.Code)
}

func TestAddressService_Lookup_NotFound(t *testing.T) {
	svc := service.NewAddressService(fakeViaCEP{
		fn: func(ctx context.Context, cep8 string) (domain.Address, error) {
			return domain.Address{}, viacep.ErrNotFound
		},
	})

	_, err := svc.Lookup(context.Background(), "58000-000")
	require.Error(t, err)

	var appErr *errors.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, 404, appErr.Status)
	require.Equal(t, "CEP_NOT_FOUND", appErr.Code)
}

func TestAddressService_Lookup_Unavailable(t *testing.T) {
	svc := service.NewAddressService(fakeViaCEP{
		fn: func(ctx context.Context, cep8 string) (domain.Address, error) {
			return domain.Address{}, viacep.ErrUnavailable
		},
	})

	_, err := svc.Lookup(context.Background(), "58000000")
	require.Error(t, err)

	var appErr *errors.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, 503, appErr.Status)
	require.Equal(t, "VIA_CEP_UNAVAILABLE", appErr.Code)
}

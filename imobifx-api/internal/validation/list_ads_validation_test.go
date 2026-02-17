package validation_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/josinaldojr/imobifx-api/internal/usecase"
	"github.com/josinaldojr/imobifx-api/internal/validation"
)

func TestValidateListAdsInput_OK_Default(t *testing.T) {
	in := usecase.ListAdsInput{Page: 1, PageSize: 10}
	require.NoError(t, validation.ValidateListAdsInput(in))
}

func TestValidateListAdsInput_InvalidPage(t *testing.T) {
	in := usecase.ListAdsInput{Page: 0, PageSize: 10}
	require.Error(t, validation.ValidateListAdsInput(in))
}

func TestValidateListAdsInput_InvalidPageSize(t *testing.T) {
	in := usecase.ListAdsInput{Page: 1, PageSize: 100}
	require.Error(t, validation.ValidateListAdsInput(in))
}

func TestValidateListAdsInput_InvalidType(t *testing.T) {
	typ := "FOO"
	in := usecase.ListAdsInput{Page: 1, PageSize: 10, Type: &typ}
	require.Error(t, validation.ValidateListAdsInput(in))
}

func TestValidateListAdsInput_InvalidPriceRange(t *testing.T) {
	min := 100.0
	max := 50.0
	in := usecase.ListAdsInput{Page: 1, PageSize: 10, MinPrice: &min, MaxPrice: &max}
	require.Error(t, validation.ValidateListAdsInput(in))
}

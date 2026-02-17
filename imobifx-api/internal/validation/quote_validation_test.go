package validation_test

import (
	"testing"

	"github.com/josinaldojr/imobifx-api/internal/usecase"
	"github.com/josinaldojr/imobifx-api/internal/validation"
	"github.com/stretchr/testify/require"
)

func TestValidateCreateQuoteInput_OK_NoEffectiveAt(t *testing.T) {
	in := usecase.CreateQuoteInput{BrlToUsd: 0.19, EffectiveAt: ""}
	eff, err := validation.ValidateCreateQuoteInput(in)
	require.NoError(t, err)
	require.Nil(t, eff)
}

func TestValidateCreateQuoteInput_OK_WithEffectiveAt(t *testing.T) {
	in := usecase.CreateQuoteInput{BrlToUsd: 0.19, EffectiveAt: "2026-02-16T10:00:00Z"}
	eff, err := validation.ValidateCreateQuoteInput(in)
	require.NoError(t, err)
	require.NotNil(t, eff)
}

func TestValidateCreateQuoteInput_Invalid(t *testing.T) {
	in := usecase.CreateQuoteInput{BrlToUsd: 0, EffectiveAt: "invalid"}
	_, err := validation.ValidateCreateQuoteInput(in)
	require.Error(t, err)
}

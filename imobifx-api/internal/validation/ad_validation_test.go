package validation_test

import (
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"github.com/josinaldojr/imobifx-api/internal/errors"
	"github.com/josinaldojr/imobifx-api/internal/usecase"
	"github.com/josinaldojr/imobifx-api/internal/validation"
)

func TestNormalizeCEP_OK(t *testing.T) {
	cep, ok := validation.NormalizeCEP("58000-000")
	require.True(t, ok)
	require.Equal(t, "58000000", cep)

	cep, ok = validation.NormalizeCEP("58000000")
	require.True(t, ok)
	require.Equal(t, "58000000", cep)

	cep, ok = validation.NormalizeCEP(" 58000.000 ")
	require.True(t, ok)
	require.Equal(t, "58000000", cep)
}

func TestNormalizeCEP_Invalid(t *testing.T) {
	_, ok := validation.NormalizeCEP("123")
	require.False(t, ok)

	_, ok = validation.NormalizeCEP("abcdefgh")
	require.False(t, ok)

	_, ok = validation.NormalizeCEP("58000-00")
	require.False(t, ok)
}

func TestFormatCEP(t *testing.T) {
	require.Equal(t, "58000-000", validation.FormatCEP("58000000"))
}

func TestValidateCreateAdInput_OK_FormatsCEP(t *testing.T) {
	in := &usecase.CreateAdInput{
		Type:         "SALE",
		PriceBRL:     250000,
		CEP:          "58000000",
		Street:       "Rua A",
		Neighborhood: "Centro",
		City:         "João Pessoa",
		State:        "PB",
	}

	err := validation.ValidateCreateAdInput(in)
	require.NoError(t, err)

	require.Equal(t, "58000-000", in.CEP)
}

func TestValidateCreateAdInput_Invalid_ReturnsAppErrorWithDetails(t *testing.T) {
	in := &usecase.CreateAdInput{
		Type:         "INVALID",
		PriceBRL:     -1,
		CEP:          "123",
		Street:       "",
		Neighborhood: "",
		City:         "",
		State:        "P",
	}

	err := validation.ValidateCreateAdInput(in)
	require.Error(t, err)

	var appErr *errors.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, 400, appErr.Status)
	require.Equal(t, "VALIDATION_ERROR", appErr.Code)

	details, ok := appErr.Details.(fiber.Map)
	require.True(t, ok)

	require.Contains(t, details, "type")
	require.Contains(t, details, "price_brl")
	require.Contains(t, details, "cep")
	require.Contains(t, details, "street")
	require.Contains(t, details, "neighborhood")
	require.Contains(t, details, "city")
	require.Contains(t, details, "state")
}

func TestValidateImage_Nil_OK(t *testing.T) {
	require.NoError(t, validation.ValidateImage(nil, 5*1024*1024))
}

func TestValidateImage_TooLarge(t *testing.T) {
	fh := fileHeader("image/jpeg", 6*1024*1024)
	err := validation.ValidateImage(fh, 5*1024*1024)
	require.Error(t, err)

	var appErr *errors.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, 400, appErr.Status)
	require.Equal(t, "IMAGE_TOO_LARGE", appErr.Code)
}

func TestValidateImage_UnsupportedType(t *testing.T) {
	fh := fileHeader("application/pdf", 100)
	err := validation.ValidateImage(fh, 5*1024*1024)
	require.Error(t, err)

	var appErr *errors.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, 400, appErr.Status)
	require.Equal(t, "UNSUPPORTED_IMAGE_TYPE", appErr.Code)
}

func TestValidateImage_SupportedTypes(t *testing.T) {
	cases := []string{"image/jpeg", "image/png", "image/webp"}
	for _, ct := range cases {
		fh := fileHeader(ct, 100)
		require.NoError(t, validation.ValidateImage(fh, 5*1024*1024), "content-type=%s", ct)
	}
}

func fileHeader(contentType string, size int64) *multipart.FileHeader {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Type", contentType)

	return &multipart.FileHeader{
		Filename: "file",
		Header:   h,
		Size:     size,
	}
}

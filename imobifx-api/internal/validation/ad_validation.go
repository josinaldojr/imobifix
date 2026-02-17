package validation

import (
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/josinaldojr/imobifix-api/internal/domain"
	"github.com/josinaldojr/imobifix-api/internal/errors"
)

var cepRe = regexp.MustCompile(`^\d{8}$`)

func NormalizeCEP(raw string) (string, bool) {
	d := strings.NewReplacer("-", "", ".", "", " ", "").Replace(raw)
	return d, cepRe.MatchString(d)
}

func FormatCEP(cep8 string) string { return cep8[:5] + "-" + cep8[5:] }

func ValidateCreateAdInput(in *domain.CreateAdInput) error {
	details := fiber.Map{}

	if in.Type != "SALE" && in.Type != "RENT" {
		details["type"] = "must be SALE or RENT"
	}
	if in.PriceBRL < 0 {
		details["price_brl"] = "must be >= 0"
	}

	cep8, ok := NormalizeCEP(in.CEP)
	if !ok {
		details["cep"] = "must have 8 digits"
	} else {
		in.CEP = FormatCEP(cep8)
	}

	if strings.TrimSpace(in.Street) == "" {
		details["street"] = "required"
	}
	if strings.TrimSpace(in.Neighborhood) == "" {
		details["neighborhood"] = "required"
	}
	if strings.TrimSpace(in.City) == "" {
		details["city"] = "required"
	}
	if len(strings.TrimSpace(in.State)) != 2 {
		details["state"] = "must have 2 letters (UF)"
	}

	if len(details) > 0 {
		return errors.New(http.StatusBadRequest, "VALIDATION_ERROR", "Dados inválidos.", details)
	}
	return nil
}

func ValidateImage(file *multipart.FileHeader, maxBytes int64) error {
	if file == nil {
		return nil
	}
	if file.Size > maxBytes {
		return errors.New(http.StatusBadRequest, "IMAGE_TOO_LARGE", "Imagem excede o tamanho máximo.", fiber.Map{"max_bytes": maxBytes})
	}

	ct := strings.ToLower(file.Header.Get("Content-Type"))
	switch ct {
	case "image/jpeg", "image/png", "image/webp":
		return nil
	default:
		return errors.New(http.StatusBadRequest, "UNSUPPORTED_IMAGE_TYPE", "Tipo de imagem não suportado (jpeg/png/webp).", fiber.Map{"content_type": ct})
	}
}

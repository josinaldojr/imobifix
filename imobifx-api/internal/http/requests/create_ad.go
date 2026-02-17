package requests

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/josinaldojr/imobifx-api/internal/errors"
	"github.com/josinaldojr/imobifx-api/internal/usecase"
)

func BindCreateAd(c *fiber.Ctx) (usecase.CreateAdInput, *multipart.FileHeader, error) {
	typ := strings.ToUpper(strings.TrimSpace(c.FormValue("type")))
	priceStr := strings.TrimSpace(c.FormValue("price_brl"))

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return usecase.CreateAdInput{}, nil, errors.New(http.StatusBadRequest, "VALIDATION_ERROR", "Dados inválidos.", fiber.Map{"price_brl": "must be a number"})
	}

	cepRaw := strings.TrimSpace(c.FormValue("cep"))

	in := usecase.CreateAdInput{
		Type:         typ,
		PriceBRL:     price,
		CEP:          cepRaw,
		Street:       strings.TrimSpace(c.FormValue("street")),
		Neighborhood: strings.TrimSpace(c.FormValue("neighborhood")),
		City:         strings.TrimSpace(c.FormValue("city")),
		State:        strings.ToUpper(strings.TrimSpace(c.FormValue("state"))),
		Number:       optStr(c.FormValue("number")),
		Complement:   optStr(c.FormValue("complement")),
	}

	file, ferr := c.FormFile("image")
	if ferr != nil {
		file = nil
	}

	return in, file, nil
}

func optStr(v string) *string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	return &v
}

package requests

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/josinaldojr/imobifx-api/internal/errors"
	"github.com/josinaldojr/imobifx-api/internal/usecase"
)

func BindListAds(c *fiber.Ctx) (usecase.ListAdsInput, error) {
	in := usecase.ListAdsInput{
		Page:     parseInt(c.Query("page"), 1),
		PageSize: parseInt(c.Query("page_size"), 10),
	}

	if v := strings.TrimSpace(c.Query("type")); v != "" {
		vv := strings.ToUpper(v)
		in.Type = &vv
	}
	if v := strings.TrimSpace(c.Query("city")); v != "" {
		in.City = &v
	}
	if v := strings.TrimSpace(c.Query("state")); v != "" {
		vv := strings.ToUpper(v)
		in.State = &vv
	}

	if v := strings.TrimSpace(c.Query("min_price")); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return usecase.ListAdsInput{}, errors.New(http.StatusBadRequest, "VALIDATION_ERROR", "Dados inválidos.", fiber.Map{"min_price": "must be a number"})
		}
		in.MinPrice = &f
	}
	if v := strings.TrimSpace(c.Query("max_price")); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return usecase.ListAdsInput{}, errors.New(http.StatusBadRequest, "VALIDATION_ERROR", "Dados inválidos.", fiber.Map{"max_price": "must be a number"})
		}
		in.MaxPrice = &f
	}

	return in, nil
}

func parseInt(v string, def int) int {
	if strings.TrimSpace(v) == "" {
		return def
	}
	x, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return x
}

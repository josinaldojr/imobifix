package requests

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/josinaldojr/imobifx-api/internal/errors"
	"github.com/josinaldojr/imobifx-api/internal/usecase"
)

func BindCreateQuote(c *fiber.Ctx) (usecase.CreateQuoteInput, error) {
	var in usecase.CreateQuoteInput
	if err := c.BodyParser(&in); err != nil {
		return usecase.CreateQuoteInput{}, errors.New(http.StatusBadRequest, "VALIDATION_ERROR", "JSON inválido.", nil)
	}
	in.EffectiveAt = strings.TrimSpace(in.EffectiveAt)
	return in, nil
}

package handlers

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/josinaldojr/imobifix-api/internal/errors"
	"github.com/josinaldojr/imobifix-api/internal/service"
)

type createQuoteReq struct {
	BrlToUsd    float64    `json:"brl_to_usd"`
	EffectiveAt *time.Time `json:"effective_at"`
}

func CreateQuote(svc *service.QuotesService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req createQuoteReq
		if err := c.BodyParser(&req); err != nil {
			return errors.New(http.StatusBadRequest, "VALIDATION_ERROR", "JSON inválido.", nil)
		}
		if req.BrlToUsd <= 0 {
			return errors.New(http.StatusBadRequest, "VALIDATION_ERROR", "Dados inválidos.", fiber.Map{"brl_to_usd": "must be > 0"})
		}

		q, err := svc.Create(c.Context(), req.BrlToUsd, req.EffectiveAt)
		if err != nil {
			return err
		}
		return c.Status(http.StatusCreated).JSON(q)
	}
}

func CurrentQuote(svc *service.QuotesService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		q, err := svc.Current(c.Context())
		if err != nil {
			return err
		}
		return c.JSON(q)
	}
}

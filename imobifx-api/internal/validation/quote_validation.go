package validation

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/josinaldojr/imobifix-api/internal/errors"
	"github.com/josinaldojr/imobifix-api/internal/usecase"
)

func ValidateCreateQuoteInput(in usecase.CreateQuoteInput) (*time.Time, error) {
	details := fiber.Map{}

	if in.BrlToUsd <= 0 {
		details["brl_to_usd"] = "must be > 0"
	}

	var eff *time.Time
	if in.EffectiveAt != "" {
		t, err := time.Parse(time.RFC3339, in.EffectiveAt)
		if err != nil {
			details["effective_at"] = "must be RFC3339 (e.g. 2026-02-16T10:00:00Z)"
		} else {
			eff = &t
		}
	}

	if len(details) > 0 {
		return nil, errors.New(http.StatusBadRequest, "VALIDATION_ERROR", "Dados inválidos.", details)
	}
	return eff, nil
}

package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/josinaldojr/imobifix-api/internal/http/requests"
	"github.com/josinaldojr/imobifix-api/internal/service"
	"github.com/josinaldojr/imobifix-api/internal/validation"
)

func CreateQuote(svc *service.QuotesService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		in, err := requests.BindCreateQuote(c)
		if err != nil {
			return err
		}

		eff, err := validation.ValidateCreateQuoteInput(in)
		if err != nil {
			return err
		}

		q, err := svc.Create(c.UserContext(), in.BrlToUsd, eff)
		if err != nil {
			return err
		}

		return c.Status(http.StatusCreated).JSON(q)
	}
}

func CurrentQuote(svc *service.QuotesService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		q, err := svc.Current(c.UserContext())
		if err != nil {
			return err
		}
		return c.JSON(q)
	}
}

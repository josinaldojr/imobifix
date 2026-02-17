package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/josinaldojr/imobifx-api/internal/config"
	"github.com/josinaldojr/imobifx-api/internal/domain"
	"github.com/josinaldojr/imobifx-api/internal/http/requests"
	"github.com/josinaldojr/imobifx-api/internal/service"
)

func CreateAd(cfg config.Config, ads *service.AdsService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		in, file, err := requests.BindCreateAd(c)
		if err != nil {
			return err
		}

		created, err := ads.Create(c.UserContext(), in, file)
		if err != nil {
			return err
		}

		item := domain.ToAdItem(created)
		return c.Status(http.StatusCreated).JSON(item)
	}
}

func ListAds(ads *service.AdsService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		in, err := requests.BindListAds(c)
		if err != nil {
			return err
		}
		resp, err := ads.List(c.UserContext(), in)
		if err != nil {
			return err
		}
		return c.JSON(resp)
	}
}

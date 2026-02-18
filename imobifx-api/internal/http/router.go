package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/josinaldojr/imobifx-api/internal/config"
	"github.com/josinaldojr/imobifx-api/internal/http/handlers"
	"github.com/josinaldojr/imobifx-api/internal/repo"
	"github.com/josinaldojr/imobifx-api/internal/service"
)

type Deps struct {
	Config  config.Config
	Ads     *service.AdsService
	Quotes  *service.QuotesService
	DB      *repo.DB
	Address *service.AddressService
}

func RegisterRoutes(app *fiber.App, d Deps) {
	app.Get("/health", handlers.Health())
	app.Get("/swagger", handlers.SwaggerUI())
	app.Get("/swagger/", handlers.SwaggerUI())
	app.Get("/swagger/openapi.yaml", handlers.SwaggerSpec())

	api := app.Group("/api")

	api.Get("/addresses/:cep", handlers.Address(d.Address))
	api.Post("/quotes", handlers.CreateQuote(d.Quotes))
	api.Get("/quotes/current", handlers.CurrentQuote(d.Quotes))

	api.Post("/ads", handlers.CreateAd(d.Config, d.Ads))
	api.Get("/ads", handlers.ListAds(d.Ads))
}

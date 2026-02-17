package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/josinaldojr/imobifix-api/internal/config"
	"github.com/josinaldojr/imobifix-api/internal/http/handlers"
	"github.com/josinaldojr/imobifix-api/internal/integrations/viacep"
	"github.com/josinaldojr/imobifix-api/internal/repo"
	"github.com/josinaldojr/imobifix-api/internal/service"
)

type Deps struct {
	Config    config.Config
	ViaCEP    *viacep.Client
	Ads       *service.AdsService
	Quotes    *service.QuotesService
	DB        *repo.DB
	StartTime time.Time
}

func RegisterRoutes(app *fiber.App, d Deps) {
	app.Get("/health", handlers.Health())

	api := app.Group("/api")

	api.Get("/addresses/:cep", handlers.Address(d.ViaCEP))
	api.Post("/quotes", handlers.CreateQuote(d.Quotes))
	api.Get("/quotes/current", handlers.CurrentQuote(d.Quotes))

	api.Post("/ads", handlers.CreateAd(d.Config, d.Ads))
	api.Get("/ads", handlers.ListAds(d.Ads, d.Quotes))
}

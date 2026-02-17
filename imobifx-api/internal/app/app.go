package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/josinaldojr/imobifx-api/internal/config"
	"github.com/josinaldojr/imobifx-api/internal/errors"
	"github.com/josinaldojr/imobifx-api/internal/http"
	middlewares "github.com/josinaldojr/imobifx-api/internal/http/midlewares"

	"github.com/josinaldojr/imobifx-api/internal/integrations/viacep"
	"github.com/josinaldojr/imobifx-api/internal/logging"
	"github.com/josinaldojr/imobifx-api/internal/repo"
	"github.com/josinaldojr/imobifx-api/internal/service"
)

func Run(cfg config.Config) error {
	db, err := repo.NewPostgres(cfg.DBDSN)
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}
	defer db.Close()

	viaCEP := viacep.NewClient(cfg.ViaCepBaseURL, cfg.ViaCepTimeout)

	addressSvc := service.NewAddressService(viaCEP)
	adsSvc := service.NewAdsService(db, cfg.ImagesDir, cfg.MaxImageBytes)
	quotesSvc := service.NewQuotesService(db)

	log := logging.New(cfg)
	slog.SetDefault(log)

	app := fiber.New(fiber.Config{
		AppName:      "ImobiFX",
		BodyLimit:    cfg.BodyLimitBytes,
		ErrorHandler: errors.FiberErrorHandler,
	})

	app.Use(requestid.New())
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(middlewares.AccessLog(log))

	app.Static("/static/images", cfg.ImagesDir)

	http.RegisterRoutes(app, http.Deps{
		Config:  cfg,
		Address: addressSvc,
		Ads:     adsSvc,
		Quotes:  quotesSvc,
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Listen(":" + cfg.Port)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-sigCh:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return app.ShutdownWithContext(ctx)
	}
}

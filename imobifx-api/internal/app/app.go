package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/josinaldojr/imobifix-api/internal/config"
	"github.com/josinaldojr/imobifix-api/internal/errors"
	"github.com/josinaldojr/imobifix-api/internal/http"
	"github.com/josinaldojr/imobifix-api/internal/integrations/viacep"
	"github.com/josinaldojr/imobifix-api/internal/repo"
	"github.com/josinaldojr/imobifix-api/internal/service"
)

func Run(cfg config.Config) error {
	db, err := repo.NewPostgres(cfg.DBDSN)
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}
	defer db.Close()

	viaCEP := viacep.NewClient(cfg.ViaCepBaseURL, cfg.ViaCepTimeout)

	adsSvc := service.NewAdsService(db, cfg.ImagesDir)
	quotesSvc := service.NewQuotesService(db)

	app := fiber.New(fiber.Config{
		AppName:      "ImobiFX",
		BodyLimit:    cfg.BodyLimitBytes,
		ErrorHandler: errors.FiberErrorHandler,
	})

	app.Static("/static/images", cfg.ImagesDir)

	http.RegisterRoutes(app, http.Deps{
		Config: cfg,
		ViaCEP: viaCEP,
		Ads:    adsSvc,
		Quotes: quotesSvc,
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

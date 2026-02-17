package middlewares

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AccessLog(log *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		lat := time.Since(start)

		rid, _ := c.Locals("requestid").(string)

		attrs := []slog.Attr{
			slog.String("request_id", rid),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Int64("latency_ms", lat.Milliseconds()),
			slog.String("ip", c.IP()),
			slog.String("ua", string(c.Context().UserAgent())),
		}

		if c.Response().StatusCode() >= 500 {
			log.Error("http_request", attrsToArgs(attrs)...)
		} else {
			log.Info("http_request", attrsToArgs(attrs)...)
		}
		return err
	}
}

func attrsToArgs(attrs []slog.Attr) []any {
	args := make([]any, len(attrs))
	for i, a := range attrs {
		args[i] = a
	}
	return args
}

package logging

import (
	"log/slog"
	"os"
	"strings"

	"github.com/josinaldojr/imobifix-api/internal/config"
)

func New(cfg config.Config) *slog.Logger {
	level := parseLevel(getenv("LOG_LEVEL", defaultLevel(cfg)))
	format := strings.ToLower(getenv("LOG_FORMAT", defaultFormat(cfg)))

	var h slog.Handler
	opts := &slog.HandlerOptions{Level: level}

	if format == "json" {
		h = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		h = slog.NewTextHandler(os.Stdout, opts)
	}
	return slog.New(h).With(
		slog.String("app", "imobifx"),
		slog.String("env", string(cfg.Env)),
	)
}

func defaultLevel(cfg config.Config) string {
	if cfg.Env == config.EnvProd {
		return "info"
	}
	return "debug"
}
func defaultFormat(cfg config.Config) string {
	if cfg.Env == config.EnvProd {
		return "json"
	}
	return "text"
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

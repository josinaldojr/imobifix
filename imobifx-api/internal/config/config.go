package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Env string

const (
	EnvDev  Env = "dev"
	EnvProd Env = "prod"
	EnvTest Env = "test"
)

type Config struct {
	Env  Env
	Port string
	DBDSN string
	ViaCepBaseURL string
	ViaCepTimeout time.Duration
	ImagesDir string
	BodyLimitBytes int
	MaxImageBytes  int64
	LogLevel string
	LogFormat string
}

func Load() (Config, error) {
	env := Env(strings.ToLower(getenv("APP_ENV", string(EnvDev))))
	if env != EnvDev && env != EnvProd && env != EnvTest {
		return Config{}, fmt.Errorf("invalid APP_ENV: %q (use dev|prod|test)", env)
	}

	if env == EnvDev {
		_ = godotenv.Load(".env")
		_ = godotenv.Overload(".env.local")
	}

	cfg := Config{
		Env:            env,
		Port:           getenv("APP_PORT", "8080"),
		DBDSN:          getenv("DB_DSN", ""),
		ViaCepBaseURL:  getenv("VIA_CEP_BASE_URL", "https://viacep.com.br"),
		ImagesDir:      getenv("IMAGES_DIR", "./data/images"),
		BodyLimitBytes: mustInt(getenv("BODY_LIMIT_BYTES", "10485760")),
		MaxImageBytes:  mustInt64(getenv("MAX_IMAGE_BYTES", "5242880")),
		LogLevel:       getenv("LOG_LEVEL", "debug"),
		LogFormat:      getenv("LOG_FORMAT", "text"),
	}

	timeoutStr := getenv("VIA_CEP_TIMEOUT", "2500ms")
	tout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid VIA_CEP_TIMEOUT=%q: %w", timeoutStr, err)
	}
	cfg.ViaCepTimeout = tout

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) IsProd() bool { return c.Env == EnvProd }

func (c Config) Validate() error {
	var errs []string

	if strings.TrimSpace(c.Port) == "" {
		errs = append(errs, "APP_PORT is required")
	}

	if strings.TrimSpace(c.DBDSN) == "" {
		errs = append(errs, "DB_DSN is required")
	}

	if strings.TrimSpace(c.ImagesDir) == "" {
		errs = append(errs, "IMAGES_DIR is required")
	}

	if c.BodyLimitBytes <= 0 {
		errs = append(errs, "BODY_LIMIT_BYTES must be > 0")
	}
	if c.MaxImageBytes <= 0 {
		errs = append(errs, "MAX_IMAGE_BYTES must be > 0")
	}
	if c.ViaCepTimeout <= 0 {
		errs = append(errs, "VIA_CEP_TIMEOUT must be > 0")
	}

	if len(errs) > 0 {
		return errors.New("config error: " + strings.Join(errs, "; "))
	}
	return nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func mustInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic("invalid int env value: " + s)
	}
	return v
}

func mustInt64(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic("invalid int64 env value: " + s)
	}
	return v
}

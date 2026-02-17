package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port          string
	DBDSN         string
	ViaCepBaseURL string
	ViaCepTimeout time.Duration
	ImagesDir     string

	BodyLimitBytes int
	MaxImageBytes  int64
}

func Load() Config {
	port := getenv("APP_PORT", "8080")
	dsn := getenv("DB_DSN", "")
	if dsn == "" {
	}

	base := getenv("VIA_CEP_BASE_URL", "https://viacep.com.br")
	timeoutMS := atoi(getenv("VIA_CEP_TIMEOUT_MS", "2500"), 2500)

	imagesDir := getenv("IMAGES_DIR", "./data/images")

	bodyLimit := atoi(getenv("BODY_LIMIT_BYTES", "10485760"), 10*1024*1024)  
	maxImg := int64(atoi(getenv("MAX_IMAGE_BYTES", "5242880"), 5*1024*1024))

	return Config{
		Port:           port,
		DBDSN:          dsn,
		ViaCepBaseURL:  base,
		ViaCepTimeout:  time.Duration(timeoutMS) * time.Millisecond,
		ImagesDir:      imagesDir,
		BodyLimitBytes: bodyLimit,
		MaxImageBytes:  maxImg,
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func atoi(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

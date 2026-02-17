//go:build integration

package repo_test

import (
	"os"

	"github.com/joho/godotenv"
)

func testDSN() string {
	if v := os.Getenv("TEST_DB_DSN"); v != "" {
		return v
	}
	if v := os.Getenv("DB_DSN"); v != "" {
		return v
	}

	_ = godotenv.Load(".env.test")
	_ = godotenv.Load("../.env.test")
	_ = godotenv.Load("../../.env.test")

	if v := os.Getenv("TEST_DB_DSN"); v != "" {
		return v
	}
	return os.Getenv("DB_DSN")
}

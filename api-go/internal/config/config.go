package config

import (
	"net"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr        string
	DatabaseURL string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return Config{}, err
	}

	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(envOr("POSTGRES_USER", "food_delivery"), os.Getenv("POSTGRES_PASSWORD")),
		Host:     net.JoinHostPort(envOr("POSTGRES_HOST", "127.0.0.1"), envOr("POSTGRES_PORT", "5432")),
		Path:     "/" + envOr("POSTGRES_DB", "food_delivery"),
		RawQuery: "sslmode=disable",
	}
	return Config{
		Addr:        net.JoinHostPort(envOr("API_HOST", "127.0.0.1"), envOr("API_PORT", "8080")),
		DatabaseURL: dsn.String(),
	}, nil
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

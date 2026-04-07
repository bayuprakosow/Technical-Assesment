package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	HTTPAddr          string
	DatabaseURL       string
	JWTSecret         string
	BasicAuthUser     string
	BasicAuthPassword string
	ServiceName       string
	ServiceVersion    string
}

func Load() (*Config, error) {
	jwt := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if len(jwt) < 32 {
		return nil, fmt.Errorf("JWT_SECRET wajib di-set dan minimal 32 karakter")
	}

	dbURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL wajib di-set")
	}

	addr := strings.TrimSpace(os.Getenv("HTTP_ADDR"))
	if addr == "" {
		addr = ":8080"
	}

	return &Config{
		HTTPAddr:          addr,
		DatabaseURL:       dbURL,
		JWTSecret:         jwt,
		BasicAuthUser:     strings.TrimSpace(os.Getenv("BASIC_AUTH_USER")),
		BasicAuthPassword: strings.TrimSpace(os.Getenv("BASIC_AUTH_PASSWORD")),
		ServiceName:       strings.TrimSpace(getenvDefault("SERVICE_NAME", "findings-api")),
		ServiceVersion:    strings.TrimSpace(getenvDefault("SERVICE_VERSION", "0.1.0")),
	}, nil
}

func getenvDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

package config

import (
	"strings"
	"testing"
)

func TestLoad_MissingDatabaseURL(t *testing.T) {
	t.Setenv("JWT_SECRET", strings.Repeat("a", 32))
	t.Setenv("DATABASE_URL", "")
	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Fatalf("expected DATABASE_URL error, got %v", err)
	}
}

func TestLoad_ShortJWTSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "tooshort")
	t.Setenv("DATABASE_URL", "postgres://u:p@localhost:5432/db?sslmode=disable")
	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "JWT_SECRET") {
		t.Fatalf("expected JWT_SECRET error, got %v", err)
	}
}

func TestLoad_OK(t *testing.T) {
	t.Setenv("JWT_SECRET", strings.Repeat("x", 32))
	t.Setenv("DATABASE_URL", "postgres://u:p@localhost:5432/db?sslmode=disable")
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("SERVICE_NAME", "")
	t.Setenv("SERVICE_VERSION", "")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Errorf("HTTPAddr = %q, want :8080", cfg.HTTPAddr)
	}
	if cfg.ServiceName != "findings-api" {
		t.Errorf("ServiceName = %q", cfg.ServiceName)
	}
	if cfg.DatabaseURL == "" || cfg.JWTSecret == "" {
		t.Fatal("expected non-empty DatabaseURL and JWTSecret")
	}
}

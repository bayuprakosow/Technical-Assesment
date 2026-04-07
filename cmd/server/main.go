package main

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/zentara/technical_assesment/internal/config"
	"github.com/zentara/technical_assesment/internal/db"
	"github.com/zentara/technical_assesment/internal/handlers"
	"github.com/zentara/technical_assesment/internal/router"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	migrationsDir, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatalf("migrations path: %v", err)
	}
	if err := db.RunMigrations(cfg.DatabaseURL, migrationsDir); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	sqlDB, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer sqlDB.Close()

	h := handlers.New(cfg, sqlDB)
	engine := router.New(cfg, h)

	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := engine.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("server: %v", err)
	}
}

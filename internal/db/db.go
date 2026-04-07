package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Open(databaseURL string) (*sql.DB, error) {
	d, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	if err := d.Ping(); err != nil {
		_ = d.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	d.SetMaxOpenConns(25)
	d.SetMaxIdleConns(5)
	return d, nil
}

func RunMigrations(databaseURL, migrationsDir string) error {
	d, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return err
	}
	defer d.Close()

	driver, err := postgres.WithInstance(d, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migrate postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

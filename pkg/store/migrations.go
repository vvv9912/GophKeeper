package store

import (
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

//go:embed sqllite/migrations/*
var migrationsSqlite embed.FS

//go:embed postgresql/migrations/*
var migrationsPostgresql embed.FS

func MigratePostgres(db *sqlx.DB) error {
	goose.SetBaseFS(migrationsPostgresql)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("postgres migrate set dialect postgres: %w", err)
	}

	if err := goose.Up(db.DB, "postgresql/migrations"); err != nil {

		return fmt.Errorf("postgres migrate up: %w", err)
	}

	return nil
}
func MigrateSQLITE(db *sqlx.DB) error {
	goose.SetBaseFS(migrationsSqlite)
	if err := goose.SetDialect("sqlite"); err != nil {
		return fmt.Errorf("sqlite migrate set dialect sqlite: %w", err)
	}

	if err := goose.Up(db.DB, "sqllite/migrations"); err != nil {

		return fmt.Errorf("sqlite migrate up: %w", err)
	}

	return nil
}

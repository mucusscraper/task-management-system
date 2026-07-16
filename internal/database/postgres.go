package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func NewPostgres() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=task_management sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Connected to PostgreSQL")
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	return db, nil
}

func runMigrations(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	fmt.Println("Running database migrations...")
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	fmt.Println("Migrations applied successfully")
	return nil
}

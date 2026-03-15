package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// InitDatabase initializes SQLite database connection and applies migrations
func InitDatabase(dbPath string) (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Apply migrations
	if err := applyMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Database migrations applied successfully")
	return db, nil
}

// CloseDatabase closes database connection
func CloseDatabase(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// applyMigrations reads and executes SQL schema from migrations/schema.sql
func applyMigrations(db *sql.DB) error {
	schemaPath := filepath.Join("migrations", "schema.sql")

	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute schema
	if _, err := db.Exec(string(schema)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

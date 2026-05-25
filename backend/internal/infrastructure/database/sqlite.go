package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// NewSQLite creates a new SQLite database connection and runs migrations
func NewSQLite(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// runMigrations executes the database schema migrations
func runMigrations(db *sql.DB) error {
	// Create sources table
	sourcesTableSQL := `
CREATE TABLE IF NOT EXISTS sources (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	url TEXT NOT NULL,
	fetch_type TEXT NOT NULL,
	llm_rules TEXT NOT NULL DEFAULT '',
	is_active BOOLEAN DEFAULT 1,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

	if _, err := db.Exec(sourcesTableSQL); err != nil {
		return fmt.Errorf("failed to create sources table: %w", err)
	}

	// Create articles table
	articlesTableSQL := `
CREATE TABLE IF NOT EXISTS articles (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	source_id INTEGER NOT NULL,
	title TEXT NOT NULL,
	url TEXT UNIQUE NOT NULL,
	summary TEXT DEFAULT '',
	published_at DATETIME,
	sent_to_discord INTEGER NOT NULL DEFAULT 0,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (source_id) REFERENCES sources (id)
);`

	if _, err := db.Exec(articlesTableSQL); err != nil {
		return fmt.Errorf("failed to create articles table: %w", err)
	}

	// Add updated_at if it doesn't exist (for existing databases)
	// We ignore the error because SQLite doesn't have "IF NOT EXISTS" for ADD COLUMN
	if _, err := db.Exec("ALTER TABLE sources ADD COLUMN updated_at DATETIME DEFAULT CURRENT_TIMESTAMP"); err != nil {
		log.Printf("Note: ALTER TABLE sources ADD COLUMN updated_at failed (might already exist): %v", err)
	} else {
		log.Println("Added updated_at column to sources table")
	}

	if _, err := db.Exec("ALTER TABLE articles ADD COLUMN updated_at DATETIME DEFAULT CURRENT_TIMESTAMP"); err != nil {
		log.Printf("Note: ALTER TABLE articles ADD COLUMN updated_at failed (might already exist): %v", err)
	} else {
		log.Println("Added updated_at column to articles table")
	}

	log.Println("Database migrations completed successfully")
	return nil
}
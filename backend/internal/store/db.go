package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Store wraps the SQLite database connection and provides data access methods.
type Store struct {
	db      *sql.DB
	dataDir string
}

// New creates a new Store, ensuring the data directory and database file exist.
// It initializes the schema (modules and settings tables) if they don't already exist.
func New(dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("creating data directory %q: %w", dataDir, err)
	}

	dbPath := filepath.Join(dataDir, "aegis.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening database %q: %w", dbPath, err)
	}

	// Enable WAL mode for better concurrent read performance.
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("setting WAL mode: %w", err)
	}

	s := &Store{db: db, dataDir: dataDir}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}
	return s, nil
}

// migrate creates the initial schema tables if they do not exist.
func (s *Store) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS modules (
		id      TEXT PRIMARY KEY,
		name    TEXT NOT NULL,
		domain  TEXT NOT NULL,
		enabled BOOLEAN NOT NULL DEFAULT 0,
		status  TEXT NOT NULL DEFAULT 'stopped'
	);

	CREATE TABLE IF NOT EXISTS settings (
		key   TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);
	`
	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("executing schema migration: %w", err)
	}
	return nil
}

// GetDB returns the underlying *sql.DB for advanced queries.
func (s *Store) GetDB() *sql.DB {
	return s.db
}

// Close closes the database connection.
func (s *Store) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("closing database: %w", err)
	}
	return nil
}

// GetSetting retrieves a setting value by key. Returns ("", nil) if not found.
func (s *Store) GetSetting(key string) (string, error) {
	var value string
	err := s.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("getting setting %q: %w", key, err)
	}
	return value, nil
}

// SetSetting upserts a setting key-value pair.
func (s *Store) SetSetting(key, value string) error {
	_, err := s.db.Exec(
		"INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value",
		key, value,
	)
	if err != nil {
		return fmt.Errorf("setting %q = %q: %w", key, value, err)
	}
	return nil
}

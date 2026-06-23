package notes

import (
	"database/sql"
	"fmt"
	"time"
)

// Note represents a field note entry.
type Note struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Tags      string `json:"tags"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Pinned    bool   `json:"pinned"`
}

// NotesManager handles CRUD operations for notes using SQLite.
type NotesManager struct {
	db *sql.DB
}

// NewNotesManager creates a new NotesManager and ensures the notes table exists.
func NewNotesManager(db *sql.DB) (*NotesManager, error) {
	nm := &NotesManager{db: db}
	if err := nm.migrate(); err != nil {
		return nil, err
	}
	return nm, nil
}

func (nm *NotesManager) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS notes (
		id         TEXT PRIMARY KEY,
		title      TEXT NOT NULL DEFAULT '',
		content    TEXT NOT NULL DEFAULT '',
		tags       TEXT NOT NULL DEFAULT '',
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		pinned     BOOLEAN NOT NULL DEFAULT 0
	);
	`
	_, err := nm.db.Exec(schema)
	return err
}

// List returns all notes, ordered by pinned first then most recently updated.
func (nm *NotesManager) List() ([]Note, error) {
	rows, err := nm.db.Query(
		"SELECT id, title, content, tags, created_at, updated_at, pinned FROM notes ORDER BY pinned DESC, updated_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("listing notes: %w", err)
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Tags, &n.CreatedAt, &n.UpdatedAt, &n.Pinned); err != nil {
			return nil, fmt.Errorf("scanning note: %w", err)
		}
		notes = append(notes, n)
	}
	if notes == nil {
		notes = make([]Note, 0)
	}
	return notes, nil
}

// Get returns a single note by ID.
func (nm *NotesManager) Get(id string) (*Note, error) {
	var n Note
	err := nm.db.QueryRow(
		"SELECT id, title, content, tags, created_at, updated_at, pinned FROM notes WHERE id = ?", id,
	).Scan(&n.ID, &n.Title, &n.Content, &n.Tags, &n.CreatedAt, &n.UpdatedAt, &n.Pinned)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("note %q not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("getting note: %w", err)
	}
	return &n, nil
}

// Create creates a new note.
func (nm *NotesManager) Create(id, title, content, tags string) (*Note, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := nm.db.Exec(
		"INSERT INTO notes (id, title, content, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		id, title, content, tags, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("creating note: %w", err)
	}
	return nm.Get(id)
}

// Update updates an existing note.
func (nm *NotesManager) Update(id, title, content, tags string) (*Note, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	result, err := nm.db.Exec(
		"UPDATE notes SET title = ?, content = ?, tags = ?, updated_at = ? WHERE id = ?",
		title, content, tags, now, id,
	)
	if err != nil {
		return nil, fmt.Errorf("updating note: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, fmt.Errorf("note %q not found", id)
	}
	return nm.Get(id)
}

// Delete deletes a note by ID.
func (nm *NotesManager) Delete(id string) error {
	result, err := nm.db.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deleting note: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("note %q not found", id)
	}
	return nil
}

// TogglePin toggles the pinned status of a note.
func (nm *NotesManager) TogglePin(id string) (*Note, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := nm.db.Exec(
		"UPDATE notes SET pinned = NOT pinned, updated_at = ? WHERE id = ?", now, id,
	)
	if err != nil {
		return nil, fmt.Errorf("toggling pin: %w", err)
	}
	return nm.Get(id)
}

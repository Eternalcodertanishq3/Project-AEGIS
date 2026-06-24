package p2p

import (
	"crypto/ed25519"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

// KeyPair holds the local node's identity keys.
type KeyPair struct {
	PublicKey  string `json:"public_key"`
	CreatedAt string `json:"created_at"`
}

// Contact represents a known P2P peer.
type Contact struct {
	ID        string `json:"id"`
	Alias     string `json:"alias"`
	PublicKey string `json:"public_key"`
	Trusted   bool   `json:"trusted"`
	LastSeen  string `json:"last_seen,omitempty"`
	AddedAt   string `json:"added_at"`
}

// SecureMessage represents an encrypted P2P message.
type SecureMessage struct {
	ID        string `json:"id"`
	ContactID string `json:"contact_id"`
	Direction string `json:"direction"` // "sent" or "received"
	Content   string `json:"content"`
	Encrypted bool   `json:"encrypted"`
	Timestamp string `json:"timestamp"`
}

// P2PStatus describes the status of the P2P module.
type P2PStatus struct {
	KeyGenerated  bool   `json:"key_generated"`
	PublicKey     string `json:"public_key"`
	ContactCount  int    `json:"contact_count"`
	MessageCount  int    `json:"message_count"`
	Listening     bool   `json:"listening"`
	ListenAddress string `json:"listen_address,omitempty"`
}

// P2PManager manages encrypted P2P communications.
type P2PManager struct {
	db        *sql.DB
	publicKey string
	privateKey ed25519.PrivateKey
}

// NewP2PManager creates a new P2P communications manager.
func NewP2PManager(db *sql.DB) (*P2PManager, error) {
	mgr := &P2PManager{db: db}

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS p2p_keys (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			public_key TEXT NOT NULL,
			private_key TEXT NOT NULL,
			created_at TEXT DEFAULT (datetime('now'))
		);
		CREATE TABLE IF NOT EXISTS p2p_contacts (
			id TEXT PRIMARY KEY,
			alias TEXT NOT NULL,
			public_key TEXT NOT NULL UNIQUE,
			trusted INTEGER DEFAULT 0,
			last_seen TEXT,
			added_at TEXT DEFAULT (datetime('now'))
		);
		CREATE TABLE IF NOT EXISTS p2p_messages (
			id TEXT PRIMARY KEY,
			contact_id TEXT NOT NULL,
			direction TEXT NOT NULL,
			content TEXT NOT NULL,
			encrypted INTEGER DEFAULT 1,
			timestamp TEXT DEFAULT (datetime('now')),
			FOREIGN KEY (contact_id) REFERENCES p2p_contacts(id)
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("p2p tables init: %w", err)
	}

	// Load or generate keys
	if err := mgr.loadOrGenerateKeys(); err != nil {
		return nil, fmt.Errorf("p2p key init: %w", err)
	}

	return mgr, nil
}

func (m *P2PManager) loadOrGenerateKeys() error {
	var pubHex, privHex string
	err := m.db.QueryRow(`SELECT public_key, private_key FROM p2p_keys WHERE id = 1`).Scan(&pubHex, &privHex)
	if err == sql.ErrNoRows {
		// Generate new keypair
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return fmt.Errorf("keygen: %w", err)
		}
		pubHex = hex.EncodeToString(pub)
		privHex = hex.EncodeToString(priv)
		_, err = m.db.Exec(`INSERT INTO p2p_keys (id, public_key, private_key) VALUES (1, ?, ?)`, pubHex, privHex)
		if err != nil {
			return fmt.Errorf("store keys: %w", err)
		}
		m.publicKey = pubHex
		m.privateKey = priv
		return nil
	}
	if err != nil {
		return err
	}
	m.publicKey = pubHex
	privBytes, _ := hex.DecodeString(privHex)
	m.privateKey = ed25519.PrivateKey(privBytes)
	return nil
}

// GetStatus returns the P2P module status.
func (m *P2PManager) GetStatus() *P2PStatus {
	var contactCount, msgCount int
	m.db.QueryRow(`SELECT COUNT(*) FROM p2p_contacts`).Scan(&contactCount)
	m.db.QueryRow(`SELECT COUNT(*) FROM p2p_messages`).Scan(&msgCount)

	return &P2PStatus{
		KeyGenerated: m.publicKey != "",
		PublicKey:    m.publicKey,
		ContactCount: contactCount,
		MessageCount: msgCount,
		Listening:    false,
	}
}

// GetKeyPair returns the public key info.
func (m *P2PManager) GetKeyPair() *KeyPair {
	var createdAt string
	m.db.QueryRow(`SELECT created_at FROM p2p_keys WHERE id = 1`).Scan(&createdAt)
	return &KeyPair{PublicKey: m.publicKey, CreatedAt: createdAt}
}

// AddContact adds a known contact.
func (m *P2PManager) AddContact(alias, publicKey string) (*Contact, error) {
	if alias == "" || publicKey == "" {
		return nil, fmt.Errorf("alias and public_key are required")
	}
	id := generateP2PID()
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := m.db.Exec(`INSERT INTO p2p_contacts (id, alias, public_key, added_at) VALUES (?, ?, ?, ?)`,
		id, alias, publicKey, now)
	if err != nil {
		return nil, fmt.Errorf("add contact: %w", err)
	}
	return &Contact{ID: id, Alias: alias, PublicKey: publicKey, AddedAt: now}, nil
}

// GetContacts returns all contacts.
func (m *P2PManager) GetContacts() []Contact {
	rows, err := m.db.Query(`SELECT id, alias, public_key, trusted, last_seen, added_at FROM p2p_contacts ORDER BY alias`)
	if err != nil {
		return []Contact{}
	}
	defer rows.Close()
	var contacts []Contact
	for rows.Next() {
		var c Contact
		var lastSeen sql.NullString
		rows.Scan(&c.ID, &c.Alias, &c.PublicKey, &c.Trusted, &lastSeen, &c.AddedAt)
		if lastSeen.Valid {
			c.LastSeen = lastSeen.String
		}
		contacts = append(contacts, c)
	}
	if contacts == nil {
		contacts = make([]Contact, 0)
	}
	return contacts
}

// SendMessage stores an outgoing encrypted message.
func (m *P2PManager) SendMessage(contactID, content string) (*SecureMessage, error) {
	if content == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}
	msg := &SecureMessage{
		ID:        generateP2PID(),
		ContactID: contactID,
		Direction: "sent",
		Content:   content,
		Encrypted: true,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	_, err := m.db.Exec(`INSERT INTO p2p_messages (id, contact_id, direction, content, encrypted, timestamp) VALUES (?, ?, ?, ?, ?, ?)`,
		msg.ID, msg.ContactID, msg.Direction, msg.Content, msg.Encrypted, msg.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("send p2p message: %w", err)
	}
	return msg, nil
}

// GetMessages returns messages for a contact.
func (m *P2PManager) GetMessages(contactID string) []SecureMessage {
	rows, err := m.db.Query(`SELECT id, contact_id, direction, content, encrypted, timestamp FROM p2p_messages WHERE contact_id = ? ORDER BY timestamp`, contactID)
	if err != nil {
		return []SecureMessage{}
	}
	defer rows.Close()
	var msgs []SecureMessage
	for rows.Next() {
		var msg SecureMessage
		rows.Scan(&msg.ID, &msg.ContactID, &msg.Direction, &msg.Content, &msg.Encrypted, &msg.Timestamp)
		msgs = append(msgs, msg)
	}
	if msgs == nil {
		msgs = make([]SecureMessage, 0)
	}
	return msgs
}

// DeleteContact removes a contact and their messages.
func (m *P2PManager) DeleteContact(id string) error {
	m.db.Exec(`DELETE FROM p2p_messages WHERE contact_id = ?`, id)
	res, err := m.db.Exec(`DELETE FROM p2p_contacts WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("contact not found")
	}
	return nil
}

func generateP2PID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

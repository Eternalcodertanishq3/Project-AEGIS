package mesh

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

// Message represents a single mesh message.
type Message struct {
	ID        string `json:"id"`
	Channel   string `json:"channel"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	Delivered bool   `json:"delivered"`
	Via       string `json:"via"` // "local", "lan", "lora"
}

// Channel represents a messaging channel.
type Channel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MessageCount int   `json:"message_count"`
	LastActivity string `json:"last_activity,omitempty"`
}

// MeshStatus describes the current status of mesh radios and transports.
type MeshStatus struct {
	LoRaConnected   bool   `json:"lora_connected"`
	LANAvailable    bool   `json:"lan_available"`
	NodeID          string `json:"node_id"`
	ActiveChannels  int    `json:"active_channels"`
	TotalMessages   int    `json:"total_messages"`
	TransportMode   string `json:"transport_mode"` // "offline", "lan", "lora", "lora+lan"
}

// MeshManager manages mesh messaging state.
type MeshManager struct {
	db     *sql.DB
	nodeID string
}

// NewMeshManager creates a new mesh messaging manager.
func NewMeshManager(db *sql.DB) (*MeshManager, error) {
	mgr := &MeshManager{db: db, nodeID: generateNodeID()}

	// Create tables
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS mesh_channels (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			created_at TEXT DEFAULT (datetime('now'))
		);
		CREATE TABLE IF NOT EXISTS mesh_messages (
			id TEXT PRIMARY KEY,
			channel_id TEXT NOT NULL,
			sender TEXT NOT NULL,
			content TEXT NOT NULL,
			timestamp TEXT DEFAULT (datetime('now')),
			delivered INTEGER DEFAULT 0,
			via TEXT DEFAULT 'local',
			FOREIGN KEY (channel_id) REFERENCES mesh_channels(id)
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("mesh tables init: %w", err)
	}

	// Seed default channels
	mgr.seedDefaults()
	return mgr, nil
}

func (m *MeshManager) seedDefaults() {
	defaults := []struct{ id, name, desc string }{
		{"emergency", "Emergency", "Emergency broadcasts — highest priority"},
		{"general", "General", "General communications channel"},
		{"logistics", "Logistics", "Supply coordination and resource sharing"},
		{"intel", "Intelligence", "Situational awareness and reports"},
	}
	for _, ch := range defaults {
		m.db.Exec(`INSERT OR IGNORE INTO mesh_channels (id, name, description) VALUES (?, ?, ?)`,
			ch.id, ch.name, ch.desc)
	}
}

// GetStatus returns the current mesh network status.
func (m *MeshManager) GetStatus() *MeshStatus {
	var msgCount int
	m.db.QueryRow(`SELECT COUNT(*) FROM mesh_messages`).Scan(&msgCount)
	var chCount int
	m.db.QueryRow(`SELECT COUNT(*) FROM mesh_channels`).Scan(&chCount)

	return &MeshStatus{
		LoRaConnected:  false, // No hardware detected
		LANAvailable:   true,  // Assume LAN is available
		NodeID:         m.nodeID,
		ActiveChannels: chCount,
		TotalMessages:  msgCount,
		TransportMode:  "lan",
	}
}

// GetChannels returns all channels with message counts.
func (m *MeshManager) GetChannels() []Channel {
	rows, err := m.db.Query(`
		SELECT c.id, c.name, c.description,
			COUNT(m.id) as msg_count,
			MAX(m.timestamp) as last_activity
		FROM mesh_channels c
		LEFT JOIN mesh_messages m ON m.channel_id = c.id
		GROUP BY c.id ORDER BY c.name
	`)
	if err != nil {
		return []Channel{}
	}
	defer rows.Close()

	var channels []Channel
	for rows.Next() {
		var ch Channel
		var lastAct sql.NullString
		rows.Scan(&ch.ID, &ch.Name, &ch.Description, &ch.MessageCount, &lastAct)
		if lastAct.Valid {
			ch.LastActivity = lastAct.String
		}
		channels = append(channels, ch)
	}
	if channels == nil {
		channels = make([]Channel, 0)
	}
	return channels
}

// GetMessages returns messages for a channel.
func (m *MeshManager) GetMessages(channelID string, limit int) []Message {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := m.db.Query(`
		SELECT id, channel_id, sender, content, timestamp, delivered, via
		FROM mesh_messages WHERE channel_id = ?
		ORDER BY timestamp DESC LIMIT ?
	`, channelID, limit)
	if err != nil {
		return []Message{}
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		rows.Scan(&msg.ID, &msg.Channel, &msg.Sender, &msg.Content, &msg.Timestamp, &msg.Delivered, &msg.Via)
		messages = append(messages, msg)
	}
	if messages == nil {
		messages = make([]Message, 0)
	}
	// Reverse so oldest first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages
}

// SendMessage sends a message to a channel.
func (m *MeshManager) SendMessage(channelID, content string) (*Message, error) {
	if content == "" {
		return nil, fmt.Errorf("message content cannot be empty")
	}
	msg := &Message{
		ID:        generateMsgID(),
		Channel:   channelID,
		Sender:    m.nodeID,
		Content:   content,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Delivered: true,
		Via:       "local",
	}
	_, err := m.db.Exec(`
		INSERT INTO mesh_messages (id, channel_id, sender, content, timestamp, delivered, via)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, msg.ID, msg.Channel, msg.Sender, msg.Content, msg.Timestamp, msg.Delivered, msg.Via)
	if err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}
	return msg, nil
}

func generateNodeID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return "AEGIS-" + hex.EncodeToString(b)
}

func generateMsgID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

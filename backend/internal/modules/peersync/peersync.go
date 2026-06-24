package peersync

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

// Peer represents a known AEGIS peer node.
type Peer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"` // IP:port
	NodeID    string `json:"node_id"`
	Status    string `json:"status"` // "online", "offline", "syncing", "unknown"
	LastSync  string `json:"last_sync,omitempty"`
	AddedAt   string `json:"added_at"`
}

// ContentManifest describes what content this node has available to share.
type ContentManifest struct {
	NodeID       string          `json:"node_id"`
	Hostname     string          `json:"hostname"`
	Modules      []ModuleContent `json:"modules"`
	TotalItems   int             `json:"total_items"`
	LastUpdated  string          `json:"last_updated"`
}

// ModuleContent describes content available from a specific module.
type ModuleContent struct {
	ModuleID    string `json:"module_id"`
	ModuleName  string `json:"module_name"`
	ItemCount   int    `json:"item_count"`
	Description string `json:"description"`
}

// SyncStatus describes the current sync state.
type SyncStatus struct {
	NodeID       string `json:"node_id"`
	PeerCount    int    `json:"peer_count"`
	OnlinePeers  int    `json:"online_peers"`
	LastSync     string `json:"last_sync,omitempty"`
	SyncEnabled  bool   `json:"sync_enabled"`
	ListenPort   int    `json:"listen_port"`
}

// SyncManager manages peer discovery and content synchronization.
type SyncManager struct {
	db     *sql.DB
	nodeID string
}

// NewSyncManager creates a new peer sync manager.
func NewSyncManager(db *sql.DB) (*SyncManager, error) {
	mgr := &SyncManager{db: db, nodeID: generateSyncID()}

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS sync_peers (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			address TEXT NOT NULL,
			node_id TEXT DEFAULT '',
			status TEXT DEFAULT 'unknown',
			last_sync TEXT,
			added_at TEXT DEFAULT (datetime('now'))
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("sync tables init: %w", err)
	}
	return mgr, nil
}

// GetStatus returns the current sync status.
func (m *SyncManager) GetStatus() *SyncStatus {
	var total, online int
	m.db.QueryRow(`SELECT COUNT(*) FROM sync_peers`).Scan(&total)
	m.db.QueryRow(`SELECT COUNT(*) FROM sync_peers WHERE status = 'online'`).Scan(&online)
	var lastSync sql.NullString
	m.db.QueryRow(`SELECT MAX(last_sync) FROM sync_peers`).Scan(&lastSync)

	status := &SyncStatus{
		NodeID:      m.nodeID,
		PeerCount:   total,
		OnlinePeers: online,
		SyncEnabled: true,
		ListenPort:  8080,
	}
	if lastSync.Valid {
		status.LastSync = lastSync.String
	}
	return status
}

// GetManifest returns this node's content manifest.
func (m *SyncManager) GetManifest() *ContentManifest {
	hostname := "AEGIS-Node"
	// Count notes, etc.
	var noteCount int
	m.db.QueryRow(`SELECT COUNT(*) FROM notes`).Scan(&noteCount)

	return &ContentManifest{
		NodeID:   m.nodeID,
		Hostname: hostname,
		Modules: []ModuleContent{
			{ModuleID: "notes", ModuleName: "Notes", ItemCount: noteCount, Description: "Field notes and documents"},
			{ModuleID: "mesh", ModuleName: "Mesh Messages", ItemCount: 0, Description: "Mesh message history"},
		},
		TotalItems:  noteCount,
		LastUpdated: time.Now().UTC().Format(time.RFC3339),
	}
}

// AddPeer adds a peer node by address.
func (m *SyncManager) AddPeer(name, address string) (*Peer, error) {
	if name == "" || address == "" {
		return nil, fmt.Errorf("name and address are required")
	}
	peer := &Peer{
		ID:      generateSyncID(),
		Name:    name,
		Address: address,
		Status:  "unknown",
		AddedAt: time.Now().UTC().Format(time.RFC3339),
	}
	_, err := m.db.Exec(`INSERT INTO sync_peers (id, name, address, status, added_at) VALUES (?, ?, ?, ?, ?)`,
		peer.ID, peer.Name, peer.Address, peer.Status, peer.AddedAt)
	if err != nil {
		return nil, fmt.Errorf("add peer: %w", err)
	}
	return peer, nil
}

// GetPeers returns all known peers.
func (m *SyncManager) GetPeers() []Peer {
	rows, err := m.db.Query(`SELECT id, name, address, node_id, status, last_sync, added_at FROM sync_peers ORDER BY name`)
	if err != nil {
		return []Peer{}
	}
	defer rows.Close()
	var peers []Peer
	for rows.Next() {
		var p Peer
		var lastSync sql.NullString
		rows.Scan(&p.ID, &p.Name, &p.Address, &p.NodeID, &p.Status, &lastSync, &p.AddedAt)
		if lastSync.Valid {
			p.LastSync = lastSync.String
		}
		peers = append(peers, p)
	}
	if peers == nil {
		peers = make([]Peer, 0)
	}
	return peers
}

// RemovePeer deletes a peer.
func (m *SyncManager) RemovePeer(id string) error {
	res, err := m.db.Exec(`DELETE FROM sync_peers WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("peer not found")
	}
	return nil
}

func generateSyncID() string {
	b := make([]byte, 6)
	rand.Read(b)
	return hex.EncodeToString(b)
}

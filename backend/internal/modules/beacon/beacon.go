package beacon

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math"
	"time"
)

// Position represents a logged position.
type Position struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude,omitempty"` // meters
	Accuracy  float64 `json:"accuracy,omitempty"` // meters
	Source    string  `json:"source"` // "manual", "gps", "network"
	Comment   string  `json:"comment,omitempty"`
	Timestamp string  `json:"timestamp"`
}

// APRSBeacon represents an APRS-format position beacon.
type APRSBeacon struct {
	Callsign  string  `json:"callsign"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Comment   string  `json:"comment"`
	Raw       string  `json:"raw"` // APRS-format string
	Timestamp string  `json:"timestamp"`
}

// DistanceResult holds the result of a distance/bearing calculation.
type DistanceResult struct {
	FromLat      float64 `json:"from_lat"`
	FromLon      float64 `json:"from_lon"`
	ToLat        float64 `json:"to_lat"`
	ToLon        float64 `json:"to_lon"`
	DistanceKm   float64 `json:"distance_km"`
	DistanceMi   float64 `json:"distance_mi"`
	DistanceNm   float64 `json:"distance_nm"`
	BearingDeg   float64 `json:"bearing_deg"`
	BearingCard  string  `json:"bearing_cardinal"`
}

// BeaconStatus describes the position beacon module status.
type BeaconStatus struct {
	Callsign       string `json:"callsign"`
	BeaconActive   bool   `json:"beacon_active"`
	PositionCount  int    `json:"position_count"`
	LastPosition   string `json:"last_position,omitempty"`
	TransportMode  string `json:"transport_mode"` // "off", "lan", "lora"
}

// BeaconManager manages position logging and beacon broadcasting.
type BeaconManager struct {
	db       *sql.DB
	callsign string
}

// NewBeaconManager creates a new position beacon manager.
func NewBeaconManager(db *sql.DB) (*BeaconManager, error) {
	mgr := &BeaconManager{db: db, callsign: "AEGIS-1"}

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS beacon_positions (
			id TEXT PRIMARY KEY,
			latitude REAL NOT NULL,
			longitude REAL NOT NULL,
			altitude REAL DEFAULT 0,
			accuracy REAL DEFAULT 0,
			source TEXT DEFAULT 'manual',
			comment TEXT DEFAULT '',
			timestamp TEXT DEFAULT (datetime('now'))
		);
		CREATE TABLE IF NOT EXISTS beacon_config (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("beacon tables init: %w", err)
	}

	// Load callsign
	var cs string
	err = db.QueryRow(`SELECT value FROM beacon_config WHERE key = 'callsign'`).Scan(&cs)
	if err == nil && cs != "" {
		mgr.callsign = cs
	}

	return mgr, nil
}

// GetStatus returns the beacon module status.
func (m *BeaconManager) GetStatus() *BeaconStatus {
	var count int
	m.db.QueryRow(`SELECT COUNT(*) FROM beacon_positions`).Scan(&count)
	var lastTs sql.NullString
	m.db.QueryRow(`SELECT MAX(timestamp) FROM beacon_positions`).Scan(&lastTs)

	status := &BeaconStatus{
		Callsign:      m.callsign,
		BeaconActive:  false,
		PositionCount: count,
		TransportMode: "off",
	}
	if lastTs.Valid {
		status.LastPosition = lastTs.String
	}
	return status
}

// LogPosition stores a new position fix.
func (m *BeaconManager) LogPosition(lat, lon, alt, acc float64, source, comment string) (*Position, error) {
	if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		return nil, fmt.Errorf("invalid coordinates: lat=%f lon=%f", lat, lon)
	}
	if source == "" {
		source = "manual"
	}
	pos := &Position{
		ID:        generateBeaconID(),
		Latitude:  lat,
		Longitude: lon,
		Altitude:  alt,
		Accuracy:  acc,
		Source:    source,
		Comment:   comment,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	_, err := m.db.Exec(`INSERT INTO beacon_positions (id, latitude, longitude, altitude, accuracy, source, comment, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		pos.ID, pos.Latitude, pos.Longitude, pos.Altitude, pos.Accuracy, pos.Source, pos.Comment, pos.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("log position: %w", err)
	}
	return pos, nil
}

// GetPositions returns the position history.
func (m *BeaconManager) GetPositions(limit int) []Position {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	rows, err := m.db.Query(`SELECT id, latitude, longitude, altitude, accuracy, source, comment, timestamp
		FROM beacon_positions ORDER BY timestamp DESC LIMIT ?`, limit)
	if err != nil {
		return []Position{}
	}
	defer rows.Close()
	var positions []Position
	for rows.Next() {
		var p Position
		rows.Scan(&p.ID, &p.Latitude, &p.Longitude, &p.Altitude, &p.Accuracy, &p.Source, &p.Comment, &p.Timestamp)
		positions = append(positions, p)
	}
	if positions == nil {
		positions = make([]Position, 0)
	}
	return positions
}

// GenerateAPRS generates an APRS-format position beacon string.
func (m *BeaconManager) GenerateAPRS(lat, lon float64, comment string) *APRSBeacon {
	// Convert decimal degrees to APRS DMM format
	latDeg := int(math.Abs(lat))
	latMin := (math.Abs(lat) - float64(latDeg)) * 60
	latDir := "N"
	if lat < 0 {
		latDir = "S"
	}
	lonDeg := int(math.Abs(lon))
	lonMin := (math.Abs(lon) - float64(lonDeg)) * 60
	lonDir := "E"
	if lon < 0 {
		lonDir = "W"
	}

	raw := fmt.Sprintf("%s>APRS,TCPIP*:=%02d%05.2f%s/%03d%05.2f%s_%s",
		m.callsign, latDeg, latMin, latDir, lonDeg, lonMin, lonDir, comment)

	return &APRSBeacon{
		Callsign:  m.callsign,
		Latitude:  lat,
		Longitude: lon,
		Comment:   comment,
		Raw:       raw,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// CalculateDistance computes great-circle distance and bearing between two points.
func CalculateDistance(lat1, lon1, lat2, lon2 float64) *DistanceResult {
	const R = 6371.0 // Earth radius in km

	lat1r := lat1 * math.Pi / 180
	lat2r := lat2 * math.Pi / 180
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1r)*math.Cos(lat2r)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distKm := R * c

	// Bearing
	y := math.Sin(dLon) * math.Cos(lat2r)
	x := math.Cos(lat1r)*math.Sin(lat2r) - math.Sin(lat1r)*math.Cos(lat2r)*math.Cos(dLon)
	bearing := math.Atan2(y, x) * 180 / math.Pi
	if bearing < 0 {
		bearing += 360
	}

	return &DistanceResult{
		FromLat:     lat1,
		FromLon:     lon1,
		ToLat:       lat2,
		ToLon:       lon2,
		DistanceKm:  math.Round(distKm*100) / 100,
		DistanceMi:  math.Round(distKm*0.621371*100) / 100,
		DistanceNm:  math.Round(distKm*0.539957*100) / 100,
		BearingDeg:  math.Round(bearing*100) / 100,
		BearingCard: degreesToCardinal(bearing),
	}
}

func degreesToCardinal(deg float64) string {
	dirs := []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE",
		"S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}
	idx := int(math.Round(deg/22.5)) % 16
	return dirs[idx]
}

func generateBeaconID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

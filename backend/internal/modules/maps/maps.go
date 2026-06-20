package maps

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// PMTilesFile represents a discovered .pmtiles map file.
type PMTilesFile struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	SizeBytes int64  `json:"size_bytes"`
	SizeHuman string `json:"size_human"`
}

// MapManager manages offline map files and serves PMTiles data.
type MapManager struct {
	mu       sync.RWMutex
	dataDir  string
	tileFiles []PMTilesFile
}

// NewMapManager creates a new MapManager.
func NewMapManager(dataDir string) *MapManager {
	return &MapManager{
		dataDir:   dataDir,
		tileFiles: make([]PMTilesFile, 0),
	}
}

// DiscoverMapFiles scans for .pmtiles files in the data and content-packs directories.
func (mm *MapManager) DiscoverMapFiles() ([]PMTilesFile, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	mm.tileFiles = make([]PMTilesFile, 0)

	searchPaths := []string{
		mm.dataDir,
		filepath.Join(filepath.Dir(mm.dataDir), "content-packs"),
		filepath.Join(filepath.Dir(mm.dataDir), "content-packs", "maps-regional"),
		filepath.Join(mm.dataDir, "maps"),
	}

	seen := make(map[string]bool)

	for _, searchPath := range searchPaths {
		err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".pmtiles" {
				absPath, _ := filepath.Abs(path)
				if seen[absPath] {
					return nil
				}
				seen[absPath] = true

				id := strings.TrimSuffix(info.Name(), ext)
				mm.tileFiles = append(mm.tileFiles, PMTilesFile{
					ID:        id,
					Name:      info.Name(),
					Path:      absPath,
					SizeBytes: info.Size(),
					SizeHuman: humanizeBytes(info.Size()),
				})
			}
			return nil
		})
		if err != nil {
			continue
		}
	}

	log.Printf("  Maps: discovered %d PMTiles files", len(mm.tileFiles))
	return mm.tileFiles, nil
}

// GetMapFiles returns currently discovered PMTiles files.
func (mm *MapManager) GetMapFiles() []PMTilesFile {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	result := make([]PMTilesFile, len(mm.tileFiles))
	copy(result, mm.tileFiles)
	return result
}

// GetMapFile returns a specific PMTiles file by ID.
func (mm *MapManager) GetMapFile(id string) (*PMTilesFile, error) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	for _, f := range mm.tileFiles {
		if f.ID == id {
			return &f, nil
		}
	}
	return nil, fmt.Errorf("map file %q not found", id)
}

// Handlers holds the HTTP handler functions for the maps module.
type Handlers struct {
	manager *MapManager
}

// NewHandlers creates new maps API handlers.
func NewHandlers(manager *MapManager) *Handlers {
	return &Handlers{manager: manager}
}

// RegisterRoutes registers all maps module routes on the given mux.
func (h *Handlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/maps/status", h.statusHandler)
	mux.HandleFunc("GET /api/maps/tiles", h.tilesListHandler)
	mux.HandleFunc("GET /api/maps/tiles/{id}/{rest...}", h.tileServeHandler)
}

// statusHandler returns the status of the maps module.
func (h *Handlers) statusHandler(w http.ResponseWriter, r *http.Request) {
	tileFiles := h.manager.GetMapFiles()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"module":        "maps",
		"status":        "ready",
		"pmtiles_count": len(tileFiles),
		"pmtiles_files": tileFiles,
	})
}

// tilesListHandler returns discovered PMTiles files.
func (h *Handlers) tilesListHandler(w http.ResponseWriter, r *http.Request) {
	tileFiles := h.manager.GetMapFiles()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"tiles": tileFiles,
		"count": len(tileFiles),
	})
}

// tileServeHandler serves PMTiles files with HTTP range request support.
// This allows MapLibre GL JS + pmtiles JS library to directly load tiles
// via range requests to the static file.
func (h *Handlers) tileServeHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	mapFile, err := h.manager.GetMapFile(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	// Serve the file with range request support (required for PMTiles protocol)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Range")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Range, Content-Length")

	http.ServeFile(w, r, mapFile.Path)
}

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]interface{}{
		"error":  message,
		"status": status,
	})
}

// humanizeBytes converts bytes to a human-readable string.
func humanizeBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

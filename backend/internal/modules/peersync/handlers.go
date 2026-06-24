package peersync

import (
	"encoding/json"
	"net/http"
)

// Handlers holds HTTP handlers for the peer sync module.
type Handlers struct {
	mgr *SyncManager
}

// NewHandlers creates new peer sync API handlers.
func NewHandlers(mgr *SyncManager) *Handlers {
	return &Handlers{mgr: mgr}
}

// RegisterRoutes registers all peer sync routes.
func (h *Handlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/sync/status", h.statusHandler)
	mux.HandleFunc("GET /api/sync/manifest", h.manifestHandler)
	mux.HandleFunc("GET /api/sync/peers", h.peersHandler)
	mux.HandleFunc("POST /api/sync/peers", h.addPeerHandler)
	mux.HandleFunc("DELETE /api/sync/peers/{id}", h.removePeerHandler)
}

func (h *Handlers) statusHandler(w http.ResponseWriter, r *http.Request) {
	syncWriteJSON(w, http.StatusOK, h.mgr.GetStatus())
}

func (h *Handlers) manifestHandler(w http.ResponseWriter, r *http.Request) {
	syncWriteJSON(w, http.StatusOK, h.mgr.GetManifest())
}

func (h *Handlers) peersHandler(w http.ResponseWriter, r *http.Request) {
	peers := h.mgr.GetPeers()
	syncWriteJSON(w, http.StatusOK, map[string]interface{}{
		"peers": peers,
		"count": len(peers),
	})
}

func (h *Handlers) addPeerHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		syncWriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	peer, err := h.mgr.AddPeer(body.Name, body.Address)
	if err != nil {
		syncWriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	syncWriteJSON(w, http.StatusCreated, peer)
}

func (h *Handlers) removePeerHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.mgr.RemovePeer(id); err != nil {
		syncWriteError(w, http.StatusNotFound, err.Error())
		return
	}
	syncWriteJSON(w, http.StatusOK, map[string]string{"status": "removed"})
}

func syncWriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func syncWriteError(w http.ResponseWriter, status int, message string) {
	syncWriteJSON(w, status, map[string]interface{}{"error": message, "status": status})
}

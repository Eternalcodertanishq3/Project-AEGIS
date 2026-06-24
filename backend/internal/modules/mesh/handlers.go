package mesh

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Handlers holds HTTP handlers for the mesh messaging module.
type Handlers struct {
	mgr *MeshManager
}

// NewHandlers creates new mesh messaging API handlers.
func NewHandlers(mgr *MeshManager) *Handlers {
	return &Handlers{mgr: mgr}
}

// RegisterRoutes registers all mesh messaging routes.
func (h *Handlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/mesh/status", h.statusHandler)
	mux.HandleFunc("GET /api/mesh/channels", h.channelsHandler)
	mux.HandleFunc("GET /api/mesh/channels/{id}/messages", h.messagesHandler)
	mux.HandleFunc("POST /api/mesh/channels/{id}/send", h.sendHandler)
}

func (h *Handlers) statusHandler(w http.ResponseWriter, r *http.Request) {
	meshWriteJSON(w, http.StatusOK, h.mgr.GetStatus())
}

func (h *Handlers) channelsHandler(w http.ResponseWriter, r *http.Request) {
	channels := h.mgr.GetChannels()
	meshWriteJSON(w, http.StatusOK, map[string]interface{}{
		"channels": channels,
		"count":    len(channels),
	})
}

func (h *Handlers) messagesHandler(w http.ResponseWriter, r *http.Request) {
	chID := r.PathValue("id")
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil {
			limit = n
		}
	}
	messages := h.mgr.GetMessages(chID, limit)
	meshWriteJSON(w, http.StatusOK, map[string]interface{}{
		"channel":  chID,
		"messages": messages,
		"count":    len(messages),
	})
}

func (h *Handlers) sendHandler(w http.ResponseWriter, r *http.Request) {
	chID := r.PathValue("id")
	var body struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		meshWriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	msg, err := h.mgr.SendMessage(chID, body.Content)
	if err != nil {
		meshWriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	meshWriteJSON(w, http.StatusCreated, msg)
}

func meshWriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func meshWriteError(w http.ResponseWriter, status int, message string) {
	meshWriteJSON(w, status, map[string]interface{}{"error": message, "status": status})
}

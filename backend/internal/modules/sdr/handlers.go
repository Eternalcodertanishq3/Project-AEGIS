package sdr

import (
	"encoding/json"
	"net/http"
)

// Handlers holds HTTP handlers for the SDR module.
type Handlers struct {
	db *SDRDatabase
}

// NewHandlers creates new SDR API handlers.
func NewHandlers(db *SDRDatabase) *Handlers {
	return &Handlers{db: db}
}

// RegisterRoutes registers all SDR module routes.
func (h *Handlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/sdr/status", h.statusHandler)
	mux.HandleFunc("GET /api/sdr/frequencies", h.frequenciesHandler)
	mux.HandleFunc("GET /api/sdr/frequencies/{id}", h.frequencyGroupHandler)
	mux.HandleFunc("GET /api/sdr/bandplans", h.bandPlansHandler)
	mux.HandleFunc("GET /api/sdr/search", h.searchHandler)
}

func (h *Handlers) statusHandler(w http.ResponseWriter, r *http.Request) {
	sdrWriteJSON(w, http.StatusOK, h.db.GetStatus())
}

func (h *Handlers) frequenciesHandler(w http.ResponseWriter, r *http.Request) {
	groups := h.db.GetGroups()
	sdrWriteJSON(w, http.StatusOK, map[string]interface{}{
		"groups": groups,
		"count":  len(groups),
	})
}

func (h *Handlers) frequencyGroupHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	group := h.db.GetGroup(id)
	if group == nil {
		sdrWriteError(w, http.StatusNotFound, "frequency group not found")
		return
	}
	sdrWriteJSON(w, http.StatusOK, group)
}

func (h *Handlers) bandPlansHandler(w http.ResponseWriter, r *http.Request) {
	plans := h.db.GetBandPlans()
	sdrWriteJSON(w, http.StatusOK, map[string]interface{}{
		"band_plans": plans,
		"count":      len(plans),
	})
}

func (h *Handlers) searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		sdrWriteError(w, http.StatusBadRequest, "q query parameter is required")
		return
	}
	results := h.db.SearchFrequencies(query)
	sdrWriteJSON(w, http.StatusOK, map[string]interface{}{
		"results": results,
		"count":   len(results),
		"query":   query,
	})
}

func sdrWriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sdrWriteError(w http.ResponseWriter, status int, message string) {
	sdrWriteJSON(w, status, map[string]interface{}{"error": message, "status": status})
}

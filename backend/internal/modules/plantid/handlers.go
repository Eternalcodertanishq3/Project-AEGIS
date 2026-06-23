package plantid

import (
	"encoding/json"
	"net/http"
)

// Handlers holds the HTTP handler functions for the plant ID module.
type Handlers struct {
	db *PlantDB
}

// NewHandlers creates new plant ID API handlers.
func NewHandlers(db *PlantDB) *Handlers {
	return &Handlers{db: db}
}

// RegisterRoutes registers all plant ID module routes.
func (h *Handlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/plants/groups", h.groupsHandler)
	mux.HandleFunc("GET /api/plants/groups/{id}", h.groupHandler)
	mux.HandleFunc("GET /api/plants/{id}", h.plantHandler)
	mux.HandleFunc("GET /api/plants/search", h.searchHandler)
}

func (h *Handlers) groupsHandler(w http.ResponseWriter, r *http.Request) {
	groups := h.db.GetGroups()
	pidWriteJSON(w, http.StatusOK, map[string]interface{}{
		"groups": groups,
		"count":  len(groups),
	})
}

func (h *Handlers) groupHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	group := h.db.GetGroup(id)
	if group == nil {
		pidWriteError(w, http.StatusNotFound, "group not found")
		return
	}
	pidWriteJSON(w, http.StatusOK, group)
}

func (h *Handlers) plantHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	plant := h.db.GetPlant(id)
	if plant == nil {
		pidWriteError(w, http.StatusNotFound, "plant not found")
		return
	}
	pidWriteJSON(w, http.StatusOK, plant)
}

func (h *Handlers) searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		pidWriteError(w, http.StatusBadRequest, "q query parameter is required")
		return
	}
	results := h.db.SearchPlants(query)
	pidWriteJSON(w, http.StatusOK, map[string]interface{}{
		"results": results,
		"count":   len(results),
		"query":   query,
	})
}

func pidWriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func pidWriteError(w http.ResponseWriter, status int, message string) {
	pidWriteJSON(w, status, map[string]interface{}{"error": message, "status": status})
}

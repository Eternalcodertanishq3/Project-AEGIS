package medical

import (
	"encoding/json"
	"net/http"
)

// Handlers holds the HTTP handler functions for the medical module.
type Handlers struct {
	db *MedicalDB
}

// NewHandlers creates new medical API handlers.
func NewHandlers(db *MedicalDB) *Handlers {
	return &Handlers{db: db}
}

// RegisterRoutes registers all medical module routes.
func (h *Handlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/medical/categories", h.categoriesHandler)
	mux.HandleFunc("GET /api/medical/categories/{id}", h.categoryHandler)
	mux.HandleFunc("GET /api/medical/entries/{id}", h.entryHandler)
}

func (h *Handlers) categoriesHandler(w http.ResponseWriter, r *http.Request) {
	cats := h.db.GetCategories()
	medWriteJSON(w, http.StatusOK, map[string]interface{}{
		"categories": cats,
		"count":      len(cats),
	})
}

func (h *Handlers) categoryHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cat := h.db.GetCategory(id)
	if cat == nil {
		medWriteError(w, http.StatusNotFound, "category not found")
		return
	}
	medWriteJSON(w, http.StatusOK, cat)
}

func (h *Handlers) entryHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	entry := h.db.GetEntry(id)
	if entry == nil {
		medWriteError(w, http.StatusNotFound, "entry not found")
		return
	}
	medWriteJSON(w, http.StatusOK, entry)
}

func medWriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func medWriteError(w http.ResponseWriter, status int, message string) {
	medWriteJSON(w, status, map[string]interface{}{"error": message, "status": status})
}

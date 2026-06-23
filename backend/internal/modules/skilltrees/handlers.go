package skilltrees

import (
	"encoding/json"
	"net/http"
)

// Handlers holds the HTTP handler functions for the skill trees module.
type Handlers struct {
	db *SkillTreeDB
}

// NewHandlers creates new skill trees API handlers.
func NewHandlers(db *SkillTreeDB) *Handlers {
	return &Handlers{db: db}
}

// RegisterRoutes registers all skill trees module routes.
func (h *Handlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/skills/categories", h.categoriesHandler)
	mux.HandleFunc("GET /api/skills/categories/{id}", h.categoryHandler)
	mux.HandleFunc("GET /api/skills/{id}", h.skillHandler)
}

func (h *Handlers) categoriesHandler(w http.ResponseWriter, r *http.Request) {
	cats := h.db.GetCategories()
	stWriteJSON(w, http.StatusOK, map[string]interface{}{
		"categories": cats,
		"count":      len(cats),
	})
}

func (h *Handlers) categoryHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cat := h.db.GetCategory(id)
	if cat == nil {
		stWriteError(w, http.StatusNotFound, "category not found")
		return
	}
	stWriteJSON(w, http.StatusOK, cat)
}

func (h *Handlers) skillHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	skill := h.db.GetSkill(id)
	if skill == nil {
		stWriteError(w, http.StatusNotFound, "skill not found")
		return
	}
	stWriteJSON(w, http.StatusOK, skill)
}

func stWriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func stWriteError(w http.ResponseWriter, status int, message string) {
	stWriteJSON(w, status, map[string]interface{}{"error": message, "status": status})
}

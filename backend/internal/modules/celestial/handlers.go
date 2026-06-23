package celestial

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// Handlers holds the HTTP handler functions for the celestial navigation module.
type Handlers struct{}

// NewHandlers creates new celestial navigation API handlers.
func NewHandlers() *Handlers {
	return &Handlers{}
}

// RegisterRoutes registers all celestial navigation module routes.
func (h *Handlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/celestial/calculate", h.calculateHandler)
	mux.HandleFunc("GET /api/celestial/stars", h.starsHandler)
	mux.HandleFunc("GET /api/celestial/techniques", h.techniquesHandler)
}

func (h *Handlers) calculateHandler(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")
	timeStr := r.URL.Query().Get("time")

	if latStr == "" || lonStr == "" {
		celWriteError(w, http.StatusBadRequest, "lat and lon query parameters are required")
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		celWriteError(w, http.StatusBadRequest, "invalid latitude")
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		celWriteError(w, http.StatusBadRequest, "invalid longitude")
		return
	}

	t := time.Now().UTC()
	if timeStr != "" {
		parsed, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			celWriteError(w, http.StatusBadRequest, "invalid time format, use RFC3339")
			return
		}
		t = parsed
	}

	result := Calculate(lat, lon, t)
	celWriteJSON(w, http.StatusOK, result)
}

func (h *Handlers) starsHandler(w http.ResponseWriter, r *http.Request) {
	stars := GetNavigationStars()
	celWriteJSON(w, http.StatusOK, map[string]interface{}{
		"stars": stars,
		"count": len(stars),
	})
}

func (h *Handlers) techniquesHandler(w http.ResponseWriter, r *http.Request) {
	techs := GetTechniques()
	celWriteJSON(w, http.StatusOK, map[string]interface{}{
		"techniques": techs,
		"count":      len(techs),
	})
}

func celWriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func celWriteError(w http.ResponseWriter, status int, message string) {
	celWriteJSON(w, status, map[string]interface{}{"error": message, "status": status})
}

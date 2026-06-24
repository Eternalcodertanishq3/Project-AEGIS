package beacon

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Handlers holds HTTP handlers for the position beacon module.
type Handlers struct {
	mgr *BeaconManager
}

// NewHandlers creates new beacon API handlers.
func NewHandlers(mgr *BeaconManager) *Handlers {
	return &Handlers{mgr: mgr}
}

// RegisterRoutes registers all beacon routes.
func (h *Handlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/beacon/status", h.statusHandler)
	mux.HandleFunc("GET /api/beacon/positions", h.positionsHandler)
	mux.HandleFunc("POST /api/beacon/positions", h.logPositionHandler)
	mux.HandleFunc("POST /api/beacon/aprs", h.aprsHandler)
	mux.HandleFunc("GET /api/beacon/distance", h.distanceHandler)
}

func (h *Handlers) statusHandler(w http.ResponseWriter, r *http.Request) {
	beaconWriteJSON(w, http.StatusOK, h.mgr.GetStatus())
}

func (h *Handlers) positionsHandler(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil {
			limit = n
		}
	}
	positions := h.mgr.GetPositions(limit)
	beaconWriteJSON(w, http.StatusOK, map[string]interface{}{
		"positions": positions,
		"count":     len(positions),
	})
}

func (h *Handlers) logPositionHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Altitude  float64 `json:"altitude"`
		Accuracy  float64 `json:"accuracy"`
		Source    string  `json:"source"`
		Comment   string  `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		beaconWriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	pos, err := h.mgr.LogPosition(body.Latitude, body.Longitude, body.Altitude, body.Accuracy, body.Source, body.Comment)
	if err != nil {
		beaconWriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	beaconWriteJSON(w, http.StatusCreated, pos)
}

func (h *Handlers) aprsHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Comment   string  `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		beaconWriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	aprs := h.mgr.GenerateAPRS(body.Latitude, body.Longitude, body.Comment)
	beaconWriteJSON(w, http.StatusOK, aprs)
}

func (h *Handlers) distanceHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	lat1, err1 := strconv.ParseFloat(q.Get("lat1"), 64)
	lon1, err2 := strconv.ParseFloat(q.Get("lon1"), 64)
	lat2, err3 := strconv.ParseFloat(q.Get("lat2"), 64)
	lon2, err4 := strconv.ParseFloat(q.Get("lon2"), 64)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		beaconWriteError(w, http.StatusBadRequest, "lat1, lon1, lat2, lon2 query parameters are required")
		return
	}
	result := CalculateDistance(lat1, lon1, lat2, lon2)
	beaconWriteJSON(w, http.StatusOK, result)
}

func beaconWriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func beaconWriteError(w http.ResponseWriter, status int, message string) {
	beaconWriteJSON(w, status, map[string]interface{}{"error": message, "status": status})
}

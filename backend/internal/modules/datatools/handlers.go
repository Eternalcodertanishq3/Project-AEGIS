package datatools

import (
	"encoding/json"
	"net/http"
)

// Handlers holds the HTTP handler functions for the data tools module.
type Handlers struct{}

// NewHandlers creates new data tools API handlers.
func NewHandlers() *Handlers {
	return &Handlers{}
}

// RegisterRoutes registers all data tools module routes.
func (h *Handlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/datatools/operations", h.operationsHandler)
	mux.HandleFunc("POST /api/datatools/transform", h.transformHandler)
}

func (h *Handlers) operationsHandler(w http.ResponseWriter, r *http.Request) {
	ops := ListOperations()
	dtWriteJSON(w, http.StatusOK, map[string]interface{}{
		"operations": ops,
	})
}

type transformRequest struct {
	Operation string `json:"operation"`
	Input     string `json:"input"`
}

func (h *Handlers) transformHandler(w http.ResponseWriter, r *http.Request) {
	var req transformRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dtWriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Operation == "" {
		dtWriteError(w, http.StatusBadRequest, "operation is required")
		return
	}

	result := Transform(req.Operation, req.Input)
	if result.Error != "" {
		dtWriteJSON(w, http.StatusUnprocessableEntity, result)
		return
	}
	dtWriteJSON(w, http.StatusOK, result)
}

func dtWriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func dtWriteError(w http.ResponseWriter, status int, message string) {
	dtWriteJSON(w, status, map[string]interface{}{"error": message, "status": status})
}

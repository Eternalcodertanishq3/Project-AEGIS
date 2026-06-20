package api

import (
	"encoding/json"
	"log"
	"net/http"

	"aegis/backend/internal/modules/ai"
	"aegis/backend/internal/modules/knowledge"
	"aegis/backend/internal/modules/maps"
	"aegis/backend/internal/orchestrator"
	"aegis/backend/internal/powermanager"
	"aegis/backend/internal/resourceprofiler"
)

// Deps holds all injected dependencies for API handlers.
type Deps struct {
	Profiler          resourceprofiler.Profiler
	PowerManager      powermanager.PowerManager
	Orchestrator      *orchestrator.Orchestrator
	KnowledgeHandlers *knowledge.Handlers
	MapsHandlers      *maps.Handlers
	AIHandlers        *ai.Handlers
}

// healthHandler returns system health.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": "0.1.0",
	})
}

// systemProfileHandler returns the detected hardware profile.
func systemProfileHandler(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profile, err := deps.Profiler.DetectProfile()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to detect profile: "+err.Error())
			return
		}
		writeJSON(w, http.StatusOK, profile)
	}
}

// systemPowerHandler returns the current power status.
func systemPowerHandler(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := deps.PowerManager.GetPowerStatus()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to get power status: "+err.Error())
			return
		}
		writeJSON(w, http.StatusOK, status)
	}
}

// modulesListHandler returns all registered modules.
func modulesListHandler(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		modules := deps.Orchestrator.ListModules()
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"modules": modules,
			"count":   len(modules),
		})
	}
}

// moduleEnableHandler enables a module by ID.
func moduleEnableHandler(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			writeError(w, http.StatusBadRequest, "module id is required")
			return
		}

		if err := deps.Orchestrator.EnableModule(id); err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		mod, _ := deps.Orchestrator.GetModule(id)
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"message": "module enabled",
			"module":  mod,
		})
	}
}

// moduleDisableHandler disables a module by ID.
func moduleDisableHandler(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			writeError(w, http.StatusBadRequest, "module id is required")
			return
		}

		if err := deps.Orchestrator.DisableModule(id); err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		mod, _ := deps.Orchestrator.GetModule(id)
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"message": "module disabled",
			"module":  mod,
		})
	}
}

// pluginsListHandler returns the list of installed plugins (empty for Phase 0).
func pluginsListHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"plugins": []interface{}{},
		"count":   0,
	})
}

// pluginsInstallHandler is a stub that returns 501 Not Implemented.
func pluginsInstallHandler(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "plugin installation is not yet implemented")
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error encoding JSON response: %v", err)
	}
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]interface{}{
		"error":  message,
		"status": status,
	})
}

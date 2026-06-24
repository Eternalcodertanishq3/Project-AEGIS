package knowledge

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Handlers holds the HTTP handler functions for the knowledge module.
type Handlers struct {
	kiwix *KiwixManager
}

// NewHandlers creates new knowledge API handlers.
func NewHandlers(kiwix *KiwixManager) *Handlers {
	return &Handlers{kiwix: kiwix}
}

// RegisterRoutes registers all knowledge module routes on the given mux.
func (h *Handlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/knowledge/status", h.statusHandler)
	mux.HandleFunc("GET /api/knowledge/zim", h.zimListHandler)
	mux.HandleFunc("GET /api/knowledge/search", h.searchHandler)
	mux.HandleFunc("GET /api/knowledge/article/{rest...}", h.articleHandler)
}

// statusHandler returns the status of the knowledge module.
func (h *Handlers) statusHandler(w http.ResponseWriter, r *http.Request) {
	zimFiles := h.kiwix.GetZIMFiles()
	sidecarPath, sidecarErr := h.kiwix.FindSidecar()

	status := map[string]interface{}{
		"module":          "knowledge",
		"kiwix_running":   h.kiwix.IsRunning(),
		"kiwix_port":      h.kiwix.Port(),
		"zim_files_count": len(zimFiles),
		"zim_files":       zimFiles,
	}

	if sidecarErr != nil {
		status["kiwix_sidecar"] = "not found"
		status["kiwix_note"] = "Place kiwix-serve in sidecars/kiwix-serve/<os>/"
	} else {
		status["kiwix_sidecar"] = sidecarPath
		if lastErr := h.kiwix.GetLastError(); lastErr != "" {
			status["kiwix_note"] = lastErr
		}
	}

	writeJSON(w, http.StatusOK, status)
}

// zimListHandler returns discovered ZIM files.
func (h *Handlers) zimListHandler(w http.ResponseWriter, r *http.Request) {
	zimFiles := h.kiwix.GetZIMFiles()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"zim_files": zimFiles,
		"count":     len(zimFiles),
	})
}

// searchHandler handles knowledge search requests.
func (h *Handlers) searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if n, err := strconv.Atoi(limitStr); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	results, err := h.kiwix.ProxySearch(query, limit)
	if err != nil {
		log.Printf("knowledge search error: %v", err)
		writeError(w, http.StatusInternalServerError, "search failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"query":   query,
		"results": results,
		"count":   len(results),
	})
}

// articleHandler proxies article requests to kiwix-serve.
func (h *Handlers) articleHandler(w http.ResponseWriter, r *http.Request) {
	// URL pattern: /api/knowledge/article/{zimId}/{path...}
	rest := r.PathValue("rest")
	if rest == "" {
		writeError(w, http.StatusBadRequest, "article path is required")
		return
	}

	parts := strings.SplitN(rest, "/", 2)
	zimID := parts[0]
	articlePath := ""
	if len(parts) > 1 {
		articlePath = parts[1]
	}

	body, contentType, err := h.kiwix.ProxyArticle(zimID, articlePath)
	if err != nil {
		log.Printf("knowledge article error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch article")
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]interface{}{
		"error":  message,
		"status": status,
	})
}

package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Handlers holds the HTTP handler functions for the AI module.
type Handlers struct {
	ai *AIManager
}

// NewHandlers creates new AI API handlers.
func NewHandlers(ai *AIManager) *Handlers {
	return &Handlers{ai: ai}
}

// RegisterRoutes registers all AI module routes on the given mux.
func (h *Handlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/ai/status", h.statusHandler)
	mux.HandleFunc("GET /api/ai/models", h.modelsHandler)
	mux.HandleFunc("POST /api/ai/start", h.startHandler)
	mux.HandleFunc("POST /api/ai/stop", h.stopHandler)
	mux.HandleFunc("POST /api/ai/chat", h.chatHandler)
}

// statusHandler returns the status of the AI module.
func (h *Handlers) statusHandler(w http.ResponseWriter, r *http.Request) {
	models := h.ai.GetModels()
	sidecarPath, sidecarErr := h.ai.FindSidecar()

	status := map[string]interface{}{
		"module":       "ai",
		"running":      h.ai.IsRunning(),
		"port":         h.ai.Port(),
		"active_model": h.ai.ActiveModel(),
		"models_count": len(models),
		"models":       models,
	}

	if sidecarErr != nil {
		status["sidecar"] = "not found"
		status["sidecar_note"] = "Place llama-server in sidecars/llama/<os>/"
	} else {
		status["sidecar"] = sidecarPath
	}

	aiWriteJSON(w, http.StatusOK, status)
}

// modelsHandler returns discovered model files.
func (h *Handlers) modelsHandler(w http.ResponseWriter, r *http.Request) {
	models := h.ai.GetModels()
	aiWriteJSON(w, http.StatusOK, map[string]interface{}{
		"models": models,
		"count":  len(models),
	})
}

// startRequest is the JSON body for the start endpoint.
type startRequest struct {
	ModelID string `json:"model_id"`
}

// startHandler launches llama-server with a specific model.
func (h *Handlers) startHandler(w http.ResponseWriter, r *http.Request) {
	var req startRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		aiWriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.ModelID == "" {
		aiWriteError(w, http.StatusBadRequest, "model_id is required")
		return
	}

	if err := h.ai.Start(req.ModelID); err != nil {
		aiWriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	aiWriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":      "AI server started",
		"model":        req.ModelID,
		"port":         h.ai.Port(),
	})
}

// stopHandler stops the running llama-server.
func (h *Handlers) stopHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.ai.Stop(); err != nil {
		aiWriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	aiWriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "AI server stopped",
	})
}

// chatRequest is the JSON body for the chat endpoint.
type chatRequest struct {
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

// chatMessage is a single message in the chat.
type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatHandler proxies chat requests to llama-server's OpenAI-compatible endpoint.
func (h *Handlers) chatHandler(w http.ResponseWriter, r *http.Request) {
	if !h.ai.IsRunning() {
		aiWriteError(w, http.StatusServiceUnavailable, "AI server is not running. Load a model first.")
		return
	}

	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		aiWriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if len(req.Messages) == 0 {
		aiWriteError(w, http.StatusBadRequest, "messages array is required")
		return
	}

	// Build the OpenAI-compatible request for llama-server
	llamaReq := map[string]interface{}{
		"messages":    req.Messages,
		"stream":      req.Stream,
		"temperature": 0.7,
		"max_tokens":  1024,
	}

	body, err := json.Marshal(llamaReq)
	if err != nil {
		aiWriteError(w, http.StatusInternalServerError, "failed to encode request")
		return
	}

	llamaURL := fmt.Sprintf("http://127.0.0.1:%d/v1/chat/completions", h.ai.Port())

	client := &http.Client{Timeout: 120 * time.Second}
	llamaResp, err := client.Post(llamaURL, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("AI chat proxy error: %v", err)
		aiWriteError(w, http.StatusBadGateway, "failed to reach AI server: "+err.Error())
		return
	}
	defer llamaResp.Body.Close()

	// Stream the response back to the client
	if req.Stream {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")
		w.WriteHeader(http.StatusOK)

		flusher, ok := w.(http.Flusher)

		buf := make([]byte, 1024)
		for {
			n, readErr := llamaResp.Body.Read(buf)
			if n > 0 {
				w.Write(buf[:n])
				if ok {
					flusher.Flush()
				}
			}
			if readErr != nil {
				break
			}
		}
		return
	}

	// Non-streaming: forward the complete response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(llamaResp.StatusCode)
	io.Copy(w, llamaResp.Body)
}

// aiWriteJSON writes a JSON response.
func aiWriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// aiWriteError writes a JSON error response.
func aiWriteError(w http.ResponseWriter, status int, message string) {
	aiWriteJSON(w, status, map[string]interface{}{
		"error":  message,
		"status": status,
	})
}
